package helix_test

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	. "github.com/kolosys/helix"
)

func TestNew(t *testing.T) {
	s := New(nil)
	if s == nil {
		t.Fatal("New() returned nil")
	}
	cfg := s.GetConfig()
	if cfg.Addr != ":8080" {
		t.Errorf("expected default addr :8080, got %s", cfg.Addr)
	}
}

func TestNewWithOptions(t *testing.T) {
	s := New(&Options{
		Addr:         ":3000",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 20 * time.Second,
		IdleTimeout:  60 * time.Second,
		GracePeriod:  15 * time.Second,
	})

	cfg := s.GetConfig()
	if cfg.Addr != ":3000" {
		t.Errorf("expected addr :3000, got %s", cfg.Addr)
	}
	if cfg.ReadTimeout != int64(10*time.Second) {
		t.Errorf("expected readTimeout 10s, got %v", time.Duration(cfg.ReadTimeout))
	}
	if cfg.WriteTimeout != int64(20*time.Second) {
		t.Errorf("expected writeTimeout 20s, got %v", time.Duration(cfg.WriteTimeout))
	}
	if cfg.IdleTimeout != int64(60*time.Second) {
		t.Errorf("expected idleTimeout 60s, got %v", time.Duration(cfg.IdleTimeout))
	}
	if cfg.GracePeriod != int64(15*time.Second) {
		t.Errorf("expected gracePeriod 15s, got %v", time.Duration(cfg.GracePeriod))
	}
}

func TestServerUse(t *testing.T) {
	s := New(nil)

	called := false
	mw := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true
			next.ServeHTTP(w, r)
		})
	}

	s.Use(mw)

	cfg := s.GetConfig()
	if cfg.MiddlewareLen != 1 {
		t.Errorf("expected 1 middleware, got %d", cfg.MiddlewareLen)
	}

	s.GET("/test", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, req)

	if !called {
		t.Error("middleware was not called")
	}
}

func TestRouteRegistration(t *testing.T) {
	methods := []struct {
		name   string
		method string
		fn     func(s *Server, pattern string, handler http.HandlerFunc)
	}{
		{"GET", http.MethodGet, (*Server).GET},
		{"POST", http.MethodPost, (*Server).POST},
		{"PUT", http.MethodPut, (*Server).PUT},
		{"PATCH", http.MethodPatch, (*Server).PATCH},
		{"DELETE", http.MethodDelete, (*Server).DELETE},
		{"OPTIONS", http.MethodOptions, (*Server).OPTIONS},
		{"HEAD", http.MethodHead, (*Server).HEAD},
	}

	for _, tc := range methods {
		t.Run(tc.name, func(t *testing.T) {
			s := New(nil)
			called := false

			tc.fn(s, "/test", func(w http.ResponseWriter, r *http.Request) {
				called = true
				w.WriteHeader(http.StatusOK)
			})

			req := httptest.NewRequest(tc.method, "/test", nil)
			rec := httptest.NewRecorder()
			s.ServeHTTP(rec, req)

			if !called {
				t.Errorf("%s handler was not called", tc.name)
			}
			if rec.Code != http.StatusOK {
				t.Errorf("expected status 200, got %d", rec.Code)
			}
		})
	}
}

func TestRouteNotFound(t *testing.T) {
	s := New(nil)
	s.GET("/exists", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/not-exists", nil)
	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", rec.Code)
	}
}

func TestPathParams(t *testing.T) {
	s := New(nil)
	s.GET("/users/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := Param(r, "id")
		Text(w, http.StatusOK, id)
	})

	req := httptest.NewRequest(http.MethodGet, "/users/123", nil)
	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}
	if rec.Body.String() != "123" {
		t.Errorf("expected body '123', got '%s'", rec.Body.String())
	}
}

func TestMultiplePathParams(t *testing.T) {
	s := New(nil)
	s.GET("/users/{userID}/posts/{postID}", func(w http.ResponseWriter, r *http.Request) {
		userID := Param(r, "userID")
		postID := Param(r, "postID")
		Text(w, http.StatusOK, userID+":"+postID)
	})

	req := httptest.NewRequest(http.MethodGet, "/users/42/posts/99", nil)
	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}
	if rec.Body.String() != "42:99" {
		t.Errorf("expected body '42:99', got '%s'", rec.Body.String())
	}
}

func TestWildcardParam(t *testing.T) {
	s := New(nil)
	s.GET("/files/{path...}", func(w http.ResponseWriter, r *http.Request) {
		path := Param(r, "path")
		Text(w, http.StatusOK, path)
	})

	req := httptest.NewRequest(http.MethodGet, "/files/docs/readme.md", nil)
	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}
	if rec.Body.String() != "docs/readme.md" {
		t.Errorf("expected body 'docs/readme.md', got '%s'", rec.Body.String())
	}
}

