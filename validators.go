package main

import (
	"encoding/json"
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
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
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
	if unmarshalError, ok := err.(*json.UnmarshalTypeError); ok {
		displayErrors := make([]FieldDisplayError, 1)
		displayErrors[0] = FieldDisplayError{
			Field: unmarshalError.Field,
			Error: fmt.Sprintf("Field %s must be of type %s, but recieved %s", unmarshalError.Field, unmarshalError.Type.String(), unmarshalError.Value),
		}
		return displayErrors
	}
	panic("Unknown error type")
}

var validationErrorMessages = map[string]string{
	"required": "%s is required",
	"date":     "%s must be a valid date in the format YYYY-MM-DD",
}
