package logs

import (
	"sync"
	"sync/atomic"
	"time"
)

// Sampler determines if a log entry should be emitted.
type Sampler interface {
	// Sample returns true if the entry should be logged.
	Sample(level Level, msg string) bool
}

// RateSampler limits logs to a certain rate per message.
type RateSampler struct {
	rate    int           // max logs per interval
	burst   int           // initial burst allowance
	window  time.Duration // time window
	counts  sync.Map      // message -> *rateBucket
	cleanup time.Duration // cleanup interval for old entries
}

type rateBucket struct {
	count     atomic.Int64
	lastReset atomic.Int64
}

// NewRateSampler creates a sampler that limits log rate per message.
// rate is the number of logs allowed per window.
// burst is the initial burst allowance.
func NewRateSampler(rate int, window time.Duration) *RateSampler {
	s := &RateSampler{
		rate:    rate,
		burst:   rate,
		window:  window,
		cleanup: window * 10,
	}
	return s
}

// WithBurst sets the burst allowance.
func (s *RateSampler) WithBurst(burst int) *RateSampler {
	s.burst = burst
	return s
}

// Sample implements Sampler.
func (s *RateSampler) Sample(level Level, msg string) bool {
	now := time.Now().UnixNano()

	// Get or create bucket for this message
	val, _ := s.counts.LoadOrStore(msg, &rateBucket{})
	bucket := val.(*rateBucket)

	// Check if window has passed
	lastReset := bucket.lastReset.Load()
	if now-lastReset >= int64(s.window) {
		// Try to reset the window
		if bucket.lastReset.CompareAndSwap(lastReset, now) {
			bucket.count.Store(1)
			return true
		}
	}

	// Increment count and check limit
	count := bucket.count.Add(1)
	return count <= int64(s.rate)
}

// CountSampler logs every Nth occurrence.
type CountSampler struct {
	n      int
	counts sync.Map // message -> *atomic.Int64
}

// NewCountSampler creates a sampler that logs every Nth occurrence.
func NewCountSampler(n int) *CountSampler {
	return &CountSampler{n: n}
}

// Sample implements Sampler.
func (s *CountSampler) Sample(level Level, msg string) bool {
	val, _ := s.counts.LoadOrStore(msg, &atomic.Int64{})
	count := val.(*atomic.Int64).Add(1)
	return count%int64(s.n) == 1
}

// LevelSampler applies different samplers per level.
type LevelSampler struct {
	samplers map[Level]Sampler
	fallback Sampler
}

// NewLevelSampler creates a sampler with per-level configuration.
func NewLevelSampler(fallback Sampler) *LevelSampler {
	return &LevelSampler{
		samplers: make(map[Level]Sampler),
		fallback: fallback,
	}
}

// WithLevel sets the sampler for a specific level.
func (s *LevelSampler) WithLevel(level Level, sampler Sampler) *LevelSampler {
	s.samplers[level] = sampler
	return s
}

// Sample implements Sampler.
func (s *LevelSampler) Sample(level Level, msg string) bool {
	if sampler, ok := s.samplers[level]; ok {
		return sampler.Sample(level, msg)
	}
	if s.fallback != nil {
		return s.fallback.Sample(level, msg)
	}
	return true
}

// FirstNSampler logs only the first N occurrences.
type FirstNSampler struct {
	n      int64
	counts sync.Map // message -> *atomic.Int64
}

// NewFirstNSampler creates a sampler that logs only the first N occurrences.
func NewFirstNSampler(n int) *FirstNSampler {
	return &FirstNSampler{n: int64(n)}
}

// Sample implements Sampler.
func (s *FirstNSampler) Sample(level Level, msg string) bool {
	val, _ := s.counts.LoadOrStore(msg, &atomic.Int64{})
	count := val.(*atomic.Int64).Add(1)
	return count <= s.n
}

// OncePerSampler logs a message only once per duration.
type OncePerSampler struct {
	period   time.Duration
	lastSeen sync.Map // message -> int64 (UnixNano)
}

// NewOncePerSampler creates a sampler that logs each message at most once per period.
func NewOncePerSampler(period time.Duration) *OncePerSampler {
	return &OncePerSampler{period: period}
}

// Sample implements Sampler.
func (s *OncePerSampler) Sample(level Level, msg string) bool {
	now := time.Now().UnixNano()

	val, loaded := s.lastSeen.Load(msg)
	if !loaded {
		s.lastSeen.Store(msg, now)
		return true
	}

	lastSeen := val.(int64)
	if now-lastSeen >= int64(s.period) {
		s.lastSeen.Store(msg, now)
		return true
	}

	return false
}

// CompositeSampler combines multiple samplers with AND logic.
type CompositeSampler struct {
	samplers []Sampler
}

// NewCompositeSampler creates a sampler that requires all samplers to pass.
func NewCompositeSampler(samplers ...Sampler) *CompositeSampler {
	return &CompositeSampler{samplers: samplers}
}

// Sample implements Sampler.
func (s *CompositeSampler) Sample(level Level, msg string) bool {
	for _, sampler := range s.samplers {
		if !sampler.Sample(level, msg) {
			return false
		}
	}
	return true
}

// RandomSampler samples a percentage of logs.
type RandomSampler struct {
	threshold uint32
	counter   atomic.Uint64
}

// NewRandomSampler creates a sampler that logs a percentage of entries.
// percentage should be between 0 and 100.
func NewRandomSampler(percentage int) *RandomSampler {
	if percentage < 0 {
		percentage = 0
	}
	if percentage > 100 {
		percentage = 100
	}
	return &RandomSampler{
		threshold: uint32(percentage),
	}
}

// Sample implements Sampler using a simple modulo for deterministic "random" sampling.
func (s *RandomSampler) Sample(level Level, msg string) bool {
	count := s.counter.Add(1)
	return uint32(count%100) < s.threshold
}

// AlwaysSampler always allows logging.
type AlwaysSampler struct{}

// Sample implements Sampler.
func (s *AlwaysSampler) Sample(level Level, msg string) bool {
	return true
}

// NeverSampler never allows logging.
type NeverSampler struct{}

// Sample implements Sampler.
func (s *NeverSampler) Sample(level Level, msg string) bool {
	return false
}

