package helix_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	. "github.com/kolosys/helix"
)

func TestCtx_Param(t *testing.T) {
	s := New()
	var gotID string

	s.GET("/users/{id}", HandleCtx(func(c *Ctx) error {
		gotID = c.Param("id")
		return c.OK(map[string]string{"id": gotID})
	}))

	req := httptest.NewRequest(http.MethodGet, "/users/123", nil)
	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, req)

	if gotID != "123" {
		t.Errorf("expected id '123', got '%s'", gotID)
	}
	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}
}

func TestCtx_ParamInt(t *testing.T) {
	s := New()
	var gotID int

	s.GET("/users/{id}", HandleCtx(func(c *Ctx) error {
		id, err := c.ParamInt("id")
		if err != nil {
			return BadRequestf("invalid id")
		}
		gotID = id
		return c.OK(map[string]int{"id": id})
	}))

	req := httptest.NewRequest(http.MethodGet, "/users/42", nil)
	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, req)

	if gotID != 42 {
		t.Errorf("expected id 42, got %d", gotID)
	}
}

func TestCtx_ParamInt_Invalid(t *testing.T) {
	s := New()

	s.GET("/users/{id}", HandleCtx(func(c *Ctx) error {
		_, err := c.ParamInt("id")
		if err != nil {
			return BadRequestf("invalid id")
		}
		return c.OK(nil)
	}))

	req := httptest.NewRequest(http.MethodGet, "/users/abc", nil)
	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", rec.Code)
	}
}

func TestCtx_Query(t *testing.T) {
	s := New()
	var gotQuery string

	s.GET("/search", HandleCtx(func(c *Ctx) error {
		gotQuery = c.Query("q")
		return c.OK(map[string]string{"query": gotQuery})
	}))

	req := httptest.NewRequest(http.MethodGet, "/search?q=hello", nil)
	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, req)

	if gotQuery != "hello" {
		t.Errorf("expected query 'hello', got '%s'", gotQuery)
	}
}

func TestCtx_QueryDefault(t *testing.T) {
	s := New()
	var gotQuery string

	s.GET("/search", HandleCtx(func(c *Ctx) error {
		gotQuery = c.QueryDefault("q", "default")
		return c.OK(map[string]string{"query": gotQuery})
	}))

	req := httptest.NewRequest(http.MethodGet, "/search", nil)
	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, req)

	if gotQuery != "default" {
		t.Errorf("expected query 'default', got '%s'", gotQuery)
	}
}

func TestCtx_QueryInt(t *testing.T) {
	s := New()
	var gotPage int

	s.GET("/list", HandleCtx(func(c *Ctx) error {
		gotPage = c.QueryInt("page", 1)
		return c.OK(map[string]int{"page": gotPage})
	}))

	req := httptest.NewRequest(http.MethodGet, "/list?page=5", nil)
	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, req)

	if gotPage != 5 {
		t.Errorf("expected page 5, got %d", gotPage)
	}
}

func TestCtx_QueryBool(t *testing.T) {
	s := New()
	var gotActive bool

	s.GET("/filter", HandleCtx(func(c *Ctx) error {
		gotActive = c.QueryBool("active")
		return c.OK(map[string]bool{"active": gotActive})
	}))

	req := httptest.NewRequest(http.MethodGet, "/filter?active=true", nil)
	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, req)

	if !gotActive {
		t.Error("expected active true")
	}
}

func TestCtx_QuerySlice(t *testing.T) {
	s := New()
	var gotTags []string

	s.GET("/tags", HandleCtx(func(c *Ctx) error {
		gotTags = c.QuerySlice("tag")
		return c.OK(map[string][]string{"tags": gotTags})
	}))

	req := httptest.NewRequest(http.MethodGet, "/tags?tag=go&tag=web", nil)
	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, req)

	if len(gotTags) != 2 || gotTags[0] != "go" || gotTags[1] != "web" {
		t.Errorf("expected tags [go, web], got %v", gotTags)
	}
}

func TestCtx_Header(t *testing.T) {
	s := New()
	var gotHeader string

	s.GET("/auth", HandleCtx(func(c *Ctx) error {
		gotHeader = c.Header("Authorization")
		return c.OK(map[string]string{"auth": gotHeader})
	}))

	req := httptest.NewRequest(http.MethodGet, "/auth", nil)
	req.Header.Set("Authorization", "Bearer token123")
	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, req)

	if gotHeader != "Bearer token123" {
		t.Errorf("expected header 'Bearer token123', got '%s'", gotHeader)
	}
}

