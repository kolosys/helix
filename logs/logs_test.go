package logs_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"

	. "github.com/kolosys/helix/logs"
)

func TestNewLogger(t *testing.T) {
	log := New(nil)
	if log == nil {
		t.Fatal("expected logger")
	}
	if log.GetLevel() != InfoLevel {
		t.Errorf("expected InfoLevel, got %v", log.GetLevel())
	}
}

func TestLoggerWithOptions(t *testing.T) {
	buf := &bytes.Buffer{}
	log := New(&Options{
		Output:    buf,
		Level:     DebugLevel,
		Formatter: &TextFormatter{DisableTimestamp: true, DisableColors: true},
	})

	log.Debug("test message")

	output := buf.String()
	if !strings.Contains(output, "DEBG") {
		t.Errorf("expected DEBG in output, got: %s", output)
	}
	if !strings.Contains(output, "test message") {
		t.Errorf("expected 'test message' in output, got: %s", output)
	}
}

func TestLogLevels(t *testing.T) {
	tests := []struct {
		level  Level
		method func(*Logger, string, ...Field)
	}{
		{TraceLevel, (*Logger).Trace},
		{DebugLevel, (*Logger).Debug},
		{InfoLevel, (*Logger).Info},
		{WarnLevel, (*Logger).Warn},
		{ErrorLevel, (*Logger).Error},
	}

	for _, tc := range tests {
		t.Run(tc.level.String(), func(t *testing.T) {
			buf := &bytes.Buffer{}
			log := New(&Options{
				Output:    buf,
				Level:     TraceLevel,
				Formatter: &TextFormatter{DisableTimestamp: true, DisableColors: true},
			})

			tc.method(log, "test")

			output := buf.String()
			if !strings.Contains(output, tc.level.ShortString()) {
				t.Errorf("expected %s in output, got: %s", tc.level.ShortString(), output)
			}
		})
	}
}

func TestLogLevelFiltering(t *testing.T) {
	buf := &bytes.Buffer{}
	log := New(&Options{
		Output:    buf,
		Level:     WarnLevel,
		Formatter: &TextFormatter{DisableTimestamp: true},
	})

	log.Debug("debug message")
	log.Info("info message")
	log.Warn("warn message")

	output := buf.String()
	if strings.Contains(output, "debug") {
		t.Error("debug should be filtered")
	}
	if strings.Contains(output, "info") {
		t.Error("info should be filtered")
	}
	if !strings.Contains(output, "warn") {
		t.Error("warn should not be filtered")
	}
}

func TestLoggerFields(t *testing.T) {
	buf := &bytes.Buffer{}
	log := New(&Options{
		Output:    buf,
		Formatter: &TextFormatter{DisableTimestamp: true, DisableColors: true},
	})

	log.Info("test", String("key", "value"), Int("count", 42))

	output := buf.String()
	if !strings.Contains(output, "key=value") {
		t.Errorf("expected key=value in output, got: %s", output)
	}
	if !strings.Contains(output, "count=42") {
		t.Errorf("expected count=42 in output, got: %s", output)
	}
}

func TestLoggerWith(t *testing.T) {
	buf := &bytes.Buffer{}
	log := New(&Options{
		Output:    buf,
		Formatter: &TextFormatter{DisableTimestamp: true, DisableColors: true},
	})

	child := log.With(String("component", "test"))
	child.Info("message")

	output := buf.String()
	if !strings.Contains(output, "component=test") {
		t.Errorf("expected component=test in output, got: %s", output)
	}
}

func TestLoggerWithDefaultFields(t *testing.T) {
	buf := &bytes.Buffer{}
	log := New(&Options{
		Output:    buf,
		Fields:    []Field{String("app", "helix")},
		Formatter: &TextFormatter{DisableTimestamp: true, DisableColors: true},
	})

	log.Info("test")

	output := buf.String()
	if !strings.Contains(output, "app=helix") {
		t.Errorf("expected app=helix in output, got: %s", output)
	}
}

