package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"reca-backend/internal/model"
	"reca-backend/internal/srs"
)

type TopicRepository struct {
	db *PostgresDB
}

func NewTopicRepository(db *PostgresDB) *TopicRepository {
	return &TopicRepository{db: db}
}

// CreateTopic inserts a new topic and creates its initial SRS item (due immediately).
func (r *TopicRepository) CreateTopic(ctx context.Context, userID uuid.UUID, req model.CreateTopicRequest) (*model.Topic, error) {
	tx, err := r.db.Pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	now := time.Now().UTC()
	topicID := uuid.New()

	if req.Category == "" {
		req.Category = "General"
	}
	if req.Tone == "" {
		req.Tone = "orange"
	}

	queryTopic := `
		INSERT INTO topics (id, user_id, title, category, tone, notes, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $7)
		RETURNING id, user_id, title, category, tone, notes, created_at, updated_at
	`
	var topic model.Topic
	err = tx.QueryRow(ctx, queryTopic, topicID, userID, req.Title, req.Category, req.Tone, req.Notes, now).Scan(
		&topic.ID, &topic.UserID, &topic.Title, &topic.Category, &topic.Tone, &topic.Notes, &topic.CreatedAt, &topic.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to insert topic: %w", err)
	}

	initialState := srs.DefaultState(now)
	srsID := uuid.New()
	querySRS := `
		INSERT INTO srs_items (id, topic_id, user_id, repetitions, interval_days, ease_factor, last_reviewed_at, next_review_at, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, NULL, $7, $8)
		RETURNING id, topic_id, user_id, repetitions, interval_days, ease_factor, last_reviewed_at, next_review_at
	`
	var srsItem model.SRSItem
	err = tx.QueryRow(ctx, querySRS,
		srsID, topic.ID, userID, initialState.Repetitions, initialState.IntervalDays, initialState.EaseFactor, initialState.NextReviewAt, now,
	).Scan(
		&srsItem.ID, &srsItem.TopicID, &srsItem.UserID, &srsItem.Repetitions, &srsItem.IntervalDays, &srsItem.EaseFactor, &srsItem.LastReviewedAt, &srsItem.NextReviewAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to insert srs item: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit tx: %w", err)
	}

	topic.SRS = &srsItem
	return &topic, nil
}

// GetDueRevisions returns all topics where next_review_at <= NOW().
func (r *TopicRepository) GetDueRevisions(ctx context.Context, userID uuid.UUID) ([]model.Topic, error) {
	query := `
		SELECT t.id, t.user_id, t.title, t.category, t.tone, t.notes, t.created_at, t.updated_at,
		       s.id, s.topic_id, s.user_id, s.repetitions, s.interval_days, s.ease_factor, s.last_reviewed_at, s.next_review_at
		FROM topics t
		JOIN srs_items s ON t.id = s.topic_id
		WHERE s.user_id = $1 AND s.next_review_at <= NOW()
		ORDER BY s.next_review_at ASC
	`
	rows, err := r.db.Pool.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var topics []model.Topic
	for rows.Next() {
		var t model.Topic
		var s model.SRSItem
		if err := rows.Scan(
			&t.ID, &t.UserID, &t.Title, &t.Category, &t.Tone, &t.Notes, &t.CreatedAt, &t.UpdatedAt,
			&s.ID, &s.TopicID, &s.UserID, &s.Repetitions, &s.IntervalDays, &s.EaseFactor, &s.LastReviewedAt, &s.NextReviewAt,
		); err != nil {
			return nil, err
		}
		t.SRS = &s
		topics = append(topics, t)
	}

	return topics, rows.Err()
}

