package helix_test

import (
	"errors"
	"net/http/httptest"
	"strings"
	"testing"

	. "github.com/kolosys/helix"
)

func TestBindRequired(t *testing.T) {
	type Request struct {
		Name string `query:"name,required"`
	}

	// Missing required field
	req := httptest.NewRequest("GET", "/", nil)
	_, err := Bind[Request](req)
	if err == nil {
		t.Error("expected error for missing required field")
	}
}

func TestBindTypes(t *testing.T) {
	type Request struct {
		Int     int     `query:"int"`
		Int64   int64   `query:"int64"`
		Uint    uint    `query:"uint"`
		Float64 float64 `query:"float64"`
		Bool    bool    `query:"bool"`
		String  string  `query:"string"`
	}

	req := httptest.NewRequest("GET", "/?int=42&int64=9223372036854775807&uint=100&float64=3.14&bool=true&string=hello", nil)

	result, err := Bind[Request](req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Int != 42 {
		t.Errorf("expected Int 42, got %d", result.Int)
	}
	if result.Int64 != 9223372036854775807 {
		t.Errorf("expected Int64 max, got %d", result.Int64)
	}
	if result.Uint != 100 {
		t.Errorf("expected Uint 100, got %d", result.Uint)
	}
	if result.Float64 != 3.14 {
		t.Errorf("expected Float64 3.14, got %f", result.Float64)
	}
	if !result.Bool {
		t.Error("expected Bool true")
	}
	if result.String != "hello" {
		t.Errorf("expected String 'hello', got '%s'", result.String)
	}
}

func TestBindPointer(t *testing.T) {
	type Request struct {
		Value *int `query:"value"`
	}

	req := httptest.NewRequest("GET", "/?value=42", nil)

	result, err := Bind[Request](req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Value == nil {
		t.Fatal("expected Value to be non-nil")
	}
	if *result.Value != 42 {
		t.Errorf("expected Value 42, got %d", *result.Value)
	}
}

func TestBindSlice(t *testing.T) {
	type Request struct {
		Tags []string `query:"tags"`
	}

	req := httptest.NewRequest("GET", "/?tags=a,b,c", nil)

	result, err := Bind[Request](req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.Tags) != 3 {
		t.Errorf("expected 3 tags, got %d", len(result.Tags))
	}
}

func TestBindInvalidInt(t *testing.T) {
	type Request struct {
		Value int `query:"value"`
	}

	req := httptest.NewRequest("GET", "/?value=not-a-number", nil)

	_, err := Bind[Request](req)
	if err == nil {
		t.Error("expected error for invalid int")
	}
}

func TestBindFormData(t *testing.T) {
	type Request struct {
		Name string `form:"name"`
	}

	body := strings.NewReader("name=John")
	req := httptest.NewRequest("POST", "/", body)
	req.Header.Set("Content-Type", MIMEApplicationForm)

	result, err := Bind[Request](req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Name != "John" {
		t.Errorf("expected Name 'John', got '%s'", result.Name)
	}
}

func TestBindJSONEmpty(t *testing.T) {
	type Request struct {
		Name string `json:"name"`
	}

	req := httptest.NewRequest("POST", "/", nil)

	_, err := BindJSON[Request](req)
	if err == nil {
		t.Error("expected error for nil body")
	}
}

func TestBindJSONInvalid(t *testing.T) {
	type Request struct {
		Name string `json:"name"`
	}

	req := httptest.NewRequest("POST", "/", strings.NewReader("not json"))

	_, err := BindJSON[Request](req)
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}

func TestBindPathNotFound(t *testing.T) {
	type Request struct {
		ID int `path:"id,required"`
	}

	req := httptest.NewRequest("GET", "/", nil)

	_, err := BindPath[Request](req)
	if err == nil {
		t.Error("expected error for missing path param")
	}
}

func TestBindHeaderNotFound(t *testing.T) {
	type Request struct {
		Token string `header:"Authorization,required"`
	}

	req := httptest.NewRequest("GET", "/", nil)

	_, err := BindHeader[Request](req)
	if err == nil {
		t.Error("expected error for missing header")
	}
}

func TestBindBoolVariants(t *testing.T) {
	type Request struct {
		Active bool `query:"active"`
	}

	tests := []struct {
		query    string
		expected bool
	}{
		{"?active=true", true},
		{"?active=false", false},
		{"?active=1", true},
		{"?active=0", false},
		{"?active=yes", true},
		{"?active=no", false},
		{"?active=on", true},
		{"?active=off", false},
	}

	for _, tc := range tests {
		t.Run(tc.query, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/"+tc.query, nil)
			result, err := Bind[Request](req)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if result.Active != tc.expected {
				t.Errorf("expected %v, got %v", tc.expected, result.Active)
			}
		})
	}
}

func TestBindOmitEmpty(t *testing.T) {
	type Request struct {
		Name string `query:"name,omitempty"`
	}

	req := httptest.NewRequest("GET", "/", nil)

	result, err := Bind[Request](req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Name != "" {
		t.Errorf("expected empty name, got '%s'", result.Name)
	}
}

func TestBindIgnoreUnexported(t *testing.T) {
	type Request struct {
		Name   string `query:"name"`
		secret string `query:"secret"` // unexported, should be ignored
	}

	req := httptest.NewRequest("GET", "/?name=John&secret=hidden", nil)

	result, err := Bind[Request](req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Name != "John" {
		t.Errorf("expected Name 'John', got '%s'", result.Name)
	}
	if result.secret != "" {
		t.Errorf("expected secret to be empty, got '%s'", result.secret)
	}
}

func TestBindTagIgnore(t *testing.T) {
	type Request struct {
		Name    string `query:"name"`
		Ignored string `query:"-"`
	}

	req := httptest.NewRequest("GET", "/?name=John&Ignored=value", nil)

	result, err := Bind[Request](req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Ignored != "" {
		t.Errorf("expected Ignored to be empty, got '%s'", result.Ignored)
	}
}

func TestBindUnsupportedType(t *testing.T) {
	type Nested struct {
		Value string
	}
	type Request struct {
		Nested Nested `query:"nested"`
	}

	req := httptest.NewRequest("GET", "/?nested=value", nil)

	_, err := Bind[Request](req)
	if err == nil {
		t.Error("expected error for unsupported type")
	}
}

func TestBindQueryMissingNonRequired(t *testing.T) {
	type Request struct {
		Name string `query:"name"`
		Age  int    `query:"age"`
	}

	req := httptest.NewRequest("GET", "/?name=John", nil)

	result, err := BindQuery[Request](req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Name != "John" {
		t.Errorf("expected Name 'John', got '%s'", result.Name)
	}
	if result.Age != 0 {
		t.Errorf("expected Age 0, got %d", result.Age)
	}
}

func BenchmarkBindSimple(b *testing.B) {
	type Request struct {
		Name string `query:"name"`
		Age  int    `query:"age"`
	}

	req := httptest.NewRequest("GET", "/?name=John&age=30", nil)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		Bind[Request](req)
	}
}

func BenchmarkBindComplex(b *testing.B) {
	type Request struct {
		ID     int     `query:"id"`
		Name   string  `query:"name"`
		Age    int     `query:"age"`
		Active bool    `query:"active"`
		Score  float64 `query:"score"`
	}

	req := httptest.NewRequest("GET", "/?id=1&name=John&age=30&active=true&score=95.5", nil)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		Bind[Request](req)
	}
}

// -----------------------------------------------------------------------------
// Tests for Improved Error Messages
// -----------------------------------------------------------------------------

func TestBindingError_Format(t *testing.T) {
	type Request struct {
		Name string `query:"name,required"`
	}

	req := httptest.NewRequest("GET", "/", nil)
	_, err := Bind[Request](req)

	if err == nil {
		t.Fatal("expected error for missing required field")
	}

	errStr := err.Error()
	// Check that error contains field name
	if !strings.Contains(errStr, "name") {
		t.Errorf("error should contain field name 'name', got: %s", errStr)
	}
	// Check that error contains source
	if !strings.Contains(errStr, "query") {
		t.Errorf("error should contain source 'query', got: %s", errStr)
	}
	// Check that error contains hint
	if !strings.Contains(errStr, "hint") {
		t.Errorf("error should contain hint, got: %s", errStr)
	}
}

func TestBindingError_PathRequired(t *testing.T) {
	type Request struct {
		ID int `path:"id,required"`
	}

	req := httptest.NewRequest("GET", "/", nil)
	_, err := BindPath[Request](req)

	if err == nil {
		t.Fatal("expected error for missing required path param")
	}

	errStr := err.Error()
	if !strings.Contains(errStr, "path") {
		t.Errorf("error should contain source 'path', got: %s", errStr)
	}
}

func TestBindingError_HeaderRequired(t *testing.T) {
	type Request struct {
		Token string `header:"Authorization,required"`
	}

	req := httptest.NewRequest("GET", "/", nil)
	_, err := BindHeader[Request](req)

	if err == nil {
		t.Fatal("expected error for missing required header")
	}

	errStr := err.Error()
	if !strings.Contains(errStr, "header") {
		t.Errorf("error should contain source 'header', got: %s", errStr)
	}
}

func TestBindingError_InvalidValue(t *testing.T) {
	type Request struct {
		Age int `query:"age"`
	}

	req := httptest.NewRequest("GET", "/?age=not-a-number", nil)
	_, err := Bind[Request](req)

	if err == nil {
		t.Fatal("expected error for invalid value")
	}

	errStr := err.Error()
	// Check that error contains field name
	if !strings.Contains(errStr, "age") {
		t.Errorf("error should contain field name 'age', got: %s", errStr)
	}
	// Check that error contains source
	if !strings.Contains(errStr, "query") {
		t.Errorf("error should contain source 'query', got: %s", errStr)
	}
	// Check that error contains hint about type
	if !strings.Contains(errStr, "hint") {
		t.Errorf("error should contain hint, got: %s", errStr)
	}
}

func TestBindingError_Unwrap(t *testing.T) {
	type Request struct {
		Name string `query:"name,required"`
	}

	req := httptest.NewRequest("GET", "/", nil)
	_, err := Bind[Request](req)

	if err == nil {
		t.Fatal("expected error")
	}

	// Check that we can unwrap to the sentinel error
	var bindErr *BindingError
	if !errors.As(err, &bindErr) {
		t.Error("expected error to be *BindingError")
	}

	if !errors.Is(bindErr.Unwrap(), ErrRequiredField) {
		t.Error("expected wrapped error to be ErrRequiredField")
	}
}

func TestBindingError_InvalidFieldUnwrap(t *testing.T) {
	type Request struct {
		Value int `query:"value"`
	}

	req := httptest.NewRequest("GET", "/?value=abc", nil)
	_, err := Bind[Request](req)

	if err == nil {
		t.Fatal("expected error")
	}

	var bindErr *BindingError
	if !errors.As(err, &bindErr) {
		t.Error("expected error to be *BindingError")
	}

	if !errors.Is(bindErr.Unwrap(), ErrInvalidFieldValue) {
		t.Error("expected wrapped error to be ErrInvalidFieldValue")
	}
}
