package decodini

import "reflect"

type StructDecoder struct {
	dec *Decoding
}

func NewStructDecoder(dec *Decoding) Decoder {
	if dec == nil {
		dec = &defaultDecoding
	}
	return &StructDecoder{dec: dec}
}

func (d *StructDecoder) Decode(tr *Tree, target DecodeTarget) error {
	if tr.Value.Kind() != reflect.Struct {
		return newDecodeErrorf(
			tr.Path,
			"struct decoder can't decode %s", tr.Value.Kind(),
		)
	}

	switch target.Value.Kind() {
	case reflect.Struct, reflect.Interface:
		return d.decodeIntoStruct(tr, target)
	case reflect.Map:
		return d.decodeIntoMap(tr, target)
	default:
		return newDecodeErrorf(
			tr.Path,
			"cannot decode struct into %s", target.Value.Kind(),
		)
	}
}

func (d *StructDecoder) decodeIntoStruct(tr *Tree, target DecodeTarget) error {
	var typ reflect.Type
	if target.Value.Kind() == reflect.Interface {
		typ = tr.Value.Type()
	} else {
		typ = target.Value.Type()
	}
	created := reflect.New(typ)

	for _, child := range tr.Children() {
		name, isString := child.Name().(string)
		if !isString {
			return newDecodeErrorf(tr.Path, "struct fields must be strings, but got %T", name)
		}

		sf, field := d.dec.structFieldByName(created.Elem(), name)
		subtarget := DecodeTarget{Value: field, structField: sf}

		if !field.IsValid() {
			if d.dec.ResolveUnknownField == nil {
				return newDecodeErrorf(tr.Path, "struct field %s does not exist", name)
			}

			value, err := d.dec.ResolveUnknownField(tr, subtarget)
			if err != nil {
				return err
			}
			if !value.IsValid() {
				continue
			}
			subtarget.Value = value
		}

		if err := d.dec.decode(append(tr.Path, name), child, subtarget); err != nil {
			return err
		}
	}

	target.Value.Set(created.Elem())
	return nil
}

func (d *StructDecoder) decodeIntoMap(tr *Tree, target DecodeTarget) error {
	typ := inferType(tr, target)
	created := reflect.MakeMapWithSize(typ, tr.NumChildren())

	for _, child := range tr.Children() {
		key := reflect.ValueOf(child.Name())
		val := reflect.New(typ.Elem()).Elem()

		subtarget := DecodeTarget{Value: val, mapKey: key}
		err := d.dec.decode(append(tr.Path, child.Name()), child, subtarget)
		if err != nil {
			return err
		}

		created.SetMapIndex(key, val)
	}

	target.Value.Set(created)
	return nil
}
