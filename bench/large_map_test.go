package bench

import (
	"os"
	"strconv"
)

func generateLargeMap() map[int]int {
	m := make(map[int]int)
	for i := range largeMapSize() {
		m[i] = i
	}
	return m
}

func largeMapSize() int {
	raw := os.Getenv("DECODINI_B_LARGE_MAP_SIZE")

	size, err := strconv.Atoi(raw)
	if err != nil {
		return 1_000_000
	}
	return size
}
