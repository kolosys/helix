package logs

import (
	"context"
)

// Builder provides a fluent/chainable API for building log entries.
// It accumulates fields and then emits a log entry when a level method is called.
type Builder struct {
	logger *Logger
	fields []Field
	ctx    context.Context
}

// newBuilder creates a new Builder.
func newBuilder(l *Logger) *Builder {
	return &Builder{
		logger: l,
		fields: make([]Field, 0, 8),
	}
}

// With adds a field to the builder using auto-detection.
func (b *Builder) With(key string, value any) *Builder {
	b.fields = append(b.fields, Any(key, value))
	return b
}

// WithField adds a typed field to the builder.
func (b *Builder) WithField(f Field) *Builder {
	b.fields = append(b.fields, f)
	return b
}

// WithFields adds multiple typed fields to the builder.
func (b *Builder) WithFields(fields ...Field) *Builder {
	b.fields = append(b.fields, fields...)
	return b
}

// WithContext sets the context for the log entry.
func (b *Builder) WithContext(ctx context.Context) *Builder {
	b.ctx = ctx
	return b
}

// WithError adds an error field.
func (b *Builder) WithError(err error) *Builder {
	if err != nil {
		b.fields = append(b.fields, Err(err))
	}
	return b
}

// Str adds a string field.
func (b *Builder) Str(key, value string) *Builder {
	b.fields = append(b.fields, String(key, value))
	return b
}

// Int adds an int field.
func (b *Builder) Int(key string, value int) *Builder {
	b.fields = append(b.fields, Int(key, value))
	return b
}

// Int64 adds an int64 field.
func (b *Builder) Int64(key string, value int64) *Builder {
	b.fields = append(b.fields, Int64(key, value))
	return b
}

// Uint adds a uint field.
func (b *Builder) Uint(key string, value uint) *Builder {
	b.fields = append(b.fields, Uint(key, value))
	return b
}

// Uint64 adds a uint64 field.
func (b *Builder) Uint64(key string, value uint64) *Builder {
	b.fields = append(b.fields, Uint64(key, value))
	return b
}

// Float64 adds a float64 field.
func (b *Builder) Float64(key string, value float64) *Builder {
	b.fields = append(b.fields, Float64(key, value))
	return b
}

// Bool adds a bool field.
func (b *Builder) Bool(key string, value bool) *Builder {
	b.fields = append(b.fields, Bool(key, value))
	return b
}

// Err adds an error field with key "error".
func (b *Builder) Err(err error) *Builder {
	if err != nil {
		b.fields = append(b.fields, Err(err))
	}
	return b
}

// Trace logs at trace level.
func (b *Builder) Trace(msg string) {
	b.emit(TraceLevel, msg)
}

// Debug logs at debug level.
func (b *Builder) Debug(msg string) {
	b.emit(DebugLevel, msg)
}

// Info logs at info level.
func (b *Builder) Info(msg string) {
	b.emit(InfoLevel, msg)
}

// Warn logs at warn level.
func (b *Builder) Warn(msg string) {
	b.emit(WarnLevel, msg)
}

// Error logs at error level.
func (b *Builder) Error(msg string) {
	b.emit(ErrorLevel, msg)
}

// Fatal logs at fatal level and exits.
func (b *Builder) Fatal(msg string) {
	b.emit(FatalLevel, msg)
}

// Panic logs at panic level and panics.
func (b *Builder) Panic(msg string) {
	b.emit(PanicLevel, msg)
}

// Log logs at the specified level.
func (b *Builder) Log(level Level, msg string) {
	b.emit(level, msg)
}

// emit sends the log entry.
func (b *Builder) emit(level Level, msg string) {
	if b.ctx != nil {
		b.logger.logContext(b.ctx, level, msg, b.fields)
	} else {
		b.logger.log(level, msg, b.fields)
	}
}

// Msg is an alias for Info (zerolog-style).
func (b *Builder) Msg(msg string) {
	b.emit(InfoLevel, msg)
}

// Send logs with an empty message (zerolog-style).
func (b *Builder) Send() {
	b.emit(InfoLevel, "")
}

// Logger methods to start a builder chain

// Build returns a new Builder for constructing a log entry.
func (l *Logger) Build() *Builder {
	return newBuilder(l)
}

// F returns a builder with fields added using key-value pairs.
// Keys must be strings, values are auto-detected.
// Example: log.F("user", "john", "age", 30).Info("created")
func (l *Logger) F(keyvals ...any) *Builder {
	b := newBuilder(l)
	for i := 0; i+1 < len(keyvals); i += 2 {
		if key, ok := keyvals[i].(string); ok {
			b.fields = append(b.fields, Any(key, keyvals[i+1]))
		}
	}
	return b
}

// Ctx returns a builder with context set.
func (l *Logger) Ctx(ctx context.Context) *Builder {
	b := newBuilder(l)
	b.ctx = ctx
	return b
}
