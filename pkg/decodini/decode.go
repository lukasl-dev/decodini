package decodini

import "reflect"

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
	if target.Value.Kind() == reflect.Ptr {
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

		from := node.Child(targetName)
		subTarget := DecodeTarget{
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
			uFrom, uErr := dec.Unmatched(from, subTarget)
			if uErr != nil {
				return uErr
			}
			if uFrom == nil {
				continue
			}
			from = uFrom
		}

		if err := dec.into(from, subTarget); err != nil {
			return err
		}
	}
	return nil
}

func (dec *Decoding) intoSlice(node *Tree, target DecodeTarget) error {
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
