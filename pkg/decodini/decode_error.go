package decodini

import (
	"fmt"
	"strings"
)

type DecodeError struct {
	Node *Tree
	Err  error
}

var _ error = (*DecodeError)(nil)

func newDecodeError(node *Tree, err error) *DecodeError {
	return &DecodeError{Node: node, Err: err}
}

func newDecodeErrorf(node *Tree, format string, args ...any) *DecodeError {
	return newDecodeError(node, fmt.Errorf(format, args...))
}

// Unwrap returns the underlying error.
func (e *DecodeError) Unwrap() error { return e.Err }

// Error returns the error message.
func (e *DecodeError) Error() string {
	return fmt.Sprintf("decodini: decode: failed at %s: %s", e.PathString(), e.Err)
}

// PathSTring returns a dot-separated string representation of the path.
func (e *DecodeError) PathString() string {
	path := e.Node.Path()
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
