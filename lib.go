package main

import (
	"errors"
	"reflect"
	"strconv"
)

var (
	EOF              = errors.New("EOF")
	ErrUnhandledType = errors.New("Unhandled type")
	ErrTypeMismatch  = errors.New("Type mismatch")
	ErrMissingSemi   = errors.New("Missing semicolon")
	ErrMissingBrace  = errors.New("Missing brace")
	ErrBadObject     = errors.New("Object with unhandled key types")
)

func unmarshalInt(bs []byte, term byte) (int64, []byte, error) {
	if len(bs) < 2 {
		return 0, nil, EOF
	}

	var end int = len(bs)

	for i, b := range bs {
		if b == term {
			end = i
			break
		}
	}
	if end >= len(bs) {
		return 0, nil, EOF
	}
	if bs[end] != term {
		return 0, nil, ErrMissingSemi
	}

	val, er := strconv.ParseInt(string(bs[:end]), 10, 64)
	if er != nil {
		return 0, nil, er
	}

	return val, bs[end+1:], nil
}

func unmarshalFloat(bs []byte) (float64, []byte, error) {
	if len(bs) < 2 {
		return 0, nil, EOF
	}

	var end int = len(bs)

	for i, b := range bs {
		if b == ';' {
			end = i
		}
	}
	if end >= len(bs) {
		return 0, nil, EOF
	}
	if bs[end] != ';' {
		return 0, nil, ErrMissingSemi
	}

	val, er := strconv.ParseFloat(string(bs[:end]), 64)
	if er != nil {
		return 0, nil, er
	}

	return val, bs[end+1:], nil
}

func unmarshalBool(bs []byte) (bool, []byte, error) {
	if len(bs) < 2 {
		return false, nil, EOF
	}

	val := (bs[0] == '1')

	if bs[1] != ';' {
		return false, nil, ErrMissingSemi
	}

	return val, bs[2:], nil
}

func unmarshalString(bs []byte) (string, []byte, error) {
	length, bs, er := unmarshalInt(bs, ':')
	if er != nil {
		return "", nil, er
	}
	if int64(len(bs)) < length+3 {
		return "", nil, EOF
	}

	return string(bs[1 : 1+length]), bs[3+length:], nil
}

func unmarshalArray(bs []byte) (interface{}, []byte, error) {
	length, bs, er := unmarshalInt(bs, ':')
	if er != nil {
		return "", nil, er
	}
	if len(bs) < 2 || bs[0] != '{' {
		return nil, nil, ErrMissingBrace
	}
	bs = bs[1:]

	var (
		obj                 = map[interface{}]interface{}{}
		key                 interface{}
		val                 interface{}
		onlyIntKeys         bool = true
		onlyConvertableKeys bool = true
	)

	for i := int64(0); i < length; i += 1 {
		key, bs, er = unmarshalValue(bs)
		if er != nil {
			return nil, nil, er
		}

		if _, ok := key.(int64); !ok {
			onlyIntKeys = false
		} else if _, ok := key.(string); !ok {
			onlyConvertableKeys = false
		}

		val, bs, er = unmarshalValue(bs)
		if er != nil {
			return nil, nil, er
		}

		obj[key] = val
	}
	if len(bs) < 1 {
		return nil, nil, EOF
	}
	if bs[0] != '}' {
		return nil, nil, ErrMissingBrace
	}

	if !onlyIntKeys && !onlyConvertableKeys {
		return nil, nil, ErrBadObject
	}

	if onlyIntKeys {
		maxKey := int64(0)

		for k := range obj {
			ki := k.(int64)
			if ki > maxKey {
				maxKey = ki
			}
		}

		out := make([]interface{}, maxKey+1)

		for k, v := range obj {
			out[k.(int64)] = v
		}

		return out, bs[1:], nil
	}

	// onlyConvertableKeys
	out := map[string]interface{}{}

	for k, v := range obj {
		if ks, ok := k.(string); ok {
			out[ks] = v
		} else if ki, ok := k.(int64); ok {
			out[strconv.FormatInt(ki, 10)] = v
		}
	}

	return out, bs[1:], nil
}

func unmarshalValue(bs []byte) (interface{}, []byte, error) {
	if len(bs) < 2 {
		return nil, nil, EOF
	}

	switch bs[0] {
	case 'i':
		return unmarshalInt(bs[2:], ';')
	case 'd':
		return unmarshalFloat(bs[2:])
	case 'b':
		return unmarshalBool(bs[2:])
	case 's':
		return unmarshalString(bs[2:])
	case 'a':
		return unmarshalArray(bs[2:])
	default:
		return nil, nil, ErrUnhandledType
	}
}

func Unmarshal(bs []byte, v interface{}) error {
	val, _, er := unmarshalValue(bs)
	if er != nil {
		return er
	}

	valv := reflect.ValueOf(val)
	out := reflect.Indirect(reflect.ValueOf(v))

	if !out.CanSet() || !valv.Type().ConvertibleTo(out.Type()) {
		return ErrTypeMismatch
	}

	out.Set(valv.Convert(out.Type()))

	return nil
}
