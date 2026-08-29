package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConvertBytes(t *testing.T) {
	t.Run("string", func(t *testing.T) {
		bytes := []byte("Input content")
		val, err := ConvertBytes[string](bytes)
		assert.NoError(t, err)
		assert.Equal(t, "Input content", val)
	})

	t.Run("*string", func(t *testing.T) {
		bytes := []byte("Input content")
		val, err := ConvertBytes[*string](bytes)
		assert.NoError(t, err)
		assert.Equal(t, "Input content", *val)
	})

	t.Run("int", func(t *testing.T) {
		bytes := []byte("-42424242")
		val, err := ConvertBytes[int](bytes)
		assert.NoError(t, err)
		assert.Equal(t, -42424242, val)
	})

	t.Run("*int", func(t *testing.T) {
		bytes := []byte("-42424242")
		val, err := ConvertBytes[*int](bytes)
		assert.NoError(t, err)
		assert.Equal(t, -42424242, *val)
	})

	t.Run("float", func(t *testing.T) {
		bytes := []byte("3.1415926535")
		val, err := ConvertBytes[float64](bytes)
		assert.NoError(t, err)
		assert.Equal(t, 3.1415926535, val)
	})

	t.Run("*float", func(t *testing.T) {
		bytes := []byte("3.1415926535")
		val, err := ConvertBytes[*float64](bytes)
		assert.NoError(t, err)
		assert.Equal(t, 3.1415926535, *val)
	})

	t.Run("bool", func(t *testing.T) {
		bytes := []byte("true")
		val, err := ConvertBytes[bool](bytes)
		assert.NoError(t, err)
		assert.Equal(t, true, val)
	})

	t.Run("*bool", func(t *testing.T) {
		bytes := []byte("false")
		val, err := ConvertBytes[*bool](bytes)
		assert.NoError(t, err)
		assert.Equal(t, false, *val)
	})

	t.Run("uint", func(t *testing.T) {
		bytes := []byte("42424242")
		val, err := ConvertBytes[uint](bytes)
		assert.NoError(t, err)
		assert.Equal(t, uint(42424242), val)
	})

	t.Run("*uint", func(t *testing.T) {
		bytes := []byte("42424242")
		val, err := ConvertBytes[*uint](bytes)
		assert.NoError(t, err)
		assert.Equal(t, uint(42424242), *val)
	})

	t.Run("Unsupported data type", func(t *testing.T) {
		bytes := []byte("ABC")
		_, err := ConvertBytes[struct{ A string }](bytes)
		assert.Error(t, err)
	})
}
