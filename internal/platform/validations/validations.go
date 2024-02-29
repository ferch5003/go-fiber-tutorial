package validations

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"strings"
)

type (
	ErrorResponse struct {
		Error       bool
		FailedField string
		Tag         string
		Value       interface{}
	}

	XValidator struct {
		validator *validator.Validate
	}
)

func NewValidator() *XValidator {
	return &XValidator{
		validator: validator.New(),
	}
}

func (v XValidator) validate(data any) []ErrorResponse {
	var validationErrors []ErrorResponse

	errs := v.validator.Struct(data)
	if errs != nil {
		for _, err := range errs.(validator.ValidationErrors) {
			var elem ErrorResponse

			elem.FailedField = err.Field() // Export struct field name
			elem.Tag = err.Tag()           // Export struct tag
			elem.Value = err.Value()       // Export field value
			elem.Error = true

			validationErrors = append(validationErrors, elem)
		}
	}

	return validationErrors
}

func (v XValidator) GetValidations(data any) string {
	if errs := v.validate(data); len(errs) > 0 && errs[0].Error {
		errMsgs := make([]string, 0)

		for _, err := range errs {
			errMsgs = append(errMsgs, fmt.Sprintf(
				"[%s]: '%v' | Needs to implement '%s'",
				err.FailedField,
				err.Value,
				err.Tag,
			))
		}

		return strings.Join(errMsgs, " and ")
	}

	return ""
}