func TestLoggerWithCaller(t *testing.T) {
	buf := &bytes.Buffer{}
	log := New(&Options{
		Output:    buf,
		AddCaller: true,
		Formatter: &TextFormatter{DisableTimestamp: true, DisableColors: true},
	})

	log.Info("test")

	output := buf.String()
	// Caller info should be present (could be logs_test.go or testing.go depending on depth)
	if !strings.Contains(output, ".go:") {
		t.Errorf("expected caller info in output, got: %s", output)
	}
}

func TestJSONFormatter(t *testing.T) {
	buf := &bytes.Buffer{}
	log := New(&Options{
		Output:    buf,
		Formatter: &JSONFormatter{},
	})

	log.Info("test message", String("key", "value"), Int("count", 42))

	var entry map[string]any
	if err := json.Unmarshal(buf.Bytes(), &entry); err != nil {
		t.Fatalf("failed to parse JSON: %v, output: %s", err, buf.String())
	}

	if entry["msg"] != "test message" {
		t.Errorf("expected msg 'test message', got %v", entry["msg"])
	}
	if entry["key"] != "value" {
		t.Errorf("expected key 'value', got %v", entry["key"])
	}
	if entry["count"].(float64) != 42 {
		t.Errorf("expected count 42, got %v", entry["count"])
	}
}

func TestPrettyFormatter(t *testing.T) {
	buf := &bytes.Buffer{}
	log := New(&Options{
		Output:    buf,
		Formatter: &PrettyFormatter{ShowTimestamp: true},
	})

	log.Info("test message", String("key", "value"))

	output := buf.String()
	if !strings.Contains(output, "test message") {
		t.Errorf("expected message in output, got: %s", output)
	}
	// Field is present but with ANSI codes
	if !strings.Contains(output, "key") || !strings.Contains(output, "value") {
		t.Errorf("expected field in output, got: %s", output)
	}
}

func TestFieldTypes(t *testing.T) {
	tests := []struct {
		name     string
		field    Field
		expected string
	}{
		{"string", String("k", "v"), "v"},
		{"int", Int("k", 42), "42"},
		{"int64", Int64("k", 9223372036854775807), "9223372036854775807"},
		{"uint", Uint("k", 100), "100"},
		{"float64", Float64("k", 3.14), "3.14"},
		{"bool true", Bool("k", true), "true"},
		{"bool false", Bool("k", false), "false"},
		{"duration", Duration("k", time.Second), "1s"},
		{"error", Err(errors.New("test error")), "test error"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.field.StringValue() != tc.expected {
				t.Errorf("expected %s, got %s", tc.expected, tc.field.StringValue())
			}
		})
	}
}

func TestAnyField(t *testing.T) {
	tests := []struct {
		name     string
		value    any
		expected string
	}{
		{"string", "hello", "hello"},
		{"int", 42, "42"},
		{"bool", true, "true"},
		{"nil", nil, "null"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			field := Any("k", tc.value)
			if field.StringValue() != tc.expected {
				t.Errorf("expected %s, got %s", tc.expected, field.StringValue())
			}
		})
	}
}

func TestContextLogging(t *testing.T) {
	buf := &bytes.Buffer{}
	log := New(&Options{
		Output:    buf,
		Formatter: &TextFormatter{DisableTimestamp: true, DisableColors: true},
	})

	ctx := context.Background()
	ctx = WithContextFields(ctx, String("request_id", "abc123"))

	log.InfoContext(ctx, "request processed")

	output := buf.String()
	if !strings.Contains(output, "request_id=abc123") {
		t.Errorf("expected request_id in output, got: %s", output)
	}
}

