package handler

import "github.com/go-playground/validator/v10"

var validate = validator.New()

func ValidateStruct(value any) error {
	return validate.Struct(value)
}
