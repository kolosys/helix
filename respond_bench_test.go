package helix_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/kolosys/helix"
)

// BenchmarkJSONAllocations benchmarks allocations in JSON encoding.
func BenchmarkJSONAllocations(b *testing.B) {
	data := map[string]any{
		"id":    1,
		"name":  "Test",
		"value": "Test Value",
		"items": []string{"item1", "item2", "item3"},
	}

	b.ReportAllocs()
	for b.Loop() {
		w := httptest.NewRecorder()
		JSON(w, http.StatusOK, data)
	}
}

// BenchmarkJSONEncoding benchmarks JSON encoding performance.
func BenchmarkJSONEncoding(b *testing.B) {
	data := map[string]any{
		"id":    1,
		"name":  "Test",
		"value": "Test Value",
		"items": []string{"item1", "item2", "item3"},
	}

	for b.Loop() {
		w := httptest.NewRecorder()
		JSON(w, http.StatusOK, data)
	}
}

// BenchmarkJSONLargePayload benchmarks JSON encoding with large payloads.
func BenchmarkJSONLargePayload(b *testing.B) {
	data := make(map[string]any)
	items := make([]string, 1000)
	for i := range items {
		items[i] = "item" + string(rune('a'+i%26))
	}
	data["items"] = items

	b.ReportAllocs()
	for b.Loop() {
		w := httptest.NewRecorder()
		JSON(w, http.StatusOK, data)
	}
}

// BenchmarkBufferPool benchmarks buffer pool efficiency.
func BenchmarkBufferPool(b *testing.B) {
	data := map[string]string{"test": "value"}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			w := httptest.NewRecorder()
			JSON(w, http.StatusOK, data)
		}
	})
}
