package logs

import (
	"errors"
	"fmt"
	"runtime"
)

// ErrorBuilder provides a fluent API for logging errors.
type ErrorBuilder struct {
	logger *Logger
	err    error
	fields []Field
}

// IfErr returns an ErrorBuilder if err is not nil, otherwise returns a no-op builder.
// This allows for one-liner error logging:
//
//	log.IfErr(err).With("user", id).Error("failed to create user")
func (l *Logger) IfErr(err error) *ErrorBuilder {
	return &ErrorBuilder{
		logger: l,
		err:    err,
		fields: make([]Field, 0, 4),
	}
}

// With adds a field to the error builder.
func (b *ErrorBuilder) With(key string, value any) *ErrorBuilder {
	if b.err == nil {
		return b
	}
	b.fields = append(b.fields, Any(key, value))
	return b
}

// WithField adds a typed field.
func (b *ErrorBuilder) WithField(f Field) *ErrorBuilder {
	if b.err == nil {
		return b
	}
	b.fields = append(b.fields, f)
	return b
}

// WithFields adds multiple fields.
func (b *ErrorBuilder) WithFields(fields ...Field) *ErrorBuilder {
	if b.err == nil {
		return b
	}
	b.fields = append(b.fields, fields...)
	return b
}

// Trace logs at trace level if error is not nil.
func (b *ErrorBuilder) Trace(msg string) {
	if b.err == nil {
		return
	}
	b.fields = append(b.fields, Err(b.err))
	b.logger.log(TraceLevel, msg, b.fields)
}

// Debug logs at debug level if error is not nil.
func (b *ErrorBuilder) Debug(msg string) {
	if b.err == nil {
		return
	}
	b.fields = append(b.fields, Err(b.err))
	b.logger.log(DebugLevel, msg, b.fields)
}

// Info logs at info level if error is not nil.
func (b *ErrorBuilder) Info(msg string) {
	if b.err == nil {
		return
	}
	b.fields = append(b.fields, Err(b.err))
	b.logger.log(InfoLevel, msg, b.fields)
}

// Warn logs at warn level if error is not nil.
func (b *ErrorBuilder) Warn(msg string) {
	if b.err == nil {
		return
	}
	b.fields = append(b.fields, Err(b.err))
	b.logger.log(WarnLevel, msg, b.fields)
}

// Error logs at error level if error is not nil.
func (b *ErrorBuilder) Error(msg string) {
	if b.err == nil {
		return
	}
	b.fields = append(b.fields, Err(b.err))
	b.logger.log(ErrorLevel, msg, b.fields)
}

// Fatal logs at fatal level if error is not nil and exits.
func (b *ErrorBuilder) Fatal(msg string) {
	if b.err == nil {
		return
	}
	b.fields = append(b.fields, Err(b.err))
	b.logger.log(FatalLevel, msg, b.fields)
}

// WrapErr wraps an error with additional context and logs it.
// Returns the wrapped error for returning from functions.
//
//	return log.WrapErr(err, "failed to connect", logs.String("host", host))
func (l *Logger) WrapErr(err error, msg string, fields ...Field) error {
	if err == nil {
		return nil
	}

	// Create wrapped error
	wrapped := fmt.Errorf("%s: %w", msg, err)

	// Log it
	allFields := make([]Field, 0, len(fields)+1)
	allFields = append(allFields, Err(err))
	allFields = append(allFields, fields...)
	l.log(ErrorLevel, msg, allFields)

	return wrapped
}

// WrapErrLevel wraps an error and logs at a specific level.
func (l *Logger) WrapErrLevel(level Level, err error, msg string, fields ...Field) error {
	if err == nil {
		return nil
	}

	wrapped := fmt.Errorf("%s: %w", msg, err)

	allFields := make([]Field, 0, len(fields)+1)
	allFields = append(allFields, Err(err))
	allFields = append(allFields, fields...)
	l.log(level, msg, allFields)

	return wrapped
}

// LogErr logs an error at error level if not nil.
// This is a simple one-liner for common error logging.
//
//	log.LogErr(err, "operation failed")
func (l *Logger) LogErr(err error, msg string, fields ...Field) {
	if err == nil {
		return
	}
	allFields := make([]Field, 0, len(fields)+1)
	allFields = append(allFields, Err(err))
	allFields = append(allFields, fields...)
	l.log(ErrorLevel, msg, allFields)
}

// ErrChain creates a field that unwraps the error chain.
func ErrChain(err error) Field {
	if err == nil {
		return String("errors", "null")
	}

	var chain []string
	for e := err; e != nil; e = errors.Unwrap(e) {
		chain = append(chain, e.Error())
	}

	return Field{
		Key:       "errors",
		Type:      FieldTypeAny,
		Interface: chain,
	}
}

// ErrWithStack creates an error field with stack trace.
func ErrWithStack(err error) Field {
	if err == nil {
		return String("error", "")
	}

	// Capture stack
	var pcs [32]uintptr
	n := runtime.Callers(2, pcs[:])
	frames := runtime.CallersFrames(pcs[:n])

	var stack []string
	for {
		frame, more := frames.Next()
		stack = append(stack, fmt.Sprintf("%s:%d %s", frame.File, frame.Line, frame.Function))
		if !more {
			break
		}
	}

	return Field{
		Key:  "error",
		Type: FieldTypeAny,
		Interface: map[string]any{
			"message": err.Error(),
			"stack":   stack,
		},
	}
}

// Must logs and panics if error is not nil.
// Useful for initialization code.
//
//	db := log.Must(sql.Open("postgres", dsn))
func Must[T any](l *Logger, val T, err error) T {
	if err != nil {
		l.log(PanicLevel, "fatal error", []Field{Err(err)})
		panic(err)
	}
	return val
}

// MustErr panics if error is not nil, logging the error first.
func (l *Logger) MustErr(err error, msg string, fields ...Field) {
	if err != nil {
		allFields := make([]Field, 0, len(fields)+1)
		allFields = append(allFields, Err(err))
		allFields = append(allFields, fields...)
		l.log(PanicLevel, msg, allFields)
		panic(err)
	}
}

// CheckErr logs an error and returns true if err is not nil.
// Useful in if statements:
//
//	if log.CheckErr(err, "failed") {
//	    return
//	}
func (l *Logger) CheckErr(err error, msg string, fields ...Field) bool {
	if err == nil {
		return false
	}
	allFields := make([]Field, 0, len(fields)+1)
	allFields = append(allFields, Err(err))
	allFields = append(allFields, fields...)
	l.log(ErrorLevel, msg, allFields)
	return true
}
