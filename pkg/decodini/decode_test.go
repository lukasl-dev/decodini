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
		B *int
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
		B: ptr(42),
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

func TestDecode_OneEmbeddedStruct(t *testing.T) {
	type (
		Embedded struct {
			A string `decodini:"a"`
			B int
			C bool `decodini:"-"`
		}
		Outer struct {
			Embedded
		}
	)

	a := assert.New(t)

	from := map[string]any{
		"a": "foo",
		"B": 42,
	}
	tr := Encode(nil, from)

	to, err := Decode[Outer](nil, tr)
	a.NoError(err)

	expected := Outer{Embedded: Embedded{A: "foo", B: 42}}
	a.Equal(expected, to)
}

func TestDecode_TwoEmbeddedStructs(t *testing.T) {
	type (
		EmbeddedA struct {
			A string `decodini:"a"`
		}
		EmbeddedB struct {
			B int `decodini:"b"`
		}
		Outer struct {
			EmbeddedA
			EmbeddedB
		}
	)

	a := assert.New(t)

	from := map[string]any{
		"a": "foo",
		"b": 7,
	}
	tr := Encode(nil, from)

	to, err := Decode[Outer](nil, tr)
	a.NoError(err)

	expected := Outer{EmbeddedA: EmbeddedA{A: "foo"}, EmbeddedB: EmbeddedB{B: 7}}
	a.Equal(expected, to)
}

func TestDecode_Slice_to_SliceOfPointers(t *testing.T) {
	a := assert.New(t)

	from := []int{1, 2, 3}
	tr := Encode(nil, from)

	to, err := Decode[[]*int](nil, tr)
	a.NoError(err)

	a.Equal(len(from), len(to))
	for i, v := range from {
		if a.NotNil(to[i]) {
			a.Equal(v, *to[i])
		}
	}
}

func TestDecode_Map_to_MapOfPointers(t *testing.T) {
	a := assert.New(t)

	from := map[string]int{
		"a": 1,
		"b": 2,
	}
	tr := Encode(nil, from)

	to, err := Decode[map[string]*int](nil, tr)
	a.NoError(err)

	a.Equal(len(from), len(to))
	for k, v := range from {
		ptrVal, ok := to[k]
		a.True(ok)
		if a.NotNil(ptrVal) {
			a.Equal(v, *ptrVal)
		}
	}
}

func TestDecode_PointerToPointerStructField(t *testing.T) {
	type toStruct struct {
		V **int `decodini:"v"`
	}

	a := assert.New(t)

	from := map[string]any{
		"v": 42,
	}
	tr := Encode(nil, from)

	to, err := Decode[toStruct](nil, tr)
	a.NoError(err)

	a.NotNil(to.V)
	a.NotNil(*to.V)
	a.Equal(42, **to.V)
}

func ptr[T any](value T) *T {
	return &value
}
