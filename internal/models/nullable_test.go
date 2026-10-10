package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSet(t *testing.T) {
	v := Null[int]()

	v.Set(42)
	val, _ := v.Get()

	assert.False(t, v.IsNull())
	assert.Equal(t, 42, val)
}

func TestIsNull(t *testing.T) {
	notNull := NotNull(5)
	null := Null[int]()

	assert.False(t, notNull.IsNull())
	assert.True(t, null.IsNull())
}

func TestGet(t *testing.T) {
	val, err := NotNull(55).Get()

	assert.NoError(t, err)
	assert.Equal(t, 55, val)

	_, err = Null[int]().Get()

	assert.Error(t, err)
}

func TestGetOrElse(t *testing.T) {
	assert.Equal(t, 33, NotNull(33).GetOrElse(88))

	assert.Equal(t, 67, Null[int]().GetOrElse(67))
}

func TestString(t *testing.T) {
	notNull := NotNull(3.14)
	null := Null[float32]()

	assert.Equal(t, "3.14", notNull.String())
	assert.Equal(t, "null", null.String())
}

func TestScan(t *testing.T) {
	t.Run("Assigns value if same type is received", func(t *testing.T) {
		expected := 42
		n := Null[int]()
		err := n.Scan(expected)

		assert.NoError(t, err)
		assert.Equal(t, expected, n.Val)
		assert.True(t, n.Valid)
	})

	t.Run("Returns error on unsupported conversions", func(t *testing.T) {
		n := Null[int]()
		err := n.Scan(Experiment{Name: "Test"})

		assert.Error(t, err)
		assert.False(t, n.Valid)
	})

	t.Run("Converts []byte to string", func(t *testing.T) {
		expected := "Variable content"
		val := []byte(expected)

		n := Null[string]()
		err := n.Scan(val)

		assert.NoError(t, err)
		assert.Equal(t, expected, n.Val)
		assert.True(t, n.Valid)
	})

	t.Run("Converts []byte to int", func(t *testing.T) {
		val := []byte(`123456`)

		n := Null[int]()

		err := n.Scan(val)

		assert.NoError(t, err)
		assert.Equal(t, 123456, n.Val)
		assert.True(t, n.Valid)
	})

	t.Run("Fails when unsupported []byte conversion is attempted", func(t *testing.T) {
		val := []byte(`123456`)
		n := Null[Experiment]()

		err := n.Scan(val)

		assert.Error(t, err)
	})
}
