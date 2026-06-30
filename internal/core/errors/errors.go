package errors

import (
	"errors"
	"fmt"
)

// Kind defines the category of a domain error.
type Kind int

const (
	KindUnknown Kind = iota
	KindReadOnly
	KindUnsupportedResource
	KindInvalidParent
	KindNotFound
	KindAlreadyExists
	KindConstraintViolation
	KindInternal
)

// Error represents a domain-level error with a kind and optional cause.
type Error struct {
	Kind  Kind
	Msg   string
	Cause error
}

func (e *Error) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %v", e.Msg, e.Cause)
	}
	return e.Msg
}

func (e *Error) Unwrap() error {
	return e.Cause
}

// KindOf resolves an error to its domain Kind.
func KindOf(err error) Kind {
	var de *Error
	if errors.As(err, &de) {
		return de.Kind
	}
	// Fallback to sentinel checks if any remain
	return KindUnknown
}

var (
	ErrReadOnly            = &Error{Kind: KindReadOnly, Msg: "provider is read-only"}
	ErrUnsupportedResource = &Error{Kind: KindUnsupportedResource, Msg: "resource type not supported"}
	ErrInvalidParent       = &Error{Kind: KindInvalidParent, Msg: "invalid parent for this group type"}
)

// Helper constructors
func NewNotFound(resource string, id string) error {
	return &Error{
		Kind: KindNotFound,
		Msg:  fmt.Sprintf("%s %q not found", resource, id),
	}
}

func NewReadOnly(provider string) error {
	return &Error{
		Kind: KindReadOnly,
		Msg:  fmt.Sprintf("provider %q is read-only", provider),
	}
}


