package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestMockWriteToPtrScalar(t *testing.T) {
	var i int
	args := mock.Arguments{&i}
	MockWriteToPtr(0, 5)(args)
	assert.Equal(t, 5, i)
}

func TestMockWriteToPtrStruct(t *testing.T) {
	type A struct {
		Value int
	}

	dest := A{}
	args := mock.Arguments{&dest}

	val := A{
		Value: 42,
	}

	MockWriteToPtr(0, val)(args)

	assert.Equal(t, val, dest)
}
