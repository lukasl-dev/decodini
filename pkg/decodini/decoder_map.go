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
	var typ reflect.Type
	if target.Value.Kind() == reflect.Interface {
		typ = tr.Value.Type()
	} else {
		typ = target.Value.Type()
	}
	created := reflect.MakeMapWithSize(typ, len(tr.Children))

	for _, child := range tr.Children {
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
	var typ reflect.Type
	if target.Value.Kind() == reflect.Interface {
		typ = tr.Value.Type()
	} else {
		typ = target.Value.Type()
	}
	created := reflect.New(typ).Elem()

	for _, child := range tr.Children {
		sf, vf := d.dec.structFieldByName(created, fmt.Sprint(child.Name()))
		if !vf.IsValid() {
			return newDecodeErrorf(tr.Path, "no such field: %s in %s", child.Name(), typ)
		}

		val := reflect.New(vf.Type()).Elem()

		subtarget := DecodeTarget{Value: val, structField: sf}
		err := d.dec.decode(append(tr.Path, child.Name()), child, subtarget)
		if err != nil {
			return err
		}

		vf.Set(val)
	}

	target.Value.Set(created)
	return nil
}
