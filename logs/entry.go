package logs

import (
	"time"
)

// Entry represents a log entry.
type Entry struct {
	Level   Level
	Time    time.Time
	Message string
	Fields  []Field
	Caller  string
	Stack   string
}

// HasField returns true if the entry has a field with the given key.
func (e *Entry) HasField(key string) bool {
	for _, f := range e.Fields {
		if f.Key == key {
			return true
		}
	}
	return false
}

// GetField returns the field with the given key, or an empty field if not found.
func (e *Entry) GetField(key string) (Field, bool) {
	for _, f := range e.Fields {
		if f.Key == key {
			return f, true
		}
	}
	return Field{}, false
}

// GetString returns the string value of a field, or empty string if not found.
func (e *Entry) GetString(key string) string {
	if f, ok := e.GetField(key); ok {
		return f.StringValue()
	}
	return ""
}
