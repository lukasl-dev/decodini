package decodini

import (
	"fmt"
	"reflect"
)

type MapDecoder struct {
	dec *Decoding
}

func NewMapDecoder(dec *Decoding) Decoder {
	if dec == nil {
		dec = &defaultDecoding
	}
	return &MapDecoder{dec: dec}
}

func (d *MapDecoder) Decode(tr *Tree, target DecodeTarget) error {
	if tr.Value.Kind() != reflect.Map {
		return newDecodeErrorf(
			tr.Path,
			"map decoder can't decode %s", tr.Value.Kind(),
		)
	}

	switch target.Value.Kind() {
	case reflect.Map, reflect.Interface:
		return d.decodeIntoMap(tr, target)
	case reflect.Struct:
		return d.decodeIntoStruct(tr, target)
	default:
		return newDecodeErrorf(tr.Path, "cannot decode map into %s", target.Value.Type())
	}
}

func (d *MapDecoder) decodeIntoMap(tr *Tree, target DecodeTarget) error {
	typ := inferType(tr, target)
	created := reflect.MakeMapWithSize(typ, tr.NumChildren())

	for _, child := range tr.Children() {
		key := reflect.ValueOf(child.Name())
		val := reflect.New(typ.Elem()).Elem()

		target := DecodeTarget{Value: val, mapKey: key}
		err := d.dec.decode(append(tr.Path, child.Name()), child, target)
		if err != nil {
			return err
		}

		created.SetMapIndex(key, val)
	}

	target.Value.Set(created)
	return nil
}

func (d *MapDecoder) decodeIntoStruct(tr *Tree, target DecodeTarget) error {
	typ := inferType(tr, target)
	created := reflect.New(typ).Elem()

	for _, child := range tr.Children() {
		if err := d.decodeIntoStructField(tr, child, created, target); err != nil {
			return err
		}
	}

	target.Value.Set(created)
	return nil
}

func (d *MapDecoder) decodeIntoStructField(root, node *Tree, target reflect.Value, parent DecodeTarget) error {
	sf, vf := d.dec.structFieldByName(target, fmt.Sprint(node.Name()))
	subtarget := DecodeTarget{structField: sf}

	if vf.IsValid() {
		val := reflect.New(vf.Type()).Elem()
		subtarget.Value = val
	} else {
		if d.dec.ResolveUnknownField == nil {
			return newDecodeErrorf(root.Path, "no such field: %s in %s", node.Name(), target.Type())
		}

		value, err := d.dec.ResolveUnknownField(root, subtarget)
		if err != nil {
			return err
		}
		if !value.IsValid() {
			return nil
		}
		subtarget.Value = value
	}

	err := d.dec.decode(append(root.Path, node.Name()), node, subtarget)
	if err != nil {
		return err
	}

	vf.Set(subtarget.Value)
	return nil
}
