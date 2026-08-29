package utils

import (
	"fmt"
	"reflect"
	"strconv"
)

// ConvertBytes converts from []byte to native data types: string, int, float, bool, uint.
func ConvertBytes[T any](data []byte) (T, error) {
	var zero T
	targetType := reflect.TypeOf(zero)

	// If the type is a pointer, get the pointed type
	if targetType.Kind() == reflect.Pointer {
		targetType = targetType.Elem()
	}

	ptr := reflect.New(targetType)
	val := ptr.Elem()

	err := assignValue(val, data)
	if err != nil {
		return zero, err
	}

	if reflect.TypeOf(zero).Kind() == reflect.Pointer {
		return ptr.Interface().(T), nil
	}
	return val.Interface().(T), nil
}

// ConvertirBytes is an alias for ConvertBytes
func ConvertirBytes[T any](data []byte) (T, error) {
	return ConvertBytes[T](data)
}

func assignValue(destValue reflect.Value, data []byte) error {
	switch destValue.Kind() {
	case reflect.String:
		str := string(data)
		destValue.SetString(str)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		val, err := strconv.ParseInt(string(data), 10, destValue.Type().Bits())
		if err != nil {
			return err
		}
		destValue.SetInt(val)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		val, err := strconv.ParseUint(string(data), 10, destValue.Type().Bits())
		if err != nil {
			return err
		}
		destValue.SetUint(val)
	case reflect.Float32, reflect.Float64:
		val, err := strconv.ParseFloat(string(data), destValue.Type().Bits())
		if err != nil {
			return err
		}
		destValue.SetFloat(val)
	case reflect.Bool:
		val, err := strconv.ParseBool(string(data))
		if err != nil {
			return err
		}
		destValue.SetBool(val)
	default:
		return fmt.Errorf("unsupported target type: %s", destValue.Type().String())
	}
	return nil
}
