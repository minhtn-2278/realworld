package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	TagsListKey = "tags:list"
	feedKeyRoot = "articles:feed:"
)

type Store interface {
	Get(ctx context.Context, key string, destination any) (bool, error)
	Set(ctx context.Context, key string, value any, ttl time.Duration) error
	Delete(ctx context.Context, keys ...string) error
	Close() error
}

type RedisStore struct {
	client *redis.Client
}

func NewRedisStore(ctx context.Context, redisURL string) (*RedisStore, error) {
	options, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, fmt.Errorf("parse Redis URL: %w", err)
	}

	client := redis.NewClient(options)
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("ping Redis: %w", err)
	}

	return &RedisStore{client: client}, nil
}

func (s *RedisStore) Get(ctx context.Context, key string, destination any) (bool, error) {
	value, err := s.client.Get(ctx, key).Bytes()
	if errors.Is(err, redis.Nil) {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("get key %q: %w", key, err)
	}
	if err := json.Unmarshal(value, destination); err != nil {
		return false, fmt.Errorf("unmarshal key %q: %w", key, err)
	}

	return true, nil
}

func (s *RedisStore) Set(ctx context.Context, key string, value any, ttl time.Duration) error {
	valueJSON, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("marshal key %q: %w", key, err)
	}

	if err := s.client.Set(ctx, key, valueJSON, ttl).Err(); err != nil {
		return fmt.Errorf("set key %q: %w", key, err)
	}

	return nil
}

func (s *RedisStore) Delete(ctx context.Context, keys ...string) error {
	if len(keys) == 0 {
		return nil
	}
	if err := s.client.Del(ctx, keys...).Err(); err != nil {
		return fmt.Errorf("delete cache keys: %w", err)
	}

	return nil
}

func (s *RedisStore) Close() error {
	return s.client.Close()
}

func FeedKey(userID uint, limit int, page int) string {
	return fmt.Sprintf("%s%d:limit:%d:page:%d", feedKeyRoot, userID, limit, page)
}
