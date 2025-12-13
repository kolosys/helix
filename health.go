package helix

import (
	"context"
	"net/http"
	"sync"
	"time"
)

// HealthStatus represents the health status of a component.
type HealthStatus string

const (
	HealthStatusUp       HealthStatus = "up"
	HealthStatusDown     HealthStatus = "down"
	HealthStatusDegraded HealthStatus = "degraded"
)

// HealthCheck is a function that checks the health of a component.
type HealthCheck func(ctx context.Context) HealthCheckResult

// HealthCheckResult contains the result of a health check.
type HealthCheckResult struct {
	Status  HealthStatus   `json:"status"`
	Message string         `json:"message,omitempty"`
	Latency time.Duration  `json:"latency_ms,omitempty"`
	Details map[string]any `json:"details,omitempty"`
}

// HealthResponse is the response returned by the health endpoint.
type HealthResponse struct {
	Status     HealthStatus                 `json:"status"`
	Timestamp  time.Time                    `json:"timestamp"`
	Version    string                       `json:"version,omitempty"`
	Components map[string]HealthCheckResult `json:"components,omitempty"`
}

// HealthBuilder provides a fluent interface for building health check endpoints.
type HealthBuilder struct {
	checks  map[string]HealthCheck
	version string
	timeout time.Duration
}

// Health creates a new HealthBuilder.
func Health() *HealthBuilder {
	return &HealthBuilder{
		checks:  make(map[string]HealthCheck),
		timeout: 5 * time.Second,
	}
}

// Version sets the application version shown in health responses.
func (h *HealthBuilder) Version(v string) *HealthBuilder {
	h.version = v
	return h
}

// Timeout sets the timeout for health checks.
func (h *HealthBuilder) Timeout(d time.Duration) *HealthBuilder {
	h.timeout = d
	return h
}

// Check adds a health check for a named component.
func (h *HealthBuilder) Check(name string, check HealthCheck) *HealthBuilder {
	h.checks[name] = check
	return h
}

// CheckFunc adds a simple health check that returns an error.
func (h *HealthBuilder) CheckFunc(name string, check func(ctx context.Context) error) *HealthBuilder {
	h.checks[name] = func(ctx context.Context) HealthCheckResult {
		start := time.Now()
		err := check(ctx)
		latency := time.Since(start)

		if err != nil {
			return HealthCheckResult{
				Status:  HealthStatusDown,
				Message: err.Error(),
				Latency: latency,
			}
		}
		return HealthCheckResult{
			Status:  HealthStatusUp,
			Latency: latency,
		}
	}
	return h
}

// Handler returns an http.HandlerFunc for the health check endpoint.
func (h *HealthBuilder) Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), h.timeout)
		defer cancel()

		response := HealthResponse{
			Status:     HealthStatusUp,
			Timestamp:  time.Now().UTC(),
			Version:    h.version,
			Components: make(map[string]HealthCheckResult),
		}

		// Run all checks concurrently
		var wg sync.WaitGroup
		var mu sync.Mutex

		for name, check := range h.checks {
			wg.Add(1)
			go func(name string, check HealthCheck) {
				defer wg.Done()

				result := check(ctx)

				mu.Lock()
				response.Components[name] = result

				// Update overall status
				if result.Status == HealthStatusDown {
					response.Status = HealthStatusDown
				} else if result.Status == HealthStatusDegraded && response.Status != HealthStatusDown {
					response.Status = HealthStatusDegraded
				}
				mu.Unlock()
			}(name, check)
		}

		wg.Wait()

		status := http.StatusOK
		if response.Status == HealthStatusDown {
			status = http.StatusServiceUnavailable
		}

		JSON(w, status, response)
	}
}

// LivenessHandler returns a simple liveness probe handler.
// Returns 200 OK if the server is running.
func LivenessHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		JSON(w, http.StatusOK, map[string]string{
			"status": "alive",
		})
	}
}

// ReadinessHandler returns a simple readiness probe handler.
// Uses the provided checks to determine readiness.
func ReadinessHandler(checks ...func(ctx context.Context) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		for _, check := range checks {
			if err := check(ctx); err != nil {
				JSON(w, http.StatusServiceUnavailable, map[string]string{
					"status": "not_ready",
					"error":  err.Error(),
				})
				return
			}
		}

		JSON(w, http.StatusOK, map[string]string{
			"status": "ready",
		})
	}
}
