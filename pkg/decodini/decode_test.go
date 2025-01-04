package decodini

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecode(t *testing.T) {
	t.Parallel()

	t.Run("Nil", func(t *testing.T) {
		a := assert.New(t)

		tr := DefaultEncode(nil)

		var dst any
		err := DefaultDecode(tr, &dst)

		a.NoError(err)
		a.Nil(dst)
	})

	t.Run("String", func(t *testing.T) {
		val := "test"
		tr := DefaultEncode(val)

		t.Run("Interface", func(t *testing.T) {
			a := assert.New(t)

			var dst any
			err := DefaultDecode(tr, &dst)

			a.NoError(err)
			a.Equal(val, dst)
		})

		t.Run("String", func(t *testing.T) {
			a := assert.New(t)

			var dst string
			err := DefaultDecode(tr, &dst)

			a.NoError(err)
			a.Equal(val, dst)
		})
	})

	t.Run("Bool", func(t *testing.T) {
		val := bool(true)
		tr := DefaultEncode(val)

		t.Run("Interface", func(t *testing.T) {
			a := assert.New(t)

			var dst any
			err := DefaultDecode(tr, &dst)

			a.NoError(err)
			a.Equal(val, dst)
		})

		t.Run("Bool", func(t *testing.T) {
			a := assert.New(t)

			var dst bool
			err := DefaultDecode(tr, &dst)

			a.NoError(err)
			a.Equal(val, dst)
		})
	})

	t.Run("Int", func(t *testing.T) {
		val := int(42)
		tr := DefaultEncode(val)

		t.Run("Interface", func(t *testing.T) {
			a := assert.New(t)

			var dst any
			err := DefaultDecode(tr, &dst)

			a.NoError(err)
			a.Equal(val, dst)
		})

		t.Run("Int", func(t *testing.T) {
			a := assert.New(t)

			var dst int
			err := DefaultDecode(tr, &dst)

			a.NoError(err)
			a.Equal(val, dst)
		})
	})

	t.Run("Int8", func(t *testing.T) {
		val := int8(42)
		tr := DefaultEncode(val)

		t.Run("Interface", func(t *testing.T) {
			a := assert.New(t)

			var dst any
			err := DefaultDecode(tr, &dst)

			a.NoError(err)
			a.Equal(val, dst)
		})

		t.Run("Int8", func(t *testing.T) {
			a := assert.New(t)

			var dst int8
			err := DefaultDecode(tr, &dst)

			a.NoError(err)
			a.Equal(val, dst)
		})
	})

	t.Run("Int16", func(t *testing.T) {
		val := int16(42)
		tr := DefaultEncode(val)

		t.Run("Interface", func(t *testing.T) {
			a := assert.New(t)

			var dst any
			err := DefaultDecode(tr, &dst)

			a.NoError(err)
			a.Equal(val, dst)
		})

		t.Run("Int16", func(t *testing.T) {
			a := assert.New(t)

			var dst int16
			err := DefaultDecode(tr, &dst)

			a.NoError(err)
			a.Equal(val, dst)
		})
	})

	t.Run("Int32", func(t *testing.T) {
		val := int32(42)
		tr := DefaultEncode(val)

		t.Run("Interface", func(t *testing.T) {
			a := assert.New(t)

			var dst any
			err := DefaultDecode(tr, &dst)

			a.NoError(err)
			a.Equal(val, dst)
		})

		t.Run("Int32", func(t *testing.T) {
			a := assert.New(t)

			var dst int32
			err := DefaultDecode(tr, &dst)

			a.NoError(err)
			a.Equal(val, dst)
		})
	})

	t.Run("Int64", func(t *testing.T) {
		val := int64(42)
		tr := DefaultEncode(val)

		t.Run("Interface", func(t *testing.T) {
			a := assert.New(t)

			var dst any
			err := DefaultDecode(tr, &dst)

			a.NoError(err)
			a.Equal(val, dst)
		})

		t.Run("Int64", func(t *testing.T) {
			a := assert.New(t)

			var dst int64
			err := DefaultDecode(tr, &dst)

			a.NoError(err)
			a.Equal(val, dst)
		})
	})

	t.Run("Uint", func(t *testing.T) {
		val := uint(42)
		tr := DefaultEncode(val)

		t.Run("Interface", func(t *testing.T) {
			a := assert.New(t)

			var dst any
			err := DefaultDecode(tr, &dst)

			a.NoError(err)
			a.Equal(val, dst)
		})

		t.Run("Uint", func(t *testing.T) {
			a := assert.New(t)

			var dst uint
			err := DefaultDecode(tr, &dst)

			a.NoError(err)
			a.Equal(val, dst)
		})
	})

	t.Run("Uint8", func(t *testing.T) {
		val := uint8(42)
		tr := DefaultEncode(val)

		t.Run("Interface", func(t *testing.T) {
			a := assert.New(t)

			var dst any
			err := DefaultDecode(tr, &dst)

			a.NoError(err)
			a.Equal(val, dst)
		})

		t.Run("Uint8", func(t *testing.T) {
			a := assert.New(t)

			var dst uint8
			err := DefaultDecode(tr, &dst)

			a.NoError(err)
			a.Equal(val, dst)
		})
	})

	t.Run("Uint16", func(t *testing.T) {
		val := uint16(42)
		tr := DefaultEncode(val)

		t.Run("Interface", func(t *testing.T) {
			a := assert.New(t)

			var dst any
			err := DefaultDecode(tr, &dst)

			a.NoError(err)
			a.Equal(val, dst)
		})

		t.Run("Uint16", func(t *testing.T) {
			a := assert.New(t)

			var dst uint16
			err := DefaultDecode(tr, &dst)

			a.NoError(err)
			a.Equal(val, dst)
		})
	})

	t.Run("Uint32", func(t *testing.T) {
		val := uint32(42)
		tr := DefaultEncode(val)

		t.Run("Interface", func(t *testing.T) {
			a := assert.New(t)

			var dst any
			err := DefaultDecode(tr, &dst)

			a.NoError(err)
			a.Equal(val, dst)
		})

		t.Run("Uint32", func(t *testing.T) {
			a := assert.New(t)

			var dst uint32
			err := DefaultDecode(tr, &dst)

			a.NoError(err)
			a.Equal(val, dst)
		})
	})

	t.Run("Uint64", func(t *testing.T) {
		val := uint64(42)
		tr := DefaultEncode(val)

		t.Run("Interface", func(t *testing.T) {
			a := assert.New(t)

			var dst any
			err := DefaultDecode(tr, &dst)

			a.NoError(err)
			a.Equal(val, dst)
		})

		t.Run("Uint64", func(t *testing.T) {
			a := assert.New(t)

			var dst uint64
			err := DefaultDecode(tr, &dst)

			a.NoError(err)
			a.Equal(val, dst)
		})
	})

	t.Run("Float32", func(t *testing.T) {
		val := float32(42)
		tr := DefaultEncode(val)

		t.Run("Interface", func(t *testing.T) {
			a := assert.New(t)

			var dst any
			err := DefaultDecode(tr, &dst)

			a.NoError(err)
			a.Equal(val, dst)
		})

		t.Run("Float32", func(t *testing.T) {
			a := assert.New(t)

			var dst float32
			err := DefaultDecode(tr, &dst)

			a.NoError(err)
			a.Equal(val, dst)
		})
	})

	t.Run("Float64", func(t *testing.T) {
		val := float64(42)
		tr := DefaultEncode(val)

		t.Run("Interface", func(t *testing.T) {
			a := assert.New(t)

			var dst any
			err := DefaultDecode(tr, &dst)

			a.NoError(err)
			a.Equal(val, dst)
		})

		t.Run("Float64", func(t *testing.T) {
			a := assert.New(t)

			var dst float64
			err := DefaultDecode(tr, &dst)

			a.NoError(err)
			a.Equal(val, dst)
		})
	})

	t.Run("Struct", func(t *testing.T) {
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
		tr := DefaultEncode(val)

		t.Run("Interface", func(t *testing.T) {
			a := assert.New(t)

			var dst any
			err := DefaultDecode(tr, &dst)

			a.NoError(err, "should not error")
			a.Equal(val, dst, "should be equal")
		})

		t.Run("Struct", func(t *testing.T) {
			a := assert.New(t)

			var dst testStruct
			err := DefaultDecode(tr, &dst)

			a.NoError(err, "should not error")
			a.Equal(val, dst, "should be equal")
		})

		t.Run("Map", func(t *testing.T) {
			t.Skip("TODO")

			a := assert.New(t)

			var dst map[string]any
			err := DefaultDecode(tr, &dst)

			a.NoError(err, "should not error")
			a.Equal(val, dst, "should be equal")
		})
	})

	t.Run("Slice", func(t *testing.T) {
		val := []string{"foo", "bar"}
		tr := DefaultEncode(val)

		t.Run("Interface", func(t *testing.T) {
			a := assert.New(t)

			var dst any
			err := DefaultDecode(tr, &dst)

			a.NoError(err, "should not error")
			a.Equal(val, dst, "should be equal")
		})

		t.Run("Slice", func(t *testing.T) {
			a := assert.New(t)

			var dst []string
			err := DefaultDecode(tr, &dst)

			a.NoError(err, "should not error")
			a.Equal(val, dst, "should be equal")
		})
	})

	t.Run("Array", func(t *testing.T) {
		t.Skip("TODO")

		val := [2]string{"foo", "bar"}
		tr := DefaultEncode(val)

		t.Run("Interface", func(t *testing.T) {
			a := assert.New(t)

			var dst any
			err := DefaultDecode(tr, &dst)

			a.NoError(err, "should not error")
			a.Equal(val, dst, "should be equal")
		})

		t.Run("Array", func(t *testing.T) {
			a := assert.New(t)

			var dst [2]string
			err := DefaultDecode(tr, &dst)

			a.NoError(err, "should not error")
			a.Equal(val, dst, "should be equal")
		})
	})

	t.Run("Map", func(t *testing.T) {
		val := map[string]int{"Foo": 42, "Bar": 1337}
		tr := DefaultEncode(val)

		t.Run("Interface", func(t *testing.T) {
			a := assert.New(t)

			var dst any
			err := DefaultDecode(tr, &dst)

			a.NoError(err, "should not error")
			a.Equal(val, dst, "should be equal")
		})

		t.Run("Map", func(t *testing.T) {
			a := assert.New(t)

			var dst map[string]int
			err := DefaultDecode(tr, &dst)

			a.NoError(err, "should not error")
			a.Equal(val, dst, "should be equal")
		})

		t.Run("Struct", func(t *testing.T) {
			t.Skip("TODO")

			type testStruct struct {
				Foo int
				Bar int
			}

			a := assert.New(t)

			var dst testStruct
			err := DefaultDecode(tr, &dst)

			a.NoError(err, "should not error")
			a.Equal(val, dst, "should be equal")
		})
	})
}