func TestLoggerFromContext(t *testing.T) {
	buf := &bytes.Buffer{}
	log := New(&Options{
		Output:    buf,
		Formatter: &TextFormatter{DisableTimestamp: true, DisableColors: true},
	})

	ctx := context.Background()
	ctx = WithLogger(ctx, log)

	CtxInfo(ctx, "from context")

	output := buf.String()
	if !strings.Contains(output, "from context") {
		t.Errorf("expected message in output, got: %s", output)
	}
}

func TestHook(t *testing.T) {
	buf := &bytes.Buffer{}
	hookBuf := &bytes.Buffer{}

	hook := NewWriterHook(hookBuf, &TextFormatter{DisableTimestamp: true, DisableColors: true}, ErrorLevel)

	log := New(&Options{
		Output:    buf,
		Formatter: &TextFormatter{DisableTimestamp: true},
		Level:     TraceLevel,
		Hooks:     []Hook{hook},
	})

	log.Info("info message")
	log.Error("error message")

	hookOutput := hookBuf.String()
	if strings.Contains(hookOutput, "info message") {
		t.Errorf("hook should not capture info, got: %s", hookOutput)
	}
	if !strings.Contains(hookOutput, "error message") {
		t.Errorf("hook should capture error, got: %s", hookOutput)
	}
}

func TestMetricsHook(t *testing.T) {
	metrics := NewMetricsHook()

	log := New(&Options{
		Output: &bytes.Buffer{},
		Hooks:  []Hook{metrics},
	})

	log.Info("one")
	log.Info("two")
	log.Warn("warn")
	log.Error("error")

	if metrics.Count(InfoLevel) != 2 {
		t.Errorf("expected 2 info logs, got %d", metrics.Count(InfoLevel))
	}
	if metrics.Count(WarnLevel) != 1 {
		t.Errorf("expected 1 warn log, got %d", metrics.Count(WarnLevel))
	}
	if metrics.Count(ErrorLevel) != 1 {
		t.Errorf("expected 1 error log, got %d", metrics.Count(ErrorLevel))
	}
}

func TestErrorHook(t *testing.T) {
	errorHook := NewErrorHook(10)

	log := New(&Options{
		Output: &bytes.Buffer{},
		Level:  TraceLevel,
		Hooks:  []Hook{errorHook},
	})

	log.Info("info")
	log.Error("error 1")
	log.Error("error 2")

	errs := errorHook.Errors()
	if len(errs) != 2 {
		t.Errorf("expected 2 errors, got %d", len(errs))
	}
}

func TestSampler(t *testing.T) {
	buf := &bytes.Buffer{}
	log := New(&Options{
		Output:    buf,
		Formatter: &TextFormatter{DisableTimestamp: true},
		Sampler:   NewCountSampler(3), // Log every 3rd message
	})

	for i := 0; i < 9; i++ {
		log.Info("message")
	}

	// Should have 3 log entries (1st, 4th, 7th)
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 3 {
		t.Errorf("expected 3 lines, got %d", len(lines))
	}
}

func TestRateSampler(t *testing.T) {
	sampler := NewRateSampler(2, 100*time.Millisecond)

	// First 2 should pass
	if !sampler.Sample(InfoLevel, "msg") {
		t.Error("first should pass")
	}
	if !sampler.Sample(InfoLevel, "msg") {
		t.Error("second should pass")
	}

	// Third should fail
	if sampler.Sample(InfoLevel, "msg") {
		t.Error("third should be sampled out")
	}

	// Wait for window to reset
	time.Sleep(150 * time.Millisecond)

	// Should pass again
	if !sampler.Sample(InfoLevel, "msg") {
		t.Error("after reset should pass")
	}
}

func TestFirstNSampler(t *testing.T) {
	sampler := NewFirstNSampler(3)

	for i := 0; i < 3; i++ {
		if !sampler.Sample(InfoLevel, "msg") {
			t.Errorf("first %d should pass", i+1)
		}
	}

	if sampler.Sample(InfoLevel, "msg") {
		t.Error("4th should be sampled out")
	}
}

