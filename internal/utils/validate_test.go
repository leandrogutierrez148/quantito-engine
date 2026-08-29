package utils

import (
	"reflect"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

type myStruct struct {
	StringValue string `validate:"required,contains=S"`
	IntValue    int    `validate:"required,min=10"`
}

func TestValidateStruct(t *testing.T) {
	t.Run("Success validate struct", validateStructSuccess)
	t.Run("Error validate struct", validateStructError)
}

func TestDecimalValue(t *testing.T) {
	t.Run("Success decimal value", successDecimalValue)
	t.Run("Error decimal value", errorDecimalValue)
}

func TestDecimalGreaterThanParam(t *testing.T) {
	t.Run("Success condition", successDecimalGreaterThan)
	t.Run("Condition failure", falseDecimalGreaterThan)
	t.Run("Error param data type", errorParamTypeDecimalGreaterThan)
	t.Run("Error value data type", errorValueTypeDecimalGreaterThan)
	t.Run("Error param not string", errorValueNotStringDecimalGreaterThan)
}

func TestDecimalGreaterThanField(t *testing.T) {
	t.Run("Condition met", successDecimalGreaterThanField)
	t.Run("Condition not met", falseDecimalGreaterThanField)
	t.Run("Error value not parseable to decimal", errorParseValueDecimalGreaterThanField)
	t.Run("Error base not parseable to decimal", errorParseBaseDecimalGreaterThanField)
	t.Run("Error value not string", errorValueNotStringDecimalGreaterThanField)
	t.Run("Error base not string", errorBaseNotStringDecimalGreaterThanField)
	t.Run("Error base does not exist", errorBaseDoesNotExistDecimalGreaterThanField)
}

func validateStructSuccess(t *testing.T) {
	s := myStruct{
		StringValue: "Si",
		IntValue:    42,
	}

	assert.NoError(t, ValidateStruct(s))
}

func validateStructError(t *testing.T) {
	s := myStruct{
		StringValue: "No",
		IntValue:    5,
	}

	err := ValidateStruct(s)
	assert.Error(t, err)

	validationErrors, ok := err.(validator.ValidationErrors)
	assert.True(t, ok)
	assert.Len(t, validationErrors, 2)
}

func successDecimalGreaterThan(t *testing.T) {
	type decimalStruct struct {
		Value1 decimal.Decimal `validate:"required,dgt=9.5"`
	}

	d := decimalStruct{
		Value1: decimal.NewFromInt(10),
	}

	assert.NoError(t, ValidateStruct(d), "test: success in `dgt` condition")
}

func falseDecimalGreaterThan(t *testing.T) {
	type decimalStruct struct {
		Value1 decimal.Decimal `validate:"required,dgt=11"`
	}

	d := decimalStruct{
		Value1: decimal.NewFromInt(10),
	}

	assert.Error(t, ValidateStruct(d), "test: failure in `dgt` condition")
}

func errorParamTypeDecimalGreaterThan(t *testing.T) {
	type decimalStruct struct {
		Value1 decimal.Decimal `validate:"required,dgt=ABC"`
	}

	d := decimalStruct{
		Value1: decimal.NewFromInt(10),
	}

	assert.Error(t, ValidateStruct(d), "test: data type error in `param` for `dgt` condition")
}

func errorValueTypeDecimalGreaterThan(t *testing.T) {
	type decimalStruct struct {
		Value1 any `validate:"required,dgt=10"`
	}

	d := decimalStruct{
		Value1: "ABC",
	}

	assert.Error(t, ValidateStruct(d), "test: data type error in `value` for `dgt` condition")
}

func errorValueNotStringDecimalGreaterThan(t *testing.T) {
	type decimalStruct struct {
		Value1 any `validate:"required,dgt=10"`
	}

	d := decimalStruct{
		Value1: 15,
	}

	assert.Error(t, ValidateStruct(d), "test: `value` cannot be cast to string in `dgt` condition")
}

func successDecimalGreaterThanField(t *testing.T) {
	type decimalStruct struct {
		Value1 decimal.Decimal `validate:"required"`
		Value2 decimal.Decimal `validate:"required,dgtf=Value1"`
	}

	val1, err := decimal.NewFromString("100.50")
	assert.NoError(t, err)
	val2, err := decimal.NewFromString("100.55")
	assert.NoError(t, err)

	d := decimalStruct{
		Value1: val1,
		Value2: val2,
	}

	assert.NoError(t, ValidateStruct(d), "test: success in `dgtf` condition")
}

func falseDecimalGreaterThanField(t *testing.T) {
	type decimalStruct struct {
		Value1 decimal.Decimal `validate:"required"`
		Value2 decimal.Decimal `validate:"required,dgtf=Value1"`
	}

	d := decimalStruct{
		Value1: decimal.NewFromInt(100),
		Value2: decimal.NewFromInt(99),
	}

	assert.Error(t, ValidateStruct(d), "test: failure in `dgtf` condition")
}

func errorParseValueDecimalGreaterThanField(t *testing.T) {
	type decimalStruct struct {
		Value1 decimal.Decimal `validate:"required"`
		Value2 any             `validate:"required,dgtf=Value1"`
	}

	d := decimalStruct{
		Value1: decimal.NewFromInt(10),
		Value2: "11.50.50",
	}

	assert.Error(t, ValidateStruct(d), "test: data type error in `value` for `dgtf` condition")
}

func errorParseBaseDecimalGreaterThanField(t *testing.T) {
	type decimalStruct struct {
		Value1 any             `validate:"required"`
		Value2 decimal.Decimal `validate:"required,dgtf=Value1"`
	}

	d := decimalStruct{
		Value1: "11.50.50",
		Value2: decimal.NewFromInt(1000),
	}

	assert.Error(t, ValidateStruct(d), "test: data type error in `value` for `dgtf` condition")
}

func errorBaseNotStringDecimalGreaterThanField(t *testing.T) {
	type decimalStruct struct {
		Value1 any             `validate:"required"`
		Value2 decimal.Decimal `validate:"required,dgtf=Value1"`
	}

	d := decimalStruct{
		Value1: 10,
		Value2: decimal.NewFromInt(99),
	}

	assert.Error(t, ValidateStruct(d), "test: data type error in `base` for `dgtf` condition")
}

func errorValueNotStringDecimalGreaterThanField(t *testing.T) {
	type decimalStruct struct {
		Value1 any `validate:"required"`
		Value2 any `validate:"required,dgtf=Value1"`
	}

	d := decimalStruct{
		Value1: decimal.NewFromInt(100),
		Value2: 101,
	}

	assert.Error(t, ValidateStruct(d), "test: `value` cannot be cast to string in `dgtf` condition")
}

func errorBaseDoesNotExistDecimalGreaterThanField(t *testing.T) {
	type decimalStruct struct {
		Value1 any `validate:"required"`
		Value2 any `validate:"required,dgtf=Value3"`
	}

	d := decimalStruct{
		Value1: decimal.NewFromInt(100),
		Value2: decimal.NewFromInt(101),
	}

	assert.Error(t, ValidateStruct(d), "test: `base` does not exist in `dgtf` condition")
}

func successDecimalValue(t *testing.T) {
	val := reflect.ValueOf(decimal.NewFromInt(120))
	assert.NotNil(t, decimalValue(val))
}

func errorDecimalValue(t *testing.T) {
	val := reflect.ValueOf(120)
	assert.Nil(t, decimalValue(val))
}
