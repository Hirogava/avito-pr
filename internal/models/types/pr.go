// Package types defines types
package types

// PRStatus - PR status
type PRStatus string

const (
	// PRStatusOpen - PR открыт
	PRStatusOpen PRStatus = "open"
	// PRStatusMerged - PR закрыт
	PRStatusMerged PRStatus = "merged"
)
