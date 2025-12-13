package middleware

import (
	"net/http"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

// RateLimitConfig configures the RateLimit middleware.
type RateLimitConfig struct {
	// Rate is the number of requests allowed per second.
	// Default: 100
	Rate float64

	// Burst is the maximum number of requests allowed in a burst.
	// Default: 10
	Burst int

	// KeyFunc extracts the rate limit key from the request.
	// Default: uses client IP address
	KeyFunc func(r *http.Request) string

	// Handler is called when the rate limit is exceeded.
	// If nil, a default 429 Too Many Requests response is sent.
	Handler http.HandlerFunc

	// SkipFunc determines if rate limiting should be skipped.
	SkipFunc func(r *http.Request) bool

	// CleanupInterval is the interval for cleaning up expired entries.
	// Default: 1 minute
	CleanupInterval time.Duration

	// ExpirationTime is how long to keep entries after last access.
	// Default: 5 minutes
	ExpirationTime time.Duration
}

// DefaultRateLimitConfig returns the default RateLimit configuration.
func DefaultRateLimitConfig() RateLimitConfig {
	return RateLimitConfig{
		Rate:            100,
		Burst:           10,
		KeyFunc:         getClientIP,
		CleanupInterval: time.Minute,
		ExpirationTime:  5 * time.Minute,
	}
}

// RateLimit returns a rate limiting middleware with the given rate and burst.
func RateLimit(rate float64, burst int) Middleware {
	config := DefaultRateLimitConfig()
	config.Rate = rate
	config.Burst = burst
	return RateLimitWithConfig(config)
}

// RateLimitWithConfig returns a RateLimit middleware with the given configuration.
func RateLimitWithConfig(config RateLimitConfig) Middleware {
	if config.Rate <= 0 {
		config.Rate = 100
	}
	if config.Burst <= 0 {
		config.Burst = 10
	}
	if config.KeyFunc == nil {
		config.KeyFunc = getClientIP
	}
	if config.CleanupInterval <= 0 {
		config.CleanupInterval = time.Minute
	}
	if config.ExpirationTime <= 0 {
		config.ExpirationTime = 5 * time.Minute
	}

	store := newRateLimitStore(config)

	// Start cleanup goroutine
	go store.cleanup(config.CleanupInterval, config.ExpirationTime)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Skip if configured
			if config.SkipFunc != nil && config.SkipFunc(r) {
				next.ServeHTTP(w, r)
				return
			}

			key := config.KeyFunc(r)
			limiter := store.get(key, config.Rate, config.Burst)

			if !limiter.Allow() {
				// Rate limit exceeded
				retryAfter := limiter.RetryAfter()

				w.Header().Set("X-RateLimit-Limit", strconv.FormatFloat(config.Rate, 'f', 0, 64))
				w.Header().Set("X-RateLimit-Remaining", "0")
				w.Header().Set("Retry-After", strconv.FormatInt(int64(retryAfter.Seconds()), 10))

				if config.Handler != nil {
					config.Handler(w, r)
				} else {
					w.WriteHeader(http.StatusTooManyRequests)
					w.Write([]byte("Too Many Requests"))
				}
				return
			}

			// Set rate limit headers
			remaining := limiter.Remaining()
			w.Header().Set("X-RateLimit-Limit", strconv.FormatFloat(config.Rate, 'f', 0, 64))
			w.Header().Set("X-RateLimit-Remaining", strconv.Itoa(remaining))

			next.ServeHTTP(w, r)
		})
	}
}

// rateLimitStore stores rate limiters per key.
type rateLimitStore struct {
	mu       sync.RWMutex
	limiters map[string]*tokenBucket
	done     chan struct{}
}

func newRateLimitStore(config RateLimitConfig) *rateLimitStore {
	return &rateLimitStore{
		limiters: make(map[string]*tokenBucket),
		done:     make(chan struct{}),
	}
}

func (s *rateLimitStore) get(key string, rate float64, burst int) *tokenBucket {
	s.mu.RLock()
	limiter, ok := s.limiters[key]
	s.mu.RUnlock()

	if ok {
		limiter.touch()
		return limiter
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	// Double check after acquiring write lock
	if limiter, ok = s.limiters[key]; ok {
		limiter.touch()
		return limiter
	}

	limiter = newTokenBucket(rate, burst)
	s.limiters[key] = limiter
	return limiter
}

func (s *rateLimitStore) cleanup(interval, expiration time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.mu.Lock()
			now := time.Now()
			for key, limiter := range s.limiters {
				if now.Sub(limiter.lastAccess()) > expiration {
					delete(s.limiters, key)
				}
			}
			s.mu.Unlock()
		case <-s.done:
			return
		}
	}
}

// tokenBucket implements the token bucket algorithm.
type tokenBucket struct {
	rate       float64      // tokens per second
	burst      int          // max tokens
	tokens     float64      // current tokens
	lastUpdate time.Time    // last token update
	lastTouch  atomic.Value // time.Time
	mu         sync.Mutex
}

func newTokenBucket(rate float64, burst int) *tokenBucket {
	tb := &tokenBucket{
		rate:       rate,
		burst:      burst,
		tokens:     float64(burst),
		lastUpdate: time.Now(),
	}
	tb.lastTouch.Store(time.Now())
	return tb
}

func (tb *tokenBucket) Allow() bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(tb.lastUpdate).Seconds()
	tb.lastUpdate = now

	// Add tokens based on elapsed time
	tb.tokens += elapsed * tb.rate
	if tb.tokens > float64(tb.burst) {
		tb.tokens = float64(tb.burst)
	}

	if tb.tokens >= 1 {
		tb.tokens--
		return true
	}

	return false
}

func (tb *tokenBucket) Remaining() int {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(tb.lastUpdate).Seconds()

	tokens := tb.tokens + elapsed*tb.rate
	if tokens > float64(tb.burst) {
		tokens = float64(tb.burst)
	}

	return int(tokens)
}

func (tb *tokenBucket) RetryAfter() time.Duration {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	if tb.tokens >= 1 {
		return 0
	}

	// Calculate time until next token
	needed := 1.0 - tb.tokens
	return time.Duration(needed/tb.rate) * time.Second
}

func (tb *tokenBucket) touch() {
	tb.lastTouch.Store(time.Now())
}

func (tb *tokenBucket) lastAccess() time.Time {
	return tb.lastTouch.Load().(time.Time)
}

// getClientIP extracts the client IP from the request.
func getClientIP(r *http.Request) string {
	// Check X-Forwarded-For
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		// Get first IP
		for i := 0; i < len(xff); i++ {
			if xff[i] == ',' {
				return xff[:i]
			}
		}
		return xff
	}

	// Check X-Real-IP
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}

	// Use RemoteAddr
	return r.RemoteAddr
}
