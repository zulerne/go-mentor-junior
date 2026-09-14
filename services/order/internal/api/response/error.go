package response

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/zulerne/go-mentor-junior/order/internal/domain"
)

type Error struct {
	ErrorData ErrorData `json:"error"`
}

type ErrorData struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details any    `json:"details"`
}

func NewError(err *domain.Error) Error {
	return Error{
		ErrorData: ErrorData{
			Code:    err.Code,
			Message: err.Message,
			Details: err.Details,
		},
	}
}

func NewBaseError(msg string) Error {
	return Error{
		ErrorData: ErrorData{
			Code:    domain.BaseErrorCode,
			Message: msg,
		},
	}
}

// NewValidationError creates a new validation error from the given validator.ValidationErrors ().
func NewValidationError(errs validator.ValidationErrors) Error {
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
			Code:    domain.ValidationErrorCode,
			Message: strings.Join(msgs, ", "),
		},
	}
}
