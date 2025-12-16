package helix

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"slices"
	"strconv"
	"strings"
	"sync"
)

// Binding errors
var (
	ErrBindingFailed     = errors.New("helix: binding failed")
	ErrUnsupportedType   = errors.New("helix: unsupported type for binding")
	ErrInvalidJSON       = errors.New("helix: invalid JSON body")
	ErrRequiredField     = errors.New("helix: required field missing")
	ErrBodyAlreadyRead   = errors.New("helix: request body already read")
	ErrInvalidFieldValue = errors.New("helix: invalid field value")
)

// Struct tag names for binding sources
const (
	tagPath   = "path"
	tagQuery  = "query"
	tagHeader = "header"
	tagJSON   = "json"
	tagForm   = "form"
)

// bindingCache caches reflected struct information for performance.
var bindingCache sync.Map

type fieldInfo struct {
	index     int
	name      string
	source    string // path, query, header, json, form
	required  bool
	omitEmpty bool
	fieldType reflect.Type
}

type structInfo struct {
	fields []fieldInfo
}

// Bind binds path parameters, query parameters, headers, and JSON body to a struct.
// The binding sources are determined by struct tags:
//   - `path:"name"` - binds from URL path parameters
//   - `query:"name"` - binds from URL query parameters
//   - `header:"name"` - binds from HTTP headers
//   - `json:"name"` - binds from JSON body
//   - `form:"name"` - binds from form data
func Bind[T any](r *http.Request) (T, error) {
	var result T

	resultVal := reflect.ValueOf(&result).Elem()
	resultType := resultVal.Type()

	// Get or create struct info
	info := getStructInfo(resultType)

	// First bind non-body fields
	for _, field := range info.fields {
		if field.source == tagJSON {
			continue // Handle JSON separately
		}

		var value string
		switch field.source {
		case tagPath:
			value = Param(r, field.name)
		case tagQuery:
			value = Query(r, field.name)
		case tagHeader:
			value = r.Header.Get(field.name)
		case tagForm:
			if err := r.ParseForm(); err == nil {
				value = r.FormValue(field.name)
			}
		}

		if value == "" {
			if field.required {
				return result, fmt.Errorf("%w: %s", ErrRequiredField, field.name)
			}
			continue
		}

		if err := setFieldValue(resultVal.Field(field.index), value); err != nil {
			return result, fmt.Errorf("%w: field %s: %v", ErrInvalidFieldValue, field.name, err)
		}
	}

	// Check if there are any JSON fields
	hasJSONFields := false
	for _, field := range info.fields {
		if field.source == tagJSON {
			hasJSONFields = true
			break
		}
	}

	// Bind JSON body if there are JSON fields
	if hasJSONFields && r.Body != nil && r.ContentLength != 0 {
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&result); err != nil && err != io.EOF {
			return result, fmt.Errorf("%w: %v", ErrInvalidJSON, err)
		}
	}

	return result, nil
}

// BindJSON binds the JSON request body to a struct.
func BindJSON[T any](r *http.Request) (T, error) {
	var result T

	if r.Body == nil {
		return result, ErrInvalidJSON
	}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&result); err != nil {
		return result, fmt.Errorf("%w: %v", ErrInvalidJSON, err)
	}

	return result, nil
}

// BindQuery binds URL query parameters to a struct.
// Uses the `query` struct tag to determine field names.
func BindQuery[T any](r *http.Request) (T, error) {
	var result T

	resultVal := reflect.ValueOf(&result).Elem()
	resultType := resultVal.Type()

	info := getStructInfo(resultType)

	for _, field := range info.fields {
		if field.source != tagQuery {
			continue
		}

		value := Query(r, field.name)
		if value == "" {
			if field.required {
				return result, fmt.Errorf("%w: %s", ErrRequiredField, field.name)
			}
			continue
		}

		if err := setFieldValue(resultVal.Field(field.index), value); err != nil {
			return result, fmt.Errorf("%w: field %s: %v", ErrInvalidFieldValue, field.name, err)
		}
	}

	return result, nil
}

