package helix_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/kolosys/helix"
)

func TestRouterStaticRoutes(t *testing.T) {
	r := NewRouter()

	r.Handle(http.MethodGet, "/", func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte("root"))
	})

	r.Handle(http.MethodGet, "/users", func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte("users"))
	})

	r.Handle(http.MethodGet, "/users/list", func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte("list"))
	})

	tests := []struct {
		path     string
		expected string
	}{
		{"/", "root"},
		{"/users", "users"},
		{"/users/list", "list"},
	}

	for _, tc := range tests {
		t.Run(tc.path, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tc.path, nil)
			rec := httptest.NewRecorder()
			r.ServeHTTP(rec, req)

			if rec.Body.String() != tc.expected {
				t.Errorf("expected '%s', got '%s'", tc.expected, rec.Body.String())
			}
		})
	}
}

func TestRouterParamRoutes(t *testing.T) {
	r := NewRouter()

	r.Handle(http.MethodGet, "/users/{id}", func(w http.ResponseWriter, req *http.Request) {
		id := Param(req, "id")
		w.Write([]byte("user:" + id))
	})

	r.Handle(http.MethodGet, "/posts/{postID}/comments/{commentID}", func(w http.ResponseWriter, req *http.Request) {
		postID := Param(req, "postID")
		commentID := Param(req, "commentID")
		w.Write([]byte("post:" + postID + ",comment:" + commentID))
	})

	tests := []struct {
		path     string
		expected string
	}{
		{"/users/123", "user:123"},
		{"/users/abc", "user:abc"},
		{"/posts/1/comments/2", "post:1,comment:2"},
	}

	for _, tc := range tests {
		t.Run(tc.path, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tc.path, nil)
			rec := httptest.NewRecorder()
			r.ServeHTTP(rec, req)

			if rec.Body.String() != tc.expected {
				t.Errorf("expected '%s', got '%s'", tc.expected, rec.Body.String())
			}
		})
	}
}

func TestRouterCatchAllRoutes(t *testing.T) {
	r := NewRouter()

	r.Handle(http.MethodGet, "/files/{path...}", func(w http.ResponseWriter, req *http.Request) {
		path := Param(req, "path")
		w.Write([]byte("file:" + path))
	})

	tests := []struct {
		path     string
		expected string
	}{
		{"/files/readme.txt", "file:readme.txt"},
		{"/files/docs/api.md", "file:docs/api.md"},
		{"/files/a/b/c/d", "file:a/b/c/d"},
	}

	for _, tc := range tests {
		t.Run(tc.path, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tc.path, nil)
			rec := httptest.NewRecorder()
			r.ServeHTTP(rec, req)

			if rec.Body.String() != tc.expected {
				t.Errorf("expected '%s', got '%s'", tc.expected, rec.Body.String())
			}
		})
	}
}

func TestRouterMixedRoutes(t *testing.T) {
	r := NewRouter()

	r.Handle(http.MethodGet, "/api/users", func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte("all-users"))
	})

	r.Handle(http.MethodGet, "/api/users/{id}", func(w http.ResponseWriter, req *http.Request) {
		id := Param(req, "id")
		w.Write([]byte("user:" + id))
	})

	r.Handle(http.MethodGet, "/api/users/{id}/posts", func(w http.ResponseWriter, req *http.Request) {
		id := Param(req, "id")
		w.Write([]byte("posts-for:" + id))
	})

	tests := []struct {
		path     string
		expected string
	}{
		{"/api/users", "all-users"},
		{"/api/users/123", "user:123"},
		{"/api/users/123/posts", "posts-for:123"},
	}

	for _, tc := range tests {
		t.Run(tc.path, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tc.path, nil)
			rec := httptest.NewRecorder()
			r.ServeHTTP(rec, req)

			if rec.Body.String() != tc.expected {
				t.Errorf("expected '%s', got '%s'", tc.expected, rec.Body.String())
			}
		})
	}
}

func TestRouterNotFound(t *testing.T) {
	r := NewRouter()

	r.Handle(http.MethodGet, "/exists", func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte("ok"))
	})

	req := httptest.NewRequest(http.MethodGet, "/not-exists", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", rec.Code)
	}
}

