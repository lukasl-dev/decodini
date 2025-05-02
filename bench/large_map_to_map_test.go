package bench

import (
	"testing"

	"github.com/lukasl-dev/decodini/pkg/decodini"
	"github.com/mitchellh/mapstructure"
)

func BenchmarkLargeMap_to_Map(b *testing.B) {
	m := generateLargeMap()

	b.Run("Decodini", func(b *testing.B) {
		for b.Loop() {
			var res map[int]int
			_ = decodini.TransmuteInto(nil, m, &res)
		}
	})

	b.Run("Mapstructure", func(b *testing.B) {
		for b.Loop() {
			var res map[int]int
			_ = mapstructure.Decode(m, &res)
		}
	})
}
