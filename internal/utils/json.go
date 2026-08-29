package utils

import "encoding/json"

// UnmarshalJson unmarshals JSON into a new instance of the type parameter T.
// Useful to execute alongside `db.QueryJson`.
//
// Example:
//
//	utils.UnmarshalJson[[]*MyStruct](jsonBytes)
func UnmarshalJson[T any](source []byte) (T, error) {
	var target T

	if err := json.Unmarshal(source, &target); err != nil {
		return target, err
	}

	return target, nil
}