func TestAsyncLogging(t *testing.T) {
	buf := &safeBuffer{}
	log := New(&Options{
		Output:          buf,
		Formatter:       &TextFormatter{DisableTimestamp: true, DisableColors: true},
		AsyncBufferSize: 100,
	})

	for i := 0; i < 10; i++ {
		log.Info("async message", Int("i", i))
	}

	// Close waits for all async entries to be processed
	log.Close()

	output := buf.String()
	if output == "" {
		t.Error("expected output from async logging")
	}
	if !strings.Contains(output, "async message") {
		t.Errorf("expected async message in output, got: %s", output)
	}
}

// safeBuffer is a thread-safe buffer for testing async logging.
type safeBuffer struct {
	buf bytes.Buffer
	mu  sync.Mutex
}

func (b *safeBuffer) Write(p []byte) (n int, err error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.buf.Write(p)
}

func (b *safeBuffer) String() string {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.buf.String()
}

func TestLevelParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected Level
	}{
		{"trace", TraceLevel},
		{"debug", DebugLevel},
		{"info", InfoLevel},
		{"warn", WarnLevel},
		{"warning", WarnLevel},
		{"error", ErrorLevel},
		{"err", ErrorLevel},
		{"fatal", FatalLevel},
		{"panic", PanicLevel},
		{"TRACE", TraceLevel},
		{"INFO", InfoLevel},
		{"unknown", InfoLevel},
	}

	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			level := ParseLevel(tc.input)
			if level != tc.expected {
				t.Errorf("expected %v, got %v", tc.expected, level)
			}
		})
	}
}

func TestLevelString(t *testing.T) {
	for _, level := range AllLevels() {
		if level.String() == "" {
			t.Errorf("level %d has empty string", level)
		}
		if level.ShortString() == "" {
			t.Errorf("level %d has empty short string", level)
		}
	}
}

func TestIsEnabled(t *testing.T) {
	log := New(&Options{
		Level: WarnLevel,
	})

	if log.IsEnabled(InfoLevel) {
		t.Error("info should not be enabled")
	}
	if !log.IsEnabled(WarnLevel) {
		t.Error("warn should be enabled")
	}
	if !log.IsEnabled(ErrorLevel) {
		t.Error("error should be enabled")
	}
}

func TestDefaultLogger(t *testing.T) {
	buf := &bytes.Buffer{}
	log := New(&Options{
		Output:    buf,
		Formatter: &TextFormatter{DisableTimestamp: true},
	})
	SetDefault(log)

	Info("default logger")

	if !strings.Contains(buf.String(), "default logger") {
		t.Errorf("expected message in output, got: %s", buf.String())
	}

	// Reset default
	SetDefault(New(nil))
}

func TestWithHelper(t *testing.T) {
	buf := &bytes.Buffer{}
	log := New(&Options{
		Output:    buf,
		Formatter: &TextFormatter{DisableTimestamp: true, DisableColors: true},
	})
	SetDefault(log)

	child := With(String("component", "test"))
	child.Info("message")

	if !strings.Contains(buf.String(), "component=test") {
		t.Errorf("expected component in output, got: %s", buf.String())
	}

	SetDefault(New(nil))
}

func TestFilterHook(t *testing.T) {
	buf := &bytes.Buffer{}
	hook := NewFilterHook(
		NewWriterHook(buf, &TextFormatter{DisableTimestamp: true, DisableColors: true}),
		func(e *Entry) bool {
			return strings.Contains(e.Message, "important")
		},
	)

	log := New(&Options{
		Output: &bytes.Buffer{},
		Hooks:  []Hook{hook},
	})

	log.Info("normal message")
	log.Info("important message")

	if strings.Contains(buf.String(), "normal") {
		t.Error("normal should be filtered")
	}
	if !strings.Contains(buf.String(), "important") {
		t.Error("important should not be filtered")
	}
}