func TestParamInt(t *testing.T) {
	s := New(nil)
	s.GET("/users/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, err := ParamInt(r, "id")
		if err != nil {
			BadRequest(w, err.Error())
			return
		}
		JSON(w, http.StatusOK, map[string]int{"id": id})
	})

	tests := []struct {
		name     string
		path     string
		wantCode int
	}{
		{"valid int", "/users/123", http.StatusOK},
		{"invalid int", "/users/abc", http.StatusBadRequest},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tc.path, nil)
			rec := httptest.NewRecorder()
			s.ServeHTTP(rec, req)

			if rec.Code != tc.wantCode {
				t.Errorf("expected status %d, got %d", tc.wantCode, rec.Code)
			}
		})
	}
}

func TestParamInt64(t *testing.T) {
	s := New(nil)
	s.GET("/items/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, err := ParamInt64(r, "id")
		if err != nil {
			BadRequest(w, err.Error())
			return
		}
		JSON(w, http.StatusOK, map[string]int64{"id": id})
	})

	req := httptest.NewRequest(http.MethodGet, "/items/9223372036854775807", nil)
	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}
}

func TestParamUUID(t *testing.T) {
	s := New(nil)
	s.GET("/resources/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, err := ParamUUID(r, "id")
		if err != nil {
			BadRequest(w, err.Error())
			return
		}
		Text(w, http.StatusOK, id)
	})

	tests := []struct {
		name     string
		uuid     string
		wantCode int
	}{
		{"valid uuid with hyphens", "550e8400-e29b-41d4-a716-446655440000", http.StatusOK},
		{"valid uuid without hyphens", "550e8400e29b41d4a716446655440000", http.StatusOK},
		{"invalid uuid", "not-a-uuid", http.StatusBadRequest},
		{"too short", "550e8400-e29b", http.StatusBadRequest},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/resources/"+tc.uuid, nil)
			rec := httptest.NewRecorder()
			s.ServeHTTP(rec, req)

			if rec.Code != tc.wantCode {
				t.Errorf("expected status %d, got %d", tc.wantCode, rec.Code)
			}
		})
	}
}

func TestQueryParams(t *testing.T) {
	tests := []struct {
		name     string
		query    string
		param    string
		expected string
	}{
		{"simple", "?name=john", "name", "john"},
		{"missing", "?other=value", "name", ""},
		{"empty", "?name=", "name", ""},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s := New(nil)
			s.GET("/search", func(w http.ResponseWriter, r *http.Request) {
				name := Query(r, tc.param)
				Text(w, http.StatusOK, name)
			})

			req := httptest.NewRequest(http.MethodGet, "/search"+tc.query, nil)
			rec := httptest.NewRecorder()
			s.ServeHTTP(rec, req)

			if rec.Body.String() != tc.expected {
				t.Errorf("expected '%s', got '%s'", tc.expected, rec.Body.String())
			}
		})
	}
}

func TestQueryInt(t *testing.T) {
	s := New(nil)
	s.GET("/page", func(w http.ResponseWriter, r *http.Request) {
		page := QueryInt(r, "page", 1)
		JSON(w, http.StatusOK, map[string]int{"page": page})
	})

	tests := []struct {
		name     string
		query    string
		expected string
	}{
		{"valid", "?page=5", `{"page":5}`},
		{"default", "", `{"page":1}`},
		{"invalid uses default", "?page=abc", `{"page":1}`},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/page"+tc.query, nil)
			rec := httptest.NewRecorder()
			s.ServeHTTP(rec, req)

			body := strings.TrimSpace(rec.Body.String())
			if body != tc.expected {
				t.Errorf("expected '%s', got '%s'", tc.expected, body)
			}
		})
	}
}

func TestQueryBool(t *testing.T) {
	s := New(nil)
	s.GET("/active", func(w http.ResponseWriter, r *http.Request) {
		active := QueryBool(r, "active")
		JSON(w, http.StatusOK, map[string]bool{"active": active})
	})

	tests := []struct {
		name     string
		query    string
		expected bool
	}{
		{"true", "?active=true", true},
		{"false", "?active=false", false},
		{"1", "?active=1", true},
		{"0", "?active=0", false},
		{"yes", "?active=yes", true},
		{"no", "?active=no", false},
		{"missing", "", false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/active"+tc.query, nil)
			rec := httptest.NewRecorder()
			s.ServeHTTP(rec, req)

			expected := `{"active":false}`
			if tc.expected {
				expected = `{"active":true}`
			}

			body := strings.TrimSpace(rec.Body.String())
			if body != expected {
				t.Errorf("expected '%s', got '%s'", expected, body)
			}
		})
	}
}

