package decodini

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDepthFirst(t *testing.T) {
	t.Run("Leaf", func(t *testing.T) {
		a := assert.New(t)

		tree := NewRootTree(reflect.ValueOf(0))

		seq := SeqToSlice(tree.DepthFirst(), 0)

		a.Len(seq, 1)
		a.Equal(tree, seq[0])
	})

	t.Run("Node", func(t *testing.T) {
		a := assert.New(t)

		tree := &Tree{
			Children: []*Tree{
				NewTree([]any{reflect.ValueOf("child1")}, reflect.ValueOf(0)),
				NewTree([]any{reflect.ValueOf("child2")}, reflect.ValueOf(0)),
			},
		}

		seq := SeqToSlice(tree.DepthFirst(), 0)

		a.Len(seq, 3)
		a.Equal(tree, seq[0])
		a.Equal(tree.Children[0], seq[1])
		a.Equal(tree.Children[1], seq[2])
	})
}
