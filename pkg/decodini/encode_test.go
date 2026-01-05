package decodini

import (
	"reflect"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncode_Nil(t *testing.T) {
	a := assert.New(t)

	tr := Encode(nil, nil)

	a.NotNil(tr)
	a.Nil(tr.Parent())
	a.EqualValues(0, tr.NumChildren())

	a.True(tr.IsNil())
}

func TestEncode_String(t *testing.T) {
	a := assert.New(t)

	tr := Encode(nil, "decodini")

	a.NotNil(tr)
	a.EqualValues(0, tr.NumChildren())

	a.Equal(reflect.String, tr.Value().Kind())
	a.Equal("decodini", tr.Value().String())
}

func TestEncode_ByteSlice(t *testing.T) {
	a := assert.New(t)

	val := []byte("héllö 🌍")
	tr := Encode(nil, val)

	a.NotNil(tr)
	a.Equal(uint(len(val)), tr.NumChildren())

	child0 := tr.Child(0)
	a.NotNil(child0)
	a.Equal(reflect.Uint8, child0.Value().Kind())
	a.Equal(val[0], byte(child0.Value().Uint()))
}

func TestEncode_RuneSlice(t *testing.T) {
	a := assert.New(t)

	val := []rune("héllö 🌍")
	tr := Encode(nil, val)

	a.NotNil(tr)
	a.Equal(uint(len(val)), tr.NumChildren())

	child0 := tr.Child(0)
	a.NotNil(child0)
	a.Equal(reflect.Int32, child0.Value().Kind())
	a.Equal(val[0], rune(child0.Value().Int()))
}

func TestEncode_ShallowStruct(t *testing.T) {
	type testStruct struct {
		A string `decodini:"a"`
		B int
		C bool `decodini:"-"`
		d float64
	}

	a := assert.New(t)

	val := testStruct{
		A: "foo",
		B: 42,
		C: true,
		d: 420.0,
	}
	tr := Encode(nil, val)

	a.NotNil(tr)

	a.Nil(tr.Parent())
	a.Nil(tr.Name())

	a.Equal(uint(2), tr.NumChildren())

	{
		a.Nil(tr.Child("A"))

		childA := tr.Child("a")
		a.NotNil(childA)

		a.Equal(tr, childA.Parent())
		a.Equal("a", childA.Name())

		a.Equal(reflect.String, childA.Value().Kind())
		a.Equal(val.A, childA.Value().String())

		a.Equal(uint(0), childA.NumChildren())
	}

	{
		a.Nil(tr.Child("b"))

		childB := tr.Child("B")
		a.NotNil(childB)

		a.Equal(tr, childB.Parent())
		a.Equal("B", childB.Name())

		a.Equal(reflect.Int, childB.Value().Kind())
		a.Equal(val.B, int(childB.Value().Int()))

		a.Equal(uint(0), childB.NumChildren())
	}

	a.Nil(tr.Child("C"))

	a.Nil(tr.Child("d"))
}

func TestEncode_ShallowSlice(t *testing.T) {
	a := assert.New(t)

	val := []string{"foo", "bar"}
	tr := Encode(nil, val)

	a.NotNil(tr)

	a.Nil(tr.Parent())
	a.Nil(tr.Name())

	a.Equal(uint(len(val)), tr.NumChildren())

	a.Nil(tr.Child(-1))
	a.Nil(tr.Child(len(val)))

	{
		child0 := tr.Child(0)
		a.NotNil(child0)

		a.Equal(tr, child0.Parent())
		a.Equal(0, child0.Name())

		a.Equal(reflect.String, child0.Value().Kind())
		a.Equal(val[0], child0.Value().String())

		a.Equal(uint(0), child0.NumChildren())
	}

	{
		child1 := tr.Child(1)
		a.NotNil(child1)

		a.Equal(tr, child1.Parent())
		a.Equal(1, child1.Name())

		a.Equal(reflect.String, child1.Value().Kind())
		a.Equal(val[1], child1.Value().String())

		a.Equal(uint(0), child1.NumChildren())
	}
}

func TestEncode_ShallowMap(t *testing.T) {
	a := assert.New(t)

	val := map[string]string{
		"foo": "bar",
		"baz": "baz",
	}
	tr := Encode(nil, val)

	a.NotNil(tr)

	a.Nil(tr.Parent())
	a.Nil(tr.Name())

	a.Equal(uint(len(val)), tr.NumChildren())

	{
		childFoo := tr.Child("foo")
		a.NotNil(childFoo)

		a.Equal(tr, childFoo.Parent())
		a.Equal("foo", childFoo.Name())

		a.Equal(reflect.String, childFoo.Value().Kind())
		a.Equal(val["foo"], childFoo.Value().String())

		a.Equal(uint(0), childFoo.NumChildren())
	}

	{
		childBaz := tr.Child("baz")
		a.NotNil(childBaz)

		a.Equal(tr, childBaz.Parent())
		a.Equal("baz", childBaz.Name())

		a.Equal(reflect.String, childBaz.Value().Kind())
		a.Equal(val["baz"], childBaz.Value().String())

		a.Equal(uint(0), childBaz.NumChildren())
	}
}

func TestEncode_NilPointerAndInterface(t *testing.T) {
	a := assert.New(t)

	type S struct{ X int }

	var p *S = nil
	trPtr := Encode(nil, p)
	a.NotNil(trPtr)
	a.True(trPtr.IsNil())
	a.Equal(uint(0), trPtr.NumChildren())

	var i any = (*S)(nil)
	trIface := Encode(nil, i)
	a.NotNil(trIface)
	a.True(trIface.IsNil())
	a.Equal(uint(0), trIface.NumChildren())
}

func TestEncode_EmbeddedStruct_NumChildrenMatchesFlattened(t *testing.T) {
	type Embedded struct {
		A string `decodini:"a"`
		B int
	}
	type Outer struct{ Embedded }

	a := assert.New(t)

	tr := Encode(nil, Outer{Embedded: Embedded{A: "x", B: 1}})
	a.NotNil(tr)

	// NumChildren should equal the number of yielded children (flattened)
	children := slices.Collect(tr.Children())
	a.Equal(uint(len(children)), tr.NumChildren())
}

func TestEncode_OneEmbeddedStruct(t *testing.T) {
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

	val := Outer{
		Embedded: Embedded{
			A: "foo",
			B: 42,
			C: true,
		},
	}

	tr := Encode(nil, val)
	a.NotNil(tr)

	a.Nil(tr.Child("Embedded"))

	children := slices.Collect(tr.Children())
	a.Len(children, 2)
	a.Equal("a", children[0].Name())
	a.Equal("B", children[1].Name())

	{
		childA := tr.Child("a")
		a.NotNil(childA)
		a.Equal(reflect.String, childA.Value().Kind())
		a.Equal(val.A, childA.Value().String())
		a.Equal(uint(0), childA.NumChildren())
	}

	{
		childB := tr.Child("B")
		a.NotNil(childB)
		a.Equal(reflect.Int, childB.Value().Kind())
		a.Equal(val.B, int(childB.Value().Int()))
		a.Equal(uint(0), childB.NumChildren())
	}
}

func TestEncode_TwoEmbeddedStructs(t *testing.T) {
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

	val := Outer{
		EmbeddedA: EmbeddedA{A: "foo"},
		EmbeddedB: EmbeddedB{B: 7},
	}

	tr := Encode(nil, val)
	a.NotNil(tr)

	a.Nil(tr.Child("EmbeddedA"))
	a.Nil(tr.Child("EmbeddedB"))

	children := slices.Collect(tr.Children())
	a.Len(children, 2)
	a.Equal("a", children[0].Name())
	a.Equal("b", children[1].Name())

	{
		childA := tr.Child("a")
		a.NotNil(childA)
		a.Equal(reflect.String, childA.Value().Kind())
		a.Equal(val.A, childA.Value().String())
	}

	{
		childB := tr.Child("b")
		a.NotNil(childB)
		a.Equal(reflect.Int, childB.Value().Kind())
		a.Equal(val.B, int(childB.Value().Int()))
	}
}

func TestTree_DepthFirst(t *testing.T) {
	t.Run("Singleton", func(t *testing.T) {
		type testStruct struct {
			A string `decodini:"a"`
		}

		a := assert.New(t)

		val := testStruct{
			A: "foo",
		}

		tr := Encode(nil, val)
		a.NotNil(tr)

		df := slices.Collect(tr.DepthFirst())
		a.Len(df, 2)
		a.Equal(val, df[0].Value().Interface())
		a.Equal(val.A, df[1].Value().Interface())
	})

	t.Run("Shallow", func(t *testing.T) {
		type testStruct struct {
			A string `decodini:"a"`
			B int
			C bool `decodini:"-"` // ignored
		}

		a := assert.New(t)

		val := testStruct{
			A: "foo",
			B: 42,
			C: true,
		}

		tr := Encode(nil, val)
		a.NotNil(tr)

		df := slices.Collect(tr.DepthFirst())
		a.Len(df, 3)

		a.Equal(val, df[0].Value().Interface())
		a.Equal([]any(nil), df[0].Path())

		a.Equal(val.A, df[1].Value().Interface())
		a.Equal([]any{"a"}, df[1].Path())

		a.Equal(val.B, df[2].Value().Interface())
		a.Equal([]any{"B"}, df[2].Path())
	})

	t.Run("Nested", func(t *testing.T) {
		type (
			innerStruct struct {
				A string `decodini:"a"`
				B int
				C bool `decodini:"-"` // ignored
			}

			testStruct struct {
				Inner innerStruct
			}
		)

		a := assert.New(t)

		val := testStruct{
			Inner: innerStruct{
				A: "foo",
				B: 42,
				C: true,
			},
		}

		tr := Encode(nil, val)
		a.NotNil(tr)

		df := slices.Collect(tr.DepthFirst())
		a.Len(df, 4)

		a.Equal(val, df[0].Value().Interface())
		a.Equal([]any(nil), df[0].Path())

		a.Equal(val.Inner, df[1].Value().Interface())
		a.Equal([]any{"Inner"}, df[1].Path())

		a.Equal(val.Inner.A, df[2].Value().Interface())
		a.Equal([]any{"Inner", "a"}, df[2].Path())

		a.Equal(val.Inner.B, df[3].Value().Interface())
		a.Equal([]any{"Inner", "B"}, df[3].Path())
	})

	t.Run("Backtracking", func(t *testing.T) {
		type (
			innerStruct struct {
				A string `decodini:"a"`
				B int
				C bool `decodini:"-"` // ignored
			}

			testStruct struct {
				Inner innerStruct
				D     string
				E     int `decodini:"-"` // ignored
			}
		)

		a := assert.New(t)

		val := testStruct{
			Inner: innerStruct{
				A: "foo",
				B: 42,
				C: true,
			},
			D: "bar",
		}

		tr := Encode(nil, val)
		a.NotNil(tr)

		df := slices.Collect(tr.DepthFirst())
		a.Len(df, 5)

		a.Equal(val, df[0].Value().Interface())
		a.Equal([]any(nil), df[0].Path())

		a.Equal(val.Inner, df[1].Value().Interface())
		a.Equal([]any{"Inner"}, df[1].Path())

		a.Equal(val.Inner.A, df[2].Value().Interface())
		a.Equal([]any{"Inner", "a"}, df[2].Path())

		a.Equal(val.Inner.B, df[3].Value().Interface())
		a.Equal([]any{"Inner", "B"}, df[3].Path())

		a.Equal(val.D, df[4].Value().Interface())
		a.Equal([]any{"D"}, df[4].Path())
	})

	t.Run("Embedded", func(t *testing.T) {
		type (
			EmbeddedStruct struct {
				A string `decodini:"a"`
				B int
				C bool `decodini:"-"` // ignored
			}

			testStruct struct {
				EmbeddedStruct
			}
		)

		a := assert.New(t)

		val := testStruct{
			EmbeddedStruct: EmbeddedStruct{
				A: "foo",
				B: 42,
				C: true,
			},
		}

		tr := Encode(nil, val)
		a.NotNil(tr)

		df := slices.Collect(tr.DepthFirst())
		a.Len(df, 3)

		a.Equal(val, df[0].Value().Interface())
		a.Equal([]any(nil), df[0].Path())

		a.Equal(val.A, df[1].Value().Interface())
		a.Equal([]any{"a"}, df[1].Path())

		a.Equal(val.B, df[2].Value().Interface())
		a.Equal([]any{"B"}, df[2].Path())
	})
}

func TestTree_BreadthFirst(t *testing.T) {
	t.Run("Singleton", func(t *testing.T) {
		type testStruct struct {
			A string `decodini:"a"`
		}

		a := assert.New(t)

		val := testStruct{
			A: "foo",
		}

		tr := Encode(nil, val)
		a.NotNil(tr)

		bf := slices.Collect(tr.BreadthFirst())
		a.Len(bf, 2)
		a.Equal(val, bf[0].Value().Interface())
		a.Equal(val.A, bf[1].Value().Interface())
	})

	t.Run("Shallow", func(t *testing.T) {
		type testStruct struct {
			A string `decodini:"a"`
			B int
			C bool `decodini:"-"` // ignored
		}

		a := assert.New(t)

		val := testStruct{
			A: "foo",
			B: 42,
			C: true,
		}

		tr := Encode(nil, val)
		a.NotNil(tr)

		bf := slices.Collect(tr.BreadthFirst())
		a.Len(bf, 3)

		a.Equal(val, bf[0].Value().Interface())
		a.Equal([]any(nil), bf[0].Path())

		a.Equal(val.A, bf[1].Value().Interface())
		a.Equal([]any{"a"}, bf[1].Path())

		a.Equal(val.B, bf[2].Value().Interface())
		a.Equal([]any{"B"}, bf[2].Path())
	})

	t.Run("Nested", func(t *testing.T) {
		type (
			innerStruct struct {
				A string `decodini:"a"`
				B int
				C bool `decodini:"-"` // ignored
			}

			testStruct struct {
				Inner innerStruct
			}
		)

		a := assert.New(t)

		val := testStruct{
			Inner: innerStruct{
				A: "foo",
				B: 42,
				C: true,
			},
		}

		tr := Encode(nil, val)
		a.NotNil(tr)

		bf := slices.Collect(tr.BreadthFirst())
		a.Len(bf, 4)

		a.Equal(val, bf[0].Value().Interface())
		a.Equal([]any(nil), bf[0].Path())

		a.Equal(val.Inner, bf[1].Value().Interface())
		a.Equal([]any{"Inner"}, bf[1].Path())

		a.Equal(val.Inner.A, bf[2].Value().Interface())
		a.Equal([]any{"Inner", "a"}, bf[2].Path())

		a.Equal(val.Inner.B, bf[3].Value().Interface())
		a.Equal([]any{"Inner", "B"}, bf[3].Path())
	})

	t.Run("Backtracking", func(t *testing.T) {
		type (
			innerStruct struct {
				A string `decodini:"a"`
				B int
				C bool `decodini:"-"` // ignored
			}

			testStruct struct {
				Inner innerStruct
				D     string
				E     int `decodini:"-"` // ignored
			}
		)

		a := assert.New(t)

		val := testStruct{
			Inner: innerStruct{
				A: "foo",
				B: 42,
				C: true,
			},
			D: "bar",
		}

		tr := Encode(nil, val)
		a.NotNil(tr)

		bf := slices.Collect(tr.BreadthFirst())
		a.Len(bf, 5)

		a.Equal(val, bf[0].Value().Interface())
		a.Equal([]any(nil), bf[0].Path())

		a.Equal(val.Inner, bf[1].Value().Interface())
		a.Equal([]any{"Inner"}, bf[1].Path())

		a.Equal(val.D, bf[2].Value().Interface())
		a.Equal([]any{"D"}, bf[2].Path())

		a.Equal(val.Inner.A, bf[3].Value().Interface())
		a.Equal([]any{"Inner", "a"}, bf[3].Path())

		a.Equal(val.Inner.B, bf[4].Value().Interface())
		a.Equal([]any{"Inner", "B"}, bf[4].Path())
	})

	t.Run("Embedded", func(t *testing.T) {
		type (
			EmbeddedStruct struct {
				A string `decodini:"a"`
				B int
				C bool `decodini:"-"` // ignored
			}

			testStruct struct {
				EmbeddedStruct
			}
		)

		a := assert.New(t)

		val := testStruct{
			EmbeddedStruct: EmbeddedStruct{
				A: "foo",
				B: 42,
				C: true,
			},
		}

		tr := Encode(nil, val)
		a.NotNil(tr)

		bf := slices.Collect(tr.BreadthFirst())
		a.Len(bf, 3)

		a.Equal(val, bf[0].Value().Interface())
		a.Equal([]any(nil), bf[0].Path())

		a.Equal(val.A, bf[1].Value().Interface())
		a.Equal([]any{"a"}, bf[1].Path())

		a.Equal(val.B, bf[2].Value().Interface())
		a.Equal([]any{"B"}, bf[2].Path())
	})
}