func TestQuerySlice(t *testing.T) {
	s := New(nil)
	s.GET("/tags", func(w http.ResponseWriter, r *http.Request) {
		tags := QuerySlice(r, "tag")
		JSON(w, http.StatusOK, map[string][]string{"tags": tags})
	})

	req := httptest.NewRequest(http.MethodGet, "/tags?tag=go&tag=web&tag=api", nil)
	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, req)

	body := strings.TrimSpace(rec.Body.String())
	expected := `{"tags":["go","web","api"]}`
	if body != expected {
		t.Errorf("expected '%s', got '%s'", expected, body)
	}
}

func TestServerShutdown(t *testing.T) {
	s := New(&Options{Addr: ":0"})

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	// Shutdown should not error when server hasn't started
	err := s.Shutdown(ctx)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestServerAddr(t *testing.T) {
	s := New(&Options{Addr: ":9999"})
	if s.Addr() != ":9999" {
		t.Errorf("expected addr :9999, got %s", s.Addr())
	}
}

func TestMiddlewareOrder(t *testing.T) {
	s := New(nil)

	var order []string

	s.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			order = append(order, "mw1-before")
			next.ServeHTTP(w, r)
			order = append(order, "mw1-after")
		})
	})

	s.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			order = append(order, "mw2-before")
			next.ServeHTTP(w, r)
			order = append(order, "mw2-after")
		})
	})

	s.GET("/test", func(w http.ResponseWriter, r *http.Request) {
		order = append(order, "handler")
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, req)

	expected := []string{"mw1-before", "mw2-before", "handler", "mw2-after", "mw1-after"}
	if len(order) != len(expected) {
		t.Fatalf("expected %d calls, got %d", len(expected), len(order))
	}
	for i, v := range expected {
		if order[i] != v {
			t.Errorf("expected order[%d] = %s, got %s", i, v, order[i])
		}
	}
}

func TestJSONResponse(t *testing.T) {
	s := New(nil)
	s.GET("/json", func(w http.ResponseWriter, r *http.Request) {
		JSON(w, http.StatusOK, map[string]string{"message": "hello"})
	})

	req := httptest.NewRequest(http.MethodGet, "/json", nil)
	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}

	contentType := rec.Header().Get("Content-Type")
	if !strings.HasPrefix(contentType, MIMEApplicationJSONCharsetUTF8) {
		t.Errorf("expected Content-Type %s, got %s", MIMEApplicationJSONCharsetUTF8, contentType)
	}

	expected := `{"message":"hello"}`
	body := strings.TrimSpace(rec.Body.String())
	if body != expected {
		t.Errorf("expected '%s', got '%s'", expected, body)
	}
}

func TestTextResponse(t *testing.T) {
	s := New(nil)
	s.GET("/text", func(w http.ResponseWriter, r *http.Request) {
		Text(w, http.StatusOK, "hello world")
	})

	req := httptest.NewRequest(http.MethodGet, "/text", nil)
	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}

	contentType := rec.Header().Get("Content-Type")
	if !strings.HasPrefix(contentType, MIMETextPlainCharsetUTF8) {
		t.Errorf("expected Content-Type %s, got %s", MIMETextPlainCharsetUTF8, contentType)
	}

	if rec.Body.String() != "hello world" {
		t.Errorf("expected 'hello world', got '%s'", rec.Body.String())
	}
}

func TestNoContentResponse(t *testing.T) {
	s := New(nil)
	s.DELETE("/resource", func(w http.ResponseWriter, r *http.Request) {
		NoContent(w)
	})

	req := httptest.NewRequest(http.MethodDelete, "/resource", nil)
	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Errorf("expected status 204, got %d", rec.Code)
	}

	if rec.Body.Len() != 0 {
		t.Errorf("expected empty body, got %d bytes", rec.Body.Len())
	}
}

func TestHTMLResponse(t *testing.T) {
	s := New(nil)
	s.GET("/html", func(w http.ResponseWriter, r *http.Request) {
		HTML(w, http.StatusOK, "<h1>Hello</h1>")
	})

	req := httptest.NewRequest(http.MethodGet, "/html", nil)
	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, req)

	contentType := rec.Header().Get("Content-Type")
	if !strings.HasPrefix(contentType, "text/html") {
		t.Errorf("expected Content-Type text/html, got %s", contentType)
	}
}

