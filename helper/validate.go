package helper

import (
	"bookstore/model"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func ValidateUser(user model.User) []*model.ErrorResponse {
	var errors []*model.ErrorResponse
	err := validate.Struct(user)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element model.ErrorResponse
			element.FailedField = err.StructNamespace()
			element.Tag = err.Tag()
			element.Value = err.Param()
			errors = append(errors, &element)
		}
	}
	return errors
}
