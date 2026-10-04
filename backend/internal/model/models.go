package model

import (
	"time"

	"github.com/google/uuid"
	"reca-backend/internal/srs"
)

// Topic represents a subject or unit of knowledge to revise.
type Topic struct {
	ID        uuid.UUID  `json:"id"`
	UserID    uuid.UUID  `json:"user_id"`
	Title     string     `json:"title"`
	Category  string     `json:"category"`
	Tone      string     `json:"tone"` // e.g. "orange", "neutral"
	Notes     string     `json:"notes"`
	SRS       *SRSItem   `json:"srs,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

// SRSItem holds the spaced repetition scheduling state.
type SRSItem struct {
	ID             uuid.UUID  `json:"id"`
	TopicID        uuid.UUID  `json:"topic_id"`
	UserID         uuid.UUID  `json:"user_id"`
	Repetitions    int        `json:"repetitions"`
	IntervalDays   int        `json:"interval_days"`
	EaseFactor     float64    `json:"ease_factor"`
	LastReviewedAt *time.Time `json:"last_reviewed_at"`
	NextReviewAt   time.Time  `json:"next_review_at"`
}

// ReviewLog records an audit of every revision performed.
type ReviewLog struct {
	ID               uuid.UUID  `json:"id"`
	SRSItemID        uuid.UUID  `json:"srs_item_id"`
	TopicID          uuid.UUID  `json:"topic_id"`
	Rating           srs.Rating `json:"rating"`
	IntervalBefore   int        `json:"interval_before"`
	IntervalAfter    int        `json:"interval_after"`
	EaseFactorAfter  float64    `json:"ease_factor_after"`
	ReviewedAt       time.Time  `json:"reviewed_at"`
}

// CreateTopicRequest defines the payload for creating a new revision topic.
type CreateTopicRequest struct {
	Title    string `json:"title"`
	Category string `json:"category"`
	Tone     string `json:"tone"`
	Notes    string `json:"notes"`
}

// ReviewRequest defines the payload when rating a review.
type ReviewRequest struct {
	Rating string `json:"rating"` // "again", "hard", "medium", "easy"
}

// TimelineDay represents one day card on the frontend carousel.
type TimelineDay struct {
	Offset    int      `json:"offset"`
	Label     string   `json:"label"` // e.g. "Today", "Yesterday", "In 3 days"
	Date      string   `json:"date"`  // e.g. "Oct 4, 2026"
	Side      string   `json:"side"`  // "left", "center", "right"
	Total     int      `json:"total"`
	Completed int      `json:"completed"`
	Topics    []string `json:"topics"`
}

// ReviewSummary is returned after completing a revision.
type ReviewResponse struct {
	TopicID      uuid.UUID `json:"topic_id"`
	Repetitions  int       `json:"repetitions"`
	IntervalDays int       `json:"interval_days"`
	EaseFactor   float64   `json:"ease_factor"`
	NextReviewAt time.Time `json:"next_review_at"`
	Message      string    `json:"message"`
}
