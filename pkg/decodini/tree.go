package decodini

import (
	"iter"
	"reflect"
)

type Tree struct {
	Name     any
	Value    reflect.Value
	Children []*Tree
}

func NewTree(name any, value reflect.Value, children ...*Tree) *Tree {
	return &Tree{Name: name, Value: value, Children: children}
}

func NewRootTree(value reflect.Value, children ...*Tree) *Tree {
	return NewTree(reflect.Value{}, value, children...)
}

// Root returns whether the tree is a root node.
func (t *Tree) Root() bool {
	return t.Name == nil
}

// Leaf returns whether the tree is a leaf node.
func (t *Tree) Leaf() bool {
	return len(t.Children) == 0
}

// Nil returns whether the tree's value is nil.
func (t *Tree) Nil() bool {
	return t.Value == reflect.Value{}
}

func (t *Tree) Child(name any) *Tree {
	for _, child := range t.Children {
		if child.Name == name {
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
			if !yield(child) {
				return
			}
			child.depthFirst()(yield)
		}
	}
}

// depthFirst returns a sequence of all tree nodes in the tree in depth-first
// order. The given tree is not included in the sequence.
func (t *Tree) depthFirst() iter.Seq[*Tree] {
	return func(yield func(*Tree) bool) {
		for _, child := range t.Children {
			child.depthFirst()(yield)
		}
	}
}
