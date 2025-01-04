package decodini

import (
	"fmt"
	"strings"
)

type EncodeError struct {
	Path []any
	Err  error
}

var _ error = (*EncodeError)(nil)

func newEncodeError(path []any, err error) *EncodeError {
	return &EncodeError{Path: path, Err: err}
}

func newEncodeErrorf(path []any, format string, args ...any) *EncodeError {
	return newEncodeError(path, fmt.Errorf(format, args...))
}

// Unwrap returns the underlying error.
func (e *EncodeError) Unwrap() error { return e.Err }

// Error returns the error message.
func (e *EncodeError) Error() string {
	return fmt.Sprintf("decodini: encode: failed at %s: %s", e.PathString(), e.Err)
}

// PathSTring returns a dot-separated string representation of the path.
func (e *EncodeError) PathString() string {
	path := make([]string, len(e.Path))
	for i, item := range e.Path {
		path[i] = fmt.Sprintf("%v", item)
	}
	return strings.Join(path, ".")
}
