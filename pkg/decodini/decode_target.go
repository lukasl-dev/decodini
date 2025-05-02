package decodini

import "reflect"

type DecodeTarget struct {
	Value reflect.Value
}

func (d DecodeTarget) IsPrimitive() bool {
	return isPrimitive(d.Value.Kind())
}
