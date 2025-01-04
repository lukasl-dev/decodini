package decodini

import (
	"fmt"
	"strings"
)

type DecodeError struct {
	Path []any
	Err  error
}

var _ error = (*DecodeError)(nil)

func newDecodeError(path []any, err error) *DecodeError {
	return &DecodeError{Path: path, Err: err}
}

func newDecodeErrorf(path []any, format string, args ...any) *DecodeError {
	return newDecodeError(path, fmt.Errorf(format, args...))
}

// Unwrap returns the underlying error.
func (e *DecodeError) Unwrap() error { return e.Err }

// Error returns the error message.
func (e *DecodeError) Error() string {
	return fmt.Sprintf("decodini: encode: failed at %s: %s", e.PathString(), e.Err)
}

// PathSTring returns a dot-separated string representation of the path.
func (e *DecodeError) PathString() string {
	if len(e.Path) == 0 {
		return "<root>"
	}
	path := make([]string, len(e.Path))
	for i, item := range e.Path {
		path[i] = fmt.Sprintf("%v", item)
	}
	return strings.Join(path, ".")
}
