package decodini

type Transmutation struct {
	Encoding *Encoding
	Decoding *Decoding
}

// TransmuteInto encodes the given `from` value into a tree and decodes the tree
// directly into the given `to` value.
func TransmuteInto(tr *Transmutation, from, to any) error {
	if tr == nil {
		tr = new(Transmutation)
	}
	return DecodeInto(tr.Decoding, Encode(tr.Encoding, from), to)
}

// Transmute encodes the given `from` value into a tree, and decodes the tree
// into a variable of type `T`.
func Transmute[T any](tr *Transmutation, from any) (T, error) {
	var to T
	return to, TransmuteInto(tr, from, &to)
}