func TestRedirect(t *testing.T) {
	s := New(nil)
	s.GET("/old", func(w http.ResponseWriter, r *http.Request) {
		Redirect(w, r, "/new", http.StatusMovedPermanently)
	})

	req := httptest.NewRequest(http.MethodGet, "/old", nil)
	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, req)

	if rec.Code != http.StatusMovedPermanently {
		t.Errorf("expected status 301, got %d", rec.Code)
	}

	location := rec.Header().Get("Location")
	if location != "/new" {
		t.Errorf("expected Location /new, got %s", location)
	}
}

func TestStream(t *testing.T) {
	s := New(nil)
	s.GET("/stream", func(w http.ResponseWriter, r *http.Request) {
		reader := strings.NewReader("streaming content")
		Stream(w, "text/plain", reader)
	})

	req := httptest.NewRequest(http.MethodGet, "/stream", nil)
	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, req)

	if rec.Body.String() != "streaming content" {
		t.Errorf("expected 'streaming content', got '%s'", rec.Body.String())
	}
}

func TestBlob(t *testing.T) {
	s := New(nil)
	s.GET("/blob", func(w http.ResponseWriter, r *http.Request) {
		Blob(w, http.StatusOK, MIMEApplicationOctetStream, []byte{0x01, 0x02, 0x03})
	})

	req := httptest.NewRequest(http.MethodGet, "/blob", nil)
	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}

	contentType := rec.Header().Get("Content-Type")
	if contentType != MIMEApplicationOctetStream {
		t.Errorf("expected Content-Type %s, got %s", MIMEApplicationOctetStream, contentType)
	}
}

func TestErrorResponses(t *testing.T) {
	tests := []struct {
		name     string
		handler  func(w http.ResponseWriter)
		expected int
	}{
		{"BadRequest", func(w http.ResponseWriter) { BadRequest(w, "bad") }, 400},
		{"Unauthorized", func(w http.ResponseWriter) { Unauthorized(w, "unauth") }, 401},
		{"Forbidden", func(w http.ResponseWriter) { Forbidden(w, "forbidden") }, 403},
		{"NotFound", func(w http.ResponseWriter) { NotFound(w, "not found") }, 404},
		{"InternalServerError", func(w http.ResponseWriter) { InternalServerError(w, "error") }, 500},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s := New(nil)
			s.GET("/test", func(w http.ResponseWriter, r *http.Request) {
				tc.handler(w)
			})

			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			rec := httptest.NewRecorder()
			s.ServeHTTP(rec, req)

			if rec.Code != tc.expected {
				t.Errorf("expected status %d, got %d", tc.expected, rec.Code)
			}
		})
	}
}

func TestGenericHandler(t *testing.T) {
	type GetUserRequest struct {
		ID int `path:"id"`
	}

	type User struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	s := New(nil)
	s.GET("/users/{id}", Handle(func(ctx context.Context, req GetUserRequest) (*User, error) {
		return &User{ID: req.ID, Name: "John"}, nil
	}))

	req := httptest.NewRequest(http.MethodGet, "/users/42", nil)
	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}

	body := strings.TrimSpace(rec.Body.String())
	if !strings.Contains(body, `"id":42`) {
		t.Errorf("expected body to contain id:42, got %s", body)
	}
}

func TestGenericHandlerWithError(t *testing.T) {
	type Request struct{}

	s := New(nil)
	s.GET("/error", Handle(func(ctx context.Context, req Request) (any, error) {
		return nil, ErrNotFound.WithDetailf("resource not found")
	}))

	req := httptest.NewRequest(http.MethodGet, "/error", nil)
	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", rec.Code)
	}

	contentType := rec.Header().Get("Content-Type")
	if !strings.Contains(contentType, "application/problem+json") {
		t.Errorf("expected Content-Type application/problem+json, got %s", contentType)
	}
}

