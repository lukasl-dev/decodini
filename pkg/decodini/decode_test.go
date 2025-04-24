package decodini

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecode_Nil(t *testing.T) {
	t.Parallel()

	a := assert.New(t)

	tr := Encode(nil, nil)

	var dst any
	err := Decode(nil, tr, &dst)

	a.NoError(err)
	a.Nil(dst)
}

func TestDecode_String(t *testing.T) {
	t.Parallel()

	val := "test"
	tr := Encode(nil, val)

	t.Run("Interface", func(t *testing.T) {
		a := assert.New(t)

		var dst any
		err := Decode(nil, tr, &dst)

		a.NoError(err)
		a.Equal(val, dst)
	})

	t.Run("String", func(t *testing.T) {
		a := assert.New(t)

		var dst string
		err := Decode(nil, tr, &dst)

		a.NoError(err)
		a.Equal(val, dst)
	})
}

func TestDecode_Bool(t *testing.T) {
	t.Parallel()

	val := bool(true)
	tr := Encode(nil, val)

	t.Run("Interface", func(t *testing.T) {
		a := assert.New(t)

		var dst any
		err := Decode(nil, tr, &dst)

		a.NoError(err)
		a.Equal(val, dst)
	})

	t.Run("Bool", func(t *testing.T) {
		a := assert.New(t)

		var dst bool
		err := Decode(nil, tr, &dst)

		a.NoError(err)
		a.Equal(val, dst)
	})
}

func TestDecode_Int(t *testing.T) {
	t.Parallel()

	val := int(42)
	tr := Encode(nil, val)

	t.Run("Interface", func(t *testing.T) {
		a := assert.New(t)

		var dst any
		err := Decode(nil, tr, &dst)

		a.NoError(err)
		a.Equal(val, dst)
	})

	t.Run("Int", func(t *testing.T) {
		a := assert.New(t)

		var dst int
		err := Decode(nil, tr, &dst)

		a.NoError(err)
		a.Equal(val, dst)
	})
}

func TestDecode_Int8(t *testing.T) {
	t.Parallel()

	val := int8(42)
	tr := Encode(nil, val)

	t.Run("Interface", func(t *testing.T) {
		a := assert.New(t)

		var dst any
		err := Decode(nil, tr, &dst)

		a.NoError(err)
		a.Equal(val, dst)
	})

	t.Run("Int8", func(t *testing.T) {
		a := assert.New(t)

		var dst int8
		err := Decode(nil, tr, &dst)

		a.NoError(err)
		a.Equal(val, dst)
	})
}

func TestDecode_Int16(t *testing.T) {
	t.Parallel()

	val := int16(42)
	tr := Encode(nil, val)

	t.Run("Interface", func(t *testing.T) {
		a := assert.New(t)

		var dst any
		err := Decode(nil, tr, &dst)

		a.NoError(err)
		a.Equal(val, dst)
	})

	t.Run("Int16", func(t *testing.T) {
		a := assert.New(t)

		var dst int16
		err := Decode(nil, tr, &dst)

		a.NoError(err)
		a.Equal(val, dst)
	})
}

func TestDecode_Int32(t *testing.T) {
	t.Parallel()

	val := int32(42)
	tr := Encode(nil, val)

	t.Run("Interface", func(t *testing.T) {
		a := assert.New(t)

		var dst any
		err := Decode(nil, tr, &dst)

		a.NoError(err)
		a.Equal(val, dst)
	})

	t.Run("Int32", func(t *testing.T) {
		a := assert.New(t)

		var dst int32
		err := Decode(nil, tr, &dst)

		a.NoError(err)
		a.Equal(val, dst)
	})
}

func TestDecode_Int64(t *testing.T) {
	t.Parallel()

	val := int64(42)
	tr := Encode(nil, val)

	t.Run("Interface", func(t *testing.T) {
		a := assert.New(t)

		var dst any
		err := Decode(nil, tr, &dst)

		a.NoError(err)
		a.Equal(val, dst)
	})

	t.Run("Int64", func(t *testing.T) {
		a := assert.New(t)

		var dst int64
		err := Decode(nil, tr, &dst)

		a.NoError(err)
		a.Equal(val, dst)
	})
}

func TestDecode_Uint(t *testing.T) {
	t.Parallel()

	val := uint(42)
	tr := Encode(nil, val)

	t.Run("Interface", func(t *testing.T) {
		a := assert.New(t)

		var dst any
		err := Decode(nil, tr, &dst)

		a.NoError(err)
		a.Equal(val, dst)
	})

	t.Run("Uint", func(t *testing.T) {
		a := assert.New(t)

		var dst uint
		err := Decode(nil, tr, &dst)

		a.NoError(err)
		a.Equal(val, dst)
	})
}

func TestDecode_Uint8(t *testing.T) {
	t.Parallel()

	val := uint8(42)
	tr := Encode(nil, val)

	t.Run("Interface", func(t *testing.T) {
		a := assert.New(t)

		var dst any
		err := Decode(nil, tr, &dst)

		a.NoError(err)
		a.Equal(val, dst)
	})

	t.Run("Uint8", func(t *testing.T) {
		a := assert.New(t)

		var dst uint8
		err := Decode(nil, tr, &dst)

		a.NoError(err)
		a.Equal(val, dst)
	})
}

func TestDecode_Uint16(t *testing.T) {
	t.Parallel()

	val := uint16(42)
	tr := Encode(nil, val)

	t.Run("Interface", func(t *testing.T) {
		a := assert.New(t)

		var dst any
		err := Decode(nil, tr, &dst)

		a.NoError(err)
		a.Equal(val, dst)
	})

	t.Run("Uint16", func(t *testing.T) {
		a := assert.New(t)

		var dst uint16
		err := Decode(nil, tr, &dst)

		a.NoError(err)
		a.Equal(val, dst)
	})
}

