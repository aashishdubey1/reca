package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"reca-backend/internal/model"
	"reca-backend/internal/repository"
	"reca-backend/internal/srs"
)

// Default demo user ID for initial development until full auth is wired up
var DemoUserID = uuid.MustParse("00000000-0000-0000-0000-000000000001")

type TopicHandler struct {
	repo  *repository.TopicRepository
	redis *repository.RedisClient
}

func NewTopicHandler(repo *repository.TopicRepository, redis *repository.RedisClient) *TopicHandler {
	return &TopicHandler{repo: repo, redis: redis}
}

// CreateTopic handles POST /api/v1/topics
func (h *TopicHandler) CreateTopic(w http.ResponseWriter, r *http.Request) {
	var req model.CreateTopicRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.Title == "" {
		http.Error(w, "title is required", http.StatusBadRequest)
		return
	}

	topic, err := h.repo.CreateTopic(r.Context(), DemoUserID, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Invalidate timeline cache
	h.redis.InvalidateUserCache(r.Context(), DemoUserID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(topic)
}

// GetDueRevisions handles GET /api/v1/reviews/due
func (h *TopicHandler) GetDueRevisions(w http.ResponseWriter, r *http.Request) {
	topics, err := h.repo.GetDueRevisions(r.Context(), DemoUserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if topics == nil {
		topics = []model.Topic{}
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(topics)
}

// SubmitReview handles POST /api/v1/reviews/{topicID}
func (h *TopicHandler) SubmitReview(w http.ResponseWriter, r *http.Request) {
	topicIDStr := chi.URLParam(r, "topicID")
	topicID, err := uuid.Parse(topicIDStr)
	if err != nil {
		http.Error(w, "invalid topic ID", http.StatusBadRequest)
		return
	}

	var req model.ReviewRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	rating, err := srs.ParseRating(req.Rating)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	res, err := h.repo.RecordReview(r.Context(), DemoUserID, topicID, rating)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Invalidate timeline cache and record daily activity
	h.redis.InvalidateUserCache(r.Context(), DemoUserID)
	h.redis.RecordActivity(r.Context(), DemoUserID, time.Now())

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(res)
}

// GetTimeline handles GET /api/v1/timeline
func (h *TopicHandler) GetTimeline(w http.ResponseWriter, r *http.Request) {
	// Try Redis cache first
	if cached, ok := h.redis.GetTimelineCache(r.Context(), DemoUserID); ok {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-Cache", "HIT")
		_ = json.NewEncoder(w).Encode(cached)
		return
	}

	timeline, err := h.repo.GetTimeline(r.Context(), DemoUserID, 3, 20)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Cache in Redis for 10 minutes
	h.redis.SetTimelineCache(r.Context(), DemoUserID, timeline, 10*time.Minute)

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Cache", "MISS")
	_ = json.NewEncoder(w).Encode(timeline)
}