func TestGenericHandlerWithJSONBody(t *testing.T) {
	type CreateUserRequest struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	type User struct {
		ID    int    `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	s := New(nil)
	s.POST("/users", Handle(func(ctx context.Context, req CreateUserRequest) (*User, error) {
		return &User{ID: 1, Name: req.Name, Email: req.Email}, nil
	}))

	body := `{"name":"John","email":"john@example.com"}`
	req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}

	respBody := strings.TrimSpace(rec.Body.String())
	if !strings.Contains(respBody, `"name":"John"`) {
		t.Errorf("expected body to contain name:John, got %s", respBody)
	}
}

func TestHandleNoRequest(t *testing.T) {
	type Health struct {
		Status string `json:"status"`
	}

	s := New(nil)
	s.GET("/health", HandleNoRequest(func(ctx context.Context) (*Health, error) {
		return &Health{Status: "ok"}, nil
	}))

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}
}

func TestHandleNoResponse(t *testing.T) {
	type DeleteRequest struct {
		ID int `path:"id"`
	}

	s := New(nil)
	s.DELETE("/items/{id}", HandleNoResponse(func(ctx context.Context, req DeleteRequest) error {
		return nil
	}))

	req := httptest.NewRequest(http.MethodDelete, "/items/123", nil)
	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Errorf("expected status 204, got %d", rec.Code)
	}
}

func TestHandleEmpty(t *testing.T) {
	s := New(nil)
	s.POST("/ping", HandleEmpty(func(ctx context.Context) error {
		return nil
	}))

	req := httptest.NewRequest(http.MethodPost, "/ping", nil)
	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Errorf("expected status 204, got %d", rec.Code)
	}
}

func TestBindJSON(t *testing.T) {
	type Data struct {
		Name  string `json:"name"`
		Value int    `json:"value"`
	}

	body := `{"name":"test","value":42}`
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))

	data, err := BindJSON[Data](req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if data.Name != "test" {
		t.Errorf("expected name 'test', got '%s'", data.Name)
	}
	if data.Value != 42 {
		t.Errorf("expected value 42, got %d", data.Value)
	}
}

func TestBindQuery(t *testing.T) {
	type Search struct {
		Query  string `query:"q"`
		Page   int    `query:"page"`
		Active bool   `query:"active"`
	}

	req := httptest.NewRequest(http.MethodGet, "/?q=hello&page=2&active=true", nil)

	search, err := BindQuery[Search](req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if search.Query != "hello" {
		t.Errorf("expected query 'hello', got '%s'", search.Query)
	}
	if search.Page != 2 {
		t.Errorf("expected page 2, got %d", search.Page)
	}
	if !search.Active {
		t.Error("expected active true")
	}
}

func TestBindPath(t *testing.T) {
	type GetUser struct {
		ID int `path:"id"`
	}

	s := New(nil)
	var boundReq GetUser

	s.GET("/users/{id}", func(w http.ResponseWriter, r *http.Request) {
		req, err := BindPath[GetUser](r)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		boundReq = req
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/users/99", nil)
	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, req)

	if boundReq.ID != 99 {
		t.Errorf("expected ID 99, got %d", boundReq.ID)
	}
}

func TestBindHeader(t *testing.T) {
	type Headers struct {
		APIKey string `header:"X-API-Key"`
		Auth   string `header:"Authorization"`
	}

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-API-Key", "secret123")
	req.Header.Set("Authorization", "Bearer token")

	headers, err := BindHeader[Headers](req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if headers.APIKey != "secret123" {
		t.Errorf("expected APIKey 'secret123', got '%s'", headers.APIKey)
	}
	if headers.Auth != "Bearer token" {
		t.Errorf("expected Auth 'Bearer token', got '%s'", headers.Auth)
	}
}

func TestBindWithAllSources(t *testing.T) {
	type Request struct {
		ID     int    `path:"id"`
		Page   int    `query:"page"`
		APIKey string `header:"X-API-Key"`
		Name   string `json:"name"`
	}

	s := New(nil)
	var boundReq Request

	s.POST("/items/{id}", func(w http.ResponseWriter, r *http.Request) {
		req, err := Bind[Request](r)
		if err != nil {
			BadRequest(w, err.Error())
			return
		}
		boundReq = req
		w.WriteHeader(http.StatusOK)
	})

	body := `{"name":"test"}`
	req := httptest.NewRequest(http.MethodPost, "/items/42?page=5", strings.NewReader(body))
	req.Header.Set("X-API-Key", "key123")
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	if boundReq.ID != 42 {
		t.Errorf("expected ID 42, got %d", boundReq.ID)
	}
	if boundReq.Page != 5 {
		t.Errorf("expected Page 5, got %d", boundReq.Page)
	}
	if boundReq.APIKey != "key123" {
		t.Errorf("expected APIKey 'key123', got '%s'", boundReq.APIKey)
	}
	if boundReq.Name != "test" {
		t.Errorf("expected Name 'test', got '%s'", boundReq.Name)
	}
}

func TestRouteGroups(t *testing.T) {
	s := New(nil)

	api := s.Group("/api")
	api.GET("/users", func(w http.ResponseWriter, r *http.Request) {
		Text(w, http.StatusOK, "users")
	})

	v1 := api.Group("/v1")
	v1.GET("/posts", func(w http.ResponseWriter, r *http.Request) {
		Text(w, http.StatusOK, "posts")
	})

	tests := []struct {
		path     string
		expected string
	}{
		{"/api/users", "users"},
		{"/api/v1/posts", "posts"},
	}

	for _, tc := range tests {
		t.Run(tc.path, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tc.path, nil)
			rec := httptest.NewRecorder()
			s.ServeHTTP(rec, req)

			if rec.Code != http.StatusOK {
				t.Errorf("expected status 200, got %d", rec.Code)
			}
			if rec.Body.String() != tc.expected {
				t.Errorf("expected '%s', got '%s'", tc.expected, rec.Body.String())
			}
		})
	}
}

func TestRouteGroupWithMiddleware(t *testing.T) {
	s := New(nil)

	headerValue := ""
	authMW := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			headerValue = r.Header.Get("Authorization")
			next.ServeHTTP(w, r)
		})
	}

	api := s.Group("/api", authMW)
	api.GET("/protected", func(w http.ResponseWriter, r *http.Request) {
		Text(w, http.StatusOK, "ok")
	})

	req := httptest.NewRequest(http.MethodGet, "/api/protected", nil)
	req.Header.Set("Authorization", "Bearer token123")
	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, req)

	if headerValue != "Bearer token123" {
		t.Errorf("expected header 'Bearer token123', got '%s'", headerValue)
	}
}

func TestProblem(t *testing.T) {
	p := ErrNotFound.WithDetailf("user %d not found", 123)

	if p.Status != 404 {
		t.Errorf("expected status 404, got %d", p.Status)
	}
	if p.Detail != "user 123 not found" {
		t.Errorf("expected detail 'user 123 not found', got '%s'", p.Detail)
	}
	if !strings.Contains(p.Error(), "Not Found") {
		t.Errorf("expected error to contain 'Not Found', got '%s'", p.Error())
	}
}

func TestProblemWithInstance(t *testing.T) {
	p := ErrBadRequest.WithInstance("/api/users/123")

	if p.Instance != "/api/users/123" {
		t.Errorf("expected instance '/api/users/123', got '%s'", p.Instance)
	}
}

func TestWriteProblem(t *testing.T) {
	rec := httptest.NewRecorder()
	p := ErrNotFound.WithDetailf("not found")

	WriteProblem(rec, p)

	if rec.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", rec.Code)
	}

	contentType := rec.Header().Get("Content-Type")
	if !strings.Contains(contentType, "application/problem+json") {
		t.Errorf("expected Content-Type application/problem+json, got %s", contentType)
	}
}

func TestValidatable(t *testing.T) {
	type ValidatedRequest struct {
		Email string `json:"email"`
	}

	// Note: This test shows the pattern - actual validation would require
	// implementing Validatable on a pointer receiver

	body := `{"email":"test@example.com"}`
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))

	data, err := BindAndValidate[ValidatedRequest](req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if data.Email != "test@example.com" {
		t.Errorf("expected email 'test@example.com', got '%s'", data.Email)
	}
}

func BenchmarkRouterLookup(b *testing.B) {
	s := New(nil)
	s.GET("/users", func(w http.ResponseWriter, r *http.Request) {})
	s.GET("/users/{id}", func(w http.ResponseWriter, r *http.Request) {})
	s.GET("/users/{id}/posts", func(w http.ResponseWriter, r *http.Request) {})
	s.GET("/users/{id}/posts/{postID}", func(w http.ResponseWriter, r *http.Request) {})

	req := httptest.NewRequest(http.MethodGet, "/users/123/posts/456", nil)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		rec := httptest.NewRecorder()
		s.ServeHTTP(rec, req)
	}
}

func BenchmarkJSONResponse(b *testing.B) {
	s := New(nil)
	s.GET("/json", func(w http.ResponseWriter, r *http.Request) {
		JSON(w, http.StatusOK, map[string]string{"message": "hello", "status": "ok"})
	})

	req := httptest.NewRequest(http.MethodGet, "/json", nil)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		rec := httptest.NewRecorder()
		s.ServeHTTP(rec, req)
	}
}

func BenchmarkMiddlewareChain(b *testing.B) {
	s := New(nil)

	// Add 5 middleware
	for i := 0; i < 5; i++ {
		s.Use(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				next.ServeHTTP(w, r)
			})
		})
	}

	s.GET("/test", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		rec := httptest.NewRecorder()
		s.ServeHTTP(rec, req)
	}
}

func BenchmarkBind(b *testing.B) {
	type Request struct {
		ID   int    `query:"id"`
		Name string `query:"name"`
		Page int    `query:"page"`
	}

	req := httptest.NewRequest(http.MethodGet, "/?id=123&name=test&page=5", nil)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		Bind[Request](req)
	}
}

func BenchmarkBindJSON(b *testing.B) {
	type Request struct {
		Name  string `json:"name"`
		Email string `json:"email"`
		Age   int    `json:"age"`
	}

	body := `{"name":"John","email":"john@example.com","age":30}`

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
		BindJSON[Request](req)
	}
}

func BenchmarkGenericHandler(b *testing.B) {
	type Req struct {
		ID int `path:"id"`
	}
	type Res struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	s := New(nil)
	s.GET("/users/{id}", Handle(func(ctx context.Context, req Req) (*Res, error) {
		return &Res{ID: req.ID, Name: "John"}, nil
	}))

	req := httptest.NewRequest(http.MethodGet, "/users/123", nil)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		rec := httptest.NewRecorder()
		s.ServeHTTP(rec, req)
	}
}

func BenchmarkParamExtraction(b *testing.B) {
	s := New(nil)
	s.GET("/users/{id}/posts/{postID}/comments/{commentID}", func(w http.ResponseWriter, r *http.Request) {
		Param(r, "id")
		Param(r, "postID")
		Param(r, "commentID")
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/users/1/posts/2/comments/3", nil)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		rec := httptest.NewRecorder()
		s.ServeHTTP(rec, req)
	}
}

// Test for TLS options
func TestWithTLS(t *testing.T) {
	s := New(&Options{
		TLSCertFile: "cert.pem",
		TLSKeyFile:  "key.pem",
	})

	cfg := s.GetConfig()
	if cfg.TLSCertFile != "cert.pem" {
		t.Errorf("expected cert.pem, got %s", cfg.TLSCertFile)
	}
	if cfg.TLSKeyFile != "key.pem" {
		t.Errorf("expected key.pem, got %s", cfg.TLSKeyFile)
	}
}

// Test for MaxHeaderBytes option
func TestWithMaxHeaderBytes(t *testing.T) {
	s := New(&Options{
		MaxHeaderBytes: 1 << 20,
	})

	cfg := s.GetConfig()
	if cfg.MaxHeaderBytes != 1<<20 {
		t.Errorf("expected maxHeaderBytes 1MB, got %d", cfg.MaxHeaderBytes)
	}
}

// Test static file serving pattern
func TestStaticRoutePattern(t *testing.T) {
	s := New(nil)

	// Just test that Static doesn't panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Static panicked: %v", r)
		}
	}()

	s.Static("/assets/", ".")
}

// Test Any method
func TestAny(t *testing.T) {
	s := New(nil)
	called := 0

	s.Any("/any", func(w http.ResponseWriter, r *http.Request) {
		called++
		w.WriteHeader(http.StatusOK)
	})

	methods := []string{
		http.MethodGet,
		http.MethodPost,
		http.MethodPut,
		http.MethodPatch,
		http.MethodDelete,
	}

	for _, method := range methods {
		req := httptest.NewRequest(method, "/any", nil)
		rec := httptest.NewRecorder()
		s.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("%s: expected status 200, got %d", method, rec.Code)
		}
	}

	if called != len(methods) {
		t.Errorf("expected %d calls, got %d", len(methods), called)
	}
}

// Test panic on invalid route patterns
func TestInvalidRoutePatterns(t *testing.T) {
	tests := []struct {
		name    string
		pattern string
	}{
		{"empty pattern", ""},
		{"no leading slash", "users"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r == nil {
					t.Error("expected panic")
				}
			}()

			s := New(nil)
			s.GET(tc.pattern, func(w http.ResponseWriter, r *http.Request) {})
		})
	}
}

// Test nil handler panic
func TestNilHandlerPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for nil handler")
		}
	}()

	s := New(nil)
	s.GET("/test", nil)
}

func TestJSONPretty(t *testing.T) {
	rec := httptest.NewRecorder()
	data := map[string]string{"key": "value"}

	JSONPretty(rec, http.StatusOK, data, "  ")

	body := rec.Body.String()
	if !strings.Contains(body, "  ") {
		t.Errorf("expected indented JSON, got '%s'", body)
	}
}

func TestCreatedResponse(t *testing.T) {
	rec := httptest.NewRecorder()
	Created(rec, map[string]int{"id": 1})

	if rec.Code != http.StatusCreated {
		t.Errorf("expected status 201, got %d", rec.Code)
	}
}

func TestAcceptedResponse(t *testing.T) {
	rec := httptest.NewRecorder()
	Accepted(rec, map[string]string{"status": "processing"})

	if rec.Code != http.StatusAccepted {
		t.Errorf("expected status 202, got %d", rec.Code)
	}
}

func TestOKResponse(t *testing.T) {
	rec := httptest.NewRecorder()
	OK(rec, map[string]string{"status": "ok"})

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}
}

func TestAttachment(t *testing.T) {
	rec := httptest.NewRecorder()
	Attachment(rec, "file.pdf")

	disposition := rec.Header().Get("Content-Disposition")
	if disposition != `attachment; filename="file.pdf"` {
		t.Errorf("expected attachment disposition, got '%s'", disposition)
	}
}

func TestInline(t *testing.T) {
	rec := httptest.NewRecorder()
	Inline(rec, "image.png")

	disposition := rec.Header().Get("Content-Disposition")
	if disposition != `inline; filename="image.png"` {
		t.Errorf("expected inline disposition, got '%s'", disposition)
	}
}

// Test for param not found
func TestParamNotFound(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/test", nil)

	_, err := ParamInt(req, "nonexistent")
	if err == nil {
		t.Error("expected error for nonexistent param")
	}

	_, err = ParamInt64(req, "nonexistent")
	if err == nil {
		t.Error("expected error for nonexistent param")
	}

	_, err = ParamUUID(req, "nonexistent")
	if err == nil {
		t.Error("expected error for nonexistent param")
	}
}

// Test QueryDefault
func TestQueryDefault(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/?existing=value", nil)

	val := QueryDefault(req, "existing", "default")
	if val != "value" {
		t.Errorf("expected 'value', got '%s'", val)
	}

	val = QueryDefault(req, "missing", "default")
	if val != "default" {
		t.Errorf("expected 'default', got '%s'", val)
	}
}

// Test QueryInt64
func TestQueryInt64(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/?big=9223372036854775807", nil)

	val := QueryInt64(req, "big", 0)
	if val != 9223372036854775807 {
		t.Errorf("expected max int64, got %d", val)
	}

	val = QueryInt64(req, "missing", 42)
	if val != 42 {
		t.Errorf("expected 42, got %d", val)
	}
}

// Test QueryFloat64
func TestQueryFloat64(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/?price=19.99", nil)

	val := QueryFloat64(req, "price", 0)
	if val != 19.99 {
		t.Errorf("expected 19.99, got %f", val)
	}

	val = QueryFloat64(req, "missing", 9.99)
	if val != 9.99 {
		t.Errorf("expected 9.99, got %f", val)
	}
}

// Integration test for full request flow
func TestFullRequestFlow(t *testing.T) {
	type CreateOrderRequest struct {
		CustomerID int    `path:"customerID"`
		Product    string `json:"product"`
		Quantity   int    `json:"quantity"`
		Priority   string `query:"priority"`
		APIKey     string `header:"X-API-Key"`
	}

	type Order struct {
		ID         int    `json:"id"`
		CustomerID int    `json:"customer_id"`
		Product    string `json:"product"`
		Quantity   int    `json:"quantity"`
		Priority   string `json:"priority"`
	}

	s := New(nil)

	s.POST("/customers/{customerID}/orders", Handle(func(ctx context.Context, req CreateOrderRequest) (*Order, error) {
		if req.APIKey == "" {
			return nil, ErrUnauthorized.WithDetailf("API key required")
		}
		return &Order{
			ID:         1,
			CustomerID: req.CustomerID,
			Product:    req.Product,
			Quantity:   req.Quantity,
			Priority:   req.Priority,
		}, nil
	}))

	body := `{"product":"Widget","quantity":5}`
	req := httptest.NewRequest(http.MethodPost, "/customers/42/orders?priority=high", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", "secret123")

	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
		t.Logf("body: %s", rec.Body.String())
	}

	respBody := rec.Body.String()
	if !strings.Contains(respBody, `"customer_id":42`) {
		t.Errorf("expected customer_id:42, got %s", respBody)
	}
	if !strings.Contains(respBody, `"product":"Widget"`) {
		t.Errorf("expected product:Widget, got %s", respBody)
	}
	if !strings.Contains(respBody, `"priority":"high"`) {
		t.Errorf("expected priority:high, got %s", respBody)
	}
}

// Test file download helper
func TestFileServing(t *testing.T) {
	s := New(nil)
	s.GET("/download", func(w http.ResponseWriter, r *http.Request) {
		Attachment(w, "report.pdf")
		Blob(w, http.StatusOK, "application/pdf", []byte("fake pdf content"))
	})

	req := httptest.NewRequest(http.MethodGet, "/download", nil)
	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, req)

	disposition := rec.Header().Get("Content-Disposition")
	if disposition != `attachment; filename="report.pdf"` {
		t.Errorf("expected attachment disposition, got '%s'", disposition)
	}
}

// Test streaming response
func TestStreamingResponse(t *testing.T) {
	s := New(nil)
	s.GET("/stream", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		for i := 0; i < 3; i++ {
			io.WriteString(w, "data: event\n\n")
		}
	})

	req := httptest.NewRequest(http.MethodGet, "/stream", nil)
	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}

	body := rec.Body.String()
	if strings.Count(body, "data: event") != 3 {
		t.Errorf("expected 3 events, got body: %s", body)
	}
}