func TestCtx_Bind(t *testing.T) {
	type CreateUser struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	s := New()
	var gotUser CreateUser

	s.POST("/users", HandleCtx(func(c *Ctx) error {
		if err := c.Bind(&gotUser); err != nil {
			return BadRequestf("invalid json")
		}
		return c.Created(gotUser)
	}))

	body := `{"name":"John","email":"john@example.com"}`
	req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(body))
	req.Header.Set("Content-Type", MIMEApplicationJSONCharsetUTF8)
	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Errorf("expected status 201, got %d", rec.Code)
	}
	if gotUser.Name != "John" {
		t.Errorf("expected name 'John', got '%s'", gotUser.Name)
	}
}

func TestCtx_SetHeader(t *testing.T) {
	s := New()

	s.GET("/custom", HandleCtx(func(c *Ctx) error {
		return c.SetHeader("X-Custom", "value").OK(map[string]string{"status": "ok"})
	}))

	req := httptest.NewRequest(http.MethodGet, "/custom", nil)
	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, req)

	if rec.Header().Get("X-Custom") != "value" {
		t.Errorf("expected X-Custom header 'value', got '%s'", rec.Header().Get("X-Custom"))
	}
}

func TestCtx_SetCookie(t *testing.T) {
	s := New()

	s.GET("/login", HandleCtx(func(c *Ctx) error {
		cookie := &http.Cookie{Name: "session", Value: "abc123"}
		return c.SetCookie(cookie).OK(map[string]string{"status": "logged in"})
	}))

	req := httptest.NewRequest(http.MethodGet, "/login", nil)
	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, req)

	cookies := rec.Result().Cookies()
	if len(cookies) != 1 || cookies[0].Name != "session" || cookies[0].Value != "abc123" {
		t.Errorf("expected session cookie, got %v", cookies)
	}
}

func TestCtx_JSON(t *testing.T) {
	s := New()

	s.GET("/json", HandleCtx(func(c *Ctx) error {
		return c.JSON(http.StatusOK, map[string]string{"message": "hello"})
	}))

	req := httptest.NewRequest(http.MethodGet, "/json", nil)
	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}
	contentType := rec.Header().Get("Content-Type")
	if !strings.HasPrefix(contentType, MIMEApplicationJSONCharsetUTF8) {
		t.Errorf("expected %s, got %s", MIMEApplicationJSONCharsetUTF8, contentType)
	}
}

func TestCtx_OK(t *testing.T) {
	s := New()

	s.GET("/ok", HandleCtx(func(c *Ctx) error {
		return c.OK(map[string]string{"status": "ok"})
	}))

	req := httptest.NewRequest(http.MethodGet, "/ok", nil)
	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}
}

func TestCtx_Created(t *testing.T) {
	s := New()

	s.POST("/create", HandleCtx(func(c *Ctx) error {
		return c.Created(map[string]int{"id": 1})
	}))

	req := httptest.NewRequest(http.MethodPost, "/create", nil)
	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Errorf("expected status 201, got %d", rec.Code)
	}
}

func TestCtx_NoContent(t *testing.T) {
	s := New()

	s.DELETE("/delete", HandleCtx(func(c *Ctx) error {
		return c.NoContent()
	}))

	req := httptest.NewRequest(http.MethodDelete, "/delete", nil)
	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Errorf("expected status 204, got %d", rec.Code)
	}
}

func TestCtx_Text(t *testing.T) {
	s := New()

	s.GET("/text", HandleCtx(func(c *Ctx) error {
		return c.Text(http.StatusOK, "hello world")
	}))

	req := httptest.NewRequest(http.MethodGet, "/text", nil)
	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, req)

	if rec.Body.String() != "hello world" {
		t.Errorf("expected 'hello world', got '%s'", rec.Body.String())
	}
}

func TestCtx_HTML(t *testing.T) {
	s := New()

	s.GET("/html", HandleCtx(func(c *Ctx) error {
		return c.HTML(http.StatusOK, "<h1>Hello</h1>")
	}))

	req := httptest.NewRequest(http.MethodGet, "/html", nil)
	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, req)

	contentType := rec.Header().Get("Content-Type")
	if !strings.HasPrefix(contentType, "text/html") {
		t.Errorf("expected text/html, got %s", contentType)
	}
}

