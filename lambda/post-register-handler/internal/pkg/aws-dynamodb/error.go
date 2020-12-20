package awsdynamodb

import (
	"net/http"
	customerrors "post-register-handler/internal/custom-errors"
)

var (
	ErrInternalServer *customerrors.CustomError = &customerrors.CustomError{
		HTTPStatusCode: http.StatusInternalServerError,
		Message:        "Internal server error",
	}
)
