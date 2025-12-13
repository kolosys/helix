package helix

import (
	"encoding/json"
	"io"
)

// UnmarshalJSON unmarshals a JSON request body into a struct.
// Returns an error if the request body is not a valid JSON or if the struct cannot be unmarshalled.
// The error will contain the stack trace of the error.
func UnmarshalJSON[T any](reader io.Reader) (*T, error) {
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, ErrBadRequest.WithDetailf("failed to read request body").WithErr(err)
	}

	var v T
	if err := json.Unmarshal(data, &v); err != nil {
		return nil, ErrBadRequest.WithDetailf("failed to unmarshal JSON body").WithErr(err)
	}
	return &v, nil
}