func TestFuncHook(t *testing.T) {
	var called bool
	hook := NewFuncHook(func(e *Entry) {
		called = true
	})

	log := New(&Options{
		Output: &bytes.Buffer{},
		Hooks:  []Hook{hook},
	})

	log.Info("test")

	if !called {
		t.Error("hook should have been called")
	}
}

func TestContextHelpers(t *testing.T) {
	ctx := context.Background()
	ctx = WithRequestID(ctx, "req-123")
	ctx = WithTraceID(ctx, "trace-456")
	ctx = WithUserID(ctx, "user-789")

	fields := FieldsFromContext(ctx)
	if len(fields) != 3 {
		t.Errorf("expected 3 fields, got %d", len(fields))
	}
}

func TestTimeField(t *testing.T) {
	now := time.Now()
	field := Time("created", now)

	if field.Type != FieldTypeTime {
		t.Errorf("expected FieldTypeTime, got %v", field.Type)
	}

	value := field.StringValue()
	if !strings.Contains(value, now.Format("2006")) {
		t.Errorf("expected time in output, got: %s", value)
	}
}

func TestBytesField(t *testing.T) {
	data := []byte("hello world")
	field := Bytes("data", data)

	if field.StringValue() != "hello world" {
		t.Errorf("expected 'hello world', got '%s'", field.StringValue())
	}
}

func TestJSONField(t *testing.T) {
	type Data struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
	field := JSON("data", Data{Name: "John", Age: 30})

	value := field.StringValue()
	if !strings.Contains(value, "John") {
		t.Errorf("expected JSON with John, got: %s", value)
	}
}

func TestEntryMethods(t *testing.T) {
	e := &Entry{
		Level:   InfoLevel,
		Message: "test",
		Fields:  []Field{String("key", "value")},
	}

	if !e.HasField("key") {
		t.Error("expected HasField to return true")
	}
	if e.HasField("missing") {
		t.Error("expected HasField to return false for missing")
	}

	field, ok := e.GetField("key")
	if !ok {
		t.Error("expected GetField to return true")
	}
	if field.StringValue() != "value" {
		t.Errorf("expected 'value', got '%s'", field.StringValue())
	}

	if e.GetString("key") != "value" {
		t.Error("expected GetString to return value")
	}
	if e.GetString("missing") != "" {
		t.Error("expected GetString to return empty for missing")
	}
}

func BenchmarkLogNoFields(b *testing.B) {
	log := New(&Options{
		Output:    &bytes.Buffer{},
		Formatter: &NoopFormatter{},
	})

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		log.Info("benchmark message")
	}
}

func BenchmarkLogWithFields(b *testing.B) {
	log := New(&Options{
		Output:    &bytes.Buffer{},
		Formatter: &NoopFormatter{},
	})

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		log.Info("benchmark message",
			String("key1", "value1"),
			Int("key2", 42),
			Bool("key3", true),
		)
	}
}

func BenchmarkLogTextFormatter(b *testing.B) {
	log := New(&Options{
		Output:    &bytes.Buffer{},
		Formatter: &TextFormatter{DisableTimestamp: true, DisableColors: true},
	})

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		log.Info("benchmark message", String("key", "value"))
	}
}

func BenchmarkLogJSONFormatter(b *testing.B) {
	log := New(&Options{
		Output:    &bytes.Buffer{},
		Formatter: &JSONFormatter{},
	})

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		log.Info("benchmark message", String("key", "value"))
	}
}

func BenchmarkFieldCreation(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_ = String("key", "value")
		_ = Int("count", 42)
		_ = Bool("active", true)
		_ = Duration("latency", time.Millisecond*100)
	}
}

// Tests for new features

