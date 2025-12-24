package helix_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	. "github.com/kolosys/helix"
)

func TestResourceBuilder_List(t *testing.T) {
	s := New(nil)
	called := false

	s.Resource("/users").List(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, req)

	if !called {
		t.Error("List handler was not called")
	}
	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}
}

func TestResourceBuilder_Create(t *testing.T) {
	s := New(nil)
	called := false

	s.Resource("/users").Create(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusCreated)
	})

	req := httptest.NewRequest(http.MethodPost, "/users", nil)
	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, req)

	if !called {
		t.Error("Create handler was not called")
	}
	if rec.Code != http.StatusCreated {
		t.Errorf("expected status 201, got %d", rec.Code)
	}
}

func TestResourceBuilder_Get(t *testing.T) {
	s := New(nil)
	var gotID string

	s.Resource("/users").Get(func(w http.ResponseWriter, r *http.Request) {
		gotID = Param(r, "id")
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/users/123", nil)
	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, req)

	if gotID != "123" {
		t.Errorf("expected id '123', got '%s'", gotID)
	}
}

func TestResourceBuilder_Update(t *testing.T) {
	s := New(nil)
	var gotID string

	s.Resource("/users").Update(func(w http.ResponseWriter, r *http.Request) {
		gotID = Param(r, "id")
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodPut, "/users/456", nil)
	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, req)

	if gotID != "456" {
		t.Errorf("expected id '456', got '%s'", gotID)
	}
}

func TestResourceBuilder_Patch(t *testing.T) {
	s := New(nil)
	called := false

	s.Resource("/users").Patch(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodPatch, "/users/789", nil)
	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, req)

	if !called {
		t.Error("Patch handler was not called")
	}
}

func TestResourceBuilder_Delete(t *testing.T) {
	s := New(nil)
	var gotID string

	s.Resource("/users").Delete(func(w http.ResponseWriter, r *http.Request) {
		gotID = Param(r, "id")
		w.WriteHeader(http.StatusNoContent)
	})

	req := httptest.NewRequest(http.MethodDelete, "/users/999", nil)
	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, req)

	if gotID != "999" {
		t.Errorf("expected id '999', got '%s'", gotID)
	}
	if rec.Code != http.StatusNoContent {
		t.Errorf("expected status 204, got %d", rec.Code)
	}
}

func TestResourceBuilder_Custom(t *testing.T) {
	s := New(nil)
	called := false

	s.Resource("/users").Custom(http.MethodPost, "/{id}/archive", func(w http.ResponseWriter, r *http.Request) {
		called = true
		id := Param(r, "id")
		Text(w, http.StatusOK, "archived "+id)
	})

	req := httptest.NewRequest(http.MethodPost, "/users/123/archive", nil)
	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, req)

	if !called {
		t.Error("Custom handler was not called")
	}
	if !strings.Contains(rec.Body.String(), "archived 123") {
		t.Errorf("expected body to contain 'archived 123', got '%s'", rec.Body.String())
	}
}

func TestResourceBuilder_CRUD(t *testing.T) {
	s := New(nil)
	calls := make(map[string]bool)

	s.Resource("/posts").CRUD(
		func(w http.ResponseWriter, r *http.Request) { calls["list"] = true; w.WriteHeader(http.StatusOK) },
		func(w http.ResponseWriter, r *http.Request) {
			calls["create"] = true
			w.WriteHeader(http.StatusCreated)
		},
		func(w http.ResponseWriter, r *http.Request) { calls["get"] = true; w.WriteHeader(http.StatusOK) },
		func(w http.ResponseWriter, r *http.Request) { calls["update"] = true; w.WriteHeader(http.StatusOK) },
		func(w http.ResponseWriter, r *http.Request) {
			calls["delete"] = true
			w.WriteHeader(http.StatusNoContent)
		},
	)

	tests := []struct {
		method   string
		path     string
		expected string
	}{
		{http.MethodGet, "/posts", "list"},
		{http.MethodPost, "/posts", "create"},
		{http.MethodGet, "/posts/1", "get"},
		{http.MethodPut, "/posts/1", "update"},
		{http.MethodDelete, "/posts/1", "delete"},
	}

	for _, tc := range tests {
		req := httptest.NewRequest(tc.method, tc.path, nil)
		rec := httptest.NewRecorder()
		s.ServeHTTP(rec, req)

		if !calls[tc.expected] {
			t.Errorf("%s %s: expected '%s' handler to be called", tc.method, tc.path, tc.expected)
		}
	}
}

func TestResourceBuilder_ReadOnly(t *testing.T) {
	s := New(nil)
	calls := make(map[string]bool)

	s.Resource("/items").ReadOnly(
		func(w http.ResponseWriter, r *http.Request) { calls["list"] = true; w.WriteHeader(http.StatusOK) },
		func(w http.ResponseWriter, r *http.Request) { calls["get"] = true; w.WriteHeader(http.StatusOK) },
	)

	// Test list
	req := httptest.NewRequest(http.MethodGet, "/items", nil)
	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, req)
	if !calls["list"] {
		t.Error("list handler was not called")
	}

	// Test get
	req = httptest.NewRequest(http.MethodGet, "/items/1", nil)
	rec = httptest.NewRecorder()
	s.ServeHTTP(rec, req)
	if !calls["get"] {
		t.Error("get handler was not called")
	}

	// Test that POST is not registered
	req = httptest.NewRequest(http.MethodPost, "/items", nil)
	rec = httptest.NewRecorder()
	s.ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Errorf("expected POST to return 404, got %d", rec.Code)
	}
}

