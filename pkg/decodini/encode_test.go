package decodini

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncode(t *testing.T) {
	t.Parallel()

	t.Run("Nil", func(t *testing.T) {
		a := assert.New(t)

		tr := Encode(nil, nil)

		a.True(tr.IsPrimitive(), "should be leaf")
		a.True(tr.IsNil(), "should be nil")
	})

	t.Run("String", func(t *testing.T) {
		a := assert.New(t)

		val := "test"
		tr := Encode(nil, val)

		a.True(tr.IsPrimitive(), "should be leaf")
		a.Equal(val, tr.Value.String())
	})

	t.Run("Bool", func(t *testing.T) {
		a := assert.New(t)

		val := bool(true)
		tr := Encode(nil, val)

		fmt.Println("OIDAAAAAAAAAAAA", tr.Value.Kind())

		a.True(tr.IsPrimitive(), "should be leaf")
		a.Equal(val, tr.Value.Bool())
	})

	t.Run("Int", func(t *testing.T) {
		a := assert.New(t)

		val := int(42)
		tr := Encode(nil, val)

		a.True(tr.IsPrimitive(), "should be leaf")
		a.Equal(val, int(tr.Value.Int()))
	})

	t.Run("Int8", func(t *testing.T) {
		a := assert.New(t)

		val := int8(42)
		tr := Encode(nil, val)

		a.True(tr.IsPrimitive(), "should be leaf")
		a.Equal(val, int8(tr.Value.Int()))
	})

	t.Run("Int16", func(t *testing.T) {
		a := assert.New(t)

		val := int16(42)
		tr := Encode(nil, val)

		a.True(tr.IsPrimitive(), "should be leaf")
		a.Equal(val, int16(tr.Value.Int()))
	})

	t.Run("Int32", func(t *testing.T) {
		a := assert.New(t)

		val := int32(42)
		tr := Encode(nil, val)

		a.True(tr.IsPrimitive(), "should be leaf")
		a.Equal(val, int32(tr.Value.Int()))
	})

	t.Run("Int64", func(t *testing.T) {
		a := assert.New(t)

		val := int64(42)
		tr := Encode(nil, val)

		a.True(tr.IsPrimitive(), "should be leaf")
		a.Equal(val, int64(tr.Value.Int()))
	})

	t.Run("Uint", func(t *testing.T) {
		a := assert.New(t)

		val := uint(42)
		tr := Encode(nil, val)

		a.True(tr.IsPrimitive(), "should be leaf")
		a.Equal(val, uint(tr.Value.Uint()))
	})

	t.Run("Uint8", func(t *testing.T) {
		a := assert.New(t)

		val := uint8(42)
		tr := Encode(nil, val)

		a.True(tr.IsPrimitive(), "should be leaf")
		a.Equal(val, uint8(tr.Value.Uint()))
	})

	t.Run("Uint16", func(t *testing.T) {
		a := assert.New(t)

		val := uint16(42)
		tr := Encode(nil, val)

		a.True(tr.IsPrimitive(), "should be leaf")
		a.Equal(val, uint16(tr.Value.Uint()))
	})

	t.Run("Uint32", func(t *testing.T) {
		a := assert.New(t)

		val := uint32(42)
		tr := Encode(nil, val)

		a.True(tr.IsPrimitive(), "should be leaf")
		a.Equal(val, uint32(tr.Value.Uint()))
	})

	t.Run("Uint64", func(t *testing.T) {
		a := assert.New(t)

		val := uint64(42)
		tr := Encode(nil, val)

		a.True(tr.IsPrimitive(), "should be leaf")
		a.Equal(val, uint64(tr.Value.Uint()))
	})

	t.Run("Float32", func(t *testing.T) {
		a := assert.New(t)

		val := float32(42)
		tr := Encode(nil, val)

		a.True(tr.IsPrimitive(), "should be leaf")
		a.Equal(val, float32(tr.Value.Float()))
	})

	t.Run("Float64", func(t *testing.T) {
		a := assert.New(t)

		val := float64(42)
		tr := Encode(nil, val)

		a.True(tr.IsPrimitive(), "should be leaf")
		a.Equal(val, float64(tr.Value.Float()))
	})

	t.Run("Ptr", func(t *testing.T) {
		a := assert.New(t)

		val := "test"
		tr := Encode(nil, &val)

		a.True(tr.IsPrimitive(), "should be leaf")
		a.Equal(val, tr.Value.String())
	})

	t.Run("Struct", func(t *testing.T) {
		a := assert.New(t)

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

		a.False(tr.IsPrimitive(), "should not be leaf")
		a.Equal(val, tr.Value.Interface())

		a.Equal(2, tr.NumChildren())

		a.Equal(val.A, tr.Child("A").Value.String())
		a.True(tr.Child("A").IsStructField())
		a.NotZero(tr.Child("A").StructField())

		a.Equal(val.B, tr.Child("B").Value.Interface())
		a.True(tr.Child("B").IsStructField())
		a.NotZero(tr.Child("B").StructField())

		a.Equal(val.B.C, int(tr.Child("B").Child("C").Value.Int()))
		a.True(tr.Child("B").Child("C").IsStructField())
		a.NotZero(tr.Child("B").Child("C").StructField())
	})

	t.Run("StructTags", func(t *testing.T) {
		a := assert.New(t)

		type testStruct struct {
			A string `decodini:"-"`
			B int    `decodini:"B"`
			C int    `decodini:"-"`
		}

		val := testStruct{
			A: "test",
			B: 42,
			C: 1337,
		}
		tr := Encode(nil, val)

		a.False(tr.IsPrimitive(), "should not be leaf")

		a.Equal(1, tr.NumChildren())
		a.Nil(tr.Child("A"))
		a.Equal(val.B, int(tr.Child("B").Value.Int()))
		a.Nil(tr.Child("C"))
	})

	t.Run("Slice", func(t *testing.T) {
		a := assert.New(t)

		val := []string{"foo", "bar"}
		tr := Encode(nil, val)

		a.False(tr.IsPrimitive(), "should not be leaf")
		a.Equal(val, tr.Value.Interface())

		a.Equal(2, tr.NumChildren())
		a.Equal(val[0], tr.Child(0).Value.String())
		a.Equal(val[1], tr.Child(1).Value.String())
	})

	t.Run("Array", func(t *testing.T) {
		a := assert.New(t)

		val := [2]string{"foo", "bar"}
		tr := Encode(nil, val)

		a.False(tr.IsPrimitive(), "should not be leaf")
		a.Equal(val, tr.Value.Interface())

		a.Equal(2, tr.NumChildren())
		a.Equal(val[0], tr.Child(0).Value.String())
		a.Equal(val[1], tr.Child(1).Value.String())
	})

	t.Run("Map", func(t *testing.T) {
		a := assert.New(t)

		val := map[string]int{"foo": 42, "bar": 1337}
		tr := Encode(nil, val)

		a.False(tr.IsPrimitive(), "should not be leaf")
		a.Equal(val, tr.Value.Interface())

		a.Equal(2, tr.NumChildren())
		a.Equal(val["foo"], int(tr.Child("foo").Value.Int()))
		a.Equal(val["bar"], int(tr.Child("bar").Value.Int()))
	})

	t.Run("Map/Empty", func(t *testing.T) {
		a := assert.New(t)

		val := map[string]int{}
		tr := Encode(nil, val)

		a.False(tr.IsPrimitive(), "should not be leaf")
		a.Equal(val, tr.Value.Interface())

		a.Equal(0, tr.NumChildren())
	})
}
