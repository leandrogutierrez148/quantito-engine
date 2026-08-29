package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
)

type JSONColumn[T any] struct {
	Val T
}

// Get returns a pointer to the contained value
func (j JSONColumn[T]) Get() *T {
	return &j.Val
}

// Value implements driver.Valuer interface
func (j JSONColumn[T]) Value() (driver.Value, error) {
	return json.Marshal(j.Val)
}

// Scan implements sql.Scanner interface
func (j *JSONColumn[T]) Scan(src any) error {
	if src == nil {
		return errors.New("cannot scan null value into JSONColumn")
	}

	var data []byte
	switch s := src.(type) {
	case string:
		data = []byte(s)
	case []byte:
		data = s
	default:
		return fmt.Errorf("incompatible type for JSONColumn: %T", src)
	}

	return json.Unmarshal(data, &j.Val)
}

func NewJSONColumn[T any](val T) JSONColumn[T] {
	return JSONColumn[T]{Val: val}
}

func NewJsonColumn[T any](val T) JSONColumn[T] {
	return JSONColumn[T]{Val: val}
}
