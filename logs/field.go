package logs

import (
	"encoding/json"
	"fmt"
	"time"
)

// FieldType represents the type of a field value.
type FieldType uint8

const (
	// FieldTypeUnknown is an unknown field type.
	FieldTypeUnknown FieldType = iota
	// FieldTypeString is a string field.
	FieldTypeString
	// FieldTypeInt is an int64 field.
	FieldTypeInt
	// FieldTypeUint is a uint64 field.
	FieldTypeUint
	// FieldTypeFloat is a float64 field.
	FieldTypeFloat
	// FieldTypeBool is a bool field.
	FieldTypeBool
	// FieldTypeTime is a time.Time field.
	FieldTypeTime
	// FieldTypeDuration is a time.Duration field.
	FieldTypeDuration
	// FieldTypeError is an error field.
	FieldTypeError
	// FieldTypeAny is any other type.
	FieldTypeAny
	// FieldTypeStringer is a fmt.Stringer field.
	FieldTypeStringer
	// FieldTypeBytes is a []byte field.
	FieldTypeBytes
)

// Field represents a structured log field.
type Field struct {
	Key       string
	Type      FieldType
	Int       int64
	Uint      uint64
	Float     float64
	String    string
	Interface any
}

// String creates a string field.
func String(key, value string) Field {
	return Field{Key: key, Type: FieldTypeString, String: value}
}

// Strings creates a string slice field.
func Strings(key string, values []string) Field {
	return Field{Key: key, Type: FieldTypeAny, Interface: values}
}

// Stringer creates a field from a fmt.Stringer.
func Stringer(key string, value fmt.Stringer) Field {
	return Field{Key: key, Type: FieldTypeStringer, Interface: value}
}

// Int creates an int field.
func Int(key string, value int) Field {
	return Field{Key: key, Type: FieldTypeInt, Int: int64(value)}
}

// Int8 creates an int8 field.
func Int8(key string, value int8) Field {
	return Field{Key: key, Type: FieldTypeInt, Int: int64(value)}
}

// Int16 creates an int16 field.
func Int16(key string, value int16) Field {
	return Field{Key: key, Type: FieldTypeInt, Int: int64(value)}
}

// Int32 creates an int32 field.
func Int32(key string, value int32) Field {
	return Field{Key: key, Type: FieldTypeInt, Int: int64(value)}
}

// Int64 creates an int64 field.
func Int64(key string, value int64) Field {
	return Field{Key: key, Type: FieldTypeInt, Int: value}
}

// Uint creates a uint field.
func Uint(key string, value uint) Field {
	return Field{Key: key, Type: FieldTypeUint, Uint: uint64(value)}
}

// Uint8 creates a uint8 field.
func Uint8(key string, value uint8) Field {
	return Field{Key: key, Type: FieldTypeUint, Uint: uint64(value)}
}

// Uint16 creates a uint16 field.
func Uint16(key string, value uint16) Field {
	return Field{Key: key, Type: FieldTypeUint, Uint: uint64(value)}
}

// Uint32 creates a uint32 field.
func Uint32(key string, value uint32) Field {
	return Field{Key: key, Type: FieldTypeUint, Uint: uint64(value)}
}

// Uint64 creates a uint64 field.
func Uint64(key string, value uint64) Field {
	return Field{Key: key, Type: FieldTypeUint, Uint: value}
}

// Float32 creates a float32 field.
func Float32(key string, value float32) Field {
	return Field{Key: key, Type: FieldTypeFloat, Float: float64(value)}
}

// Float64 creates a float64 field.
func Float64(key string, value float64) Field {
	return Field{Key: key, Type: FieldTypeFloat, Float: value}
}

// Bool creates a bool field.
func Bool(key string, value bool) Field {
	f := Field{Key: key, Type: FieldTypeBool}
	if value {
		f.Int = 1
	}
	return f
}

// Time creates a time.Time field.
func Time(key string, value time.Time) Field {
	return Field{Key: key, Type: FieldTypeTime, Int: value.UnixNano(), Interface: value}
}

// Duration creates a time.Duration field.
func Duration(key string, value time.Duration) Field {
	return Field{Key: key, Type: FieldTypeDuration, Int: int64(value)}
}

// Err creates an error field with key "error".
func Err(err error) Field {
	if err == nil {
		return Field{Key: "error", Type: FieldTypeString, String: ""}
	}
	return Field{Key: "error", Type: FieldTypeError, Interface: err, String: err.Error()}
}

// NamedErr creates an error field with a custom key.
func NamedErr(key string, err error) Field {
	if err == nil {
		return Field{Key: key, Type: FieldTypeString, String: ""}
	}
	return Field{Key: key, Type: FieldTypeError, Interface: err, String: err.Error()}
}

