package decodini

import (
	"reflect"
	"unicode/utf16"
)

type Decoder func(tr *Tree, target DecodeTarget) error

func DecodeIgnoreUnmatched(*Tree, DecodeTarget) (*Tree, error) {
	return nil, nil
}

type Decoding struct {
	StructTag string

	// Decoder is a custom decoder. If nil is returned, the default decoding
	// mechanism is used.
	Decoder func(tr *Tree, target DecodeTarget) Decoder

	Unmatched func(tr *Tree, target DecodeTarget) (*Tree, error)
}

var defaultDecoding = Decoding{
	StructTag: "decodini",
	Unmatched: nil,
}

func DecodeInto(dec *Decoding, tr *Tree, into any) error {
	if dec == nil {
		dec = &defaultDecoding
	}
	if dec.StructTag == "" {
		dec.StructTag = defaultDecoding.StructTag
	}

	if tr == nil {
		panic("decodini: cannot decode from nil tree")
	}

	rVal, isVal := into.(reflect.Value)
	if !isVal {
		rVal = reflect.ValueOf(into)
	}

	return dec.into(tr, DecodeTarget{Value: rVal})
}

func Decode[T any](dec *Decoding, tr *Tree) (T, error) {
	var to T
	return to, DecodeInto(dec, tr, &to)
}

func (dec *Decoding) into(node *Tree, target DecodeTarget) error {
	if target.Value.Kind() == reflect.Pointer {
		if node.IsNil() {
			if target.Value.CanSet() {
				target.Value.Set(reflect.Zero(target.Value.Type()))
				return nil
			}
			if target.Value.IsNil() {
				return newDecodeErrorf(node, target, "cannot decode into unsettable value")
			}

			target.Value = target.Value.Elem()
			return dec.into(node, target)
		}

		if target.Value.IsNil() {
			if !target.Value.CanSet() {
				return newDecodeErrorf(node, target, "cannot decode into unsettable value")
			}

			target.Value.Set(reflect.New(target.Value.Type().Elem()))
		}

		target.Value = target.Value.Elem()
		return dec.into(node, target)
	}

	if !target.Value.CanSet() {
		return newDecodeErrorf(node, target, "cannot decode into unsettable value")
	}

	if node.IsNil() {
		target.Value.Set(reflect.Zero(target.Value.Type()))
		return nil
	}

	if dec.Decoder != nil {
		if fn := dec.Decoder(node, target); fn != nil {
			return fn(node, target)
		}
	}

	if target.IsPrimitive() {
		return dec.intoScalar(node, target)
	}

	switch target.Value.Kind() {
	case reflect.Struct:
		return dec.intoStruct(node, target)

	case reflect.Slice:
		return dec.intoSlice(node, target)

	case reflect.Array:
		return dec.intoArray(node, target)

	case reflect.Map:
		return dec.intoMap(node, target)
	default:
		return newDecodeErrorf(
			node,
			target,
			"cannot decode into %s", target.Value.Type(),
		)
	}
}

func (dec *Decoding) intoScalar(node *Tree, target DecodeTarget) error {
	if target.Value.Kind() == reflect.String {
		switch node.Value().Kind() {
		case reflect.Slice, reflect.Array:
			elemKind := node.Value().Type().Elem().Kind()
			n := node.Value().Len()
			switch elemKind {
			case reflect.Uint8, reflect.Int8:
				b := make([]byte, n)
				for i := range n {
					v := node.Value().Index(i)
					if elemKind == reflect.Uint8 {
						b[i] = byte(v.Uint())
					} else {
						b[i] = byte(v.Int())
					}
				}
				target.Value.SetString(string(b))
				return nil
			case reflect.Uint16, reflect.Int16:
				u := make([]uint16, n)
				for i := range n {
					v := node.Value().Index(i)
					if elemKind == reflect.Uint16 {
						u[i] = uint16(v.Uint())
					} else {
						u[i] = uint16(v.Int())
					}
				}
				target.Value.SetString(string(utf16.Decode(u)))
				return nil
			case reflect.Int, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
				r := make([]rune, n)
				for i := range n {
					v := node.Value().Index(i)
					switch elemKind {
					case reflect.Int, reflect.Int32, reflect.Int64:
						r[i] = rune(v.Int())
					default:
						r[i] = rune(v.Uint())
					}
				}
				target.Value.SetString(string(r))
				return nil
			}
		}
	}

	if !node.IsPrimitive() {
		return newDecodeErrorf(
			node,
			target,
			"leaf decoder can't decode non-leaf %s", node.Value().Kind(),
		)
	}

	target.Value.Set(node.Value())
	return nil
}

func (dec *Decoding) intoStruct(node *Tree, target DecodeTarget) error {
	switch node.Value().Kind() {
	case reflect.Struct, reflect.Map:
		return dec.intoStructFromStructOrMap(node, target)
	default:
		return newDecodeErrorf(
			node,
			target,
			"cannot decode %s into struct", node.Value().Kind(),
		)
	}
}