func TestDecode_Uint32(t *testing.T) {
	t.Parallel()

	val := uint32(42)
	tr := Encode(nil, val)

	t.Run("Interface", func(t *testing.T) {
		a := assert.New(t)

		var dst any
		err := Decode(nil, tr, &dst)

		a.NoError(err)
		a.Equal(val, dst)
	})

	t.Run("Uint32", func(t *testing.T) {
		a := assert.New(t)

		var dst uint32
		err := Decode(nil, tr, &dst)

		a.NoError(err)
		a.Equal(val, dst)
	})
}

func TestDecode_Uint64(t *testing.T) {
	t.Parallel()

	val := uint64(42)
	tr := Encode(nil, val)

	t.Run("Interface", func(t *testing.T) {
		a := assert.New(t)

		var dst any
		err := Decode(nil, tr, &dst)

		a.NoError(err)
		a.Equal(val, dst)
	})

	t.Run("Uint64", func(t *testing.T) {
		a := assert.New(t)

		var dst uint64
		err := Decode(nil, tr, &dst)

		a.NoError(err)
		a.Equal(val, dst)
	})
}

func TestDecode_Float32(t *testing.T) {
	t.Parallel()

	val := float32(42)
	tr := Encode(nil, val)

	t.Run("Interface", func(t *testing.T) {
		a := assert.New(t)

		var dst any
		err := Decode(nil, tr, &dst)

		a.NoError(err)
		a.Equal(val, dst)
	})

	t.Run("Float32", func(t *testing.T) {
		a := assert.New(t)

		var dst float32
		err := Decode(nil, tr, &dst)

		a.NoError(err)
		a.Equal(val, dst)
	})
}

func TestDecode_Float64(t *testing.T) {
	t.Parallel()

	val := float64(42)
	tr := Encode(nil, val)

	t.Run("Interface", func(t *testing.T) {
		a := assert.New(t)

		var dst any
		err := Decode(nil, tr, &dst)

		a.NoError(err)
		a.Equal(val, dst)
	})

	t.Run("Float64", func(t *testing.T) {
		a := assert.New(t)

		var dst float64
		err := Decode(nil, tr, &dst)

		a.NoError(err)
		a.Equal(val, dst)
	})
}

func TestDecode_Struct(t *testing.T) {
	t.Parallel()

	type testStruct struct {
		A string
		B struct {
			C int
		}
	}

	val := testStruct{
		A: "test",
		B: struct{ C int }{C: 42},
	}
	tr := Encode(nil, val)

	t.Run("Interface", func(t *testing.T) {
		a := assert.New(t)

		var dst any
		err := Decode(nil, tr, &dst)

		a.NoError(err, "should not error")
		a.Equal(val, dst, "should be equal")
	})

	t.Run("Struct", func(t *testing.T) {
		a := assert.New(t)

		var dst testStruct
		err := Decode(nil, tr, &dst)

		a.NoError(err, "should not error")
		a.Equal(val, dst, "should be equal")
	})

	t.Run("Map", func(t *testing.T) {
		a := assert.New(t)

		var dst map[string]any
		err := Decode(nil, tr, &dst)

		exp := map[string]any{
			"A": val.A,
			"B": val.B,
		}

		a.NoError(err, "should not error")
		a.Equal(exp, dst, "should be equal")
	})
}

func TestDecode_Slice(t *testing.T) {
	t.Parallel()

	val := []string{"foo", "bar"}
	tr := Encode(nil, val)

	t.Run("Interface", func(t *testing.T) {
		a := assert.New(t)

		var dst any
		err := Decode(nil, tr, &dst)

		a.NoError(err, "should not error")
		a.Equal(val, dst, "should be equal")
	})

	t.Run("Slice", func(t *testing.T) {
		a := assert.New(t)

		var dst []string
		err := Decode(nil, tr, &dst)

		a.NoError(err, "should not error")
		a.Equal(val, dst, "should be equal")
	})
}

func TestDecode_Array(t *testing.T) {
	t.Parallel()

	t.Skip("TODO")

	val := [2]string{"foo", "bar"}
	tr := Encode(nil, val)

	t.Run("Interface", func(t *testing.T) {
		a := assert.New(t)

		var dst any
		err := Decode(nil, tr, &dst)

		a.NoError(err, "should not error")
		a.Equal(val, dst, "should be equal")
	})

	t.Run("Array", func(t *testing.T) {
		a := assert.New(t)

		var dst [2]string
		err := Decode(nil, tr, &dst)

		a.NoError(err, "should not error")
		a.Equal(val, dst, "should be equal")
	})
}

func TestDecode_Map(t *testing.T) {
	t.Parallel()

	val := map[string]int{"Foo": 42, "Bar": 1337}
	tr := Encode(nil, val)

	t.Run("Interface", func(t *testing.T) {
		a := assert.New(t)

		var dst any
		err := Decode(nil, tr, &dst)

		a.NoError(err, "should not error")
		a.Equal(val, dst, "should be equal")
	})

	t.Run("Map", func(t *testing.T) {
		a := assert.New(t)

		var dst map[string]int
		err := Decode(nil, tr, &dst)

		a.NoError(err, "should not error")
		a.Equal(val, dst, "should be equal")
	})

	t.Run("Struct", func(t *testing.T) {
		type testStruct struct {
			Foo int
			Bar int
		}

		a := assert.New(t)

		var dst testStruct
		err := Decode(nil, tr, &dst)

		exp := testStruct{
			Foo: val["Foo"],
			Bar: val["Bar"],
		}

		a.NoError(err, "should not error")
		a.Equal(exp, dst, "should be equal")
	})
}
