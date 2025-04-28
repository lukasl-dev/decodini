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

func (e *Encoding) encode(path []any, val reflect.Value) *Tree {
	switch val.Kind() {
	case reflect.Ptr:
		return e.encode(path, val.Elem())
	case reflect.Struct:
		return e.encodeStruct(path, val)
	case reflect.Slice, reflect.Array:
		return e.encodeSlice(path, val)
	case reflect.Map:
		return e.encodeMap(path, val)
	case reflect.Invalid, reflect.String, reflect.Bool,
		reflect.Float32, reflect.Float64,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return e.encodeScalar(path, val)
	case reflect.Interface:
		if val.IsNil() {
			return NewTree(path, val)
		}
		return e.encode(path, val.Elem())
	default:
		if val.CanInterface() {
			panic(fmt.Sprintf("unsupported kind %s for value %v", val.Kind(), val.Interface()))
		}
		panic(fmt.Sprintf("unsupported kind %s", val.Kind()))
	}
}

func (e *Encoding) encodeStruct(path []any, val reflect.Value) *Tree {
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

		enc := e.encode(append(path, name), vf)
		enc.structField = tf
		children = append(children, enc)
	}
	return NewTree(path, val, children...)
}

func (e *Encoding) encodeSlice(path []any, val reflect.Value) *Tree {
	var children []*Tree
	for i := 0; i < val.Len(); i++ {
		children = append(children, e.encode(append(path, i), val.Index(i)))
	}
	return NewTree(path, val, children...)
}

func (e *Encoding) encodeMap(path []any, val reflect.Value) *Tree {
	var children []*Tree
	for _, key := range val.MapKeys() {
		children = append(children, e.encode(append(path, key.Interface()), val.MapIndex(key)))
	}
	return NewTree(path, val, children...)
}

func (e *Encoding) encodeScalar(path []any, val reflect.Value) *Tree {
	return NewTree(path, val)
}
