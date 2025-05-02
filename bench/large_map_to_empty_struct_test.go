package bench

import (
	"testing"

	"github.com/lukasl-dev/decodini/pkg/decodini"
	"github.com/mitchellh/mapstructure"
)

type emptyStruct struct{}

func BenchmarkDecodini_LargeMap_to_EmptyStruct(b *testing.B) {
	m := generateLargeMap()

	b.ResetTimer()
	for b.Loop() {
		var res emptyStruct
		_ = decodini.TransmuteInto(nil, m, &res)
	}
}

func BenchmarkMapstructure_LargeMap_to_EmptyStruct(b *testing.B) {
	m := generateLargeMap()

	b.ResetTimer()
	for b.Loop() {
		var res emptyStruct
		_ = mapstructure.Decode(m, &res)
	}
}
