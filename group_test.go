package helix_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/kolosys/helix"
)

func TestGroupBasic(t *testing.T) {
	s := New(nil)

	api := s.Group("/api")
	api.GET("/users", func(w http.ResponseWriter, r *http.Request) {
		Text(w, http.StatusOK, "users")
	})

	req := httptest.NewRequest(http.MethodGet, "/api/users", nil)
	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}
	if rec.Body.String() != "users" {
		t.Errorf("expected 'users', got '%s'", rec.Body.String())
	}
}

func TestGroupNested(t *testing.T) {
	s := New(nil)

	api := s.Group("/api")
	v1 := api.Group("/v1")
	v1.GET("/posts", func(w http.ResponseWriter, r *http.Request) {
		Text(w, http.StatusOK, "v1-posts")
	})

	v2 := api.Group("/v2")
	v2.GET("/posts", func(w http.ResponseWriter, r *http.Request) {
		Text(w, http.StatusOK, "v2-posts")
	})

	tests := []struct {
		path     string
		expected string
	}{
		{"/api/v1/posts", "v1-posts"},
		{"/api/v2/posts", "v2-posts"},
	}

	for _, tc := range tests {
		t.Run(tc.path, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tc.path, nil)
			rec := httptest.NewRecorder()
			s.ServeHTTP(rec, req)

			if rec.Body.String() != tc.expected {
				t.Errorf("expected '%s', got '%s'", tc.expected, rec.Body.String())
			}
		})
	}
}

func TestGroupMiddleware(t *testing.T) {
	s := New(nil)

	var middlewareCalled bool
	mw := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			middlewareCalled = true
			next.ServeHTTP(w, r)
		})
	}

	api := s.Group("/api", mw)
	api.GET("/test", func(w http.ResponseWriter, r *http.Request) {
		Text(w, http.StatusOK, "ok")
	})

	req := httptest.NewRequest(http.MethodGet, "/api/test", nil)
	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, req)

	if !middlewareCalled {
		t.Error("middleware was not called")
	}
}

func TestGroupUse(t *testing.T) {
	s := New(nil)

	var order []string

	api := s.Group("/api")
	api.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			order = append(order, "mw1")
			next.ServeHTTP(w, r)
		})
	})
	api.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			order = append(order, "mw2")
			next.ServeHTTP(w, r)
		})
	})

	api.GET("/test", func(w http.ResponseWriter, r *http.Request) {
		order = append(order, "handler")
		Text(w, http.StatusOK, "ok")
	})

	req := httptest.NewRequest(http.MethodGet, "/api/test", nil)
	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, req)

	expected := []string{"mw1", "mw2", "handler"}
	if len(order) != len(expected) {
		t.Fatalf("expected %d calls, got %d", len(expected), len(order))
	}
	for i, v := range expected {
		if order[i] != v {
			t.Errorf("expected order[%d] = %s, got %s", i, v, order[i])
		}
	}
}

func TestGroupNestedMiddleware(t *testing.T) {
	s := New(nil)

	var order []string

	api := s.Group("/api")
	api.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			order = append(order, "api-mw")
			next.ServeHTTP(w, r)
		})
	})

	v1 := api.Group("/v1")
	v1.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			order = append(order, "v1-mw")
			next.ServeHTTP(w, r)
		})
	})

	v1.GET("/test", func(w http.ResponseWriter, r *http.Request) {
		order = append(order, "handler")
		Text(w, http.StatusOK, "ok")
	})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/test", nil)
	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, req)

	expected := []string{"api-mw", "v1-mw", "handler"}
	if len(order) != len(expected) {
		t.Fatalf("expected %d calls, got %d", len(expected), len(order))
	}
	for i, v := range expected {
		if order[i] != v {
			t.Errorf("expected order[%d] = %s, got %s", i, v, order[i])
		}
	}
}

