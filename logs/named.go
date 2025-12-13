package logs

import (
	"strings"
)

// Named creates a child logger with a name prefix.
// Names are joined with dots to create a hierarchy.
//
//	userLog := log.Named("users")           // [users]
//	authLog := userLog.Named("auth")        // [users.auth]
//	log.Info("action")                      // [users.auth] action
func (l *Logger) Named(name string) *Logger {
	child := l.clone()

	// Build the full name
	if existingName := l.getName(); existingName != "" {
		child.setName(existingName + "." + name)
	} else {
		child.setName(name)
	}

	return child
}

// getName returns the logger's name from its default fields.
func (l *Logger) getName() string {
	for _, f := range l.fields {
		if f.Key == loggerNameKey {
			return f.String
		}
	}
	return ""
}

// setName sets the logger's name in its default fields.
func (l *Logger) setName(name string) {
	// Remove existing name field
	newFields := make([]Field, 0, len(l.fields)+1)
	for _, f := range l.fields {
		if f.Key != loggerNameKey {
			newFields = append(newFields, f)
		}
	}

	// Add new name field at the beginning
	l.fields = append([]Field{{Key: loggerNameKey, Type: FieldTypeString, String: name}}, newFields...)
}

// loggerNameKey is the field key for logger names.
const loggerNameKey = "_logger"

// clone creates a shallow copy of the logger.
func (l *Logger) clone() *Logger {
	child := &Logger{
		output:      l.output,
		formatter:   l.formatter,
		hooks:       l.hooks,
		callerDepth: l.callerDepth,
		addCaller:   l.addCaller,
		addStack:    l.addStack,
		async:       l.async,
		asyncCh:     l.asyncCh,
		entryPool:   l.entryPool,
		sampler:     l.sampler,
		fields:      make([]Field, len(l.fields)),
	}
	child.level.Store(l.level.Load())
	copy(child.fields, l.fields)
	return child
}

// Component creates a named logger for a specific component.
// This is an alias for Named with a more semantic name.
func (l *Logger) Component(name string) *Logger {
	return l.Named(name)
}

// Module creates a named logger for a module.
// This is an alias for Named with a more semantic name.
func (l *Logger) Module(name string) *Logger {
	return l.Named(name)
}

// Service creates a named logger for a service.
// This is an alias for Named with a more semantic name.
func (l *Logger) Service(name string) *Logger {
	return l.Named(name)
}

// Package-level named logger creation

// Named creates a named child of the default logger.
func Named(name string) *Logger {
	return defaultLogger.Named(name)
}

// Component creates a component logger from the default logger.
func Component(name string) *Logger {
	return defaultLogger.Component(name)
}

// NamedFormatter is a formatter that includes the logger name in output.
type NamedFormatter struct {
	// Inner is the formatter to wrap.
	Inner Formatter
	// Separator is placed between the name and message.
	// Default: " "
	Separator string
	// Brackets wraps the name. Default: "[]"
	Brackets string
}

// Format formats an entry, prefixing with the logger name if present.
func (f *NamedFormatter) Format(entry *Entry) ([]byte, error) {
	// Find the logger name field
	var name string
	var filteredFields []Field
	for _, field := range entry.Fields {
		if field.Key == loggerNameKey {
			name = field.String
		} else {
			filteredFields = append(filteredFields, field)
		}
	}

	// If there's a name, modify the message
	if name != "" {
		brackets := f.Brackets
		if brackets == "" {
			brackets = "[]"
		}
		sep := f.Separator
		if sep == "" {
			sep = " "
		}

		// Modify entry for formatting
		modifiedEntry := *entry
		modifiedEntry.Message = string(brackets[0]) + name + string(brackets[1]) + sep + entry.Message
		modifiedEntry.Fields = filteredFields

		return f.Inner.Format(&modifiedEntry)
	}

	return f.Inner.Format(entry)
}

// WithNamedFormat wraps the current formatter to include logger names.
func WithNamedFormat() Option {
	return func(l *Logger) {
		l.formatter = &NamedFormatter{
			Inner:     l.formatter,
			Separator: " ",
			Brackets:  "[]",
		}
	}
}

// GetName returns the logger's name, or empty string if not named.
func (l *Logger) GetName() string {
	return l.getName()
}

// FullName returns the complete hierarchical name of the logger.
func (l *Logger) FullName() string {
	return l.getName()
}

// NameParts returns the name split into parts.
func (l *Logger) NameParts() []string {
	name := l.getName()
	if name == "" {
		return nil
	}
	return strings.Split(name, ".")
}

