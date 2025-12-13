package helix_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	. "github.com/kolosys/helix"
)

func TestHandleWithStatus(t *testing.T) {
	type Request struct{}
	type Response struct {
		ID int `json:"id"`
	}

	s := New()
	s.POST("/create", HandleWithStatus(http.StatusCreated, func(ctx context.Context, req Request) (*Response, error) {
		return &Response{ID: 1}, nil
	}))

	req := httptest.NewRequest(http.MethodPost, "/create", nil)
	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Errorf("expected status 201, got %d", rec.Code)
	}
}

func TestHandleWithStatusError(t *testing.T) {
	type Request struct{}

	s := New()
	s.POST("/error", HandleWithStatus(http.StatusCreated, func(ctx context.Context, req Request) (any, error) {
		return nil, ErrBadRequest.WithDetailf("invalid input")
	}))

	req := httptest.NewRequest(http.MethodPost, "/error", nil)
	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", rec.Code)
	}
}

func TestHandleNoRequestError(t *testing.T) {
	s := New()
	s.GET("/error", HandleNoRequest(func(ctx context.Context) (any, error) {
		return nil, ErrNotFound.WithDetailf("not found")
	}))

	req := httptest.NewRequest(http.MethodGet, "/error", nil)
	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", rec.Code)
	}
}

func TestHandleNoResponseError(t *testing.T) {
	type Request struct {
		ID int `path:"id"`
	}

	s := New()
	s.DELETE("/items/{id}", HandleNoResponse(func(ctx context.Context, req Request) error {
		return ErrForbidden.WithDetailf("cannot delete")
	}))

	req := httptest.NewRequest(http.MethodDelete, "/items/123", nil)
	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Errorf("expected status 403, got %d", rec.Code)
	}
}

func TestHandleEmptyError(t *testing.T) {
	s := New()
	s.POST("/fail", HandleEmpty(func(ctx context.Context) error {
		return ErrInternal.WithDetailf("something went wrong")
	}))

	req := httptest.NewRequest(http.MethodPost, "/fail", nil)
	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("expected status 500, got %d", rec.Code)
	}
}

func TestHandleBindingError(t *testing.T) {
	type Request struct {
		ID int `path:"id,required"`
	}

	s := New()
	s.GET("/items/{id}", Handle(func(ctx context.Context, req Request) (any, error) {
		return map[string]int{"id": req.ID}, nil
	}))

	// The ID is valid so no binding error should occur
	req := httptest.NewRequest(http.MethodGet, "/items/123", nil)
	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}
}

func TestHandleGenericError(t *testing.T) {
	type Request struct{}

	s := New()
	s.GET("/error", Handle(func(ctx context.Context, req Request) (any, error) {
		return nil, context.DeadlineExceeded // generic error
	}))

	req := httptest.NewRequest(http.MethodGet, "/error", nil)
	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, req)

	// Generic errors should return 500
	if rec.Code != http.StatusInternalServerError {
		t.Errorf("expected status 500, got %d", rec.Code)
	}
}

func TestHandleWithValidatable(t *testing.T) {
	s := New()

	type ValidatableRequest struct {
		Email string `json:"email"`
	}

	// Note: Without implementing Validatable on ValidatableRequest,
	// this test just verifies the code path works

	s.POST("/validate", Handle(func(ctx context.Context, req ValidatableRequest) (any, error) {
		return map[string]string{"email": req.Email}, nil
	}))

	body := `{"email":"test@example.com"}`
	req := httptest.NewRequest(http.MethodPost, "/validate", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}
}

func TestIsBindingError(t *testing.T) {
	tests := []struct {
		err      error
		expected bool
	}{
		{ErrBindingFailed, true},
		{ErrUnsupportedType, true},
		{ErrInvalidJSON, true},
		{ErrRequiredField, true},
		{ErrInvalidFieldValue, true},
		{context.DeadlineExceeded, false},
	}

	for _, tc := range tests {
		t.Run(tc.err.Error(), func(t *testing.T) {
			result := IsBindingError(tc.err)
			if result != tc.expected {
				t.Errorf("expected %v, got %v", tc.expected, result)
			}
		})
	}
}

func TestHandleContextCancellation(t *testing.T) {
	type Request struct{}

	s := New()
	s.GET("/slow", Handle(func(ctx context.Context, req Request) (any, error) {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			return map[string]string{"status": "ok"}, nil
		}
	}))

	req := httptest.NewRequest(http.MethodGet, "/slow", nil)
	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}
}

func BenchmarkHandle(b *testing.B) {
	type Request struct {
		ID int `path:"id"`
	}
	type Response struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	s := New()
	s.GET("/users/{id}", Handle(func(ctx context.Context, req Request) (*Response, error) {
		return &Response{ID: req.ID, Name: "John"}, nil
	}))

	req := httptest.NewRequest(http.MethodGet, "/users/123", nil)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		rec := httptest.NewRecorder()
		s.ServeHTTP(rec, req)
	}
}

func BenchmarkHandleWithStatus(b *testing.B) {
	type Request struct{}
	type Response struct {
		ID int `json:"id"`
	}

	s := New()
	s.POST("/create", HandleWithStatus(http.StatusCreated, func(ctx context.Context, req Request) (*Response, error) {
		return &Response{ID: 1}, nil
	}))

	req := httptest.NewRequest(http.MethodPost, "/create", nil)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		rec := httptest.NewRecorder()
		s.ServeHTTP(rec, req)
	}
}

func BenchmarkHandleNoRequest(b *testing.B) {
	type Response struct {
		Status string `json:"status"`
	}

	s := New()
	s.GET("/health", HandleNoRequest(func(ctx context.Context) (*Response, error) {
		return &Response{Status: "ok"}, nil
	}))

	req := httptest.NewRequest(http.MethodGet, "/health", nil)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		rec := httptest.NewRecorder()
		s.ServeHTTP(rec, req)
	}
}