func TestChainableBuilder(t *testing.T) {
	buf := &bytes.Buffer{}
	log := New(&Options{
		Output:    buf,
		Formatter: &TextFormatter{DisableTimestamp: true, DisableColors: true},
	})

	log.Build().
		Str("user", "john").
		Int("age", 30).
		Bool("admin", true).
		Info("user created")

	output := buf.String()
	if !strings.Contains(output, "user=john") {
		t.Errorf("expected user=john in output, got: %s", output)
	}
	if !strings.Contains(output, "age=30") {
		t.Errorf("expected age=30 in output, got: %s", output)
	}
	if !strings.Contains(output, "admin=true") {
		t.Errorf("expected admin=true in output, got: %s", output)
	}
}

func TestBuilderWithAuto(t *testing.T) {
	buf := &bytes.Buffer{}
	log := New(&Options{
		Output:    buf,
		Formatter: &TextFormatter{DisableTimestamp: true, DisableColors: true},
	})

	log.Build().
		With("user", "john").
		With("count", 42).
		With("active", true).
		Info("test")

	output := buf.String()
	if !strings.Contains(output, "user=john") {
		t.Errorf("expected user=john in output, got: %s", output)
	}
	if !strings.Contains(output, "count=42") {
		t.Errorf("expected count=42 in output, got: %s", output)
	}
}

func TestFHelper(t *testing.T) {
	buf := &bytes.Buffer{}
	log := New(&Options{
		Output:    buf,
		Formatter: &TextFormatter{DisableTimestamp: true, DisableColors: true},
	})

	log.F("name", "alice", "age", 25).Info("created")

	output := buf.String()
	if !strings.Contains(output, "name=alice") {
		t.Errorf("expected name=alice in output, got: %s", output)
	}
	if !strings.Contains(output, "age=25") {
		t.Errorf("expected age=25 in output, got: %s", output)
	}
}

func TestStructField(t *testing.T) {
	type User struct {
		ID    string `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email,omitempty"`
	}

	user := User{ID: "123", Name: "John"}
	field := Struct("user", user)

	if field.Key != "user" {
		t.Errorf("expected key 'user', got '%s'", field.Key)
	}
}

func TestStructFlat(t *testing.T) {
	type Request struct {
		Method string `json:"method"`
		Path   string `json:"path"`
		Status int    `json:"status"`
	}

	req := Request{Method: "GET", Path: "/api/users", Status: 200}
	fields := StructFlat(req)

	if len(fields) != 3 {
		t.Errorf("expected 3 fields, got %d", len(fields))
	}
}

func TestVField(t *testing.T) {
	tests := []struct {
		name  string
		value any
	}{
		{"string", "hello"},
		{"int", 42},
		{"bool", true},
		{"float", 3.14},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			field := V("key", tc.value)
			if field.Key != "key" {
				t.Errorf("expected key 'key', got '%s'", field.Key)
			}
		})
	}
}

func TestIfErr(t *testing.T) {
	buf := &bytes.Buffer{}
	log := New(&Options{
		Output:    buf,
		Formatter: &TextFormatter{DisableTimestamp: true, DisableColors: true},
	})

	// Should not log when error is nil
	log.IfErr(nil).With("key", "value").Error("should not log")
	if buf.Len() > 0 {
		t.Error("should not log when error is nil")
	}

	// Should log when error is not nil
	log.IfErr(errors.New("test error")).With("key", "value").Error("operation failed")
	output := buf.String()
	if !strings.Contains(output, "test error") {
		t.Errorf("expected error in output, got: %s", output)
	}
	if !strings.Contains(output, "key=value") {
		t.Errorf("expected field in output, got: %s", output)
	}
}

func TestWrapErr(t *testing.T) {
	buf := &bytes.Buffer{}
	log := New(&Options{
		Output:    buf,
		Formatter: &TextFormatter{DisableTimestamp: true, DisableColors: true},
	})

	originalErr := errors.New("connection refused")
	wrapped := log.WrapErr(originalErr, "failed to connect", String("host", "localhost"))

	if wrapped == nil {
		t.Error("wrapped error should not be nil")
	}
	if !strings.Contains(wrapped.Error(), "failed to connect") {
		t.Errorf("wrapped error should contain message, got: %s", wrapped.Error())
	}
	if !errors.Is(wrapped, originalErr) {
		t.Error("wrapped error should unwrap to original")
	}

	output := buf.String()
	if !strings.Contains(output, "connection refused") {
		t.Errorf("expected error in output, got: %s", output)
	}
}

