package decodini

import (
	"fmt"
	"reflect"
)

type DecoderFunc func(tr *Tree, dst reflect.Value) error

type Decoding struct {
	// StructTag is the name of the struct tag used to specify the name of a
	// field. The default value is "decodini".
	//
	// If the value of a struct field is "-", the field is ignored.
	StructTag string

	// LeafDecoder is a custom leaf decoder function. If nil, the default leaf
	// decoder is used.
	LeafDecoder func(leaf *Tree, dst reflect.Value) DecoderFunc

	// StructDecoder is a custom struct decoder function. If nil, the default
	// struct decoder is used.
	StructDecoder func(str *Tree, dst reflect.Value) DecoderFunc

	// SliceDecoder is a custom slice decoder function. If nil, the default slice
	// decoder is used.
	SliceDecoder func(sl *Tree, dst reflect.Value) DecoderFunc

	// ArrayDecoder is a custom array decoder function. If nil, the default array
	// decoder is used.
	ArrayDecoder func(ar *Tree, dst reflect.Value) DecoderFunc

	// MapDecoder is a custom map decoder function. If nil, the default map
	// decoder is used.
	MapDecoder func(ma *Tree, dst reflect.Value) DecoderFunc
}

var defaultDecoding = Decoding{
	StructTag: "decodini",
}

func Decode(dec *Decoding, tr *Tree, dst any) error {
	if dec == nil {
		dec = &defaultDecoding
	}
	rVal := reflect.ValueOf(dst)
	return dec.decode(nil, tr, rVal)
}

func DefaultDecode(tr *Tree, dst any) error {
	return Decode(nil, tr, dst)
}

func (d *Decoding) decode(path []any, tr *Tree, dst reflect.Value) error {
	if tr == nil {
		return newDecodeErrorf(path, "decodini: cannot decode nil tree")
	}

	if dst.Kind() == reflect.Ptr {
		return d.decode(path, tr, dst.Elem())
	}
	if !dst.CanSet() {
		return newDecodeErrorf(path, "decodini: cannot decode into unsettable value")
	}

	if tr.Nil() {
		dst.Set(reflect.Zero(dst.Type()))
		return nil
	}

	if tr.Leaf() {
		return d.decodeLeaf(path, tr, dst)
	}
	switch tr.Value.Kind() {
	case reflect.Struct:
		return d.decodeStruct(path, tr, dst)
	case reflect.Slice:
		return d.decodeSlice(path, tr, dst)
	case reflect.Array:
		return d.decodeArray(path, tr, dst)
	case reflect.Map:
		return d.decodeMap(path, tr, dst)
	// TODO: pointers
	default:
		return fmt.Errorf("decodini: cannot decode into %s", dst.Type())
	}
}

func (d *Decoding) decodeLeaf(path []any, tr *Tree, dst reflect.Value) error {
	if d.LeafDecoder != nil {
		fn := d.LeafDecoder(tr, dst)
		if fn != nil {
			return fn(tr, dst)
		}
	}
	dst.Set(tr.Value)
	return nil
}

func (d *Decoding) decodeStruct(path []any, tr *Tree, dst reflect.Value) error {
	if d.StructDecoder != nil {
		fn := d.StructDecoder(tr, dst)
		if fn != nil {
			return fn(tr, dst)
		}
	}

	switch dst.Kind() {
	case reflect.Struct, reflect.Interface:
		return d.decodeStructIntoStruct(path, tr, dst)
	default:
		return newDecodeErrorf(path, "decodini: cannot decode struct into %s", dst.Kind())
	}
}

func (d *Decoding) decodeSlice(path []any, tr *Tree, dst reflect.Value) error {
	if d.SliceDecoder != nil {
		fn := d.SliceDecoder(tr, dst)
		if fn != nil {
			return fn(tr, dst)
		}
	}

	switch dst.Kind() {
	case reflect.Slice, reflect.Interface:
		return d.decodeSliceIntoSlice(path, tr, dst)
	default:
		return fmt.Errorf("decodini: cannot decode slice into %s", dst.Type())
	}
}

func (d *Decoding) decodeArray(path []any, tr *Tree, dst reflect.Value) error {
	if d.ArrayDecoder != nil {
		fn := d.ArrayDecoder(tr, dst)
		if fn != nil {
			return fn(tr, dst)
		}
	}

	switch dst.Kind() {
	// TODO: implement
	// case reflect.Array, reflect.Interface:
	// 	return d.decodeArrayIntoArray(tr, dst)
	default:
		return fmt.Errorf("decodini: cannot decode array into %s", dst.Type())
	}
}

func (d *Decoding) decodeMap(path []any, tr *Tree, dst reflect.Value) error {
	if d.MapDecoder != nil {
		fn := d.MapDecoder(tr, dst)
		if fn != nil {
			return fn(tr, dst)
		}
	}

	switch dst.Kind() {
	case reflect.Map, reflect.Interface:
		return d.decodeMapIntoMap(path, tr, dst)
	default:
		return fmt.Errorf("decodini: cannot decode map into %s", dst.Type())
	}
}

func (d *Decoding) decodeStructIntoStruct(
	path []any,
	tr *Tree,
	dst reflect.Value,
) error {
	var typ reflect.Type
	if dst.Kind() == reflect.Interface {
		typ = tr.Value.Type()
	} else {
		typ = dst.Type()
	}
	created := reflect.New(typ)

	for _, child := range tr.Children {
		name, isString := child.Name.(string)
		if !isString {
			return fmt.Errorf("decodini: struct fields must be strings, but got %T", name)
		}

		_, field := d.structFieldByName(created.Elem(), name)
		if !field.IsValid() {
			// TODO: allow ignoring fields
			return fmt.Errorf("decodini: struct field %s does not exist", name)
		}

		if err := d.decode(append(path, name), child, field); err != nil {
			return err
		}
	}

	dst.Set(created.Elem())
	return nil
}

func (d *Decoding) decodeSliceIntoSlice(
	path []any,
	tr *Tree,
	dst reflect.Value,
) error {
	var typ reflect.Type
	if dst.Kind() == reflect.Interface {
		typ = tr.Value.Type()
	} else {
		typ = dst.Type()
	}
	created := reflect.MakeSlice(typ, len(tr.Children), len(tr.Children))

	for i, child := range tr.Children {
		if err := d.decode(append(path, i), child, created.Index(i)); err != nil {
			return err
		}
	}
	dst.Set(created)
	return nil
}

func (d *Decoding) decodeMapIntoMap(
	path []any,
	tr *Tree,
	dst reflect.Value,
) error {
	var typ reflect.Type
	if dst.Kind() == reflect.Interface {
		typ = tr.Value.Type()
	} else {
		typ = dst.Type()
	}
	created := reflect.MakeMapWithSize(typ, len(tr.Children))

	for _, child := range tr.Children {
		key := reflect.ValueOf(child.Name)

		val := reflect.New(typ.Elem()).Elem()
		err := d.decode(append(path, child.Name), child, val)
		if err != nil {
			return err
		}

		created.SetMapIndex(key, val)
	}

	dst.Set(created)
	return nil
}

func (d *Decoding) structFieldByName(
	val reflect.Value,
	name string,
) (reflect.StructField, reflect.Value) {
	for i := 0; i < val.NumField(); i++ {
		tf, vf := val.Type().Field(i), val.Field(i)

		fieldName := tf.Name
		if tag := tf.Tag.Get(d.StructTag); tag != "" {
			if tag == "-" {
				continue
			}
			fieldName = tag
		}

		if fieldName == name {
			return tf, vf
		}
	}
	return reflect.StructField{}, reflect.Value{}
}
