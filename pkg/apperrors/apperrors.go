package apperrors

import "errors"

var (
	ErrInvalidSource   = errors.New("invalid source body")
	ErrSourceNotFound  = errors.New("source not found")
	ErrClaimNotFound   = errors.New("claim not found")
	ErrValidationLogic = errors.New("validation logic error")
)
