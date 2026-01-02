package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/kolosys/helix/middleware"
)

func TestAPI(t *testing.T) {
	bundle := API()

	if len(bundle) != 3 {
		t.Errorf("expected 3 middleware in API bundle, got %d", len(bundle))
	}

	// Test that the bundle can be applied
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Apply middleware in reverse order
	var h http.Handler = handler
	for i := len(bundle) - 1; i >= 0; i-- {
		h = bundle[i](h)
	}

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}

	// Check that RequestID header is set
	if rec.Header().Get("X-Request-ID") == "" {
		t.Error("expected X-Request-ID header to be set")
	}
}

func TestAPIWithCORS(t *testing.T) {
	config := CORSConfig{
		AllowOrigins: []string{"https://example.com"},
	}
	bundle := APIWithCORS(config)

	if len(bundle) != 3 {
		t.Errorf("expected 3 middleware in APIWithCORS bundle, got %d", len(bundle))
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	var h http.Handler = handler
	for i := len(bundle) - 1; i >= 0; i-- {
		h = bundle[i](h)
	}

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Origin", "https://example.com")
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}

	// Check CORS header
	if rec.Header().Get("Access-Control-Allow-Origin") != "https://example.com" {
		t.Errorf("expected Access-Control-Allow-Origin to be 'https://example.com', got '%s'",
			rec.Header().Get("Access-Control-Allow-Origin"))
	}
}

func TestWeb(t *testing.T) {
	bundle := Web()

	if len(bundle) != 3 {
		t.Errorf("expected 3 middleware in Web bundle, got %d", len(bundle))
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	var h http.Handler = handler
	for i := len(bundle) - 1; i >= 0; i-- {
		h = bundle[i](h)
	}

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}
}

func TestMinimal(t *testing.T) {
	bundle := Minimal()

	if len(bundle) != 1 {
		t.Errorf("expected 1 middleware in Minimal bundle, got %d", len(bundle))
	}

	// Test that Recover is included by causing a panic
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("test panic")
	})

	var h http.Handler = handler
	for i := len(bundle) - 1; i >= 0; i-- {
		h = bundle[i](h)
	}

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()

	// This should not panic due to Recover middleware
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("expected status 500 after panic, got %d", rec.Code)
	}
}

func TestProduction(t *testing.T) {
	bundle := Production()

	if len(bundle) != 2 {
		t.Errorf("expected 2 middleware in Production bundle, got %d", len(bundle))
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	var h http.Handler = handler
	for i := len(bundle) - 1; i >= 0; i-- {
		h = bundle[i](h)
	}

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}
}

func TestDevelopment(t *testing.T) {
	bundle := Development()

	if len(bundle) != 2 {
		t.Errorf("expected 2 middleware in Development bundle, got %d", len(bundle))
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	var h http.Handler = handler
	for i := len(bundle) - 1; i >= 0; i-- {
		h = bundle[i](h)
	}

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}
}

func TestSecure(t *testing.T) {
	bundle := Secure(100.0, 10)

	if len(bundle) != 3 {
		t.Errorf("expected 3 middleware in Secure bundle, got %d", len(bundle))
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	var h http.Handler = handler
	for i := len(bundle) - 1; i >= 0; i-- {
		h = bundle[i](h)
	}

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}
}

func TestBundleComposition(t *testing.T) {
	// Test that bundles can be combined
	bundle1 := Minimal()
	bundle2 := []Middleware{RequestID()}

	combined := make([]Middleware, 0, len(bundle1)+len(bundle2))
	combined = append(combined, bundle1...)
	combined = append(combined, bundle2...)

	if len(combined) != 2 {
		t.Errorf("expected 2 middleware in combined bundle, got %d", len(combined))
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	var h http.Handler = handler
	for i := len(combined) - 1; i >= 0; i-- {
		h = combined[i](h)
	}

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}
}
