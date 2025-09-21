package envparse

import (
	"errors"
	"fmt"
	"reflect"
	"runtime"
	"strconv"
)

func setSlice(f *reflect.Value, value string, parserName string) error {
	var parser func(string) (interface{}, error)

	if parserName == "default" {
		parser = defaultSliceParser
	} else {
		found, err := getParser(parserName)
		if err != nil {
			return err
		}
		parser = found
	}

	parsedValue, err := parser(value)

	if err != nil {
		return err
	}

	f.Set(reflect.ValueOf(parsedValue))
	return nil
}

func setBool(f *reflect.Value, value string) error {
	boolVal, err := strconv.ParseBool(value)
	if err != nil {
		return err
	}
	f.Set(reflect.ValueOf(boolVal))
	return nil
}

func setMap(f *reflect.Value, value string, parserName string) error {

	var parser func(string) (interface{}, error)

	if parserName == "default" {
		parser = defaultStructParser
	} else {
		found, err := getParser(parserName)
		if err != nil {
			return err
		}
		parser = found
	}

	parsedValue, err := parser(value)

	if err != nil {
		return err
	}

	f.Set(reflect.ValueOf(parsedValue))
	return nil
}

func setFloat[T float32 | float64](f *reflect.Value, value string, size int) error {
	determinedSize, err := determineSize(size)
	if err != nil {
		return err
	}

	floatVal, err := strconv.ParseFloat(value, determinedSize)

	if err != nil {
		return err
	}

	f.Set(reflect.ValueOf(T(floatVal)))
	return nil
}

func setInt[T int | int8 | int16 | int32 | int64](f *reflect.Value, value string, size int) error {

	determinedSize, err := determineSize(size)
	if err != nil {
		return err
	}

	intVal, err := strconv.ParseInt(value, 10, determinedSize)

	if err != nil {
		return err
	}

	f.Set(reflect.ValueOf(T(intVal)))
	return err
}

func setUint[T uint | uint16 | uint32 | uint64 | uint8](f *reflect.Value, value string, size int) error {
	determinedSize, err := determineSize(size)

	if err != nil {
		return err
	}

	intVal, err := strconv.ParseUint(value, 10, determinedSize)

	if err != nil {
		return err
	}

	f.Set(reflect.ValueOf(T(intVal)))
	return err
}

func setField[T interface{}](obj *T, field reflect.StructField, value string, parserName string) error {
	f := reflect.ValueOf(obj).Elem().FieldByName(field.Name)
	fieldType := f.Type()

	switch fieldType.Kind() {
	case reflect.Int:
		return setInt[int](&f, value, 0)
	case reflect.Int8:
		return setInt[int8](&f, value, 8)
	case reflect.Int16:
		return setInt[int16](&f, value, 16)
	case reflect.Int32:
		return setInt[int32](&f, value, 32)
	case reflect.Int64:
		return setInt[int64](&f, value, 64)
	case reflect.Uint:
		return setUint[uint](&f, value, 0)
	case reflect.Uint8:
		return setUint[uint8](&f, value, 8)
	case reflect.Uint16:
		return setUint[uint16](&f, value, 16)
	case reflect.Uint32:
		return setUint[uint32](&f, value, 32)
	case reflect.Uint64:
		return setUint[uint64](&f, value, 64)
	case reflect.Float32:
		return setFloat[float32](&f, value, 32)
	case reflect.Float64:
		return setFloat[float64](&f, value, 64)
	case reflect.Bool:
		return setBool(&f, value)
	case reflect.Map:
		return setMap(&f, value, parserName)
	case reflect.String:
		f.SetString(value)
		return nil
	case reflect.Slice:
		return setSlice(&f, value, parserName)
	default:
		return errors.New("no registed set method for type")
	}

}

func getArchBits() (int, error) {
	switch runtime.GOARCH {
	case "amd64", "arm64", "ppc64", "s390x":
		return 64, nil
	case "386", "arm", "ppc", "s390":
		return 32, nil
	default:
		return 0, fmt.Errorf("unknown architecture: %s", runtime.GOARCH)
	}
}

func determineSize(size int) (int, error) {

	if size == 0 {
		archBits, err := getArchBits()
		if err != nil {
			return 0, err
		}
		return archBits, nil
	}

	return size, nil
}
