package customerrors

import "errors"

var (
	ErrInvalidUUIDValue = errors.New("invalid UUID value")
	ErrResourceNotFound = errors.New("resource not found")
)

const (
	CodeInternalError    = "INTERNAL_ERROR"
	CodeResourceNotFound = "RESOURCE_NOT_FOUND"
)
