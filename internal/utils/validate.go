package utils

import (
	"reflect"

	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"
	"github.com/shopspring/decimal"
)

var validate = validator.New(validator.WithRequiredStructEnabled())

func init() {
	validate.RegisterCustomTypeFunc(decimalValue, decimal.Decimal{})

	err := validate.RegisterValidation("dgt", decimalGreaterThan)
	if err != nil {
		log.Panic().Err(err).Msg("Error registering validator dgt")
	}

	err = validate.RegisterValidation("dgtf", decimalGreaterThanField)
	if err != nil {
		log.Panic().Err(err).Msg("Error registering validator dgtf")
	}
}

// ValidateStruct validates a struct using go-playground/validator.
func ValidateStruct(s any) error {
	return validate.Struct(s)
}

// ValidarStruct alias for ValidateStruct
var ValidarStruct = ValidateStruct

func decimalValue(field reflect.Value) interface{} {
	valuer, ok := field.Interface().(decimal.Decimal)
	if ok {
		return valuer.String()
	}
	return nil
}

func decimalGreaterThan(fl validator.FieldLevel) bool {
	field := fl.Field()
	valueStr, ok := field.Interface().(string)
	if !ok {
		return false
	}

	value, err := decimal.NewFromString(valueStr)
	if err != nil {
		return false
	}

	baseValue, err := decimal.NewFromString(fl.Param())
	if err != nil {
		return false
	}

	return value.GreaterThan(baseValue)
}

func decimalGreaterThanField(fl validator.FieldLevel) bool {
	valueStr, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}

	value, err := decimal.NewFromString(valueStr)
	if err != nil {
		return false
	}

	baseField, _, _, ok := fl.GetStructFieldOK2()
	if !ok {
		return false
	}

	baseStr, ok := baseField.Interface().(string)
	if !ok {
		return false
	}

	baseValue, err := decimal.NewFromString(baseStr)
	if err != nil {
		return false
	}

	return value.GreaterThan(baseValue)
}