func TestLogErr(t *testing.T) {
	buf := &bytes.Buffer{}
	log := New(&Options{
		Output:    buf,
		Formatter: &TextFormatter{DisableTimestamp: true, DisableColors: true},
	})

	// Should not log when nil
	log.LogErr(nil, "should not log")
	if buf.Len() > 0 {
		t.Error("should not log when error is nil")
	}

	// Should log when error exists
	log.LogErr(errors.New("test error"), "operation failed")
	if !strings.Contains(buf.String(), "test error") {
		t.Errorf("expected error in output, got: %s", buf.String())
	}
}

func TestCheckErr(t *testing.T) {
	buf := &bytes.Buffer{}
	log := New(&Options{
		Output:    buf,
		Formatter: &TextFormatter{DisableTimestamp: true, DisableColors: true},
	})

	if log.CheckErr(nil, "no error") {
		t.Error("CheckErr should return false for nil error")
	}

	if !log.CheckErr(errors.New("test"), "has error") {
		t.Error("CheckErr should return true for non-nil error")
	}
}

func TestPrintfMethods(t *testing.T) {
	buf := &bytes.Buffer{}
	log := New(&Options{
		Output:    buf,
		Formatter: &TextFormatter{DisableTimestamp: true, DisableColors: true},
	})

	log.Infof("user %s created with id %d", "john", 123)

	output := buf.String()
	if !strings.Contains(output, "user john created with id 123") {
		t.Errorf("expected formatted message, got: %s", output)
	}
}

func TestPrintfAllLevels(t *testing.T) {
	buf := &bytes.Buffer{}
	log := New(&Options{
		Output:    buf,
		Level:     TraceLevel,
		Formatter: &TextFormatter{DisableTimestamp: true, DisableColors: true},
	})

	log.Tracef("trace %d", 1)
	log.Debugf("debug %d", 2)
	log.Infof("info %d", 3)
	log.Warnf("warn %d", 4)
	log.Errorf("error %d", 5)

	output := buf.String()
	for i, level := range []string{"trace 1", "debug 2", "info 3", "warn 4", "error 5"} {
		if !strings.Contains(output, level) {
			t.Errorf("expected %s in output (test %d)", level, i)
		}
	}
}

func TestPrintCompat(t *testing.T) {
	buf := &bytes.Buffer{}
	log := New(&Options{
		Output:    buf,
		Formatter: &TextFormatter{DisableTimestamp: true, DisableColors: true},
	})

	log.Print("hello", " ", "world")
	log.Println("line 2")
	log.Printf("formatted %s", "message")

	output := buf.String()
	if !strings.Contains(output, "hello world") {
		t.Errorf("expected 'hello world' in output, got: %s", output)
	}
	if !strings.Contains(output, "formatted message") {
		t.Errorf("expected 'formatted message' in output, got: %s", output)
	}
}

func TestNamedLogger(t *testing.T) {
	buf := &bytes.Buffer{}
	log := New(&Options{
		Output: buf,
		Formatter: &NamedFormatter{
			Inner:     &TextFormatter{DisableTimestamp: true, DisableColors: true},
			Separator: " ",
			Brackets:  "[]",
		},
	})

	userLog := log.Named("users")
	userLog.Info("created")

	output := buf.String()
	if !strings.Contains(output, "[users]") {
		t.Errorf("expected [users] in output, got: %s", output)
	}
}

