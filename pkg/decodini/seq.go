package decodini

import "iter"

// SeqToSlice converts a sequence to a slice.
func SeqToSlice[T any](seq iter.Seq[T], maximum uint) []T {
	var slice []T
	for item := range seq {
		if maximum != 0 && len(slice) >= int(maximum) {
			break
		}
		slice = append(slice, item)
	}
	return slice
}

// SliceToSeq converts a slice to a sequence.
func SliceToSeq[T any](slice []T) iter.Seq[T] {
	return func(yield func(T) bool) {
		for _, item := range slice {
			if !yield(item) {
				return
			}
		}
	}
}
