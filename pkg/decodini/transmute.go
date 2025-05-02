package decodini

type Transmutation struct {
	Encoding *Encoding
	Decoding *Decoding
}

// TransmuteInto encodes the given `from` value into a tree and decodes the tree
// directly into the given `to` value.
func TransmuteInto(tr *Transmutation, from, to any) error {
	panic("not implemented yet")
}

// Transmute encodes the given `from` value into a tree, and decodes the tree
// into a variable of type `T`.
func Transmute[T any](tr *Transmutation, from any) (T, error) {
	var to T
	err := TransmuteInto(tr, from, &to)
	return to, err
}
