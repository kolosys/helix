package logs

import (
	"reflect"
	"strings"
	"sync"
	"time"
)

// structFieldCache caches struct field information for performance.
var structFieldCache sync.Map // map[reflect.Type][]structFieldInfo

type structFieldInfo struct {
	index    []int
	name     string
	omitempty bool
}

// Struct creates fields from a struct's exported fields.
// It uses json tags for field names if available, otherwise uses the field name.
// Nested structs are flattened with dot notation.
func Struct(key string, v any) Field {
	if v == nil {
		return String(key, "null")
	}

	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return String(key, "null")
		}
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return Any(key, v)
	}

	// For simple struct logging, marshal to a map
	fields := extractStructFields(key, val)
	if len(fields) == 0 {
		return Any(key, v)
	}

	// Return as a special struct field that the formatter can handle
	return Field{
		Key:       key,
		Type:      FieldTypeAny,
		Interface: structFields(fields),
	}
}

// structFields is a wrapper type for struct fields.
type structFields []Field

// StructFlat creates multiple top-level fields from a struct's exported fields.
// Unlike Struct, this doesn't nest under a key - fields are added directly.
func StructFlat(v any) []Field {
	if v == nil {
		return nil
	}

	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return nil
		}
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return nil
	}

	return extractStructFields("", val)
}

// extractStructFields extracts fields from a struct value.
func extractStructFields(prefix string, val reflect.Value) []Field {
	typ := val.Type()

	// Check cache
	cached, ok := structFieldCache.Load(typ)
	var infos []structFieldInfo
	if ok {
		infos = cached.([]structFieldInfo)
	} else {
		infos = buildStructFieldInfos(typ)
		structFieldCache.Store(typ, infos)
	}

	fields := make([]Field, 0, len(infos))
	for _, info := range infos {
		fieldVal := val.FieldByIndex(info.index)

		// Skip zero values if omitempty
		if info.omitempty && isZeroValue(fieldVal) {
			continue
		}

		name := info.name
		if prefix != "" {
			name = prefix + "." + name
		}

		field := valueToField(name, fieldVal)
		fields = append(fields, field)
	}

	return fields
}

// buildStructFieldInfos builds field info for a struct type.
func buildStructFieldInfos(typ reflect.Type) []structFieldInfo {
	var infos []structFieldInfo

	for i := 0; i < typ.NumField(); i++ {
		sf := typ.Field(i)

		// Skip unexported fields
		if !sf.IsExported() {
			continue
		}

		// Skip fields with json:"-"
		jsonTag := sf.Tag.Get("json")
		if jsonTag == "-" {
			continue
		}

		// Parse json tag
		name := sf.Name
		var omitempty bool
		if jsonTag != "" {
			parts := strings.Split(jsonTag, ",")
			if parts[0] != "" {
				name = parts[0]
			}
			for _, opt := range parts[1:] {
				if opt == "omitempty" {
					omitempty = true
				}
			}
		}

		// Also check log tag for custom log field names
		if logTag := sf.Tag.Get("log"); logTag != "" {
			if logTag == "-" {
				continue
			}
			parts := strings.Split(logTag, ",")
			if parts[0] != "" {
				name = parts[0]
			}
		}

		infos = append(infos, structFieldInfo{
			index:    sf.Index,
			name:     name,
			omitempty: omitempty,
		})
	}

	return infos
}

// valueToField converts a reflect.Value to a Field.
func valueToField(key string, val reflect.Value) Field {
	// Handle pointers
	if val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return String(key, "null")
		}
		val = val.Elem()
	}

	// Handle interface
	if val.Kind() == reflect.Interface {
		if val.IsNil() {
			return String(key, "null")
		}
		return Any(key, val.Interface())
	}

	switch val.Kind() {
	case reflect.String:
		return String(key, val.String())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		// Check for time.Duration
		if val.Type() == reflect.TypeOf(time.Duration(0)) {
			return Duration(key, time.Duration(val.Int()))
		}
		return Int64(key, val.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return Uint64(key, val.Uint())
	case reflect.Float32, reflect.Float64:
		return Float64(key, val.Float())
	case reflect.Bool:
		return Bool(key, val.Bool())
	case reflect.Struct:
		// Check for time.Time
		if val.Type() == reflect.TypeOf(time.Time{}) {
			return Time(key, val.Interface().(time.Time))
		}
		// Nested struct - could flatten or marshal
		return Any(key, val.Interface())
	case reflect.Slice, reflect.Array:
		// Check for []byte
		if val.Type().Elem().Kind() == reflect.Uint8 {
			return Bytes(key, val.Bytes())
		}
		return Any(key, val.Interface())
	case reflect.Map:
		return Any(key, val.Interface())
	default:
		return Any(key, val.Interface())
	}
}

// isZeroValue checks if a value is the zero value for its type.
func isZeroValue(val reflect.Value) bool {
	switch val.Kind() {
	case reflect.String:
		return val.String() == ""
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return val.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return val.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return val.Float() == 0
	case reflect.Bool:
		return !val.Bool()
	case reflect.Ptr, reflect.Interface, reflect.Slice, reflect.Map, reflect.Chan, reflect.Func:
		return val.IsNil()
	case reflect.Struct:
		// Check for time.Time zero
		if val.Type() == reflect.TypeOf(time.Time{}) {
			return val.Interface().(time.Time).IsZero()
		}
		return false
	default:
		return false
	}
}

// V creates a field with automatic type detection.
// This is a convenience function for when you don't want to specify the type.
func V(key string, value any) Field {
	return Any(key, value)
}

// Object creates a field from any value, attempting to extract structure.
// For structs, it extracts fields. For other types, it uses Any.
func Object(key string, v any) Field {
	if v == nil {
		return String(key, "null")
	}

	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return String(key, "null")
		}
		val = val.Elem()
	}

	if val.Kind() == reflect.Struct {
		return Struct(key, v)
	}

	return Any(key, v)
}

