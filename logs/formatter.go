package logs

import (
	"bytes"
	"encoding/json"
	"strconv"
	"sync"
	"time"
)

// Formatter formats log entries.
type Formatter interface {
	Format(entry *Entry) ([]byte, error)
}

// bufferPool is a pool of byte buffers for formatting.
var bufferPool = sync.Pool{
	New: func() any {
		return bytes.NewBuffer(make([]byte, 0, 256))
	},
}

func getBuffer() *bytes.Buffer {
	buf := bufferPool.Get().(*bytes.Buffer)
	buf.Reset()
	return buf
}

func putBuffer(buf *bytes.Buffer) {
	if buf.Cap() > 64*1024 {
		// Don't pool buffers larger than 64KB
		return
	}
	bufferPool.Put(buf)
}

// TextFormatter formats logs as text.
type TextFormatter struct {
	// TimestampFormat is the format for timestamps.
	// Default: "2006-01-02T15:04:05.000Z07:00"
	TimestampFormat string

	// DisableTimestamp disables timestamp output.
	DisableTimestamp bool

	// DisableColors disables ANSI colors.
	DisableColors bool

	// DisableQuoting disables quoting of string values.
	DisableQuoting bool

	// QuoteEmptyFields quotes empty field values.
	QuoteEmptyFields bool

	// FullTimestamp shows full timestamp instead of relative time.
	FullTimestamp bool

	// PadLevelText pads level text for alignment.
	PadLevelText bool

	// FieldSeparator is the separator between fields.
	// Default: " "
	FieldSeparator string

	// KeyValueSeparator is the separator between key and value.
	// Default: "="
	KeyValueSeparator string
}

// Format formats an entry as text.
func (f *TextFormatter) Format(entry *Entry) ([]byte, error) {
	buf := getBuffer()
	defer putBuffer(buf)

	timestampFormat := f.TimestampFormat
	if timestampFormat == "" {
		timestampFormat = "2006-01-02T15:04:05.000Z07:00"
	}

	fieldSep := f.FieldSeparator
	if fieldSep == "" {
		fieldSep = " "
	}

	kvSep := f.KeyValueSeparator
	if kvSep == "" {
		kvSep = "="
	}

	// Timestamp
	if !f.DisableTimestamp {
		buf.WriteString(entry.Time.Format(timestampFormat))
		buf.WriteString(fieldSep)
	}

	// Level
	levelStr := entry.Level.ShortString()
	if !f.DisableColors {
		buf.WriteString(entry.Level.Color())
		buf.WriteString(levelStr)
		buf.WriteString("\033[0m")
	} else {
		buf.WriteString(levelStr)
	}
	buf.WriteString(fieldSep)

	// Caller
	if entry.Caller != "" {
		if !f.DisableColors {
			buf.WriteString("\033[90m") // Gray
		}
		buf.WriteString(entry.Caller)
		if !f.DisableColors {
			buf.WriteString("\033[0m")
		}
		buf.WriteString(fieldSep)
	}

	// Message
	buf.WriteString(entry.Message)

	// Fields
	for _, field := range entry.Fields {
		buf.WriteString(fieldSep)

		if !f.DisableColors {
			buf.WriteString("\033[36m") // Cyan
		}
		buf.WriteString(field.Key)
		if !f.DisableColors {
			buf.WriteString("\033[0m")
		}

		buf.WriteString(kvSep)

		value := field.StringValue()
		if f.needsQuoting(value) {
			buf.WriteString(strconv.Quote(value))
		} else {
			buf.WriteString(value)
		}
	}

	buf.WriteByte('\n')

	// Stack trace
	if entry.Stack != "" {
		buf.WriteString(entry.Stack)
		buf.WriteByte('\n')
	}

	result := make([]byte, buf.Len())
	copy(result, buf.Bytes())
	return result, nil
}

// needsQuoting returns true if the value needs quoting.
func (f *TextFormatter) needsQuoting(value string) bool {
	if f.DisableQuoting {
		return false
	}
	if f.QuoteEmptyFields && value == "" {
		return true
	}
	for _, c := range value {
		if c == ' ' || c == '"' || c == '=' || c == '\n' || c == '\r' || c == '\t' {
			return true
		}
	}
	return false
}

