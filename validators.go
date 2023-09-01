package main

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin/binding"
	validator "github.com/go-playground/validator/v10"
)

func init() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("date", date)
	}
}

var date validator.Func = func(fl validator.FieldLevel) bool {
	_, err := time.Parse("2006-01-02", fl.Field().String())
	return err == nil
}

type FieldDisplayError struct {
	Field string `json:"field"`
	Error string `json:"error"`
}

func getFieldDisplayErrors(err error) []FieldDisplayError {
	validationErrors := err.(validator.ValidationErrors)
	displayErrors := make([]FieldDisplayError, len(validationErrors))
	for i, fieldError := range validationErrors {
		var displayMsg string
		if msg, ok := validationErrorMessages[fieldError.Tag()]; ok {
			displayMsg = fmt.Sprintf(msg, fieldError.Field())
		} else {
			displayMsg = fmt.Sprintf("Field %s failed validation %s", fieldError.Error(), fieldError.Field())
		}
		displayErrors[i] = FieldDisplayError{
			Field: fieldError.Field(),
			Error: displayMsg,
		}
	}
	return displayErrors
}

var validationErrorMessages = map[string]string{
	"required": "%s is required",
	"date":     "%s must be a valid date in the format YYYY-MM-DD",
}