func TestResourceBuilder_Chaining(t *testing.T) {
	s := New(nil)
	calls := make(map[string]bool)

	s.Resource("/articles").
		List(func(w http.ResponseWriter, r *http.Request) { calls["list"] = true }).
		Create(func(w http.ResponseWriter, r *http.Request) { calls["create"] = true }).
		Get(func(w http.ResponseWriter, r *http.Request) { calls["get"] = true }).
		Update(func(w http.ResponseWriter, r *http.Request) { calls["update"] = true }).
		Delete(func(w http.ResponseWriter, r *http.Request) { calls["delete"] = true }).
		Custom(http.MethodPost, "/{id}/publish", func(w http.ResponseWriter, r *http.Request) { calls["publish"] = true })

	tests := []struct {
		method   string
		path     string
		expected string
	}{
		{http.MethodGet, "/articles", "list"},
		{http.MethodPost, "/articles", "create"},
		{http.MethodGet, "/articles/1", "get"},
		{http.MethodPut, "/articles/1", "update"},
		{http.MethodDelete, "/articles/1", "delete"},
		{http.MethodPost, "/articles/1/publish", "publish"},
	}

	for _, tc := range tests {
		req := httptest.NewRequest(tc.method, tc.path, nil)
		rec := httptest.NewRecorder()
		s.ServeHTTP(rec, req)

		if !calls[tc.expected] {
			t.Errorf("%s %s: expected '%s' handler to be called", tc.method, tc.path, tc.expected)
		}
	}
}

func TestResourceBuilder_WithMiddleware(t *testing.T) {
	s := New(nil)
	middlewareCalled := false

	mw := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			middlewareCalled = true
			next.ServeHTTP(w, r)
		})
	}

	s.Resource("/protected", mw).Get(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/protected/1", nil)
	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, req)

	if !middlewareCalled {
		t.Error("middleware was not called")
	}
}

func TestResourceBuilder_Aliases(t *testing.T) {
	s := New(nil)
	calls := make(map[string]bool)

	s.Resource("/widgets").
		Index(func(w http.ResponseWriter, r *http.Request) { calls["index"] = true }).
		Store(func(w http.ResponseWriter, r *http.Request) { calls["store"] = true }).
		Show(func(w http.ResponseWriter, r *http.Request) { calls["show"] = true }).
		Destroy(func(w http.ResponseWriter, r *http.Request) { calls["destroy"] = true })

	tests := []struct {
		method   string
		path     string
		expected string
	}{
		{http.MethodGet, "/widgets", "index"},
		{http.MethodPost, "/widgets", "store"},
		{http.MethodGet, "/widgets/1", "show"},
		{http.MethodDelete, "/widgets/1", "destroy"},
	}

	for _, tc := range tests {
		req := httptest.NewRequest(tc.method, tc.path, nil)
		rec := httptest.NewRecorder()
		s.ServeHTTP(rec, req)

		if !calls[tc.expected] {
			t.Errorf("%s %s: expected '%s' handler to be called", tc.method, tc.path, tc.expected)
		}
	}
}

func TestGroupResource(t *testing.T) {
	s := New(nil)
	var gotID string

	api := s.Group("/api/v1")
	api.Resource("/users").Get(func(w http.ResponseWriter, r *http.Request) {
		gotID = Param(r, "id")
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/users/42", nil)
	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, req)

	if gotID != "42" {
		t.Errorf("expected id '42', got '%s'", gotID)
	}
	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}
}

func TestServer_Routes(t *testing.T) {
	s := New(nil)

	s.GET("/users", func(w http.ResponseWriter, r *http.Request) {})
	s.POST("/users", func(w http.ResponseWriter, r *http.Request) {})
	s.GET("/users/{id}", func(w http.ResponseWriter, r *http.Request) {})

	routes := s.Routes()

	if len(routes) != 3 {
		t.Errorf("expected 3 routes, got %d", len(routes))
	}

	// Check that routes are present
	routeMap := make(map[string]bool)
	for _, r := range routes {
		routeMap[r.Method+" "+r.Pattern] = true
	}

	expectedRoutes := []string{
		"GET /users",
		"POST /users",
		"GET /users/{id}",
	}

	for _, expected := range expectedRoutes {
		if !routeMap[expected] {
			t.Errorf("expected route '%s' not found", expected)
		}
	}
}

func TestServer_PrintRoutes(t *testing.T) {
	s := New(nil)

	s.GET("/users", func(w http.ResponseWriter, r *http.Request) {})
	s.POST("/users", func(w http.ResponseWriter, r *http.Request) {})
	s.DELETE("/users/{id}", func(w http.ResponseWriter, r *http.Request) {})

	var buf bytes.Buffer
	s.PrintRoutes(&buf)

	output := buf.String()

	if !strings.Contains(output, "GET") {
		t.Error("expected output to contain GET")
	}
	if !strings.Contains(output, "POST") {
		t.Error("expected output to contain POST")
	}
	if !strings.Contains(output, "DELETE") {
		t.Error("expected output to contain DELETE")
	}
	if !strings.Contains(output, "/users") {
		t.Error("expected output to contain /users")
	}
}