// JSONFormatter formats logs as JSON.
type JSONFormatter struct {
	// TimestampFormat is the format for timestamps.
	// Default: time.RFC3339Nano
	TimestampFormat string

	// DisableTimestamp disables timestamp output.
	DisableTimestamp bool

	// TimestampKey is the key for the timestamp field.
	// Default: "time"
	TimestampKey string

	// LevelKey is the key for the level field.
	// Default: "level"
	LevelKey string

	// MessageKey is the key for the message field.
	// Default: "msg"
	MessageKey string

	// CallerKey is the key for the caller field.
	// Default: "caller"
	CallerKey string

	// StackKey is the key for the stack trace field.
	// Default: "stack"
	StackKey string

	// PrettyPrint enables pretty-printed JSON.
	PrettyPrint bool

	// EscapeHTML escapes HTML in JSON strings.
	EscapeHTML bool
}

// jsonEntry is the structure for JSON encoding.
type jsonEntry struct {
	Time    string         `json:"time,omitempty"`
	Level   string         `json:"level"`
	Msg     string         `json:"msg"`
	Caller  string         `json:"caller,omitempty"`
	Stack   string         `json:"stack,omitempty"`
	Fields  map[string]any `json:"-"`
}

// Format formats an entry as JSON.
func (f *JSONFormatter) Format(entry *Entry) ([]byte, error) {
	buf := getBuffer()
	defer putBuffer(buf)

	timestampFormat := f.TimestampFormat
	if timestampFormat == "" {
		timestampFormat = time.RFC3339Nano
	}

	timestampKey := f.TimestampKey
	if timestampKey == "" {
		timestampKey = "time"
	}

	levelKey := f.LevelKey
	if levelKey == "" {
		levelKey = "level"
	}

	messageKey := f.MessageKey
	if messageKey == "" {
		messageKey = "msg"
	}

	callerKey := f.CallerKey
	if callerKey == "" {
		callerKey = "caller"
	}

	stackKey := f.StackKey
	if stackKey == "" {
		stackKey = "stack"
	}

	// Build JSON object
	buf.WriteByte('{')

	// Timestamp
	if !f.DisableTimestamp {
		buf.WriteByte('"')
		buf.WriteString(timestampKey)
		buf.WriteString(`":"`)
		buf.WriteString(entry.Time.Format(timestampFormat))
		buf.WriteString(`",`)
	}

	// Level
	buf.WriteByte('"')
	buf.WriteString(levelKey)
	buf.WriteString(`":"`)
	buf.WriteString(entry.Level.String())
	buf.WriteString(`",`)

	// Message
	buf.WriteByte('"')
	buf.WriteString(messageKey)
	buf.WriteString(`":`)
	f.writeJSONString(buf, entry.Message)

	// Caller
	if entry.Caller != "" {
		buf.WriteString(`,"`)
		buf.WriteString(callerKey)
		buf.WriteString(`":"`)
		buf.WriteString(entry.Caller)
		buf.WriteByte('"')
	}

	// Stack
	if entry.Stack != "" {
		buf.WriteString(`,"`)
		buf.WriteString(stackKey)
		buf.WriteString(`":`)
		f.writeJSONString(buf, entry.Stack)
	}

	// Fields
	for _, field := range entry.Fields {
		buf.WriteString(`,"`)
		buf.WriteString(field.Key)
		buf.WriteString(`":`)
		f.writeJSONValue(buf, field)
	}

	buf.WriteByte('}')
	buf.WriteByte('\n')

	result := make([]byte, buf.Len())
	copy(result, buf.Bytes())
	return result, nil
}

// writeJSONString writes a JSON-encoded string.
func (f *JSONFormatter) writeJSONString(buf *bytes.Buffer, s string) {
	data, _ := json.Marshal(s)
	buf.Write(data)
}

