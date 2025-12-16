package helix_test

import (
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	. "github.com/kolosys/helix"
)

// BenchmarkRouterLockContention benchmarks router lock contention under high concurrency.
func BenchmarkRouterLockContention(b *testing.B) {
	router := NewRouter()

	// Register routes for different methods
	router.Handle(http.MethodGet, "/users/{id}", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	router.Handle(http.MethodPost, "/users", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
	})
	router.Handle(http.MethodPut, "/users/{id}", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	router.Handle(http.MethodDelete, "/users/{id}", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		methods := []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete}
		i := 0
		for pb.Next() {
			method := methods[i%len(methods)]
			i++

			req := httptest.NewRequest(method, "/users/123", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
		}
	})
}

// BenchmarkRouterAllocations benchmarks allocations during route matching.
func BenchmarkRouterAllocations(b *testing.B) {
	router := NewRouter()

	router.Handle(http.MethodGet, "/users/{id}", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/users/123", nil)

	b.ReportAllocs()
	for b.Loop() {
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
	}
}

// BenchmarkRouterConcurrentWrites benchmarks concurrent route registration.
func BenchmarkRouterConcurrentWrites(b *testing.B) {
	router := NewRouter()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			path := "/test/" + string(rune('a'+i%26))
			router.Handle(http.MethodGet, path, func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})
			i++
		}
	})
}

// BenchmarkPerMethodLockContention compares per-method locks vs global lock.
func BenchmarkPerMethodLockContention(b *testing.B) {
	router := NewRouter()

	// Register routes
	for i := 0; i < 10; i++ {
		path := "/test" + string(rune('a'+i))
		router.Handle(http.MethodGet, path, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})
		router.Handle(http.MethodPost, path, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusCreated)
		})
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		methods := []string{http.MethodGet, http.MethodPost}
		paths := []string{"/testa", "/testb", "/testc", "/testd", "/teste"}
		i := 0
		for pb.Next() {
			method := methods[i%len(methods)]
			path := paths[i%len(paths)]
			i++

			req := httptest.NewRequest(method, path, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
		}
	})
}

// TestRouterConcurrentAccess tests that router is safe for concurrent access.
func TestRouterConcurrentAccess(t *testing.T) {
	router := NewRouter()

	// Register initial routes for different methods
	router.Handle(http.MethodGet, "/test", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	router.Handle(http.MethodPost, "/test", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
	})
	router.Handle(http.MethodPut, "/test", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	router.Handle(http.MethodDelete, "/test", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	var wg sync.WaitGroup
	concurrency := 100
	iterations := 1000

	// Concurrent reads - this is what per-method locks optimize
	wg.Add(concurrency)
	for i := range concurrency {
		go func(id int) {
			defer wg.Done()
			methods := []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete}
			for j := 0; j < iterations; j++ {
				method := methods[(id+j)%len(methods)]
				req := httptest.NewRequest(method, "/test", nil)
				w := httptest.NewRecorder()
				router.ServeHTTP(w, req)
			}
		}(i)
	}

	wg.Wait()
}
