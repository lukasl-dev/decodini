package bench

import (
	"testing"

	"github.com/lukasl-dev/decodini/pkg/decodini"
	"github.com/mitchellh/mapstructure"
)

type emptyStruct struct{}

func BenchmarkLargeMap_to_EmptyStruct(b *testing.B) {
	m := generateLargeMap()
	size := float64(len(m))

	b.Run("Decodini", func(b *testing.B) {
		for b.Loop() {
			b.ReportMetric(size, "len/op")

			var res emptyStruct
			_ = decodini.TransmuteInto(nil, m, &res)
		}
	})

	b.Run("Mapstructure", func(b *testing.B) {
		for b.Loop() {
			b.ReportMetric(size, "len/op")

			var res emptyStruct
			_ = mapstructure.Decode(m, &res)
		}
	})
}