// writeJSONValue writes a JSON-encoded field value.
func (f *JSONFormatter) writeJSONValue(buf *bytes.Buffer, field Field) {
	switch field.Type {
	case FieldTypeString:
		f.writeJSONString(buf, field.String)
	case FieldTypeInt:
		buf.WriteString(strconv.FormatInt(field.Int, 10))
	case FieldTypeUint:
		buf.WriteString(strconv.FormatUint(field.Uint, 10))
	case FieldTypeFloat:
		buf.WriteString(strconv.FormatFloat(field.Float, 'g', -1, 64))
	case FieldTypeBool:
		if field.Int == 1 {
			buf.WriteString("true")
		} else {
			buf.WriteString("false")
		}
	case FieldTypeTime:
		if t, ok := field.Interface.(time.Time); ok {
			buf.WriteByte('"')
			buf.WriteString(t.Format(time.RFC3339Nano))
			buf.WriteByte('"')
		} else {
			buf.WriteString(strconv.FormatInt(field.Int, 10))
		}
	case FieldTypeDuration:
		buf.WriteByte('"')
		buf.WriteString(time.Duration(field.Int).String())
		buf.WriteByte('"')
	case FieldTypeError:
		f.writeJSONString(buf, field.String)
	case FieldTypeStringer:
		if s, ok := field.Interface.(interface{ String() string }); ok {
			f.writeJSONString(buf, s.String())
		} else {
			buf.WriteString("null")
		}
	case FieldTypeBytes:
		if b, ok := field.Interface.([]byte); ok {
			// Check if it's already valid JSON
			if json.Valid(b) {
				buf.Write(b)
			} else {
				f.writeJSONString(buf, string(b))
			}
		} else {
			buf.WriteString("null")
		}
	default:
		if field.Interface == nil {
			buf.WriteString("null")
		} else {
			data, err := json.Marshal(field.Interface)
			if err != nil {
				f.writeJSONString(buf, field.StringValue())
			} else {
				buf.Write(data)
			}
		}
	}
}

// PrettyFormatter formats logs with colors and alignment for development.
type PrettyFormatter struct {
	// TimestampFormat is the format for timestamps.
	// Default: "15:04:05.000"
	TimestampFormat string

	// ShowCaller shows caller information.
	ShowCaller bool

	// ShowTimestamp shows timestamps.
	ShowTimestamp bool
}

// Format formats an entry in a pretty, colorful format.
func (f *PrettyFormatter) Format(entry *Entry) ([]byte, error) {
	buf := getBuffer()
	defer putBuffer(buf)

	timestampFormat := f.TimestampFormat
	if timestampFormat == "" {
		timestampFormat = "15:04:05.000"
	}

	// Timestamp
	if f.ShowTimestamp {
		buf.WriteString("\033[90m") // Gray
		buf.WriteString(entry.Time.Format(timestampFormat))
		buf.WriteString("\033[0m ")
	}

	// Level with color and emoji
	buf.WriteString(f.levelEmoji(entry.Level))
	buf.WriteByte(' ')
	buf.WriteString(entry.Level.Color())
	buf.WriteString(entry.Level.ShortString())
	buf.WriteString("\033[0m ")

	// Message
	buf.WriteString("\033[1m") // Bold
	buf.WriteString(entry.Message)
	buf.WriteString("\033[0m")

	// Fields
	if len(entry.Fields) > 0 {
		buf.WriteString(" \033[90m│\033[0m")
		for _, field := range entry.Fields {
			buf.WriteByte(' ')
			buf.WriteString("\033[36m") // Cyan
			buf.WriteString(field.Key)
			buf.WriteString("\033[0m")
			buf.WriteByte('=')
			buf.WriteString(field.StringValue())
		}
	}

	// Caller
	if f.ShowCaller && entry.Caller != "" {
		buf.WriteString(" \033[90m(")
		buf.WriteString(entry.Caller)
		buf.WriteString(")\033[0m")
	}

	buf.WriteByte('\n')

	// Stack trace
	if entry.Stack != "" {
		buf.WriteString("\033[90m")
		buf.WriteString(entry.Stack)
		buf.WriteString("\033[0m\n")
	}

	result := make([]byte, buf.Len())
	copy(result, buf.Bytes())
	return result, nil
}

// levelEmoji returns an emoji for the log level.
func (f *PrettyFormatter) levelEmoji(level Level) string {
	switch level {
	case PanicLevel:
		return "💥"
	case FatalLevel:
		return "☠️"
	case ErrorLevel:
		return "❌"
	case WarnLevel:
		return "⚠️"
	case InfoLevel:
		return "ℹ️"
	case DebugLevel:
		return "🔍"
	case TraceLevel:
		return "📍"
	default:
		return "  "
	}
}

// NoopFormatter discards all output.
type NoopFormatter struct{}

// Format returns nil.
func (f *NoopFormatter) Format(entry *Entry) ([]byte, error) {
	return nil, nil
}