// BindPath binds URL path parameters to a struct.
// Uses the `path` struct tag to determine field names.
func BindPath[T any](r *http.Request) (T, error) {
	var result T

	resultVal := reflect.ValueOf(&result).Elem()
	resultType := resultVal.Type()

	info := getStructInfo(resultType)

	for _, field := range info.fields {
		if field.source != tagPath {
			continue
		}

		value := Param(r, field.name)
		if value == "" {
			if field.required {
				return result, fmt.Errorf("%w: %s", ErrRequiredField, field.name)
			}
			continue
		}

		if err := setFieldValue(resultVal.Field(field.index), value); err != nil {
			return result, fmt.Errorf("%w: field %s: %v", ErrInvalidFieldValue, field.name, err)
		}
	}

	return result, nil
}

// BindHeader binds HTTP headers to a struct.
// Uses the `header` struct tag to determine field names.
func BindHeader[T any](r *http.Request) (T, error) {
	var result T

	resultVal := reflect.ValueOf(&result).Elem()
	resultType := resultVal.Type()

	info := getStructInfo(resultType)

	for _, field := range info.fields {
		if field.source != tagHeader {
			continue
		}

		value := r.Header.Get(field.name)
		if value == "" {
			if field.required {
				return result, fmt.Errorf("%w: %s", ErrRequiredField, field.name)
			}
			continue
		}

		if err := setFieldValue(resultVal.Field(field.index), value); err != nil {
			return result, fmt.Errorf("%w: field %s: %v", ErrInvalidFieldValue, field.name, err)
		}
	}

	return result, nil
}

// getStructInfo gets or creates cached struct information.
func getStructInfo(t reflect.Type) *structInfo {
	if cached, ok := bindingCache.Load(t); ok {
		return cached.(*structInfo)
	}

	info := &structInfo{
		fields: make([]fieldInfo, 0),
	}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if !field.IsExported() {
			continue
		}

		// Check each tag type
		for _, tagName := range []string{tagPath, tagQuery, tagHeader, tagJSON, tagForm} {
			tag := field.Tag.Get(tagName)
			if tag == "" {
				continue
			}

			name, opts := parseTag(tag)
			if name == "-" {
				continue
			}

			info.fields = append(info.fields, fieldInfo{
				index:     i,
				name:      name,
				source:    tagName,
				required:  containsOption(opts, "required"),
				omitEmpty: containsOption(opts, "omitempty"),
				fieldType: field.Type,
			})
		}
	}

	bindingCache.Store(t, info)
	return info
}

// parseTag parses a struct tag into name and options.
func parseTag(tag string) (string, []string) {
	parts := strings.Split(tag, ",")
	if len(parts) == 0 {
		return "", nil
	}
	return parts[0], parts[1:]
}

// containsOption checks if an option is present in the options slice.
func containsOption(opts []string, option string) bool {
	return slices.Contains(opts, option)
}

// setFieldValue sets a struct field value from a string.
func setFieldValue(field reflect.Value, value string) error {
	if !field.CanSet() {
		return ErrUnsupportedType
	}

	switch field.Kind() {
	case reflect.String:
		field.SetString(value)

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		v, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return err
		}
		field.SetInt(v)

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		v, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return err
		}
		field.SetUint(v)

	case reflect.Float32, reflect.Float64:
		v, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return err
		}
		field.SetFloat(v)

	case reflect.Bool:
		v, err := strconv.ParseBool(value)
		if err != nil {
			// Also accept "yes", "no", "on", "off"
			switch strings.ToLower(value) {
			case "yes", "on":
				v = true
			case "no", "off":
				v = false
			default:
				return err
			}
		}
		field.SetBool(v)

	case reflect.Slice:
		// Handle []string
		if field.Type().Elem().Kind() == reflect.String {
			values := strings.Split(value, ",")
			field.Set(reflect.ValueOf(values))
		} else {
			return ErrUnsupportedType
		}

	case reflect.Pointer:
		// Create a new value of the underlying type
		elemType := field.Type().Elem()
		newVal := reflect.New(elemType)
		if err := setFieldValue(newVal.Elem(), value); err != nil {
			return err
		}
		field.Set(newVal)

	default:
		return ErrUnsupportedType
	}

	return nil
}

// Validatable is an interface for types that can validate themselves.
type Validatable interface {
	Validate() error
}

// BindAndValidate binds and validates a request.
// If the bound type implements Validatable, Validate() is called after binding.
func BindAndValidate[T any](r *http.Request) (T, error) {
	result, err := Bind[T](r)
	if err != nil {
		return result, err
	}

	// Check if the result implements Validatable
	if v, ok := any(&result).(Validatable); ok {
		if err := v.Validate(); err != nil {
			return result, err
		}
	}

	return result, nil
}
