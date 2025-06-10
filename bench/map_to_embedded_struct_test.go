package bench

import (
	"testing"

	"github.com/lukasl-dev/decodini/pkg/decodini"
	"github.com/mitchellh/mapstructure"
)

func BenchmarkEmbeddedStruct(b *testing.B) {
	type embeddedStruct struct {
		Inner struct {
			A int
			B string
		}
		C string
	}

	m := map[string]any{
		"A": 1,
		"B": "foo",
		"C": "bar",
	}

	b.Run("Decodini", func(b *testing.B) {
		for b.Loop() {
			var res embeddedStruct
			_ = decodini.TransmuteInto(nil, m, &res)
		}
	})

	b.Run("Mapstructure", func(b *testing.B) {
		for b.Loop() {
			var res embeddedStruct
			_ = mapstructure.Decode(m, &res)
		}
	})
}
