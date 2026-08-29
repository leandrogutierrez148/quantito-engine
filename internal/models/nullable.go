package models

import (
	"errors"
	"fmt"
	"quantito-engine/internal/utils"

	"github.com/LukaGiorgadze/gonull"
)

// Nullable models nullable columns in the database and provides safe access methods to their values.
type Nullable[T any] struct {
	gonull.Nullable[T]
}

// NotNull creates a non-null Nullable with the received value.
func NotNull[T any](value T) Nullable[T] {
	return Nullable[T]{Nullable: gonull.NewNullable(value)}
}

// Null creates a Nullable with no value (null).
func Null[T any]() Nullable[T] {
	return Nullable[T]{Nullable: gonull.Nullable[T]{Valid: false, Present: true}}
}

// Set sets the value of the Nullable.
func (n *Nullable[T]) Set(val T) {
	n.Nullable.Val = val
	n.Nullable.Valid = true
}

func (n Nullable[T]) IsNull() bool {
	return !n.IsSet()
}

func (n Nullable[T]) IsSet() bool {
	return n.Valid
}

// Get returns the value or an error in case of NULL.
func (n Nullable[T]) Get() (T, error) {
	if n.IsNull() {
		return n.Val, errors.New("value is null")
	}
	return n.Val, nil
}

// GetOrElse returns the value or a default one provided as parameter.
func (n Nullable[T]) GetOrElse(def T) T {
	if n.IsNull() {
		return def
	}
	return n.Val
}

func (n Nullable[T]) String() string {
	if n.IsNull() {
		return "null"
	}
	return fmt.Sprint(n.Val)
}

func (n *Nullable[T]) Scan(value any) error {
	err := n.Nullable.Scan(value)
	if err == nil {
		return nil
	}

	bytes, isBytes := value.([]byte)
	if !isBytes {
		return err
	}

	n.Val, err = utils.ConvertBytes[T](bytes)
	if err != nil {
		return err
	}

	n.Valid = true
	return nil
}
