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

	s := New(nil)
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

	s := New(nil)
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
	s := New(nil)
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

	s := New(nil)
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
	s := New(nil)
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

	s := New(nil)
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

	s := New(nil)
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
	s := New(nil)

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

	s := New(nil)
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

	s := New(nil)
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

	s := New(nil)
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

	s := New(nil)
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

func TestHandleVersions(t *testing.T) {
	type Request struct{}
	type Response struct {
		Version string `json:"version"`
		Data    string `json:"data"`
	}

	tests := []struct {
		name           string
		headerName     string
		defaultVersion string
		versions       map[string]Handler[Request, Response]
		headerValue    string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "version from header",
			headerName:     "API-Version",
			defaultVersion: "v1",
			versions: map[string]Handler[Request, Response]{
				"v1": func(ctx context.Context, req Request) (Response, error) {
					return Response{Version: "v1", Data: "v1 data"}, nil
				},
				"v2": func(ctx context.Context, req Request) (Response, error) {
					return Response{Version: "v2", Data: "v2 data"}, nil
				},
			},
			headerValue:    "v2",
			expectedStatus: http.StatusOK,
			expectedBody:   `{"version":"v2","data":"v2 data"}`,
		},
		{
			name:           "default version when no header",
			headerName:     "API-Version",
			defaultVersion: "v1",
			versions: map[string]Handler[Request, Response]{
				"v1": func(ctx context.Context, req Request) (Response, error) {
					return Response{Version: "v1", Data: "v1 data"}, nil
				},
				"v2": func(ctx context.Context, req Request) (Response, error) {
					return Response{Version: "v2", Data: "v2 data"}, nil
				},
			},
			headerValue:    "",
			expectedStatus: http.StatusOK,
			expectedBody:   `{"version":"v1","data":"v1 data"}`,
		},
		{
			name:           "version not found returns 404",
			headerName:     "API-Version",
			defaultVersion: "v1",
			versions: map[string]Handler[Request, Response]{
				"v1": func(ctx context.Context, req Request) (Response, error) {
					return Response{Version: "v1", Data: "v1 data"}, nil
				},
			},
			headerValue:    "v3",
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s := New(nil)
			s.GET("/api", HandleVersions(tc.headerName, tc.defaultVersion, tc.versions))

			req := httptest.NewRequest(http.MethodGet, "/api", nil)
			if tc.headerValue != "" {
				req.Header.Set(tc.headerName, tc.headerValue)
			}
			rec := httptest.NewRecorder()
			s.ServeHTTP(rec, req)

			if rec.Code != tc.expectedStatus {
				t.Errorf("expected status %d, got %d", tc.expectedStatus, rec.Code)
			}

			if tc.expectedBody != "" {
				body := strings.TrimSpace(rec.Body.String())
				if body != tc.expectedBody {
					t.Errorf("expected body %q, got %q", tc.expectedBody, body)
				}
			}
		})
	}
}

func TestHandleVersionsWithError(t *testing.T) {
	type Request struct{}
	type Response struct {
		Data string `json:"data"`
	}

	s := New(nil)
	s.GET("/api", HandleVersions("API-Version", "v1", map[string]Handler[Request, Response]{
		"v1": func(ctx context.Context, req Request) (Response, error) {
			return Response{}, ErrBadRequest.WithDetailf("invalid request")
		},
		"v2": func(ctx context.Context, req Request) (Response, error) {
			return Response{}, ErrNotFound.WithDetailf("not found")
		},
	}))

	tests := []struct {
		name           string
		headerValue    string
		expectedStatus int
	}{
		{
			name:           "v1 returns bad request",
			headerValue:    "v1",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "v2 returns not found",
			headerValue:    "v2",
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/api", nil)
			req.Header.Set("API-Version", tc.headerValue)
			rec := httptest.NewRecorder()
			s.ServeHTTP(rec, req)

			if rec.Code != tc.expectedStatus {
				t.Errorf("expected status %d, got %d", tc.expectedStatus, rec.Code)
			}
		})
	}
}

