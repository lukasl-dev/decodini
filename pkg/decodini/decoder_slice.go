package decodini

import "reflect"

type SliceDecoder struct {
	dec *Decoding
}

func NewSliceDecoder(dec *Decoding) Decoder {
	if dec == nil {
		dec = &defaultDecoding
	}
	return &SliceDecoder{dec: dec}
}

func (d *SliceDecoder) Decode(tr *Tree, target DecodeTarget) error {
	if tr.Value.Kind() != reflect.Slice {
		return newDecodeErrorf(
			tr.Path,
			"map decoder can't decode %s", tr.Value.Kind(),
		)
	}

	switch target.Value.Kind() {
	case reflect.Slice, reflect.Interface:
		return d.decodeIntoSlice(tr, target)
	default:
		return newDecodeErrorf(
			tr.Path,
			"cannot decode slice into %s", target.Value.Type(),
		)
	}
}

func (d *SliceDecoder) decodeIntoSlice(tr *Tree, target DecodeTarget) error {
	typ := inferType(tr, target)
	created := reflect.MakeSlice(typ, len(tr.Children), len(tr.Children))

	for i, child := range tr.Children {
		subtarget := DecodeTarget{Parent: &target, Value: created.Index(i), sliceIndex: &i}
		if err := d.dec.decode(append(tr.Path, i), child, subtarget); err != nil {
			return err
		}
	}
	target.Value.Set(created)
	return nil
}
