package models

import (
	"fmt"
	"reflect"
	"sync"
)

// Lazy[T] allows initializing an object of type T using the fetcher function when Value is called.
type Lazy[T any] struct {
	value   *T
	err     error
	once    *sync.Once
	fetcher func() (T, error)
}

// NewLazy instantiates a Lazy[T] with the fetcher function to be called on first access.
func NewLazy[T any](fetcher func() (T, error)) Lazy[T] {
	return Lazy[T]{
		value:   new(T),
		once:    &sync.Once{},
		fetcher: fetcher,
	}
}

// Value returns the wrapped object, computing it on first call and caching thereafter.
func (l *Lazy[T]) Value() (T, error) {
	var zero T

	if l.fetcher == nil {
		tipoT := reflect.TypeOf((*T)(nil)).Elem().Name()
		return zero, fmt.Errorf("Lazy[%s].fetcher is not defined", tipoT)
	}

	l.once.Do(func() {
		val, err := l.fetcher()
		*l.value, l.err = val, err
	})

	if l.err != nil {
		return zero, l.err
	}

	return *l.value, nil
}
