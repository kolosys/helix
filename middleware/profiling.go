//go:build profile

package middleware

import (
	"net/http"
	"runtime"
	"sync"
	"time"
)

// MiddlewareProfile contains profiling information for a middleware.
type MiddlewareProfile struct {
	Name     string
	Duration time.Duration
	Allocs   uint64
	Bytes    uint64
}

// middlewareProfiler tracks profiling data for middleware.
type middlewareProfiler struct {
	mu         sync.RWMutex
	profiles   map[string]*MiddlewareProfile
	totalCalls int64
}

var profiler = &middlewareProfiler{
	profiles: make(map[string]*MiddlewareProfile),
}

// ProfileMiddleware wraps a middleware with profiling instrumentation.
func ProfileMiddleware(name string, mw Middleware) Middleware {
	return func(next http.Handler) http.Handler {
		return mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var memBefore, memAfter runtime.MemStats
			runtime.ReadMemStats(&memBefore)
			start := time.Now()

			next.ServeHTTP(w, r)

			duration := time.Since(start)
			runtime.ReadMemStats(&memAfter)

			allocs := memAfter.Mallocs - memBefore.Mallocs
			bytes := memAfter.TotalAlloc - memBefore.TotalAlloc

			profiler.mu.Lock()
			if p, ok := profiler.profiles[name]; ok {
				p.Duration += duration
				p.Allocs += allocs
				p.Bytes += bytes
			} else {
				profiler.profiles[name] = &MiddlewareProfile{
					Name:     name,
					Duration: duration,
					Allocs:   allocs,
					Bytes:    bytes,
				}
			}
			profiler.totalCalls++
			profiler.mu.Unlock()
		}))
	}
}

// GetProfiles returns all middleware profiles.
func GetProfiles() map[string]*MiddlewareProfile {
	profiler.mu.RLock()
	defer profiler.mu.RUnlock()

	result := make(map[string]*MiddlewareProfile, len(profiler.profiles))
	for k, v := range profiler.profiles {
		result[k] = &MiddlewareProfile{
			Name:     v.Name,
			Duration: v.Duration,
			Allocs:   v.Allocs,
			Bytes:    v.Bytes,
		}
	}
	return result
}

// ResetProfiles clears all profiling data.
func ResetProfiles() {
	profiler.mu.Lock()
	defer profiler.mu.Unlock()
	profiler.profiles = make(map[string]*MiddlewareProfile)
	profiler.totalCalls = 0
}
