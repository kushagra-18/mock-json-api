package services

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
)

// RedisService handles operations with Redis.
type RedisService struct {
	Client *redis.Client
	ctx    context.Context
}

// NewRedisService creates a new RedisService.
func NewRedisService(client *redis.Client) *RedisService {
	return &RedisService{
		Client: client,
		ctx:    context.Background(), // Or a more specific context if needed
	}
}

// CreateRedisKey generates a Redis key by joining parts with a colon.
func (s *RedisService) CreateRedisKey(parts ...string) string {
	return strings.Join(parts, ":")
}

// RateLimit implements a fixed window rate limiting algorithm.
// It returns true if the request is rate-limited (i.e., exceeds the limit), false otherwise.
func (s *RedisService) RateLimit(key string, limit int, windowSeconds int64) (bool, error) {
	// Increment the counter for the key
	count, err := s.Client.Incr(s.ctx, key).Result()
	if err != nil {
		return true, fmt.Errorf("failed to increment rate limit counter: %w", err)
	}

	// If this is the first request in the window, set the expiration
	if count == 1 {
		if err := s.Client.Expire(s.ctx, key, time.Duration(windowSeconds)*time.Second).Err(); err != nil {
			// If setting expire fails, it's safer to assume rate limited or handle cleanup
			return true, fmt.Errorf("failed to set expiration for rate limit key: %w", err)
		}
	}

	// Check if the count exceeds the limit
	if count > int64(limit) {
		return true, nil // Rate limited
	}

	return false, nil // Not rate limited
}

// GetValue retrieves a value from Redis.
func (s *RedisService) GetValue(key string) (string, error) {
	val, err := s.Client.Get(s.ctx, key).Result()
	if err == redis.Nil {
		return "", nil // Key does not exist
	} else if err != nil {
		return "", fmt.Errorf("failed to get value from redis: %w", err)
	}
	return val, nil
}

// SetValue sets a value in Redis with an optional expiration.
// If expiration is 0, the key does not expire.
func (s *RedisService) SetValue(key string, value interface{}, expiration time.Duration) error {
	err := s.Client.Set(s.ctx, key, value, expiration).Err()
	if err != nil {
		return fmt.Errorf("failed to set value in redis: %w", err)
	}
	return nil
}

// DeleteValue deletes a key from Redis.
func (s *RedisService) DeleteValue(key string) error {
	err := s.Client.Del(s.ctx, key).Err()
	if err != nil {
		return fmt.Errorf("failed to delete value from redis: %w", err)
	}
	return nil
}
