package common

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	FieldErrors []FieldDisplayError `json:"fieldErrors,omitempty"`
	Message     string              `json:"error,omitempty"`
}

func (e *ErrorResponse) Error() string {
	result, err := json.Marshal(e)
	if err != nil {
		panic(err.Error())
	}
	return string(result)
}

// Bind JSON into the given type
//
// If an error is returned, it will be an `ErrorResponse` containing
// developer-friendly error messages. This can be returned directly
// to the client.
func BindJson[T interface{}](c *gin.Context) (*T, error) {
	var result T
	if c.Request.ContentLength == 0 {
		return nil, &ErrorResponse{Message: "Request body must not be empty"}
	}

	if err := c.ShouldBindJSON(&result); err != nil {
		displayErrors := GetFieldDisplayErrors(err)
		if displayErrors != nil {
			return nil, &ErrorResponse{FieldErrors: displayErrors}
		}
		return nil, &ErrorResponse{Message: "Invalid JSON"}
	}
	return &result, nil
}
