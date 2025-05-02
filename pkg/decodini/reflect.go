package decodini

import "reflect"

func includeStructField(tag string, sf reflect.StructField) bool {
	return sf.IsExported() && sf.Tag.Get(tag) != "-"
}

func structFieldName(tag string, sf reflect.StructField) string {
	if tagName, hasTag := sf.Tag.Lookup(tag); hasTag {
		return tagName
	}
	return sf.Name
}

func structFieldByName(
	tag string,
	val reflect.Value,
	name string,
) (reflect.StructField, reflect.Value) {
	if val.Kind() != reflect.Struct {
		panic("decodini: cannot get struct field of non-struct")
	}

	typ := val.Type()
	for i := range val.NumField() {
		sf := typ.Field(i)
		if !includeStructField(tag, sf) {
			continue
		}
		vf := val.Field(i)

		if structFieldName(tag, sf) == name {
			return sf, vf
		}
	}

	return reflect.StructField{}, reflect.Value{}
}

func isPrimitive(kind reflect.Kind) bool {
	switch kind {
	case reflect.Array, reflect.Slice, reflect.Map, reflect.Struct:
		return false
	default:
		return true
	}
}

func inferType(from *Tree, target DecodeTarget) reflect.Type {
	if target.Value.Kind() == reflect.Interface {
		return from.Value().Type()
	}
	return target.Value.Type()
}
