package utils

import (
	"github.com/stretchr/testify/mock"
)

// MockWriteToPtr helps mock functions that write output to a pointer argument (e.g. db.QueryAll or json.Unmarshal).
func MockWriteToPtr[T any](argPos int, val T) func(mock.Arguments) {
	return func(args mock.Arguments) {
		arg := args.Get(argPos).(*T)
		*arg = val
	}
}

// MockEscribirEnPuntero is an alias for MockWriteToPtr
func MockEscribirEnPuntero[T any](argPos int, val T) func(mock.Arguments) {
	return MockWriteToPtr[T](argPos, val)
}
