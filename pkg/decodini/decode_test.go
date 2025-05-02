package decodini

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecode_String_to_String(t *testing.T) {
	a := assert.New(t)

	expected := "decodini"
	tr := Encode(nil, expected)

	actual, err := Decode[string](nil, tr)
	a.Equal(expected, actual)
	a.NoError(err)
}
