package bench

import (
	"testing"

	"github.com/lukasl-dev/decodini/pkg/decodini"
	"github.com/mitchellh/mapstructure"
)

func BenchmarkDecodini(b *testing.B) {
	m := generateMap()

	b.ResetTimer()
	for b.Loop() {
		var res map[int]int
		_ = decodini.TransmuteInto(nil, m, &res)
	}
}

func BenchmarkMapstructure(b *testing.B) {
	m := generateMap()

	b.ResetTimer()
	for b.Loop() {
		var res map[int]int
		_ = mapstructure.Decode(m, &res)
	}
}

func generateMap() map[int]int {
	m := make(map[int]int)
	for i := range 1_000_000 {
		m[i] = i
	}
	return m
}
