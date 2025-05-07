package decodini

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// http://github.com/lukasl-dev/decodini/issues/2
func TestIssue_2(t *testing.T) {
	type (
		embedded struct {
			Foo string
		}
		outer struct {
			embedded
		}
	)

	a := assert.New(t)

	from := map[string]any{
		"Foo": "bar",
	}
	tr := Encode(nil, from)

	to, err := Decode[outer](nil, tr)
	a.NoError(err)

	expected := outer{
		embedded: embedded{
			Foo: "bar",
		},
	}
	a.Equal(expected, to)
}
