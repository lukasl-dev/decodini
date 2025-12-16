package decodini

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// https://github.com/lukasl-dev/decodini/issues/2
func TestIssue_2(t *testing.T) {
	type (
		Embedded struct {
			Foo string
		}
		Outer struct {
			Embedded
		}
	)

	a := assert.New(t)

	from := map[string]any{
		"Foo": "bar",
	}
	tr := Encode(nil, from)

	to, err := Decode[Outer](nil, tr)
	a.NoError(err)

	expected := Outer{
		Embedded: Embedded{
			Foo: "bar",
		},
	}
	a.Equal(expected, to)
}
