package helix_test

import (
	"net/http/httptest"
	"strings"
	"testing"

	. "github.com/kolosys/helix"
)

func TestNewProblem(t *testing.T) {
	p := NewProblem(400, "bad_request", "Bad Request")

	if p.Status != 400 {
		t.Errorf("expected status 400, got %d", p.Status)
	}
	if p.Title != "Bad Request" {
		t.Errorf("expected title 'Bad Request', got '%s'", p.Title)
	}
	if !strings.Contains(p.Type, "bad_request") {
		t.Errorf("expected type to contain 'bad_request', got '%s'", p.Type)
	}
}

func TestProblemError(t *testing.T) {
	p := NewProblem(404, "not_found", "Not Found")

	if p.Error() != "Not Found" {
		t.Errorf("expected 'Not Found', got '%s'", p.Error())
	}

	p2 := p.WithDetailf("resource %s not found", "user")
	if !strings.Contains(p2.Error(), "resource user not found") {
		t.Errorf("expected error with detail, got '%s'", p2.Error())
	}
}

func TestProblemWithDetail(t *testing.T) {
	p := ErrNotFound.WithDetailf("user %d not found", 123)

	if p.Detail != "user 123 not found" {
		t.Errorf("expected 'user 123 not found', got '%s'", p.Detail)
	}

	// Original should be unchanged
	if ErrNotFound.Detail != "" {
		t.Error("original problem should be unchanged")
	}
}

func TestProblemWithInstanceChain(t *testing.T) {
	p := ErrBadRequest.WithInstance("/api/users/123")

	if p.Instance != "/api/users/123" {
		t.Errorf("expected '/api/users/123', got '%s'", p.Instance)
	}

	// Original should be unchanged
	if ErrBadRequest.Instance != "" {
		t.Error("original problem should be unchanged")
	}
}

func TestProblemWithType(t *testing.T) {
	p := ErrNotFound.WithType("https://example.com/errors/not-found")

	if p.Type != "https://example.com/errors/not-found" {
		t.Errorf("expected custom type, got '%s'", p.Type)
	}
}

func TestSentinelErrors(t *testing.T) {
	tests := []struct {
		problem Problem
		status  int
	}{
		{ErrBadRequest, 400},
		{ErrUnauthorized, 401},
		{ErrForbidden, 403},
		{ErrNotFound, 404},
		{ErrMethodNotAllowed, 405},
		{ErrConflict, 409},
		{ErrGone, 410},
		{ErrUnprocessableEntity, 422},
		{ErrTooManyRequests, 429},
		{ErrInternal, 500},
		{ErrNotImplemented, 501},
		{ErrBadGateway, 502},
		{ErrServiceUnavailable, 503},
		{ErrGatewayTimeout, 504},
	}

	for _, tc := range tests {
		t.Run(tc.problem.Title, func(t *testing.T) {
			if tc.problem.Status != tc.status {
				t.Errorf("expected status %d, got %d", tc.status, tc.problem.Status)
			}
		})
	}
}

func TestProblemFromStatus(t *testing.T) {
	p := ProblemFromStatus(418)

	if p.Status != 418 {
		t.Errorf("expected status 418, got %d", p.Status)
	}
}

func TestWriteProblemContentType(t *testing.T) {
	rec := httptest.NewRecorder()
	p := ErrNotFound.WithDetailf("not found")

	WriteProblem(rec, p)

	contentType := rec.Header().Get("Content-Type")
	if !strings.Contains(contentType, "application/problem+json") {
		t.Errorf("expected Content-Type application/problem+json, got '%s'", contentType)
	}

	if rec.Code != 404 {
		t.Errorf("expected status 404, got %d", rec.Code)
	}
}

func TestProblemChaining(t *testing.T) {
	p := ErrBadRequest.
		WithDetailf("invalid email").
		WithInstance("/api/users").
		WithType("https://example.com/errors/validation")

	if p.Detail != "invalid email" {
		t.Errorf("expected detail, got '%s'", p.Detail)
	}
	if p.Instance != "/api/users" {
		t.Errorf("expected instance, got '%s'", p.Instance)
	}
	if p.Type != "https://example.com/errors/validation" {
		t.Errorf("expected type, got '%s'", p.Type)
	}
}

func TestProblemWithDetailNoArgs(t *testing.T) {
	p := ErrBadRequest.WithDetailf("simple message")

	if p.Detail != "simple message" {
		t.Errorf("expected 'simple message', got '%s'", p.Detail)
	}
}

func BenchmarkProblemWithDetail(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		ErrNotFound.WithDetailf("user %d not found", 123)
	}
}

func BenchmarkWriteProblem(b *testing.B) {
	p := ErrNotFound.WithDetailf("not found")

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		rec := httptest.NewRecorder()
		WriteProblem(rec, p)
	}
}
