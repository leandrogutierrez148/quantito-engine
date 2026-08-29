package models

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type fetcher struct {
	mock.Mock
}

func (f *fetcher) fetch() (string, error) {
	args := f.Called()
	return args.String(0), args.Error(1)
}

func TestNewLazy(t *testing.T) {
	t.Run("Initializes correctly without calling fetcher", func(t *testing.T) {
		f := fetcher{}
		l := NewLazy(func() (string, error) {
			return f.fetch()
		})

		assert.Equal(t, "", *l.value)
		assert.Nil(t, l.err)
		assert.NotNil(t, l.fetcher)
		assert.NotNil(t, l.once)

		f.AssertNotCalled(t, "fetch")
	})

	t.Run("Does not copy once by value", func(t *testing.T) {
		f := fetcher{}
		var l1, l2 Lazy[string]

		l1 = NewLazy(func() (string, error) {
			return f.fetch()
		})

		func(lazyCopy Lazy[string]) {
			l2 = lazyCopy
		}(l1)

		assert.Same(t, l1.once, l2.once)
		assert.NotSame(t, l1, l2)
	})

	t.Run("Keeps reference to internal value across value copy and avoids duplicate calls", func(t *testing.T) {
		f := fetcher{}
		f.On("fetch").Return("valor", nil)

		var l1, l2 Lazy[string]

		l1 = NewLazy(func() (string, error) {
			return f.fetch()
		})

		func(lazyCopy Lazy[string]) {
			l2 = lazyCopy
		}(l1)

		val1, err1 := l1.Value()
		val2, err2 := l2.Value()

		f.AssertNumberOfCalls(t, "fetch", 1)
		assert.Same(t, l1.value, l2.value)
		assert.NotSame(t, l1, l2)
		assert.Equal(t, "valor", val1)
		assert.Equal(t, "valor", val2)
		assert.NoError(t, err1)
		assert.NoError(t, err2)
	})
}

func TestLazyValue(t *testing.T) {
	t.Run("Calls fetcher only once and caches result for multiple calls", func(t *testing.T) {
		f := fetcher{}
		l := NewLazy(func() (string, error) {
			return f.fetch()
		})

		f.On("fetch").Return("valor", nil)

		val1, err1 := l.Value()
		val2, err2 := l.Value()

		f.AssertNumberOfCalls(t, "fetch", 1)
		assert.Equal(t, "valor", val1)
		assert.Equal(t, "valor", val2)
		assert.NoError(t, err1)
		assert.NoError(t, err2)
	})

	t.Run("Returns error on multiple calls", func(t *testing.T) {
		f := fetcher{}
		l := NewLazy(func() (string, error) {
			return f.fetch()
		})

		f.On("fetch").Return("", errors.New("ERROR"))

		_, err1 := l.Value()
		_, err2 := l.Value()

		f.AssertNumberOfCalls(t, "fetch", 1)

		assert.Error(t, err1)
		assert.Error(t, err2)
	})

	t.Run("Returns error when created without NewLazy", func(t *testing.T) {
		l := Lazy[string]{}
		_, err := l.Value()
		assert.Error(t, err)
	})
}
