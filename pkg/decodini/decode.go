package decodini

import (
	"fmt"
	"reflect"
)

type Decoder interface {
	Decode(tr *Tree, target DecodeTarget) error
}

type DecoderFunc func(tr *Tree, target DecodeTarget) error

func (f DecoderFunc) Decode(tr *Tree, target DecodeTarget) error {
	return f(tr, target)
}

type Decoding struct {
	// StructTag is the name of the struct tag used to specify the name of a
	// field. The default value is "decodini".
	//
	// If the value of a struct field is "-", the field is ignored.
	StructTag string

	// Decoder is a custom decoder. If nil is returned, the default decoding
	// mechanism is used.
	Decoder func(tr *Tree, target DecodeTarget) Decoder
}

var defaultDecoding = Decoding{
	StructTag: "decodini",
}

type DecodeTarget struct {
	Value reflect.Value

	structField reflect.StructField
	mapKey      reflect.Value
	sliceIndex  *int
}

// IsStructField returns whether the target represents a struct field.
func (t DecodeTarget) IsStructField() bool {
	return t.structField.Name != ""
}

// StructField returns the struct field. If the target does not represent a struct
// field (i.e. IsStructField() is false), it panics.
func (t DecodeTarget) StructField() reflect.StructField {
	if !t.IsStructField() {
		panic("decodini: target does not represent a struct field")
	}
	return t.structField
}

// IsMapKey returns whether the target represents a map key.
func (t DecodeTarget) IsMapKey() bool {
	return t.mapKey.IsValid()
}

// MapKey returns the map key. If the target does not represent a map key (i.e.
// IsMapKey() is false), it panics.
func (t DecodeTarget) MapKey() reflect.Value {
	if !t.IsMapKey() {
		panic("decodini: target does not represent a map key")
	}
	return t.mapKey
}

// IsSliceIndex returns whether the target represents a slice index.
func (t DecodeTarget) IsSliceIndex() bool {
	return t.sliceIndex != nil
}

// SliceIndex returns the slice index. If the target does not represent a slice
// index (i.e. IsSliceIndex() is false), it panics.
func (t DecodeTarget) SliceIndex() int {
	if !t.IsSliceIndex() {
		panic("decodini: target does not represent a slice index")
	}
	return *t.sliceIndex
}

func Decode(dec *Decoding, tr *Tree, dst any) error {
	if dec == nil {
		dec = &defaultDecoding
	}

	var rVal reflect.Value
	switch v := dst.(type) {
	case reflect.Value:
		if v.CanSet() {
			rVal = v
		}
	default:
		rVal = reflect.ValueOf(dst)
	}

	return dec.decode(nil, tr, DecodeTarget{Value: rVal})
}

func (d *Decoding) decode(path []any, tr *Tree, target DecodeTarget) error {
	if tr == nil {
		return newDecodeErrorf(path, "cannot decode nil tree")
	}

	if target.Value.Kind() == reflect.Ptr {
		target.Value = target.Value.Elem()
		return d.decode(path, tr, target)
	}

	if !target.Value.CanSet() {
		return newDecodeErrorf(path, "cannot decode into unsettable value")
	}

	if tr.IsNil() {
		target.Value.Set(reflect.Zero(target.Value.Type()))
		return nil
	}

	if tr.IsLeaf() {
		return d.decodeLeaf(path, tr, target)
	}

	if d.Decoder != nil {
		// if custom decoder is specified
		dec := d.Decoder(tr, target)
		if dec != nil {
			return dec.Decode(tr, target)
		}
	}

	switch tr.Value.Kind() {
	case reflect.Struct:
		return d.decodeStruct(path, tr, target)
	case reflect.Slice:
		return d.decodeSlice(path, tr, target)
	case reflect.Array:
		return d.decodeArray(path, tr, target)
	case reflect.Map:
		return d.decodeMap(path, tr, target)
	// TODO: pointers
	default:
		return newDecodeErrorf(path, "cannot decode %s", target.Value.Type())
	}
}

func (d *Decoding) decodeLeaf(_ []any, tr *Tree, target DecodeTarget) error {
	if d.Decoder != nil {
		dec := d.Decoder(tr, target)
		if dec != nil {
			return dec.Decode(tr, target)
		}
	}
	target.Value.Set(tr.Value)
	return nil
}

func (d *Decoding) decodeStruct(path []any, tr *Tree, target DecodeTarget) error {
	switch target.Value.Kind() {
	case reflect.Struct, reflect.Interface:
		return d.decodeStructIntoStruct(path, tr, target)
	case reflect.Map:
		return d.decodeStructIntoMap(path, tr, target)
	default:
		return newDecodeErrorf(path, "cannot decode struct into %s", target.Value.Kind())
	}
}

func (d *Decoding) decodeSlice(path []any, tr *Tree, target DecodeTarget) error {
	switch target.Value.Kind() {
	case reflect.Slice, reflect.Interface:
		return d.decodeSliceIntoSlice(path, tr, target)
	default:
		return newDecodeErrorf(path, "cannot decode slice into %s", target.Value.Type())
	}
}

