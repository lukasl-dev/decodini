package bench

import (
	"testing"

	"github.com/lukasl-dev/decodini/pkg/decodini"
	"github.com/mitchellh/mapstructure"
)

func BenchmarkLargeMap_to_Map(b *testing.B) {
	m := generateLargeMap()
	size := float64(len(m))

	b.Run("Decodini", func(b *testing.B) {
		for b.Loop() {
			b.ReportMetric(size, "len/op")

			var res map[int]int
			_ = decodini.TransmuteInto(nil, m, &res)
		}
	})

	b.Run("Mapstructure", func(b *testing.B) {
		for b.Loop() {
			b.ReportMetric(size, "len/op")

			var res map[int]int
			_ = mapstructure.Decode(m, &res)
		}
	})
}
