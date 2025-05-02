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

func TestDecode_ShallowStruct_to_ShallowStruct(t *testing.T) {
	type fromStruct struct {
		A string `decodini:"a"`
		B int
		C bool `decodini:"-"`
		d float64
	}

	type toStruct struct {
		A string `decodini:"a"`
		B int
		C int `decodini:"-"`
		x bool
	}

	a := assert.New(t)

	from := fromStruct{
		A: "foo",
		B: 42,
		C: true,
		d: 420.0,
	}
	tr := Encode(nil, from)

	to, err := Decode[toStruct](nil, tr)
	a.NoError(err)

	expected := toStruct{
		A: "foo",
		B: 42,
		C: 0,
		x: false,
	}
	a.Equal(expected, to)
}

func TestDecode_ShallowStruct_to_ShallowMap(t *testing.T) {
	type fromStruct struct {
		A string `decodini:"a"`
		B int
		C bool `decodini:"-"`
		d float64
	}

	a := assert.New(t)

	from := fromStruct{
		A: "foo",
		B: 42,
		C: true,
		d: 420.0,
	}
	tr := Encode(nil, from)

	to, err := Decode[map[string]any](nil, tr)
	a.NoError(err)

	expected := map[string]any{
		"a": "foo",
		"B": 42,
	}
	a.Equal(expected, to)
}

func TestDecode_ShallowMap_to_ShallowStruct(t *testing.T) {
	type toStruct struct {
		A string `decodini:"a"`
		B int
		C int `decodini:"-"`
		x bool
	}

	a := assert.New(t)

	from := map[string]any{
		"a": "foo",
		"B": 42,
		"C": true,
	}
	tr := Encode(nil, from)

	to, err := Decode[toStruct](nil, tr)
	a.NoError(err)

	expected := toStruct{
		A: "foo",
		B: 42,
		C: 0,
		x: false,
	}
	a.Equal(expected, to)
}

func TestDecode_ShallowMap_to_ShallowMap(t *testing.T) {
	a := assert.New(t)

	expected := map[string]any{
		"a": 12,
		"B": true,
		"c": "foo",
	}
	tr := Encode(nil, expected)

	actual, err := Decode[map[string]any](nil, tr)
	a.EqualValues(expected, actual)
	a.NoError(err)
}

func TestDecode_ShallowMap_to_ShallowSlice(t *testing.T) {
	a := assert.New(t)

	from := map[string]any{
		"a": 12,
		"B": true,
		"c": "foo",
	}
	tr := Encode(nil, from)

	actual, err := Decode[[]any](nil, tr)
	a.Contains(actual, 12)
	a.Contains(actual, true)
	a.Contains(actual, "foo")
	a.NoError(err)
}

func TestDecode_ShallowSlice_to_ShallowSlice(t *testing.T) {
	a := assert.New(t)

	expected := []string{"foo", "bar", "baz"}
	tr := Encode(nil, expected)

	actual, err := Decode[[]string](nil, tr)
	a.EqualValues(expected, actual)
	a.NoError(err)
}
