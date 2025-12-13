package logs

import (
	"context"
)

// contextKey is the type for context keys.
type contextKey int

const (
	fieldsKey contextKey = iota
	loggerKey
)

// WithFields adds fields to the context that will be included in all logs.
func WithContextFields(ctx context.Context, fields ...Field) context.Context {
	existing := FieldsFromContext(ctx)
	allFields := make([]Field, 0, len(existing)+len(fields))
	allFields = append(allFields, existing...)
	allFields = append(allFields, fields...)
	return context.WithValue(ctx, fieldsKey, allFields)
}

// FieldsFromContext extracts fields from the context.
func FieldsFromContext(ctx context.Context) []Field {
	if ctx == nil {
		return nil
	}
	if fields, ok := ctx.Value(fieldsKey).([]Field); ok {
		return fields
	}
	return nil
}

// WithLogger attaches a logger to the context.
func WithLogger(ctx context.Context, logger *Logger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}

// LoggerFromContext extracts a logger from the context.
// Returns the default logger if none is set.
func LoggerFromContext(ctx context.Context) *Logger {
	if ctx == nil {
		return defaultLogger
	}
	if logger, ok := ctx.Value(loggerKey).(*Logger); ok {
		return logger
	}
	return defaultLogger
}

// CtxTrace logs at trace level using the logger from context.
func CtxTrace(ctx context.Context, msg string, fields ...Field) {
	LoggerFromContext(ctx).TraceContext(ctx, msg, fields...)
}

// CtxDebug logs at debug level using the logger from context.
func CtxDebug(ctx context.Context, msg string, fields ...Field) {
	LoggerFromContext(ctx).DebugContext(ctx, msg, fields...)
}

// CtxInfo logs at info level using the logger from context.
func CtxInfo(ctx context.Context, msg string, fields ...Field) {
	LoggerFromContext(ctx).InfoContext(ctx, msg, fields...)
}

// CtxWarn logs at warn level using the logger from context.
func CtxWarn(ctx context.Context, msg string, fields ...Field) {
	LoggerFromContext(ctx).WarnContext(ctx, msg, fields...)
}

// CtxError logs at error level using the logger from context.
func CtxError(ctx context.Context, msg string, fields ...Field) {
	LoggerFromContext(ctx).ErrorContext(ctx, msg, fields...)
}

// RequestID is a common field key for request IDs.
const RequestIDKey = "request_id"

// WithRequestID adds a request ID to the context.
func WithRequestID(ctx context.Context, requestID string) context.Context {
	return WithContextFields(ctx, String(RequestIDKey, requestID))
}

// TraceID is a common field key for trace IDs.
const TraceIDKey = "trace_id"

// WithTraceID adds a trace ID to the context.
func WithTraceID(ctx context.Context, traceID string) context.Context {
	return WithContextFields(ctx, String(TraceIDKey, traceID))
}

// UserID is a common field key for user IDs.
const UserIDKey = "user_id"

// WithUserID adds a user ID to the context.
func WithUserID(ctx context.Context, userID string) context.Context {
	return WithContextFields(ctx, String(UserIDKey, userID))
}
