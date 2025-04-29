package decodini

import (
	"iter"
	"reflect"
)

type Encoding struct {
	// StructTag is the name of the struct tag used to specify the name of a
	// field. The default value is "decodini".
	//
	// If the value of a struct field is "-", the field is ignored.
	StructTag string
}

// defaultEncoding is the default encoding used by Encode.
var defaultEncoding = Encoding{
	StructTag: "decodini",
}

// Encode encodes the given value into a tree based on the given encoding. If
// the encoding is nil, the default encoding is used.
func Encode(enc *Encoding, val any) *Tree {
	if enc == nil {
		enc = &defaultEncoding
	}
	return encode(enc, nil, reflect.ValueOf(val))
}

func encode(enc *Encoding, path []any, val reflect.Value) *Tree {
	if enc == nil {
		enc = &defaultEncoding
	}
	switch val.Kind() {
	case reflect.Ptr:
		return encode(enc, path, val.Elem())

	case reflect.Interface:
		if !val.IsNil() {
			return encode(enc, path, val.Elem())
		}
	}
	return newTree(enc, path, val)
}

type Tree struct {
	enc *Encoding

	Path  []any
	Value reflect.Value

	structField reflect.StructField
}

func newTree(enc *Encoding, path []any, value reflect.Value) *Tree {
	return &Tree{enc: enc, Path: path, Value: value}
}

// Name returns the last element of the path, or nil if the path is empty,
func (t *Tree) Name() any {
	if len(t.Path) == 0 {
		return nil
	}
	return t.Path[len(t.Path)-1]
}

// IsRoot returns whether the tree is a root node.
func (t *Tree) IsRoot() bool {
	return len(t.Path) == 0
}

// IsLeaf returns whether the tree is a leaf node, i.e. it has no children.
func (t *Tree) IsLeaf() bool {
	return t.NumChildren() == 0
}

// IsPrimitive returns whether the tree is a primitive node, i.e. int, bool,
// string, etc.
func (t *Tree) IsPrimitive() bool {
	switch t.Value.Kind() {
	case reflect.Array, reflect.Slice, reflect.Map, reflect.Struct:
		return false
	default:
		return true
	}
}

// IsNil returns whether the tree's value is nil.
func (t *Tree) IsNil() bool {
	return t.Value == reflect.Value{}
}

// IsStructField returns whether the tree represents a struct field.
func (t *Tree) IsStructField() bool {
	return t.structField.Name != ""
}

// StructField returns the struct field. If the tree does not represent a struct
// field (i.e. IsStructField() is false), it panics.
func (t *Tree) StructField() reflect.StructField {
	if !t.IsStructField() {
		panic("decodini: tree does not represent a struct field")
	}
	return t.structField
}

// NumChildren returns the number of children of this tree. Transitive children
// (i.e. grandchildren) are not counted.
func (t *Tree) NumChildren() int {
	switch t.Value.Kind() {
	case reflect.Struct:
		fields := 0
		for range t.structFields() {
			fields++
		}
		return fields

	case reflect.Slice, reflect.Array, reflect.Map:
		return t.Value.Len()

	default:
		return 0
	}
}

// Child returns the child tree that matches the given name exactly. If no
// child tree matches the given name, the returned value is nil.
func (t *Tree) Child(name any) *Tree {
	for _, child := range t.Children() {
		if child.Name() == name {
			return child
		}
	}
	return nil
}

// Children iterates over the children of this tree.
func (t *Tree) Children() iter.Seq2[any, *Tree] {
	switch t.Value.Kind() {
	case reflect.Struct:
		return t.structChildren()

	case reflect.Slice, reflect.Array:
		return t.sliceChildren()

	case reflect.Map:
		return t.mapChildren()

	default:
		return nil
	}
}

func (t *Tree) structChildren() iter.Seq2[any, *Tree] {
	val := t.Value
	typ := t.Value.Type()

	return func(yield func(any, *Tree) bool) {
		for i := range t.structFields() {
			sf, vf := typ.Field(i), val.Field(i)

			name := sf.Name
			if tag := sf.Tag.Get(t.enc.StructTag); tag != "" {
				name = tag
			}

			enc := encode(t.enc, append(t.Path, name), vf)
			enc.structField = sf
			if !yield(name, enc) {
				return
			}
		}
	}
}

func (t *Tree) sliceChildren() iter.Seq2[any, *Tree] {
	val := t.Value

	return func(yield func(any, *Tree) bool) {
		for i := range t.Value.Len() {
			enc := encode(t.enc, append(t.Path, i), val.Index(i))
			if !yield(i, enc) {
				return
			}
		}
	}
}

func (t *Tree) mapChildren() iter.Seq2[any, *Tree] {
	val := t.Value

	return func(yield func(any, *Tree) bool) {
		for _, key := range t.Value.MapKeys() {
			enc := encode(t.enc, append(t.Path, key.Interface()), val.MapIndex(key))
			if !yield(key.Interface(), enc) {
				return
			}
		}
	}
}

// structFields iterates over the struct fields of the tree, skipping unexported
// and ignored fields.
func (t *Tree) structFields() iter.Seq[int] {
	typ := t.Value.Type()

	return func(yield func(int) bool) {
		for i := range t.Value.NumField() {
			tf := typ.Field(i)
			if tf.IsExported() && tf.Tag.Get(t.enc.StructTag) != "-" {
				if !yield(i) {
					return
				}
			}
		}
	}
}
