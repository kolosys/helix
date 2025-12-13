package logs

import (
	"context"
	"fmt"
)

// Printf-style logging methods for easy migration from fmt/log packages.
// These methods format the message using fmt.Sprintf and log it.
// For structured logging, prefer the standard methods with fields.

// Tracef logs a formatted message at trace level.
func (l *Logger) Tracef(format string, args ...any) {
	if l.IsEnabled(TraceLevel) {
		l.log(TraceLevel, fmt.Sprintf(format, args...), nil)
	}
}

// Debugf logs a formatted message at debug level.
func (l *Logger) Debugf(format string, args ...any) {
	if l.IsEnabled(DebugLevel) {
		l.log(DebugLevel, fmt.Sprintf(format, args...), nil)
	}
}

// Infof logs a formatted message at info level.
func (l *Logger) Infof(format string, args ...any) {
	if l.IsEnabled(InfoLevel) {
		l.log(InfoLevel, fmt.Sprintf(format, args...), nil)
	}
}

// Warnf logs a formatted message at warn level.
func (l *Logger) Warnf(format string, args ...any) {
	if l.IsEnabled(WarnLevel) {
		l.log(WarnLevel, fmt.Sprintf(format, args...), nil)
	}
}

// Errorf logs a formatted message at error level.
func (l *Logger) Errorf(format string, args ...any) {
	if l.IsEnabled(ErrorLevel) {
		l.log(ErrorLevel, fmt.Sprintf(format, args...), nil)
	}
}

// Fatalf logs a formatted message at fatal level and exits.
func (l *Logger) Fatalf(format string, args ...any) {
	l.log(FatalLevel, fmt.Sprintf(format, args...), nil)
}

// Panicf logs a formatted message at panic level and panics.
func (l *Logger) Panicf(format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	l.log(PanicLevel, msg, nil)
	panic(msg)
}

// Printf logs a formatted message at info level (stdlib log compatibility).
func (l *Logger) Printf(format string, args ...any) {
	if l.IsEnabled(InfoLevel) {
		l.log(InfoLevel, fmt.Sprintf(format, args...), nil)
	}
}

// Print logs a message at info level (stdlib log compatibility).
func (l *Logger) Print(args ...any) {
	if l.IsEnabled(InfoLevel) {
		l.log(InfoLevel, fmt.Sprint(args...), nil)
	}
}

// Println logs a message at info level (stdlib log compatibility).
func (l *Logger) Println(args ...any) {
	if l.IsEnabled(InfoLevel) {
		l.log(InfoLevel, fmt.Sprint(args...), nil)
	}
}

// Context-aware printf methods

// TracefContext logs a formatted message at trace level with context.
func (l *Logger) TracefContext(ctx context.Context, format string, args ...any) {
	if l.IsEnabled(TraceLevel) {
		l.logContext(ctx, TraceLevel, fmt.Sprintf(format, args...), nil)
	}
}

// DebugfContext logs a formatted message at debug level with context.
func (l *Logger) DebugfContext(ctx context.Context, format string, args ...any) {
	if l.IsEnabled(DebugLevel) {
		l.logContext(ctx, DebugLevel, fmt.Sprintf(format, args...), nil)
	}
}

// InfofContext logs a formatted message at info level with context.
func (l *Logger) InfofContext(ctx context.Context, format string, args ...any) {
	if l.IsEnabled(InfoLevel) {
		l.logContext(ctx, InfoLevel, fmt.Sprintf(format, args...), nil)
	}
}

// WarnfContext logs a formatted message at warn level with context.
func (l *Logger) WarnfContext(ctx context.Context, format string, args ...any) {
	if l.IsEnabled(WarnLevel) {
		l.logContext(ctx, WarnLevel, fmt.Sprintf(format, args...), nil)
	}
}

// ErrorfContext logs a formatted message at error level with context.
func (l *Logger) ErrorfContext(ctx context.Context, format string, args ...any) {
	if l.IsEnabled(ErrorLevel) {
		l.logContext(ctx, ErrorLevel, fmt.Sprintf(format, args...), nil)
	}
}

// Package-level printf functions using the default logger

// Tracef logs a formatted message at trace level.
func Tracef(format string, args ...any) { defaultLogger.Tracef(format, args...) }

// Debugf logs a formatted message at debug level.
func Debugf(format string, args ...any) { defaultLogger.Debugf(format, args...) }

// Infof logs a formatted message at info level.
func Infof(format string, args ...any) { defaultLogger.Infof(format, args...) }

// Warnf logs a formatted message at warn level.
func Warnf(format string, args ...any) { defaultLogger.Warnf(format, args...) }

// Errorf logs a formatted message at error level.
func Errorf(format string, args ...any) { defaultLogger.Errorf(format, args...) }

// Fatalf logs a formatted message at fatal level and exits.
func Fatalf(format string, args ...any) { defaultLogger.Fatalf(format, args...) }

// Panicf logs a formatted message at panic level and panics.
func Panicf(format string, args ...any) { defaultLogger.Panicf(format, args...) }

// Printf logs a formatted message at info level.
func Printf(format string, args ...any) { defaultLogger.Printf(format, args...) }

// Print logs a message at info level.
func Print(args ...any) { defaultLogger.Print(args...) }

// Println logs a message at info level.
func Println(args ...any) { defaultLogger.Println(args...) }

