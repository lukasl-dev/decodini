package decodini

import "reflect"

type Decoding struct {
	StructTag string
}

var defaultDecoding = Decoding{
	StructTag: "decodini",
}

func DecodeInto(dec *Decoding, tr *Tree, into any) error {
	if dec == nil {
		dec = &defaultDecoding
	}

	rVal, isVal := into.(reflect.Value)
	if !isVal {
		rVal = reflect.ValueOf(into)
	}

	return dec.into(tr, DecodeTarget{Value: rVal})
}

func Decode[T any](dec *Decoding, tr *Tree) (T, error) {
	var to T
	return to, DecodeInto(dec, tr, &to)
}

func (dec *Decoding) into(node *Tree, target DecodeTarget) error {
	if target.Value.Kind() == reflect.Ptr {
		target.Value = target.Value.Elem()
		return dec.into(node, target)
	}

	if !target.Value.CanSet() {
		return newDecodeErrorf(node, "cannot decode into unsettable value")
	}

	if node.IsNil() {
		target.Value.Set(reflect.Zero(target.Value.Type()))
		return nil
	}

	if target.IsPrimitive() {
		return dec.intoScalar(node, target)
	}

	switch target.Value.Kind() {
	case reflect.Struct:
		return dec.intoStruct(node, target)

	case reflect.Slice:
		return dec.intoSlice(node, target)

	case reflect.Array:
		return dec.intoArray(node, target)

	case reflect.Map:
		return dec.intoMap(node, target)
	default:
		return newDecodeErrorf(node, "cannot decode into %s", target.Value.Type())
	}

	return nil
}

func (dec *Decoding) intoScalar(node *Tree, target DecodeTarget) error {
	if !node.IsPrimitive() {
		return newDecodeErrorf(
			node,
			"leaf decoder can't decode non-leaf %s", node.Value().Kind(),
		)
	}

	target.Value.Set(node.Value())
	return nil
}

func (dec *Decoding) intoStruct(node *Tree, target DecodeTarget) error {
	return nil
}

func (dec *Decoding) intoStructFromStruct(node *Tree, target DecodeTarget) error {
	panic("not implemented yet")
}

func (dec *Decoding) intoStructFromMap(node *Tree, target DecodeTarget) error {
	panic("not implemented yet")
}

func (dec *Decoding) intoSlice(node *Tree, target DecodeTarget) error {
	panic("not implemented yet")
}

func (dec *Decoding) intoArray(node *Tree, target DecodeTarget) error {
	panic("not implemented yet")
}

func (dec *Decoding) intoMap(node *Tree, target DecodeTarget) error {
	panic("not implemented yet")
}
