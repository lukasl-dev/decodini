package decodini

import "reflect"

func includeStructField(tag string, sf reflect.StructField) bool {
	return sf.IsExported() && sf.Tag.Get(tag) != "-"
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

		tagName, hasTag := sf.Tag.Lookup(tag)
		if hasTag {
			if tagName == name {
				return sf, vf
			}
			continue
		}

		if sf.Name == name {
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
