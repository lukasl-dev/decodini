package bench

func generateLargeMap() map[int]int {
	m := make(map[int]int)
	for i := range 1_000_000 {
		m[i] = i
	}
	return m
}
