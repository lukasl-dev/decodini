package decodini

import (
	"fmt"
	"strings"
)

type DecodeError struct {
	From *Tree
	Into DecodeTarget
	Err  error
}

var _ error = (*DecodeError)(nil)

func newDecodeError(from *Tree, into DecodeTarget, err error) *DecodeError {
	return &DecodeError{From: from, Into: into, Err: err}
}

func newDecodeErrorf(
	from *Tree,
	into DecodeTarget,
	format string,
	args ...any,
) *DecodeError {
	return newDecodeError(from, into, fmt.Errorf(format, args...))
}

// Unwrap returns the underlying error.
func (e *DecodeError) Unwrap() error { return e.Err }

// Error returns the error message.
func (e *DecodeError) Error() string {
	return fmt.Sprintf("decodini: decode: failed at %s: %s", e.PathString(), e.Err)
}

// PathSTring returns a dot-separated string representation of the path.
func (e *DecodeError) PathString() string {
	path := e.From.Path()
	if path == nil {
		return "<root>"
	}

	var sb strings.Builder
	for i, item := range path {
		sb.WriteString(fmt.Sprintf("%v", item))
		if i < len(path)-1 {
			sb.WriteString(".")
		}
	}
	return sb.String()
}