func TestRouterMethodNotFound(t *testing.T) {
	r := NewRouter()

	r.Handle(http.MethodGet, "/resource", func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte("ok"))
	})

	req := httptest.NewRequest(http.MethodPost, "/resource", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", rec.Code)
	}
}

func TestRouterDifferentMethods(t *testing.T) {
	r := NewRouter()

	r.Handle(http.MethodGet, "/resource", func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte("GET"))
	})

	r.Handle(http.MethodPost, "/resource", func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte("POST"))
	})

	r.Handle(http.MethodPut, "/resource", func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte("PUT"))
	})

	tests := []struct {
		method   string
		expected string
	}{
		{http.MethodGet, "GET"},
		{http.MethodPost, "POST"},
		{http.MethodPut, "PUT"},
	}

	for _, tc := range tests {
		t.Run(tc.method, func(t *testing.T) {
			req := httptest.NewRequest(tc.method, "/resource", nil)
			rec := httptest.NewRecorder()
			r.ServeHTTP(rec, req)

			if rec.Body.String() != tc.expected {
				t.Errorf("expected '%s', got '%s'", tc.expected, rec.Body.String())
			}
		})
	}
}

func TestRouterPanics(t *testing.T) {
	t.Run("empty pattern", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic")
			}
		}()
		r := NewRouter()
		r.Handle(http.MethodGet, "", nil)
	})

	t.Run("no leading slash", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic")
			}
		}()
		r := NewRouter()
		r.Handle(http.MethodGet, "users", nil)
	})

	t.Run("nil handler", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic")
			}
		}()
		r := NewRouter()
		r.Handle(http.MethodGet, "/users", nil)
	})

	t.Run("duplicate route", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic")
			}
		}()
		r := NewRouter()
		r.Handle(http.MethodGet, "/users", func(w http.ResponseWriter, req *http.Request) {})
		r.Handle(http.MethodGet, "/users", func(w http.ResponseWriter, req *http.Request) {})
	})
}

func TestParsePattern(t *testing.T) {
	tests := []struct {
		pattern  string
		segments int
	}{
		{"/", 0},
		{"/users", 1},
		{"/users/{id}", 2},
		{"/users/{id}/posts", 3},
		{"/files/{path...}", 2},
	}

	for _, tc := range tests {
		t.Run(tc.pattern, func(t *testing.T) {
			segs := ParsePattern(tc.pattern)
			if len(segs) != tc.segments {
				t.Errorf("expected %d segments, got %d", tc.segments, len(segs))
			}
		})
	}
}

func TestParamsPool(t *testing.T) {
	r := NewRouter()

	r.Handle(http.MethodGet, "/users/{id}", func(w http.ResponseWriter, req *http.Request) {
		id := Param(req, "id")
		w.Write([]byte(id))
	})

	// Make many requests to test pool reuse
	for i := 0; i < 100; i++ {
		req := httptest.NewRequest(http.MethodGet, "/users/123", nil)
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		if rec.Body.String() != "123" {
			t.Fatalf("iteration %d: expected '123', got '%s'", i, rec.Body.String())
		}
	}
}

func BenchmarkRouterStatic(b *testing.B) {
	r := NewRouter()
	r.Handle(http.MethodGet, "/users", func(w http.ResponseWriter, req *http.Request) {})

	req := httptest.NewRequest(http.MethodGet, "/users", nil)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)
	}
}

func BenchmarkRouterParam(b *testing.B) {
	r := NewRouter()
	r.Handle(http.MethodGet, "/users/{id}", func(w http.ResponseWriter, req *http.Request) {
		Param(req, "id")
	})

	req := httptest.NewRequest(http.MethodGet, "/users/123", nil)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)
	}
}

func BenchmarkRouterMultiParam(b *testing.B) {
	r := NewRouter()
	r.Handle(http.MethodGet, "/users/{userID}/posts/{postID}/comments/{commentID}", func(w http.ResponseWriter, req *http.Request) {
		Param(req, "userID")
		Param(req, "postID")
		Param(req, "commentID")
	})

	req := httptest.NewRequest(http.MethodGet, "/users/1/posts/2/comments/3", nil)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)
	}
}
