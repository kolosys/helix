package helix

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

// contextKey is a private type for context keys.
type contextKey int

const (
	paramsKey contextKey = iota
	servicesCtxKey
)

// setParams stores path parameters in the context.
func setParams(ctx context.Context, ps *params) context.Context {
	return context.WithValue(ctx, paramsKey, ps)
}

// getParams retrieves path parameters from the context.
func getParams(ctx context.Context) *params {
	ps, _ := ctx.Value(paramsKey).(*params)
	return ps
}

// Param returns the value of a path parameter.
// Returns an empty string if the parameter does not exist.
func Param(r *http.Request, name string) string {
	ps := getParams(r.Context())
	if ps == nil {
		return ""
	}
	return ps.get(name)
}

// ParamInt returns the value of a path parameter as an int.
// Returns an error if the parameter does not exist or cannot be parsed.
func ParamInt(r *http.Request, name string) (int, error) {
	s := Param(r, name)
	if s == "" {
		return 0, fmt.Errorf("helix: parameter %q not found", name)
	}
	return strconv.Atoi(s)
}

// ParamInt64 returns the value of a path parameter as an int64.
// Returns an error if the parameter does not exist or cannot be parsed.
func ParamInt64(r *http.Request, name string) (int64, error) {
	s := Param(r, name)
	if s == "" {
		return 0, fmt.Errorf("helix: parameter %q not found", name)
	}
	return strconv.ParseInt(s, 10, 64)
}

// ParamUUID returns the value of a path parameter validated as a UUID.
// Returns an error if the parameter does not exist or is not a valid UUID format.
func ParamUUID(r *http.Request, name string) (string, error) {
	s := Param(r, name)
	if s == "" {
		return "", fmt.Errorf("helix: parameter %q not found", name)
	}
	if !isValidUUID(s) {
		return "", fmt.Errorf("helix: parameter %q is not a valid UUID", name)
	}
	return s, nil
}

// isValidUUID validates a UUID string (accepts both with and without hyphens).
func isValidUUID(s string) bool {
	// Standard UUID format: xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx (36 chars)
	// Without hyphens: xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx (32 chars)
	switch len(s) {
	case 36:
		// Validate hyphen positions
		if s[8] != '-' || s[13] != '-' || s[18] != '-' || s[23] != '-' {
			return false
		}
		// Check hex characters
		for i, c := range s {
			if i == 8 || i == 13 || i == 18 || i == 23 {
				continue
			}
			if !isHexChar(byte(c)) {
				return false
			}
		}
		return true
	case 32:
		// Check all hex characters
		for i := 0; i < len(s); i++ {
			if !isHexChar(s[i]) {
				return false
			}
		}
		return true
	default:
		return false
	}
}

func isHexChar(c byte) bool {
	return (c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')
}

// Query returns the first value of a query parameter.
// Returns an empty string if the parameter does not exist.
func Query(r *http.Request, name string) string {
	return r.URL.Query().Get(name)
}

// QueryDefault returns the first value of a query parameter or a default value.
func QueryDefault(r *http.Request, name, defaultVal string) string {
	v := r.URL.Query().Get(name)
	if v == "" {
		return defaultVal
	}
	return v
}

// QueryInt returns the first value of a query parameter as an int.
// Returns the default value if the parameter does not exist or cannot be parsed.
func QueryInt(r *http.Request, name string, defaultVal int) int {
	s := r.URL.Query().Get(name)
	if s == "" {
		return defaultVal
	}
	v, err := strconv.Atoi(s)
	if err != nil {
		return defaultVal
	}
	return v
}

// QueryInt64 returns the first value of a query parameter as an int64.
// Returns the default value if the parameter does not exist or cannot be parsed.
func QueryInt64(r *http.Request, name string, defaultVal int64) int64 {
	s := r.URL.Query().Get(name)
	if s == "" {
		return defaultVal
	}
	v, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return defaultVal
	}
	return v
}

// QueryBool returns the first value of a query parameter as a bool.
// Returns false if the parameter does not exist or cannot be parsed.
// Accepts "1", "t", "T", "true", "TRUE", "True" as true.
// Accepts "0", "f", "F", "false", "FALSE", "False" as false.
func QueryBool(r *http.Request, name string) bool {
	s := r.URL.Query().Get(name)
	if s == "" {
		return false
	}
	v, err := strconv.ParseBool(s)
	if err != nil {
		// Also check for "yes", "no", "on", "off"
		switch strings.ToLower(s) {
		case "yes", "on":
			return true
		default:
			return false
		}
	}
	return v
}

// QuerySlice returns all values of a query parameter as a string slice.
// Returns nil if the parameter does not exist.
func QuerySlice(r *http.Request, name string) []string {
	values, ok := r.URL.Query()[name]
	if !ok {
		return nil
	}
	return values
}

// QueryFloat64 returns the first value of a query parameter as a float64.
// Returns the default value if the parameter does not exist or cannot be parsed.
func QueryFloat64(r *http.Request, name string, defaultVal float64) float64 {
	s := r.URL.Query().Get(name)
	if s == "" {
		return defaultVal
	}
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return defaultVal
	}
	return v
}
