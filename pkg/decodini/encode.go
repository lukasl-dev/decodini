package decodini

import (
	"fmt"
	"reflect"
)

type Encoding struct {
	// StructTag is the name of the struct tag used to specify the name of a
	// field. The default value is "decodini".
	//
	// If the value of a struct field is "-", the field is ignored.
	StructTag string
}

// defaultEncoding is the default encoding used by Encode.
var defaultEncoding = Encoding{
	StructTag: "decodini",
}

// Encode encodes the given value into a tree based on the given encoding. If
// the encoding is nil, the default encoding is used.
func Encode(enc *Encoding, val any) *Tree {
	if enc == nil {
		enc = &defaultEncoding
	}
	rVal := reflect.ValueOf(val)
	return enc.encode(nil, rVal)
}

func (e *Encoding) encode(name any, val reflect.Value) *Tree {
	switch val.Kind() {
	case reflect.Ptr:
		return e.encode(name, val.Elem())
	case reflect.Struct:
		return e.encodeStruct(name, val)
	case reflect.Slice, reflect.Array:
		return e.encodeSlice(name, val)
	case reflect.Map:
		return e.encodeMap(name, val)
	case reflect.Invalid, reflect.String, reflect.Bool,
		reflect.Float32, reflect.Float64,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return e.encodeScalar(name, val)
	default:
		panic(fmt.Sprintf("unsupported kind %s", val.Kind()))
	}
}

func (e *Encoding) encodeStruct(name any, val reflect.Value) *Tree {
	typ := val.Type()
	var children []*Tree
	for i := 0; i < val.NumField(); i++ {
		tf, vf := typ.Field(i), val.Field(i)
		if !tf.IsExported() {
			continue
		}
		name := tf.Name

		tag := tf.Tag.Get(e.StructTag)
		if tag != "" {
			if tag == "-" {
				continue
			}
			name = tag
		}

		children = append(children, e.encode(name, vf))
	}
	return NewTree(name, val, children...)
}

func (e *Encoding) encodeSlice(name any, val reflect.Value) *Tree {
	var children []*Tree
	for i := 0; i < val.Len(); i++ {
		children = append(children, e.encode(i, val.Index(i)))
	}
	return NewTree(name, val, children...)
}

func (e *Encoding) encodeMap(name any, val reflect.Value) *Tree {
	var children []*Tree
	for _, key := range val.MapKeys() {
		children = append(children, e.encode(key.Interface(), val.MapIndex(key)))
	}
	return NewTree(name, val, children...)
}

func (e *Encoding) encodeScalar(name any, val reflect.Value) *Tree {
	return NewTree(name, val)
}