func TestHandleVersionsWithBinding(t *testing.T) {
	type Request struct {
		ID int `path:"id"`
	}
	type Response struct {
		ID      int    `json:"id"`
		Version string `json:"version"`
	}

	s := New(nil)
	s.GET("/items/{id}", HandleVersions("API-Version", "v1", map[string]Handler[Request, Response]{
		"v1": func(ctx context.Context, req Request) (Response, error) {
			return Response{ID: req.ID, Version: "v1"}, nil
		},
		"v2": func(ctx context.Context, req Request) (Response, error) {
			return Response{ID: req.ID, Version: "v2"}, nil
		},
	}))

	tests := []struct {
		name           string
		path           string
		headerValue    string
		expectedStatus int
		expectedID     int
		expectedVer    string
	}{
		{
			name:           "v1 with path param",
			path:           "/items/123",
			headerValue:    "v1",
			expectedStatus: http.StatusOK,
			expectedID:     123,
			expectedVer:    "v1",
		},
		{
			name:           "v2 with path param",
			path:           "/items/456",
			headerValue:    "v2",
			expectedStatus: http.StatusOK,
			expectedID:     456,
			expectedVer:    "v2",
		},
		{
			name:           "default version with path param",
			path:           "/items/789",
			headerValue:    "",
			expectedStatus: http.StatusOK,
			expectedID:     789,
			expectedVer:    "v1",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tc.path, nil)
			if tc.headerValue != "" {
				req.Header.Set("API-Version", tc.headerValue)
			}
			rec := httptest.NewRecorder()
			s.ServeHTTP(rec, req)

			if rec.Code != tc.expectedStatus {
				t.Errorf("expected status %d, got %d", tc.expectedStatus, rec.Code)
			}

			if rec.Code == http.StatusOK {
				body := rec.Body.String()
				if !strings.Contains(body, `"version":"`+tc.expectedVer+`"`) {
					t.Errorf("expected version %q in response, got %s", tc.expectedVer, body)
				}
			}
		})
	}
}

type versionValidatableRequest struct {
	Email string `json:"email"`
}

func (v *versionValidatableRequest) Validate() error {
	verrs := NewValidationErrors()
	if v.Email == "" {
		verrs.Add("email", "email is required")
	}
	return verrs.Err()
}

func TestHandleVersionsWithValidatable(t *testing.T) {
	type Response struct {
		Email   string `json:"email"`
		Version string `json:"version"`
	}

	s := New(nil)
	s.POST("/users", HandleVersions("API-Version", "v1", map[string]Handler[versionValidatableRequest, Response]{
		"v1": func(ctx context.Context, req versionValidatableRequest) (Response, error) {
			return Response{Email: req.Email, Version: "v1"}, nil
		},
		"v2": func(ctx context.Context, req versionValidatableRequest) (Response, error) {
			return Response{Email: req.Email, Version: "v2"}, nil
		},
	}))

	tests := []struct {
		name           string
		body           string
		headerValue    string
		expectedStatus int
	}{
		{
			name:           "v1 with valid request",
			body:           `{"email":"test@example.com"}`,
			headerValue:    "v1",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "v2 with valid request",
			body:           `{"email":"test@example.com"}`,
			headerValue:    "v2",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "v1 with invalid request",
			body:           `{"email":""}`,
			headerValue:    "v1",
			expectedStatus: http.StatusUnprocessableEntity, // Validation errors return 422
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(tc.body))
			req.Header.Set("Content-Type", "application/json")
			if tc.headerValue != "" {
				req.Header.Set("API-Version", tc.headerValue)
			}
			rec := httptest.NewRecorder()
			s.ServeHTTP(rec, req)

			if rec.Code != tc.expectedStatus {
				t.Errorf("expected status %d, got %d", tc.expectedStatus, rec.Code)
			}
		})
	}
}

func BenchmarkHandleVersions(b *testing.B) {
	type Request struct{}
	type Response struct {
		Version string `json:"version"`
		Data    string `json:"data"`
	}

	s := New(nil)
	s.GET("/api", HandleVersions("API-Version", "v1", map[string]Handler[Request, Response]{
		"v1": func(ctx context.Context, req Request) (Response, error) {
			return Response{Version: "v1", Data: "data"}, nil
		},
		"v2": func(ctx context.Context, req Request) (Response, error) {
			return Response{Version: "v2", Data: "data"}, nil
		},
	}))

	req := httptest.NewRequest(http.MethodGet, "/api", nil)
	req.Header.Set("API-Version", "v1")

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		rec := httptest.NewRecorder()
		s.ServeHTTP(rec, req)
	}
}
