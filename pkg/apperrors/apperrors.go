package apperrors

import "errors"

var (
	ErrInvalidSource   = errors.New("invalid source body")
	ErrValidationLogic = errors.New("validation logic error")
)
