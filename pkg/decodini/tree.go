package decodini

import (
	"iter"
	"reflect"
)

type Tree struct {
	Path     []any
	Value    reflect.Value
	Children []*Tree

	structField reflect.StructField
}

func NewTree(path []any, value reflect.Value, children ...*Tree) *Tree {
	return &Tree{Path: path, Value: value, Children: children}
}

func NewRootTree(value reflect.Value, children ...*Tree) *Tree {
	return NewTree(nil, value, children...)
}

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

// IsLeaf returns whether the tree is a leaf node. A leaf node is a node
// that does not have any children.
func (t *Tree) IsLeaf() bool {
	return len(t.Children) == 0
}

// IsPrimitive returns whether the tree is a primitive node, i.e. int, bool,
// string, etc.
func (t *Tree) IsPrimitive() bool {
	switch t.Value.Kind() {
	case reflect.Array, reflect.Slice, reflect.Map, reflect.Struct:
		return false
	default:
		return len(t.Children) == 0
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

// Child returns the child tree that matches the given name exactly. If no
// child tree matches the given name, the returned value is nil.
func (t *Tree) Child(name any) *Tree {
	for _, child := range t.Children {
		if child.Name() == name {
			return child
		}
	}
	return nil
}

// DepthFirst returns a sequence of all tree nodes in the tree in depth-first
// order. The given tree is included in the sequence.
func (t *Tree) DepthFirst() iter.Seq[*Tree] {
	return func(yield func(*Tree) bool) {
		if !yield(t) {
			return
		}
		for _, child := range t.Children {
			if !child.depthFirst(yield) {
				return
			}
		}
	}
}

func (t *Tree) depthFirst(yield func(*Tree) bool) bool {
	if !yield(t) {
		return false
	}
	for _, child := range t.Children {
		if !child.depthFirst(yield) {
			return false
		}
	}
	return true
}
