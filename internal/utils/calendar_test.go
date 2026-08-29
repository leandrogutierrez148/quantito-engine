package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsLeapYear(t *testing.T) {
	t.Run("Divisible by 4 but not by 100", func(t *testing.T) {
		assert.True(t, IsLeapYear(2024))
	})

	t.Run("Divisible by 400", func(t *testing.T) {
		assert.True(t, IsLeapYear(2000))
	})

	t.Run("Divisible by 100 but not by 400", func(t *testing.T) {
		assert.False(t, IsLeapYear(2100))
	})

	t.Run("Not divisible by 4", func(t *testing.T) {
		assert.False(t, IsLeapYear(2023))
	})
}
