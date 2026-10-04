package srs

import (
	"errors"
	"math"
	"strings"
	"time"
)

// Rating represents the user's subjective assessment of recall difficulty.
type Rating int

const (
	RatingAgain Rating = 1 // Forgot or failed recall
	RatingHard  Rating = 3 // Recalled with significant difficulty
	RatingGood  Rating = 4 // Recalled with minor hesitation (standard success)
	RatingEasy  Rating = 5 // Instant, effortless recall
)

// ParseRating converts string ratings ("again", "hard", "medium", "good", "easy") into a typed Rating.
func ParseRating(val string) (Rating, error) {
	switch strings.ToLower(strings.TrimSpace(val)) {
	case "again", "fail", "1":
		return RatingAgain, nil
	case "hard", "3":
		return RatingHard, nil
	case "medium", "good", "4":
		return RatingGood, nil
	case "easy", "5":
		return RatingEasy, nil
	default:
		return 0, errors.New("invalid rating: must be 'again', 'hard', 'medium'/'good', or 'easy'")
	}
}

// ItemState holds the spaced repetition parameters for a topic or card.
type ItemState struct {
	Repetitions  int       `json:"repetitions"`
	IntervalDays int       `json:"interval_days"`
	EaseFactor   float64   `json:"ease_factor"`
	NextReviewAt time.Time `json:"next_review_at"`
}

// DefaultState returns the initial state for a newly created topic.
func DefaultState(now time.Time) ItemState {
	return ItemState{
		Repetitions:  0,
		IntervalDays: 1,
		EaseFactor:   2.50,
		NextReviewAt: now, // Due immediately
	}
}

// Calculate applies the SM-2 algorithm to compute the new interval, ease factor, and due date.
func Calculate(current ItemState, rating Rating, now time.Time) ItemState {
	q := float64(rating)
	next := current

	// 1. Ease Factor adjustment (SM-2 standard formula)
	// EF' = EF + (0.1 - (5 - q) * (0.08 + (5 - q) * 0.02))
	newEF := current.EaseFactor + (0.1 - (5.0-q)*(0.08+(5.0-q)*0.02))
	if newEF < 1.30 {
		newEF = 1.30 // Hard lower bound
	}
	next.EaseFactor = math.Round(newEF*100) / 100

	// 2. Repetition & Interval Calculation
	if rating < RatingHard {
		// Failed recall / Lapse: reset streak
		next.Repetitions = 0
		next.IntervalDays = 1
	} else {
		// Successful recall
		switch current.Repetitions {
		case 0:
			if rating == RatingEasy {
				next.IntervalDays = 2
			} else {
				next.IntervalDays = 1
			}
		case 1:
			if rating == RatingEasy {
				next.IntervalDays = 7
			} else if rating == RatingHard {
				next.IntervalDays = 3
			} else {
				next.IntervalDays = 6
			}
		default:
			// Interval = PreviousInterval * EaseFactor
			multiplier := next.EaseFactor
			if rating == RatingHard {
				multiplier = math.Max(1.2, next.EaseFactor*0.8)
			} else if rating == RatingEasy {
				multiplier = next.EaseFactor * 1.3
			}
			newInterval := int(math.Round(float64(current.IntervalDays) * multiplier))
			if newInterval <= current.IntervalDays {
				newInterval = current.IntervalDays + 1
			}
			next.IntervalDays = newInterval
		}
		next.Repetitions++
	}

	// 3. Compute next review date (midnight or exact offset)
	next.NextReviewAt = now.AddDate(0, 0, next.IntervalDays)

	return next
}
