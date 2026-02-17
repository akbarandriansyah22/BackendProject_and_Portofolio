package security

import (
	"sync"
	"time"
)

// RateLimiter implements token bucket algorithm for rate limiting
type RateLimiter struct {
	mu         sync.RWMutex
	buckets    map[string]*TokenBucket
	capacity   int64
	refillRate time.Duration
}

// TokenBucket represents a user's rate limit bucket
type TokenBucket struct {
	tokens     int64
	lastRefill time.Time
	capacity   int64
}

// NewRateLimiter creates a new rate limiter
// capacity: max requests per window
// refillRate: time window for refilling tokens
func NewRateLimiter(capacity int64, refillRate time.Duration) *RateLimiter {
	return &RateLimiter{
		buckets:    make(map[string]*TokenBucket),
		capacity:   capacity,
		refillRate: refillRate,
	}
}

// IsAllowed checks if user can make a request
func (rl *RateLimiter) IsAllowed(userID string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	bucket, exists := rl.buckets[userID]
	if !exists {
		bucket = &TokenBucket{
			tokens:     rl.capacity,
			lastRefill: time.Now(),
			capacity:   rl.capacity,
		}
		rl.buckets[userID] = bucket
	}

	now := time.Now()
	timePassed := now.Sub(bucket.lastRefill)

	if timePassed >= rl.refillRate {
		bucket.tokens = rl.capacity
		bucket.lastRefill = now
	}

	if bucket.tokens > 0 {
		bucket.tokens--
		return true
	}

	return false
}

// GetRemaining returns remaining tokens for user
func (rl *RateLimiter) GetRemaining(userID string) int64 {
	rl.mu.RLock()
	defer rl.mu.RUnlock()

	if bucket, exists := rl.buckets[userID]; exists {
		return bucket.tokens
	}
	return rl.capacity
}

// Reset clears all rate limit data
func (rl *RateLimiter) Reset() {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	rl.buckets = make(map[string]*TokenBucket)
}