func TestGroupAllMethods(t *testing.T) {
	s := New(nil)

	api := s.Group("/api")

	api.GET("/get", func(w http.ResponseWriter, r *http.Request) { Text(w, 200, "GET") })
	api.POST("/post", func(w http.ResponseWriter, r *http.Request) { Text(w, 200, "POST") })
	api.PUT("/put", func(w http.ResponseWriter, r *http.Request) { Text(w, 200, "PUT") })
	api.PATCH("/patch", func(w http.ResponseWriter, r *http.Request) { Text(w, 200, "PATCH") })
	api.DELETE("/delete", func(w http.ResponseWriter, r *http.Request) { Text(w, 200, "DELETE") })
	api.OPTIONS("/options", func(w http.ResponseWriter, r *http.Request) { Text(w, 200, "OPTIONS") })
	api.HEAD("/head", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(200) })

	tests := []struct {
		method string
		path   string
	}{
		{http.MethodGet, "/api/get"},
		{http.MethodPost, "/api/post"},
		{http.MethodPut, "/api/put"},
		{http.MethodPatch, "/api/patch"},
		{http.MethodDelete, "/api/delete"},
		{http.MethodOptions, "/api/options"},
		{http.MethodHead, "/api/head"},
	}

	for _, tc := range tests {
		t.Run(tc.method, func(t *testing.T) {
			req := httptest.NewRequest(tc.method, tc.path, nil)
			rec := httptest.NewRecorder()
			s.ServeHTTP(rec, req)

			if rec.Code != http.StatusOK {
				t.Errorf("expected status 200, got %d", rec.Code)
			}
		})
	}
}

func TestGroupAny(t *testing.T) {
	s := New(nil)

	api := s.Group("/api")
	api.Any("/any", func(w http.ResponseWriter, r *http.Request) {
		Text(w, http.StatusOK, r.Method)
	})

	methods := []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete}

	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			req := httptest.NewRequest(method, "/api/any", nil)
			rec := httptest.NewRecorder()
			s.ServeHTTP(rec, req)

			if rec.Code != http.StatusOK {
				t.Errorf("expected status 200, got %d", rec.Code)
			}
			if rec.Body.String() != method {
				t.Errorf("expected '%s', got '%s'", method, rec.Body.String())
			}
		})
	}
}

func TestGroupHandle(t *testing.T) {
	s := New(nil)

	api := s.Group("/api")
	api.Handle(http.MethodGet, "/custom", func(w http.ResponseWriter, r *http.Request) {
		Text(w, http.StatusOK, "custom")
	})

	req := httptest.NewRequest(http.MethodGet, "/api/custom", nil)
	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, req)

	if rec.Body.String() != "custom" {
		t.Errorf("expected 'custom', got '%s'", rec.Body.String())
	}
}

func TestGroupStatic(t *testing.T) {
	s := New(nil)

	api := s.Group("/api")

	// Just test that Static doesn't panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Static panicked: %v", r)
		}
	}()

	api.Static("/files/", ".")
}

func TestGroupWithParams(t *testing.T) {
	s := New(nil)

	users := s.Group("/users")
	users.GET("/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := Param(r, "id")
		Text(w, http.StatusOK, "user:"+id)
	})

	users.GET("/{id}/posts/{postID}", func(w http.ResponseWriter, r *http.Request) {
		id := Param(r, "id")
		postID := Param(r, "postID")
		Text(w, http.StatusOK, "user:"+id+",post:"+postID)
	})

	tests := []struct {
		path     string
		expected string
	}{
		{"/users/123", "user:123"},
		{"/users/123/posts/456", "user:123,post:456"},
	}

	for _, tc := range tests {
		t.Run(tc.path, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tc.path, nil)
			rec := httptest.NewRecorder()
			s.ServeHTTP(rec, req)

			if rec.Body.String() != tc.expected {
				t.Errorf("expected '%s', got '%s'", tc.expected, rec.Body.String())
			}
		})
	}
}

func TestGroupNoMiddleware(t *testing.T) {
	s := New(nil)

	// Group without middleware
	api := s.Group("/api")
	api.GET("/test", func(w http.ResponseWriter, r *http.Request) {
		Text(w, http.StatusOK, "ok")
	})

	req := httptest.NewRequest(http.MethodGet, "/api/test", nil)
	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}
}

func BenchmarkGroup(b *testing.B) {
	s := New(nil)

	api := s.Group("/api")
	v1 := api.Group("/v1")
	v1.GET("/users", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/users", nil)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		rec := httptest.NewRecorder()
		s.ServeHTTP(rec, req)
	}
}

func BenchmarkGroupWithMiddleware(b *testing.B) {
	s := New(nil)

	mw := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r)
		})
	}

	api := s.Group("/api", mw)
	v1 := api.Group("/v1", mw)
	v1.GET("/users", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/users", nil)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		rec := httptest.NewRecorder()
		s.ServeHTTP(rec, req)
	}
}