// Any creates a field with any value.
func Any(key string, value any) Field {
	switch v := value.(type) {
	case nil:
		return Field{Key: key, Type: FieldTypeString, String: "null"}
	case string:
		return String(key, v)
	case int:
		return Int(key, v)
	case int64:
		return Int64(key, v)
	case int32:
		return Int32(key, v)
	case int16:
		return Int16(key, v)
	case int8:
		return Int8(key, v)
	case uint:
		return Uint(key, v)
	case uint64:
		return Uint64(key, v)
	case uint32:
		return Uint32(key, v)
	case uint16:
		return Uint16(key, v)
	case uint8:
		return Uint8(key, v)
	case float64:
		return Float64(key, v)
	case float32:
		return Float32(key, v)
	case bool:
		return Bool(key, v)
	case time.Time:
		return Time(key, v)
	case time.Duration:
		return Duration(key, v)
	case error:
		return NamedErr(key, v)
	case fmt.Stringer:
		return Stringer(key, v)
	case []byte:
		return Bytes(key, v)
	default:
		return Field{Key: key, Type: FieldTypeAny, Interface: v}
	}
}

// Bytes creates a []byte field.
func Bytes(key string, value []byte) Field {
	return Field{Key: key, Type: FieldTypeBytes, Interface: value}
}

// JSON creates a field that will be JSON-encoded.
func JSON(key string, value any) Field {
	data, err := json.Marshal(value)
	if err != nil {
		return String(key, fmt.Sprintf("json error: %v", err))
	}
	return Field{Key: key, Type: FieldTypeBytes, Interface: data}
}

// Stack creates a field containing a stack trace.
func Stack(key string) Field {
	return String(key, getStack())
}

// Namespace creates a namespace field for grouping.
func Namespace(key string) Field {
	return Field{Key: key, Type: FieldTypeUnknown, String: "namespace"}
}

// Value returns the field value as an interface{}.
func (f Field) Value() any {
	switch f.Type {
	case FieldTypeString:
		return f.String
	case FieldTypeInt:
		return f.Int
	case FieldTypeUint:
		return f.Uint
	case FieldTypeFloat:
		return f.Float
	case FieldTypeBool:
		return f.Int == 1
	case FieldTypeTime:
		if t, ok := f.Interface.(time.Time); ok {
			return t
		}
		return time.Unix(0, f.Int)
	case FieldTypeDuration:
		return time.Duration(f.Int)
	case FieldTypeError:
		return f.Interface
	case FieldTypeStringer:
		if s, ok := f.Interface.(fmt.Stringer); ok {
			return s.String()
		}
		return f.Interface
	case FieldTypeBytes:
		return f.Interface
	default:
		return f.Interface
	}
}

// StringValue returns the field value as a string.
func (f Field) StringValue() string {
	switch f.Type {
	case FieldTypeString:
		return f.String
	case FieldTypeInt:
		return formatInt(f.Int)
	case FieldTypeUint:
		return formatUint(f.Uint)
	case FieldTypeFloat:
		return formatFloat(f.Float)
	case FieldTypeBool:
		if f.Int == 1 {
			return "true"
		}
		return "false"
	case FieldTypeTime:
		if t, ok := f.Interface.(time.Time); ok {
			return t.Format(time.RFC3339)
		}
		return time.Unix(0, f.Int).Format(time.RFC3339)
	case FieldTypeDuration:
		return time.Duration(f.Int).String()
	case FieldTypeError:
		return f.String
	case FieldTypeStringer:
		if s, ok := f.Interface.(fmt.Stringer); ok {
			return s.String()
		}
		return fmt.Sprintf("%v", f.Interface)
	case FieldTypeBytes:
		if b, ok := f.Interface.([]byte); ok {
			return string(b)
		}
		return fmt.Sprintf("%v", f.Interface)
	default:
		if f.Interface == nil {
			return "null"
		}
		return fmt.Sprintf("%v", f.Interface)
	}
}

// formatInt formats an int64 without allocation for common cases.
func formatInt(n int64) string {
	if n >= 0 && n < 100 {
		return smallInts[n]
	}
	return fmt.Sprintf("%d", n)
}

// formatUint formats a uint64 without allocation for common cases.
func formatUint(n uint64) string {
	if n < 100 {
		return smallInts[n]
	}
	return fmt.Sprintf("%d", n)
}

// formatFloat formats a float64.
func formatFloat(f float64) string {
	return fmt.Sprintf("%g", f)
}

// Small integer strings for zero-allocation formatting.
var smallInts = []string{
	"0", "1", "2", "3", "4", "5", "6", "7", "8", "9",
	"10", "11", "12", "13", "14", "15", "16", "17", "18", "19",
	"20", "21", "22", "23", "24", "25", "26", "27", "28", "29",
	"30", "31", "32", "33", "34", "35", "36", "37", "38", "39",
	"40", "41", "42", "43", "44", "45", "46", "47", "48", "49",
	"50", "51", "52", "53", "54", "55", "56", "57", "58", "59",
	"60", "61", "62", "63", "64", "65", "66", "67", "68", "69",
	"70", "71", "72", "73", "74", "75", "76", "77", "78", "79",
	"80", "81", "82", "83", "84", "85", "86", "87", "88", "89",
	"90", "91", "92", "93", "94", "95", "96", "97", "98", "99",
}