func TestCtx_Problem(t *testing.T) {
	s := New()

	s.GET("/error", HandleCtx(func(c *Ctx) error {
		return c.Problem(ErrNotFound.WithDetailf("resource not found"))
	}))

	req := httptest.NewRequest(http.MethodGet, "/error", nil)
	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", rec.Code)
	}
	contentType := rec.Header().Get("Content-Type")
	if !strings.Contains(contentType, "application/problem+json") {
		t.Errorf("expected application/problem+json, got %s", contentType)
	}
}

func TestCtx_Redirect(t *testing.T) {
	s := New()

	s.GET("/old", HandleCtx(func(c *Ctx) error {
		c.Redirect("/new", http.StatusMovedPermanently)
		return nil
	}))

	req := httptest.NewRequest(http.MethodGet, "/old", nil)
	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, req)

	if rec.Code != http.StatusMovedPermanently {
		t.Errorf("expected status 301, got %d", rec.Code)
	}
	if rec.Header().Get("Location") != "/new" {
		t.Errorf("expected Location /new, got %s", rec.Header().Get("Location"))
	}
}

func TestCtx_ErrorResponse(t *testing.T) {
	s := New()

	s.GET("/bad", HandleCtx(func(c *Ctx) error {
		return c.BadRequest("invalid input")
	}))

	req := httptest.NewRequest(http.MethodGet, "/bad", nil)
	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", rec.Code)
	}
}

func TestCtx_Context(t *testing.T) {
	s := New()
	var gotCtx context.Context

	s.GET("/ctx", HandleCtx(func(c *Ctx) error {
		gotCtx = c.Context()
		return c.OK(nil)
	}))

	req := httptest.NewRequest(http.MethodGet, "/ctx", nil)
	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, req)

	if gotCtx == nil {
		t.Error("expected context to be non-nil")
	}
}

func TestCtx_Attachment(t *testing.T) {
	s := New()

	s.GET("/download", HandleCtx(func(c *Ctx) error {
		return c.Attachment("file.pdf").Blob(http.StatusOK, MIMEApplicationPDF, []byte("content"))
	}))

	req := httptest.NewRequest(http.MethodGet, "/download", nil)
	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, req)

	disposition := rec.Header().Get("Content-Disposition")
	if disposition != `attachment; filename="file.pdf"` {
		t.Errorf("expected attachment disposition, got '%s'", disposition)
	}
}

func TestCtx_ChainedSetHeader(t *testing.T) {
	s := New()

	s.GET("/headers", HandleCtx(func(c *Ctx) error {
		return c.SetHeader("X-One", "1").
			SetHeader("X-Two", "2").
			SetHeader("X-Three", "3").
			OK(map[string]string{"status": "ok"})
	}))

	req := httptest.NewRequest(http.MethodGet, "/headers", nil)
	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, req)

	if rec.Header().Get("X-One") != "1" {
		t.Errorf("expected X-One '1', got '%s'", rec.Header().Get("X-One"))
	}
	if rec.Header().Get("X-Two") != "2" {
		t.Errorf("expected X-Two '2', got '%s'", rec.Header().Get("X-Two"))
	}
	if rec.Header().Get("X-Three") != "3" {
		t.Errorf("expected X-Three '3', got '%s'", rec.Header().Get("X-Three"))
	}
}

func TestHandleCtx_ReturnsError(t *testing.T) {
	s := New()

	s.GET("/error", HandleCtx(func(c *Ctx) error {
		return NotFoundf("user %d not found", 123)
	}))

	req := httptest.NewRequest(http.MethodGet, "/error", nil)
	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", rec.Code)
	}
	body := rec.Body.String()
	if !strings.Contains(body, "user 123 not found") {
		t.Errorf("expected body to contain 'user 123 not found', got %s", body)
	}
}

func BenchmarkCtx_ParamAccess(b *testing.B) {
	s := New()
	s.GET("/users/{id}", HandleCtx(func(c *Ctx) error {
		c.Param("id")
		return c.NoContent()
	}))

	req := httptest.NewRequest(http.MethodGet, "/users/123", nil)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		rec := httptest.NewRecorder()
		s.ServeHTTP(rec, req)
	}
}

func BenchmarkCtx_QueryAccess(b *testing.B) {
	s := New()
	s.GET("/search", HandleCtx(func(c *Ctx) error {
		c.Query("q")
		c.QueryInt("page", 1)
		c.QueryBool("active")
		return c.NoContent()
	}))

	req := httptest.NewRequest(http.MethodGet, "/search?q=test&page=5&active=true", nil)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		rec := httptest.NewRecorder()
		s.ServeHTTP(rec, req)
	}
}
