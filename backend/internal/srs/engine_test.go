package srs

import (
	"testing"
	"time"
)

func TestDefaultState(t *testing.T) {
	now := time.Date(2026, 10, 4, 12, 0, 0, 0, time.UTC)
	state := DefaultState(now)

	if state.Repetitions != 0 {
		t.Errorf("expected 0 repetitions, got %d", state.Repetitions)
	}
	if state.IntervalDays != 1 {
		t.Errorf("expected interval 1, got %d", state.IntervalDays)
	}
	if state.EaseFactor != 2.50 {
		t.Errorf("expected EF 2.50, got %f", state.EaseFactor)
	}
	if !state.NextReviewAt.Equal(now) {
		t.Errorf("expected due at %v, got %v", now, state.NextReviewAt)
	}
}

func TestSM2Progression(t *testing.T) {
	now := time.Date(2026, 10, 4, 12, 0, 0, 0, time.UTC)
	state := DefaultState(now)

	// Review 1: Good
	state = Calculate(state, RatingGood, now)
	if state.Repetitions != 1 {
		t.Errorf("step 1: expected rep 1, got %d", state.Repetitions)
	}
	if state.IntervalDays != 1 {
		t.Errorf("step 1: expected interval 1, got %d", state.IntervalDays)
	}

	// Review 2: Good
	state = Calculate(state, RatingGood, now)
	if state.Repetitions != 2 {
		t.Errorf("step 2: expected rep 2, got %d", state.Repetitions)
	}
	if state.IntervalDays != 6 {
		t.Errorf("step 2: expected interval 6, got %d", state.IntervalDays)
	}

	// Review 3: Easy (should increase interval and EF)
	prevEF := state.EaseFactor
	state = Calculate(state, RatingEasy, now)
	if state.Repetitions != 3 {
		t.Errorf("step 3: expected rep 3, got %d", state.Repetitions)
	}
	if state.EaseFactor <= prevEF {
		t.Errorf("expected EF to increase on Easy, got %f <= %f", state.EaseFactor, prevEF)
	}
	if state.IntervalDays <= 6 {
		t.Errorf("expected interval to grow beyond 6, got %d", state.IntervalDays)
	}

	// Review 4: Fail (Again) - should reset repetitions and interval to 1
	state = Calculate(state, RatingAgain, now)
	if state.Repetitions != 0 {
		t.Errorf("step 4: expected reset to rep 0 on fail, got %d", state.Repetitions)
	}
	if state.IntervalDays != 1 {
		t.Errorf("step 4: expected reset to interval 1 on fail, got %d", state.IntervalDays)
	}
}

func TestEaseFactorLowerBound(t *testing.T) {
	now := time.Now()
	state := ItemState{Repetitions: 5, IntervalDays: 20, EaseFactor: 1.35, NextReviewAt: now}

	// Multiple consecutive failures should not drop EF below 1.30
	for i := 0; i < 5; i++ {
		state = Calculate(state, RatingAgain, now)
	}
	if state.EaseFactor < 1.30 {
		t.Errorf("expected EF not to drop below 1.30, got %f", state.EaseFactor)
	}
}
