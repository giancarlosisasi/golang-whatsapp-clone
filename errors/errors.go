package customerrors

import "errors"

var (
	ErrInvalidUUIDValue = errors.New("invalid UUID value")
)

const (
	CodeInternalError = "INTERNAL_ERROR"
)
