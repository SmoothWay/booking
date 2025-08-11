package validator

import (
	"errors"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func Validate(data any) error {
	err := validate.Struct(data)
	if err == nil {
		return nil
	}

	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		return errors.New(err.Error())
	}

	var errMsg string
	for _, fieldErr := range validationErrors {
		label := fieldErr.Field()
		if labelTag := fieldErr.StructField(); labelTag != "" {
			label = labelTag
		}
		if tag := fieldErr.Param(); tag != "" {
			errMsg += label + " cannot be " + fieldErr.Tag() + " (param: " + tag + "); "
		} else {
			errMsg += label + " cannot be " + fieldErr.Tag() + "; "
		}
	}
	if errMsg == "" {
		errMsg = "validation failed"
	}
	return errors.New(errMsg)
}
