package decodini

type LeafDecoder struct {
	dec *Decoding
}

func NewLeafDecoder(dec *Decoding) Decoder {
	if dec == nil {
		dec = &defaultDecoding
	}
	return &LeafDecoder{dec: dec}
}

func (d *LeafDecoder) Decode(tr *Tree, target DecodeTarget) error {
	if !tr.IsLeaf() {
		return newDecodeErrorf(
			tr.Path,
			"leaf decoder can't decode non-leaf %s", tr.Value.Kind(),
		)
	}

	target.Value.Set(tr.Value)
	return nil
}