func (dec *Decoding) intoStructFromStructOrMap(node *Tree, target DecodeTarget) error {
	targetType := target.Value.Type()
	for i := range target.Value.NumField() {
		targetSF := targetType.Field(i)
		if !includeStructField(dec.StructTag, targetSF) {
			continue
		}

		targetName := structFieldName(dec.StructTag, targetSF)

		if targetSF.Anonymous {
			from := node
			if child := node.Child(targetName); child != nil {
				from = child
			}

			sub := DecodeTarget{
				Value:       target.Value.Field(i),
				structField: &targetSF,
			}
			if err := dec.into(from, sub); err != nil {
				return err
			}
			continue
		}

		from := node.Child(targetName)
		sub := DecodeTarget{
			Name:        targetName,
			Value:       target.Value.Field(i),
			structField: &targetSF,
		}

		if from == nil {
			if dec.Unmatched == nil {
				return newDecodeErrorf(
					node.dummyChild(targetName),
					target,
					"struct field %s is unmatched in source tree", targetName,
				)
			}
			uFrom, uErr := dec.Unmatched(from, sub)
			if uErr != nil {
				return uErr
			}
			if uFrom == nil {
				continue
			}
			from = uFrom
		}

		if err := dec.into(from, sub); err != nil {
			return err
		}
	}
	return nil
}

func (dec *Decoding) intoSlice(node *Tree, target DecodeTarget) error {
	if node.Value().Kind() == reflect.String {
		s := node.Value().String()
		elemType := target.Value.Type().Elem()
		elemKind := elemType.Kind()
		switch elemKind {
		case reflect.Uint8, reflect.Int8:
			b := []byte(s)
			dst := reflect.MakeSlice(target.Value.Type(), len(b), len(b))
			for i := range len(b) {
				dst.Index(i).Set(reflect.ValueOf(b[i]).Convert(elemType))
			}
			target.Value.Set(dst)
			return nil
		case reflect.Uint16, reflect.Int16:
			u := utf16.Encode([]rune(s))
			dst := reflect.MakeSlice(target.Value.Type(), len(u), len(u))
			for i := range len(u) {
				dst.Index(i).Set(reflect.ValueOf(u[i]).Convert(elemType))
			}
			target.Value.Set(dst)
			return nil
		case reflect.Int, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			r := []rune(s)
			dst := reflect.MakeSlice(target.Value.Type(), len(r), len(r))
			for i := range len(r) {
				dst.Index(i).Set(reflect.ValueOf(r[i]).Convert(elemType))
			}
			target.Value.Set(dst)
			return nil
		}
	}

	switch node.Value().Kind() {
	case reflect.Slice, reflect.Array:
		return dec.intoSliceFromSliceOrArray(node, target)
	case reflect.Map:
		return dec.intoSliceFromMap(node, target)
	default:
		return newDecodeErrorf(
			node,
			target,
			"cannot decode %s into slice", node.Value().Kind(),
		)
	}
}

func (dec *Decoding) intoSliceFromSliceOrArray(
	node *Tree,
	target DecodeTarget,
) error {
	nChildren := int(node.NumChildren())
	if target.Value.IsNil() {
		target.Value.Set(
			reflect.MakeSlice(target.Value.Type(), nChildren, nChildren),
		)
	}
	typ := inferType(node, target)

	for from := range node.Children() {
		val := reflect.New(typ.Elem()).Elem()

		subtarget := DecodeTarget{Name: from.Name(), Value: val}
		err := dec.into(from, subtarget)
		if err != nil {
			return err
		}

		target.Value.Index(from.Name().(int)).Set(val)
	}

	return nil
}

func (dec *Decoding) intoSliceFromMap(node *Tree, target DecodeTarget) error {
	nChildren := int(node.NumChildren())
	if target.Value.IsNil() {
		target.Value.Set(
			reflect.MakeSlice(target.Value.Type(), nChildren, nChildren),
		)
	}
	typ := inferType(node, target)

	i := 0
	for from := range node.Children() {
		val := reflect.New(typ.Elem()).Elem()

		subtarget := DecodeTarget{Name: i, Value: val}
		err := dec.into(from, subtarget)
		if err != nil {
			return err
		}

		target.Value.Index(i).Set(val)
		i++
	}

	return nil
}

func (dec *Decoding) intoArray(node *Tree, target DecodeTarget) error {
	// TODO: implement array decoding
	return newDecodeErrorf(
		node,
		target,
		"decodini does currently not support arrays",
	)
}

func (dec *Decoding) intoMap(node *Tree, target DecodeTarget) error {
	switch node.Value().Kind() {
	case reflect.Map, reflect.Struct:
		return dec.intoMapFromMapOrStruct(node, target)
	default:
		return newDecodeErrorf(
			node,
			target,
			"cannot decode %s into map", node.Value().Kind(),
		)
	}
}

func (dec *Decoding) intoMapFromMapOrStruct(node *Tree, target DecodeTarget) error {
	if target.Value.IsNil() {
		target.Value.Set(
			reflect.MakeMapWithSize(target.Value.Type(), int(node.NumChildren())),
		)
	}
	typ := inferType(node, target)

	for from := range node.Children() {
		key := reflect.ValueOf(from.Name())
		val := reflect.New(typ.Elem()).Elem()

		subtarget := DecodeTarget{Name: key.Interface(), Value: val}
		err := dec.into(from, subtarget)
		if err != nil {
			return err
		}

		target.Value.SetMapIndex(key, val)
	}

	return nil
}
