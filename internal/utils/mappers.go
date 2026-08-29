package utils

import (
	"encoding/json"
)

// MapParams converts a map[string]any to the struct specified in the type parameter.
func MapParams[T any](m map[string]any) (*T, error) {
	var out T
	data, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(data, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// MapAndValidateParams converts a map[string]any to the struct specified in the type parameter
// and validates it against go-playground/validator struct tags.
func MapAndValidateParams[T any](m map[string]any) (*T, error) {
	out, err := MapParams[T](m)
	if err != nil {
		return nil, err
	}
	if err := ValidateStruct(out); err != nil {
		return nil, err
	}
	return out, nil
}

// Backward compatibility functions
func MapearParametros[T any](m map[string]any) (*T, error) {
	return MapParams[T](m)
}

func MapearYValidarParametros[T any](m map[string]any) (*T, error) {
	return MapAndValidateParams[T](m)
}
