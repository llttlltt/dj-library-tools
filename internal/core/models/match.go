package models

// MetadataMatch pairs a source track with a target track for reconciliation.
type MetadataMatch struct {
	Source Track
	Target Track
}
