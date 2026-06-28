package djerr

import (
	"errors"
	"fmt"
)

// Sentinel errors for simple checks
var (
	ErrReadOnly            = errors.New("provider is read-only")
	ErrUnsupportedResource = errors.New("resource type not supported")
	ErrInvalidParent       = errors.New("invalid parent for this group type")
)

// DomainError represents a structured error with metadata.
type DomainError struct {
	Code    ErrorCode
	Message string
	Wrapped error
}

func (e *DomainError) Error() string {
	if e.Wrapped != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Wrapped)
	}
	return e.Message
}

func (e *DomainError) Unwrap() error {
	return e.Wrapped
}

type ErrorCode string

const (
	CodeNotFound            ErrorCode = "NOT_FOUND"
	CodeAlreadyExists       ErrorCode = "ALREADY_EXISTS"
	CodeReadOnly            ErrorCode = "READ_ONLY"
	CodeUnsupported         ErrorCode = "UNSUPPORTED"
	CodeConstraintViolation ErrorCode = "CONSTRAINT_VIOLATION"
	CodeInternal            ErrorCode = "INTERNAL"
)

// Helper constructors
func NewNotFound(resource string, id string) error {
	return &DomainError{
		Code:    CodeNotFound,
		Message: fmt.Sprintf("%s %q not found", resource, id),
	}
}

func NewReadOnly(provider string) error {
	return &DomainError{
		Code:    CodeReadOnly,
		Message: fmt.Sprintf("provider %q is read-only", provider),
	}
}

func NewConstraint(msg string) error {
	return &DomainError{
		Code:    CodeConstraintViolation,
		Message: msg,
	}
}
