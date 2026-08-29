package utils

import (
	"math"
	"testing"

	"github.com/go-playground/validator/v10"

	"github.com/stretchr/testify/assert"
)

func TestMapAndValidateParams(t *testing.T) {
	in := map[string]any{
		"StringValue": "mystring",
		"IntValue":    2,
		"FloatValue":  3.14,
		"BoolValue":   true,
	}

	type myStruct struct {
		StringValue string
		IntValue    int
		FloatValue  float32
		BoolValue   bool
	}

	expectedResult := &myStruct{
		StringValue: "mystring",
		IntValue:    2,
		FloatValue:  3.14,
		BoolValue:   true,
	}
	out, err := MapAndValidateParams[myStruct](in)

	assert.NoError(t, err)
	assert.Equal(t, expectedResult, out)
}

func TestMapAndValidateParamsErrorEncode(t *testing.T) {
	in := map[string]any{
		"invalidValue": math.Inf(1),
	}

	out, err := MapAndValidateParams[struct{ InvalidValue int }](in)

	assert.Error(t, err)
	assert.Nil(t, out)
}

func TestMapAndValidateParamsErrorDecode(t *testing.T) {
	in := map[string]any{
		"boolValue": "notABool",
	}

	out, err := MapAndValidateParams[struct{ BoolValue bool }](in)

	assert.Error(t, err)
	assert.Nil(t, out)
}

func TestMapAndValidateParamsValidationSuccess(t *testing.T) {
	type myStruct struct {
		StringValue string `validate:"required"`
		IntValue    int    `validate:"min=5"`
	}

	in := map[string]any{
		"StringValue": "value",
		"IntValue":    42,
	}

	_, err := MapAndValidateParams[myStruct](in)
	assert.NoError(t, err)
}

func TestMapAndValidateParamsValidationError(t *testing.T) {
	type myStruct struct {
		StringValue string `validate:"contains=x"`
		IntValue    int    `validate:"min=5"`
	}

	in := map[string]any{
		"StringValue": "value",
		"IntValue":    2,
	}

	_, err := MapAndValidateParams[myStruct](in)
	assert.Error(t, err)

	validationError, ok := err.(validator.ValidationErrors)
	assert.True(t, ok)
	assert.Len(t, validationError, 2)
}