// RecordReview executes SM-2 calculation atomically inside a transaction.
func (r *TopicRepository) RecordReview(ctx context.Context, userID, topicID uuid.UUID, rating srs.Rating) (*model.ReviewResponse, error) {
	tx, err := r.db.Pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	// Fetch current SRS state with row lock
	var srsItem model.SRSItem
	queryLock := `
		SELECT id, topic_id, user_id, repetitions, interval_days, ease_factor, last_reviewed_at, next_review_at
		FROM srs_items
		WHERE topic_id = $1 AND user_id = $2
		FOR UPDATE
	`
	err = tx.QueryRow(ctx, queryLock, topicID, userID).Scan(
		&srsItem.ID, &srsItem.TopicID, &srsItem.UserID, &srsItem.Repetitions, &srsItem.IntervalDays, &srsItem.EaseFactor, &srsItem.LastReviewedAt, &srsItem.NextReviewAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("topic not found or access denied")
		}
		return nil, err
	}

	now := time.Now().UTC()
	currentState := srs.ItemState{
		Repetitions:  srsItem.Repetitions,
		IntervalDays: srsItem.IntervalDays,
		EaseFactor:   srsItem.EaseFactor,
		NextReviewAt: srsItem.NextReviewAt,
	}

	// Apply SM-2 engine
	nextState := srs.Calculate(currentState, rating, now)

	// Update SRS item
	queryUpdate := `
		UPDATE srs_items
		SET repetitions = $1, interval_days = $2, ease_factor = $3, last_reviewed_at = $4, next_review_at = $5
		WHERE id = $6
	`
	_, err = tx.Exec(ctx, queryUpdate,
		nextState.Repetitions, nextState.IntervalDays, nextState.EaseFactor, now, nextState.NextReviewAt, srsItem.ID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update srs item: %w", err)
	}

	// Insert audit review log
	logID := uuid.New()
	queryLog := `
		INSERT INTO review_logs (id, srs_item_id, topic_id, user_id, rating, interval_before, interval_after, ease_factor_after, reviewed_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	_, err = tx.Exec(ctx, queryLog,
		logID, srsItem.ID, topicID, userID, int(rating), srsItem.IntervalDays, nextState.IntervalDays, nextState.EaseFactor, now,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to insert review log: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit review tx: %w", err)
	}

	return &model.ReviewResponse{
		TopicID:      topicID,
		Repetitions:  nextState.Repetitions,
		IntervalDays: nextState.IntervalDays,
		EaseFactor:   nextState.EaseFactor,
		NextReviewAt: nextState.NextReviewAt,
		Message:      fmt.Sprintf("Next review scheduled in %d days", nextState.IntervalDays),
	}, nil
}

// GetTimeline constructs day cards (past completed reviews, today, and future projected reviews).
func (r *TopicRepository) GetTimeline(ctx context.Context, userID uuid.UUID, pastDays, futureDays int) ([]model.TimelineDay, error) {
	now := time.Now().UTC()
	startOfToday := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)

	var timeline []model.TimelineDay

	for offset := -pastDays; offset <= futureDays; offset++ {
		targetDate := startOfToday.AddDate(0, 0, offset)
		nextDay := targetDate.AddDate(0, 0, 1)

		side := "right"
		label := fmt.Sprintf("In %d days", offset)
		if offset < 0 {
			side = "left"
			if offset == -1 {
				label = "Yesterday"
			} else {
				label = fmt.Sprintf("%d days ago", -offset)
			}
		} else if offset == 0 {
			side = "center"
			label = "Today"
		} else if offset == 1 {
			label = "Tomorrow"
		}

		day := model.TimelineDay{
			Offset: offset,
			Label:  label,
			Date:   targetDate.Format("Jan 2, 2006"),
			Side:   side,
			Topics: []string{},
		}

		if offset < 0 {
			// Query reviews completed on this past day
			q := `
				SELECT t.title
				FROM review_logs l
				JOIN topics t ON l.topic_id = t.id
				WHERE l.user_id = $1 AND l.reviewed_at >= $2 AND l.reviewed_at < $3
				ORDER BY l.reviewed_at DESC
				LIMIT 10
			`
			rows, err := r.db.Pool.Query(ctx, q, userID, targetDate, nextDay)
			if err == nil {
				for rows.Next() {
					var title string
					if err := rows.Scan(&title); err == nil {
						day.Topics = append(day.Topics, title)
					}
				}
				rows.Close()
			}
			day.Completed = len(day.Topics)
			day.Total = len(day.Topics)
		} else if offset == 0 {
			// Today: Due items
			q := `
				SELECT t.title
				FROM srs_items s
				JOIN topics t ON s.topic_id = t.id
				WHERE s.user_id = $1 AND s.next_review_at <= $2
				ORDER BY s.next_review_at ASC
				LIMIT 12
			`
			rows, err := r.db.Pool.Query(ctx, q, userID, nextDay)
			if err == nil {
				for rows.Next() {
					var title string
					if err := rows.Scan(&title); err == nil {
						day.Topics = append(day.Topics, title)
					}
				}
				rows.Close()
			}
			day.Total = len(day.Topics)
		} else {
			// Future: Projected due items
			q := `
				SELECT t.title
				FROM srs_items s
				JOIN topics t ON s.topic_id = t.id
				WHERE s.user_id = $1 AND s.next_review_at >= $2 AND s.next_review_at < $3
				ORDER BY s.next_review_at ASC
				LIMIT 8
			`
			rows, err := r.db.Pool.Query(ctx, q, userID, targetDate, nextDay)
			if err == nil {
				for rows.Next() {
					var title string
					if err := rows.Scan(&title); err == nil {
						day.Topics = append(day.Topics, title)
					}
				}
				rows.Close()
			}
			day.Total = len(day.Topics)
		}

		timeline = append(timeline, day)
	}

	return timeline, nil
}
