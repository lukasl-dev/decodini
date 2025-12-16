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

// https://github.com/lukasl-dev/decodini/issues/4
func TestIssue_4(t *testing.T) {
	type Payload struct {
		Value *int `decodini:"value"`
	}

	a := assert.New(t)

	from := map[string]any{
		"value": 42,
	}
	tr := Encode(nil, from)

	to, err := Decode[Payload](nil, tr)
	a.NoError(err)

	expected := Payload{
		Value: ptr(42),
	}
	a.Equal(expected, to)
}
