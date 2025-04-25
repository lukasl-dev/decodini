package decodini

type ArrayDecoder struct {
	dec *Decoding
}

func NewArrayDecoder(dec *Decoding) Decoder {
	if dec == nil {
		dec = &defaultDecoding
	}
	return &ArrayDecoder{dec: dec}
}

func (d *ArrayDecoder) Decode(tr *Tree, target DecodeTarget) error {
	panic("decodini: arrays are not implemented yet")
}