func TestNestedNamedLogger(t *testing.T) {
	buf := &bytes.Buffer{}
	log := New(&Options{
		Output: buf,
		Formatter: &NamedFormatter{
			Inner:     &TextFormatter{DisableTimestamp: true, DisableColors: true},
			Separator: " ",
			Brackets:  "[]",
		},
	})

	userLog := log.Named("users")
	authLog := userLog.Named("auth")
	authLog.Info("login")

	output := buf.String()
	if !strings.Contains(output, "[users.auth]") {
		t.Errorf("expected [users.auth] in output, got: %s", output)
	}
}

func TestLoggerGetName(t *testing.T) {
	log := New(nil)
	userLog := log.Named("users")
	authLog := userLog.Named("auth")

	if userLog.GetName() != "users" {
		t.Errorf("expected 'users', got '%s'", userLog.GetName())
	}
	if authLog.GetName() != "users.auth" {
		t.Errorf("expected 'users.auth', got '%s'", authLog.GetName())
	}
}

func TestLoggerNameParts(t *testing.T) {
	log := New(nil)
	authLog := log.Named("users").Named("auth").Named("oauth")

	parts := authLog.NameParts()
	if len(parts) != 3 {
		t.Errorf("expected 3 parts, got %d", len(parts))
	}
	if parts[0] != "users" || parts[1] != "auth" || parts[2] != "oauth" {
		t.Errorf("unexpected parts: %v", parts)
	}
}

func TestComponentAlias(t *testing.T) {
	log := New(nil)
	compLog := log.Component("database")

	if compLog.GetName() != "database" {
		t.Errorf("expected 'database', got '%s'", compLog.GetName())
	}
}

func TestErrChain(t *testing.T) {
	inner := errors.New("connection refused")
	outer := fmt.Errorf("failed to connect: %w", inner)

	field := ErrChain(outer)
	if field.Key != "errors" {
		t.Errorf("expected key 'errors', got '%s'", field.Key)
	}

	chain, ok := field.Interface.([]string)
	if !ok {
		t.Fatal("expected []string interface")
	}
	if len(chain) != 2 {
		t.Errorf("expected 2 errors in chain, got %d", len(chain))
	}
}

func TestBuilderWithContext(t *testing.T) {
	buf := &bytes.Buffer{}
	log := New(&Options{
		Output:    buf,
		Formatter: &TextFormatter{DisableTimestamp: true, DisableColors: true},
	})

	ctx := context.Background()
	ctx = WithContextFields(ctx, String("request_id", "req-123"))

	log.Build().
		WithContext(ctx).
		Str("action", "test").
		Info("message")

	output := buf.String()
	if !strings.Contains(output, "request_id=req-123") {
		t.Errorf("expected request_id in output, got: %s", output)
	}
}

func TestBuilderWithError(t *testing.T) {
	buf := &bytes.Buffer{}
	log := New(&Options{
		Output:    buf,
		Formatter: &TextFormatter{DisableTimestamp: true, DisableColors: true},
	})

	err := errors.New("test error")
	log.Build().
		WithError(err).
		Str("context", "test").
		Error("operation failed")

	output := buf.String()
	if !strings.Contains(output, "test error") {
		t.Errorf("expected error in output, got: %s", output)
	}
}

func BenchmarkBuilder(b *testing.B) {
	log := New(&Options{
		Output:    &bytes.Buffer{},
		Formatter: &NoopFormatter{},
	})

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		log.Build().
			Str("key1", "value1").
			Int("key2", 42).
			Bool("key3", true).
			Info("message")
	}
}

func BenchmarkFHelper(b *testing.B) {
	log := New(&Options{
		Output:    &bytes.Buffer{},
		Formatter: &NoopFormatter{},
	})

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		log.F("key1", "value1", "key2", 42, "key3", true).Info("message")
	}
}

func BenchmarkPrintf(b *testing.B) {
	log := New(&Options{
		Output:    &bytes.Buffer{},
		Formatter: &NoopFormatter{},
	})

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		log.Infof("user %s with id %d", "john", 123)
	}
}
