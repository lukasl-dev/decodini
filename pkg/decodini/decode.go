package decodini

import "reflect"

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

	// ResolveUnknownField is called for values for which the decoder cannot find
	// a matching field in the target struct. If it is nil, the decoder will fail
	// for such values. If the returned [reflect.Value] is invalid (i.e. the zero
	// value), the decoder will skip the field. Otherwise, decoding will continue
	// for the remaining tree at that value. Therefore, it should be a pointer.
	ResolveUnknownField func(tr *Tree, target DecodeTarget) (reflect.Value, error)
}

// DecodeIgnoreUnknownField can be used as a [Decoding.ResolveUnknownField]
// to ignore unknown fields.
func DecodeIgnoreUnknownField(_ *Tree, _ DecodeTarget) (reflect.Value, error) {
	return reflect.Value{}, nil
}

var defaultDecoding = Decoding{
	StructTag: "decodini",
}

type DecodeTarget struct {
	Value reflect.Value
	// Parent *DecodeTarget

	structField reflect.StructField
	mapKey      reflect.Value
	sliceIndex  *int
}

// IsStructField returns whether the target represents a struct field.
func (t DecodeTarget) IsStructField() bool {
	return false
	// if t.Parent == nil {
	// 	return false
	// }
	//
	// parent := t.Parent.Value
	// return parent.IsValid() && parent.Kind() == reflect.Struct
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
	return false
	// if t.Parent == nil {
	// 	return false
	// }
	//
	// parent := t.Parent.Value
	// return parent.IsValid() && parent.Kind() == reflect.Map
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
	return false
	// if t.Parent == nil {
	// 	return false
	// }
	//
	// parent := t.Parent.Value
	// return parent.IsValid() && parent.Kind() == reflect.Slice
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

	if d.Decoder != nil {
		// if custom decoder is specified
		dec := d.Decoder(tr, target)
		if dec != nil {
			return dec.Decode(tr, target)
		}
	}

	if tr.IsPrimitive() {
		return NewLeafDecoder(d).Decode(tr, target)
	}

	switch tr.Value.Kind() {
	case reflect.Struct:
		return NewStructDecoder(d).Decode(tr, target)
	case reflect.Slice:
		return NewSliceDecoder(d).Decode(tr, target)
	case reflect.Array:
		return NewArrayDecoder(d).Decode(tr, target)
	case reflect.Map:
		return NewMapDecoder(d).Decode(tr, target)
	// TODO: pointers
	default:
		return newDecodeErrorf(path, "cannot decode %s", target.Value.Type())
	}
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

func inferType(tr *Tree, target DecodeTarget) reflect.Type {
	if target.Value.Kind() == reflect.Interface {
		return tr.Value.Type()
	}
	return target.Value.Type()
}
