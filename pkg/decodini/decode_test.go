package decodini

import (
	"testing"
	"unicode/utf16"

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

func TestDecode_FromNilPointerOrInterface_SetsZero(t *testing.T) {
	a := assert.New(t)

	type S struct {
		A string
		B int
	}

	var p *S = nil
	trPtr := Encode(nil, p)
	dst1, err1 := Decode[S](nil, trPtr)
	a.NoError(err1)
	a.Equal(S{}, dst1)

	var i any = (*S)(nil)
	trIface := Encode(nil, i)
	dst2, err2 := Decode[S](nil, trIface)
	a.NoError(err2)
	a.Equal(S{}, dst2)
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

func TestDecode_String_to_ByteSlice(t *testing.T) {
	a := assert.New(t)

	from := "héllö 🌍"
	tr := Encode(nil, from)

	to, err := Decode[[]byte](nil, tr)
	a.NoError(err)
	a.Equal([]byte(from), to)
}

func TestDecode_String_to_RuneSlice(t *testing.T) {
	a := assert.New(t)

	from := "héllö 🌍"
	tr := Encode(nil, from)

	to, err := Decode[[]rune](nil, tr)
	a.NoError(err)
	a.Equal([]rune(from), to)
}

func TestDecode_ByteSlice_to_String(t *testing.T) {
	a := assert.New(t)

	from := []byte("héllö 🌍")
	tr := Encode(nil, from)

	to, err := Decode[string](nil, tr)
	a.NoError(err)
	a.Equal(string(from), to)
}

func TestDecode_RuneSlice_to_String(t *testing.T) {
	a := assert.New(t)

	from := []rune("héllö 🌍")
	tr := Encode(nil, from)

	to, err := Decode[string](nil, tr)
	a.NoError(err)
	a.Equal(string(from), to)
}

func TestDecode_String_to_Int8Slice(t *testing.T) {
	a := assert.New(t)

	from := "héllö 🌍"
	tr := Encode(nil, from)

	to, err := Decode[[]int8](nil, tr)
	a.NoError(err)

	b := []byte(from)
	expected := make([]int8, len(b))
	for i := range len(b) {
		expected[i] = int8(b[i])
	}
	a.Equal(expected, to)
}

func TestDecode_Int8Slice_to_String(t *testing.T) {
	a := assert.New(t)

	expected := "héllö 🌍"
	b := []byte(expected)
	from := make([]int8, len(b))
	for i := range len(b) {
		from[i] = int8(b[i])
	}
	tr := Encode(nil, from)

	to, err := Decode[string](nil, tr)
	a.NoError(err)
	a.Equal(expected, to)
}

func TestDecode_String_to_Uint16Slice_UTF16(t *testing.T) {
	a := assert.New(t)

	from := "héllö 🌍"
	tr := Encode(nil, from)

	to, err := Decode[[]uint16](nil, tr)
	a.NoError(err)
	a.Equal(utf16.Encode([]rune(from)), to)
}

func TestDecode_Uint16Slice_to_String_UTF16(t *testing.T) {
	a := assert.New(t)

	expected := "héllö 🌍"
	from := utf16.Encode([]rune(expected))
	tr := Encode(nil, from)

	to, err := Decode[string](nil, tr)
	a.NoError(err)
	a.Equal(expected, to)
}

func TestDecode_String_to_Int16Slice_UTF16(t *testing.T) {
	a := assert.New(t)

	from := "héllö 🌍"
	tr := Encode(nil, from)

	to, err := Decode[[]int16](nil, tr)
	a.NoError(err)

	u := utf16.Encode([]rune(from))
	expected := make([]int16, len(u))
	for i := range len(u) {
		expected[i] = int16(u[i])
	}
	a.Equal(expected, to)
}

func TestDecode_Int16Slice_to_String_UTF16(t *testing.T) {
	a := assert.New(t)

	expected := "héllö 🌍"
	u := utf16.Encode([]rune(expected))
	from := make([]int16, len(u))
	for i := range len(u) {
		from[i] = int16(u[i])
	}
	tr := Encode(nil, from)

	to, err := Decode[string](nil, tr)
	a.NoError(err)
	a.Equal(expected, to)
}

func TestDecode_String_to_Int32Slice(t *testing.T) {
	a := assert.New(t)

	from := "héllö 🌍"
	tr := Encode(nil, from)

	to, err := Decode[[]int32](nil, tr)
	a.NoError(err)
	a.Equal([]rune(from), to)
}

func TestDecode_Int32Slice_to_String(t *testing.T) {
	a := assert.New(t)

	expected := "héllö 🌍"
	from := []int32(expected)
	tr := Encode(nil, from)

	to, err := Decode[string](nil, tr)
	a.NoError(err)
	a.Equal(expected, to)
}

func TestDecode_String_to_Int64Slice(t *testing.T) {
	a := assert.New(t)

	from := "héllö 🌍"
	tr := Encode(nil, from)

	to, err := Decode[[]int64](nil, tr)
	a.NoError(err)

	r := []rune(from)
	expected := make([]int64, len(r))
	for i := range len(r) {
		expected[i] = int64(r[i])
	}
	a.Equal(expected, to)
}

func TestDecode_Int64Slice_to_String(t *testing.T) {
	a := assert.New(t)

	expected := "héllö 🌍"
	r := []rune(expected)
	from := make([]int64, len(r))
	for i := range len(r) {
		from[i] = int64(r[i])
	}
	tr := Encode(nil, from)

	to, err := Decode[string](nil, tr)
	a.NoError(err)
	a.Equal(expected, to)
}

func TestDecode_String_to_Uint32Slice(t *testing.T) {
	a := assert.New(t)

	from := "héllö 🌍"
	tr := Encode(nil, from)

	to, err := Decode[[]uint32](nil, tr)
	a.NoError(err)

	r := []rune(from)
	expected := make([]uint32, len(r))
	for i := range len(r) {
		expected[i] = uint32(r[i])
	}
	a.Equal(expected, to)
}

func TestDecode_Uint32Slice_to_String(t *testing.T) {
	a := assert.New(t)

	expected := "héllö 🌍"
	r := []rune(expected)
	from := make([]uint32, len(r))
	for i := range len(r) {
		from[i] = uint32(r[i])
	}
	tr := Encode(nil, from)

	to, err := Decode[string](nil, tr)
	a.NoError(err)
	a.Equal(expected, to)
}

func TestDecode_String_to_Uint64Slice(t *testing.T) {
	a := assert.New(t)

	from := "héllö 🌍"
	tr := Encode(nil, from)

	to, err := Decode[[]uint64](nil, tr)
	a.NoError(err)

	r := []rune(from)
	expected := make([]uint64, len(r))
	for i := range len(r) {
		expected[i] = uint64(r[i])
	}
	a.Equal(expected, to)
}

func TestDecode_Uint64Slice_to_String(t *testing.T) {
	a := assert.New(t)

	expected := "héllö 🌍"
	r := []rune(expected)
	from := make([]uint64, len(r))
	for i := range len(r) {
		from[i] = uint64(r[i])
	}
	tr := Encode(nil, from)

	to, err := Decode[string](nil, tr)
	a.NoError(err)
	a.Equal(expected, to)
}

func TestDecode_EmptyString_to_ByteSlice(t *testing.T) {
	a := assert.New(t)

	from := ""
	tr := Encode(nil, from)

	to, err := Decode[[]byte](nil, tr)
	a.NoError(err)
	a.Equal([]byte{}, to)
}

func TestDecode_EmptyByteSlice_to_String(t *testing.T) {
	a := assert.New(t)

	from := []byte{}
	tr := Encode(nil, from)

	to, err := Decode[string](nil, tr)
	a.NoError(err)
	a.Equal("", to)
}

func TestDecode_ASCIIOnly_String_to_ByteSlice(t *testing.T) {
	a := assert.New(t)

	from := "hello world"
	tr := Encode(nil, from)

	to, err := Decode[[]byte](nil, tr)
	a.NoError(err)
	a.Equal([]byte(from), to)
}

func TestDecode_ASCIIOnly_ByteSlice_to_String(t *testing.T) {
	a := assert.New(t)

	from := []byte("hello world")
	tr := Encode(nil, from)

	to, err := Decode[string](nil, tr)
	a.NoError(err)
	a.Equal(string(from), to)
}

func ptr[T any](value T) *T {
	return &value
}
