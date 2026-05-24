package apperrors

import "errors"

var (
	ErrInvalidSource            = errors.New("invalid source body")
	ErrInvalidClaim             = errors.New("invalid claim body")
	ErrInvalidProof             = errors.New("invalid proof body")
	ErrSourceNotFound           = errors.New("source not found")
	ErrClaimNotFound            = errors.New("claim not found")
	ErrProofNotFound            = errors.New("proof not found")
	ErrValidationLogic          = errors.New("validation logic error")
	ErrInvalidClaimVerification = errors.New("invalid claim verification body")
	ErrDuplicateSource          = errors.New("source already exists")
	ErrDuplicateClaim           = errors.New("claim already exists")
	ErrDuplicateProof           = errors.New("proof already exists")
	ErrSameClaimProofSource     = errors.New("claim and proofs are from the same source")
)
