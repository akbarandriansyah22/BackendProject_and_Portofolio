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
	stopOnce   sync.Once     
	stopChan   chan struct{}  
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
		stopChan:   make(chan struct{}),
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
		// Jika window sudah berlalu, token sudah full kembali
		if time.Since(bucket.lastRefill) >= rl.refillRate {
			return rl.capacity
		}
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

func NewRateLimiterWithCleanup(capacity int64, refillRate time.Duration, cleanupInterval time.Duration) *RateLimiter {
	rl := NewRateLimiter(capacity, refillRate)
	rl.StartCleanup(cleanupInterval)
	return rl
}

// StartCleanup menjalankan background goroutine yang secara berkala
// menghapus bucket expired untuk mencegah memory leak.
// Tanpa ini, setiap unique IP yang pernah request akan tersimpan selamanya.
func (rl *RateLimiter) StartCleanup(interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-rl.stopChan:
				return
			case <-ticker.C:
				rl.cleanupExpired()
			}
		}
	}()
}

// Stop menghentikan background cleanup goroutine.
// Panggil saat aplikasi shutdown.
func (rl *RateLimiter) Stop() {
	rl.stopOnce.Do(func() {
		close(rl.stopChan)
	})
}

// cleanupExpired menghapus bucket yang sudah tidak aktif selama 2x window.
func (rl *RateLimiter) cleanupExpired() {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	cutoff := rl.refillRate * 2
	now := time.Now()

	for key, bucket := range rl.buckets {
		if now.Sub(bucket.lastRefill) > cutoff {
			delete(rl.buckets, key)
		}
	}
}

// BucketCount returns jumlah bucket aktif saat ini.
// Berguna untuk monitoring memory usage rate limiter.
func (rl *RateLimiter) BucketCount() int {
	rl.mu.RLock()
	defer rl.mu.RUnlock()
	return len(rl.buckets)
}