package decodini

import "reflect"

type DecodeTarget struct {
	Name  any
	Value reflect.Value

	structField *reflect.StructField
}

func (d DecodeTarget) IsPrimitive() bool {
	return isPrimitive(d.Value.Kind())
}

func (d DecodeTarget) IsStructField() bool {
	return d.structField != nil
}

func (d DecodeTarget) StructField() reflect.StructField {
	if !d.IsStructField() {
		panic("decodini: decode target is not a struct field")
	}
	return *d.structField
}
