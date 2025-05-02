package decodini

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncode_Nil(t *testing.T) {
	a := assert.New(t)

	tr := Encode(nil, nil)

	a.NotNil(tr)
	a.True(tr.IsRoot())
	a.True(tr.IsLeaf())
	a.EqualValues(0, tr.NumChildren())

	a.True(tr.IsNil())
}

func TestEncode_String(t *testing.T) {
	a := assert.New(t)

	tr := Encode(nil, "decodini")

	a.NotNil(tr)
	a.True(tr.IsRoot())
	a.True(tr.IsLeaf())
	a.EqualValues(0, tr.NumChildren())

	a.Equal(reflect.String, tr.Value().Kind())
	a.Equal("decodini", tr.Value().String())
}

func TestEncode_ShallowStruct(t *testing.T) {
	type testStruct struct {
		A string `decodini:"a"`
		B int
		C bool `decodini:"-"`
	}

	a := assert.New(t)

	val := testStruct{
		A: "foo",
		B: 42,
		C: true,
	}
	tr := Encode(nil, val)

	a.NotNil(tr)
	a.Nil(tr.Name())
	a.True(tr.IsRoot())
	a.Equal(uint(2), tr.NumChildren())

	{
		a.Nil(tr.Child("A"))
		childA := tr.Child("a")
		a.NotNil(childA)
		a.Equal("a", childA.Name())
		a.Equal(reflect.String, childA.Value().Kind())
		a.Equal(val.A, childA.Value().String())
		a.Equal(tr, childA.Parent())
		a.False(childA.IsRoot())
		a.Equal(uint(0), childA.NumChildren())
	}

	{
		a.Nil(tr.Child("b"))
		childB := tr.Child("B")
		a.NotNil(childB)
		a.Equal("B", childB.Name())
		a.Equal(reflect.Int, childB.Value().Kind())
		a.Equal(val.B, int(childB.Value().Int()))
		a.Equal(tr, childB.Parent())
		a.False(childB.IsRoot())
		a.Equal(uint(0), childB.NumChildren())
	}

	a.Nil(tr.Child("C"))
}

func TestEncode_ShallowSlice(t *testing.T) {
	a := assert.New(t)

	val := []string{"foo", "bar"}
	tr := Encode(nil, val)

	a.NotNil(tr)
	a.Nil(tr.Name())
	a.True(tr.IsRoot())
	a.Equal(uint(len(val)), tr.NumChildren())

	a.Nil(tr.Child(-1))
	a.Nil(tr.Child(len(val)))

	{
		child0 := tr.Child(0)
		a.NotNil(child0)
		a.Equal(0, child0.Name())
		a.Equal(reflect.String, child0.Value().Kind())
		a.Equal(val[0], child0.Value().String())
		a.Equal(tr, child0.Parent())
		a.False(child0.IsRoot())
		a.Equal(uint(0), child0.NumChildren())
	}

	{
		child1 := tr.Child(1)
		a.NotNil(child1)
		a.Equal(1, child1.Name())
		a.Equal(reflect.String, child1.Value().Kind())
		a.Equal(val[1], child1.Value().String())
		a.Equal(tr, child1.Parent())
		a.False(child1.IsRoot())
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
	a.True(tr.IsRoot())
	a.Equal(uint(len(val)), tr.NumChildren())

	{
		childFoo := tr.Child("foo")
		a.NotNil(childFoo)
		a.Equal("foo", childFoo.Name())
		a.Equal(reflect.String, childFoo.Value().Kind())
		a.Equal(val["foo"], childFoo.Value().String())
		a.Equal(tr, childFoo.Parent())
		a.False(childFoo.IsRoot())
		a.Equal(uint(0), childFoo.NumChildren())
	}

	{
		childBaz := tr.Child("baz")
		a.NotNil(childBaz)
		a.Equal("baz", childBaz.Name())
		a.Equal(reflect.String, childBaz.Value().Kind())
		a.Equal(val["baz"], childBaz.Value().String())
		a.Equal(tr, childBaz.Parent())
		a.False(childBaz.IsRoot())
		a.Equal(uint(0), childBaz.NumChildren())
	}
}
