package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"reca-backend/internal/model"
)

type RedisClient struct {
	rdb *redis.Client
}

func NewRedisClient(addr string) *RedisClient {
	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Printf("[Redis] Warning: unable to connect to Redis at %s (%v). Caching will be disabled/bypassed.", addr, err)
	} else {
		log.Println("[Redis] Connected successfully")
	}

	return &RedisClient{rdb: rdb}
}

func (r *RedisClient) GetTimelineCache(ctx context.Context, userID uuid.UUID) ([]model.TimelineDay, bool) {
	if r.rdb == nil {
		return nil, false
	}
	key := fmt.Sprintf("reca:timeline:%s", userID.String())
	data, err := r.rdb.Get(ctx, key).Bytes()
	if err != nil {
		return nil, false
	}

	var days []model.TimelineDay
	if err := json.Unmarshal(data, &days); err != nil {
		return nil, false
	}
	return days, true
}

func (r *RedisClient) SetTimelineCache(ctx context.Context, userID uuid.UUID, days []model.TimelineDay, ttl time.Duration) {
	if r.rdb == nil {
		return
	}
	key := fmt.Sprintf("reca:timeline:%s", userID.String())
	data, err := json.Marshal(days)
	if err != nil {
		return
	}
	_ = r.rdb.Set(ctx, key, data, ttl).Err()
}

func (r *RedisClient) InvalidateUserCache(ctx context.Context, userID uuid.UUID) {
	if r.rdb == nil {
		return
	}
	key := fmt.Sprintf("reca:timeline:%s", userID.String())
	_ = r.rdb.Del(ctx, key).Err()
}

func (r *RedisClient) RecordActivity(ctx context.Context, userID uuid.UUID, date time.Time) {
	if r.rdb == nil {
		return
	}
	key := fmt.Sprintf("reca:streak:%s", userID.String())
	dayOfYear := date.YearDay()
	_ = r.rdb.SetBit(ctx, key, int64(dayOfYear), 1).Err()
}
