package awsdynamodb

import (
	customerrors "get-redirect-handler/internal/custom-errors"
	"net/http"
)

var (
	ErrResourceNotFound *customerrors.CustomError = &customerrors.CustomError{
		HTTPStatusCode: http.StatusInternalServerError,
		Message:        "dynamodb - ResourceNotFoundException",
	}
	ErrInternalServer *customerrors.CustomError = &customerrors.CustomError{
		HTTPStatusCode: http.StatusInternalServerError,
		Message:        "Internal server error",
	}
	ErrItemNotExist *customerrors.CustomError = &customerrors.CustomError{
		HTTPStatusCode: http.StatusNotFound,
		Message:        "shortenId doesn't exist",
	}
)
