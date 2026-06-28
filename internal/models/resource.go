package models

// Resource is the interface for any item in a music library.
type Resource interface {
	GetID() string
	GetName() string
	GetKind() string
}
