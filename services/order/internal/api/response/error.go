package response

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

type Error struct {
	ErrorData ErrorData `json:"error"`
}

type ErrorData struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details any    `json:"details"`
}

func ValidationError(errs validator.ValidationErrors) Error {
	var msgs []string

	for _, err := range errs {
		switch err.ActualTag() {
		case "required":
			msgs = append(msgs, fmt.Sprintf("'%s' is required", err.Field()))
		default:
			msgs = append(msgs, fmt.Sprintf("'%s' is invalid", err.Field()))
		}
	}

	return Error{
		ErrorData: ErrorData{
			Code:    "VALIDATION_ERROR",
			Message: strings.Join(msgs, ", "),
		},
	}
}
