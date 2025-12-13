package logs

import (
	"strings"
)

// Level represents a log level.
type Level int

// Log levels from most to least severe.
const (
	// PanicLevel logs and then panics.
	PanicLevel Level = iota
	// FatalLevel logs and then calls os.Exit(1).
	FatalLevel
	// ErrorLevel is for errors that should be noted.
	ErrorLevel
	// WarnLevel is for non-critical issues.
	WarnLevel
	// InfoLevel is for general operational information.
	InfoLevel
	// DebugLevel is for debugging information.
	DebugLevel
	// TraceLevel is for fine-grained debugging.
	TraceLevel
)

// String returns the string representation of a level.
func (l Level) String() string {
	switch l {
	case PanicLevel:
		return "panic"
	case FatalLevel:
		return "fatal"
	case ErrorLevel:
		return "error"
	case WarnLevel:
		return "warn"
	case InfoLevel:
		return "info"
	case DebugLevel:
		return "debug"
	case TraceLevel:
		return "trace"
	default:
		return "unknown"
	}
}

// ShortString returns a 4-character representation for alignment.
func (l Level) ShortString() string {
	switch l {
	case PanicLevel:
		return "PANC"
	case FatalLevel:
		return "FATL"
	case ErrorLevel:
		return "ERRO"
	case WarnLevel:
		return "WARN"
	case InfoLevel:
		return "INFO"
	case DebugLevel:
		return "DEBG"
	case TraceLevel:
		return "TRAC"
	default:
		return "UNKN"
	}
}

// Color returns the ANSI color code for the level.
func (l Level) Color() string {
	switch l {
	case PanicLevel, FatalLevel:
		return "\033[35m" // Magenta
	case ErrorLevel:
		return "\033[31m" // Red
	case WarnLevel:
		return "\033[33m" // Yellow
	case InfoLevel:
		return "\033[32m" // Green
	case DebugLevel:
		return "\033[36m" // Cyan
	case TraceLevel:
		return "\033[37m" // White
	default:
		return "\033[0m" // Reset
	}
}

// ParseLevel parses a string into a Level.
func ParseLevel(s string) Level {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "panic":
		return PanicLevel
	case "fatal":
		return FatalLevel
	case "error", "err":
		return ErrorLevel
	case "warn", "warning":
		return WarnLevel
	case "info":
		return InfoLevel
	case "debug":
		return DebugLevel
	case "trace":
		return TraceLevel
	default:
		return InfoLevel
	}
}

// AllLevels returns all log levels.
func AllLevels() []Level {
	return []Level{
		PanicLevel,
		FatalLevel,
		ErrorLevel,
		WarnLevel,
		InfoLevel,
		DebugLevel,
		TraceLevel,
	}
}