func (d *Decoding) decodeArray(path []any, tr *Tree, target DecodeTarget) error {
	switch target.Value.Kind() {
	// TODO: implement
	// case reflect.Array, reflect.Interface:
	// 	return d.decodeArrayIntoArray(tr, dst)
	default:
		return newDecodeErrorf(path, "cannot decode array into %s", target.Value.Type())
	}
}

func (d *Decoding) decodeMap(path []any, tr *Tree, target DecodeTarget) error {
	switch target.Value.Kind() {
	case reflect.Map, reflect.Interface:
		return d.decodeMapIntoMap(path, tr, target)
	case reflect.Struct:
		return d.decodeMapIntoStruct(path, tr, target)
	default:
		return newDecodeErrorf(path, "cannot decode map into %s", target.Value.Type())
	}
}

func (d *Decoding) decodeStructIntoStruct(
	path []any,
	tr *Tree,
	target DecodeTarget,
) error {
	var typ reflect.Type
	if target.Value.Kind() == reflect.Interface {
		typ = tr.Value.Type()
	} else {
		typ = target.Value.Type()
	}
	created := reflect.New(typ)

	for _, child := range tr.Children {
		name, isString := child.Name.(string)
		if !isString {
			return newDecodeErrorf(path, "struct fields must be strings, but got %T", name)
		}

		sf, field := d.structFieldByName(created.Elem(), name)
		if !field.IsValid() {
			// TODO: allow ignoring fields
			return newDecodeErrorf(path, "struct field %s does not exist", name)
		}

		subtarget := DecodeTarget{Value: field, structField: sf}
		if err := d.decode(append(path, name), child, subtarget); err != nil {
			return err
		}
	}

	target.Value.Set(created.Elem())
	return nil
}

func (d *Decoding) decodeStructIntoMap(
	path []any,
	tr *Tree,
	target DecodeTarget,
) error {
	var typ reflect.Type
	if target.Value.Kind() == reflect.Interface {
		typ = tr.Value.Type()
	} else {
		typ = target.Value.Type()
	}
	created := reflect.MakeMapWithSize(typ, len(tr.Children))

	for _, child := range tr.Children {
		key := reflect.ValueOf(child.Name)
		val := reflect.New(typ.Elem()).Elem()

		subtarget := DecodeTarget{Value: val, mapKey: key}
		err := d.decode(append(path, child.Name), child, subtarget)
		if err != nil {
			return err
		}

		created.SetMapIndex(key, val)
	}

	target.Value.Set(created)
	return nil
}

func (d *Decoding) decodeSliceIntoSlice(
	path []any,
	tr *Tree,
	target DecodeTarget,
) error {
	var typ reflect.Type
	if target.Value.Kind() == reflect.Interface {
		typ = tr.Value.Type()
	} else {
		typ = target.Value.Type()
	}
	created := reflect.MakeSlice(typ, len(tr.Children), len(tr.Children))

	for i, child := range tr.Children {
		subtarget := DecodeTarget{Value: created.Index(i), sliceIndex: &i}
		if err := d.decode(append(path, i), child, subtarget); err != nil {
			return err
		}
	}
	target.Value.Set(created)
	return nil
}

func (d *Decoding) decodeMapIntoMap(
	path []any,
	tr *Tree,
	target DecodeTarget,
) error {
	var typ reflect.Type
	if target.Value.Kind() == reflect.Interface {
		typ = tr.Value.Type()
	} else {
		typ = target.Value.Type()
	}
	created := reflect.MakeMapWithSize(typ, len(tr.Children))

	for _, child := range tr.Children {
		key := reflect.ValueOf(child.Name)
		val := reflect.New(typ.Elem()).Elem()

		target := DecodeTarget{Value: val, mapKey: key}
		err := d.decode(append(path, child.Name), child, target)
		if err != nil {
			return err
		}

		created.SetMapIndex(key, val)
	}

	target.Value.Set(created)
	return nil
}

func (d *Decoding) decodeMapIntoStruct(
	path []any,
	tr *Tree,
	target DecodeTarget,
) error {
	var typ reflect.Type
	if target.Value.Kind() == reflect.Interface {
		typ = tr.Value.Type()
	} else {
		typ = target.Value.Type()
	}
	created := reflect.New(typ).Elem()

	for _, child := range tr.Children {
		sf, vf := d.structFieldByName(created, fmt.Sprint(child.Name))
		if !vf.IsValid() {
			return newDecodeErrorf(path, "no such field: %s in %s", child.Name, typ)
		}

		val := reflect.New(vf.Type()).Elem()

		subtarget := DecodeTarget{Value: val, structField: sf}
		err := d.decode(append(path, child.Name), child, subtarget)
		if err != nil {
			return err
		}

		vf.Set(val)
	}

	target.Value.Set(created)
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
