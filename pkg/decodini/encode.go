package decodini

import (
	"iter"
	"reflect"
)

type Encoding struct {
	StructTag string
}

var defaultEncoding = Encoding{
	StructTag: "decodini",
}

// Encode encodes the given value into a (lazy) Tree.
func Encode(enc *Encoding, val any) *Tree {
	if enc == nil {
		enc = &defaultEncoding
	}

	rVal, isVal := val.(reflect.Value)
	if !isVal {
		rVal = reflect.ValueOf(val)
	}

	tr := encode(enc, nil, nil, rVal)
	tr.isNil = val == nil
	return tr
}

func encode(enc *Encoding, parent *Tree, name any, val reflect.Value) *Tree {
	switch val.Kind() {
	case reflect.Ptr:
		return encode(enc, parent, name, val.Elem())
	case reflect.Interface:
		if !val.IsNil() {
			return encode(enc, parent, name, val.Elem())
		}
	}
	return &Tree{enc: enc, name: name, parent: parent, val: val}
}

type Tree struct {
	enc    *Encoding
	name   any
	parent *Tree
	val    reflect.Value

	isNil bool
}

// Name returns the name of this node in the parent node. If this node is root,
// null is returned.
func (t *Tree) Name() any {
	return t.name
}

// Parent returns the parent of this node. If this is a root node, nil is
// returned.
func (t *Tree) Parent() *Tree {
	return t.parent
}

// Value returns the underlying reflect.Value of this node.
func (t *Tree) Value() reflect.Value {
	return t.val
}

func (t *Tree) IsPrimitive() bool {
	return isPrimitive(t.val.Kind())
}

// Path returns the path from the root to this node. The first element is the
// name of
func (t *Tree) Path() (path []any) {
	if t.parent == nil {
		return nil
	}

	curr := t
	for curr != nil {
		path = append(path, curr.name)
		curr = curr.parent
	}
	return path
}

// IsNil returns true if this node's value is nil.
func (t *Tree) IsNil() bool {
	return t.isNil
}

// DepthFirst returns a sequence of the tree nodes in depth-first order.
func (t *Tree) DepthFirst() iter.Seq[*Tree] {
	return func(yield func(*Tree) bool) {
		if !t.yieldDFS(yield) {
			return
		}
	}
}

func (t *Tree) yieldDFS(yield func(*Tree) bool) bool {
	if !yield(t) {
		return false
	}
	for child := range t.Children() {
		if !child.yieldDFS(yield) {
			return false
		}
	}
	return true
}

// BreadthFirst returns a sequence of the tree nodes in breadth-first order.
func (t *Tree) BreadthFirst() iter.Seq[*Tree] {
	queue := []*Tree{t}

	return func(yield func(*Tree) bool) {
		for len(queue) > 0 {
			frontier := queue[0]
			if !yield(frontier) {
				return
			}

			children := make([]*Tree, 0, frontier.NumChildren())
			for child := range frontier.Children() {
				children = append(children, child)
			}

			queue = append(queue[1:], children...)
		}
	}
}

// NumChildren returns the number of children of this node.
func (t *Tree) NumChildren() uint {
	switch t.val.Kind() {
	case reflect.Struct:
		typ := t.val.Type()
		fields := uint(0)
		for i := range t.val.NumField() {
			field := typ.Field(i)
			if includeStructField(t.enc.StructTag, field) {
				fields++
			}
		}
		return fields

	case reflect.Map, reflect.Slice, reflect.Array:
		return uint(t.val.Len())

	default:
		return 0
	}
}

// Child returns the child of this node with the given name.
func (t *Tree) Child(name any) *Tree {
	switch t.val.Kind() {
	case reflect.Struct:
		nameStr, ok := name.(string)
		if !ok {
			return nil
		}
		_, vf := structFieldByName(t.enc.StructTag, t.val, nameStr)
		if !vf.IsValid() {
			return nil
		}
		tr := encode(t.enc, t, name, vf)
		return tr

	case reflect.Slice, reflect.Array:
		nameInt, ok := name.(int)
		if !ok || nameInt < 0 || nameInt >= t.val.Len() {
			return nil
		}
		tr := encode(t.enc, t, nameInt, t.val.Index(nameInt))
		return tr

	case reflect.Map:
		nameVal := reflect.ValueOf(name)
		for _, key := range t.val.MapKeys() {
			if key.Equal(nameVal) {
				return encode(t.enc, t, name, t.val.MapIndex(key))
			}
		}
		return nil

	default:
		return nil
	}
}

// Children returns a sequence of the children of this node, preserving their
// order.
func (t *Tree) Children() iter.Seq[*Tree] {
	switch t.val.Kind() {
	case reflect.Struct:
		typ := t.val.Type()
		return func(yield func(*Tree) bool) {
			for i := range t.val.NumField() {
				sf := typ.Field(i)
				if !includeStructField(t.enc.StructTag, sf) {
					continue
				}
				vf := t.val.Field(i)

				name := sf.Name
				tagName, hasTag := sf.Tag.Lookup(t.enc.StructTag)
				if hasTag {
					name = tagName
				}

				tr := encode(t.enc, t, name, vf)
				if !yield(tr) {
					return
				}
			}
		}

	case reflect.Slice, reflect.Array:
		return func(yield func(*Tree) bool) {
			for i := range t.val.Len() {
				tr := encode(t.enc, t, i, t.val.Index(i))
				if !yield(tr) {
					return
				}
			}
		}

	case reflect.Map:
		return func(yield func(*Tree) bool) {
			for _, key := range t.val.MapKeys() {
				tr := encode(t.enc, t, key.Interface(), t.val.MapIndex(key))
				if !yield(tr) {
					return
				}
			}
		}

	default:
		return nil
	}
}
