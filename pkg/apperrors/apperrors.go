package apperrors

import "errors"

var (
	InvalidSource   = errors.New("invalid source body")
	ValidationLogic = errors.New("validation logic error")
)
