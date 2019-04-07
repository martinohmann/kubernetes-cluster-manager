package infra

import "github.com/martinohmann/cluster-manager/pkg/api"

// Manager is the interface for a cloud infrastructure manager.
type Manager interface {
	// Apply will changes to the infrastructure. It will automatically create
	// or update a kubernetes cluster.
	Apply() error

	// GetOutput obtains output variables from the infrastructure manager.
	// These variables are made available during kubernetes manifest
	// renderering.
	GetOutput() (*api.InfraOutput, error)

	// Destroy performs all actions needed to destroy a kubernetes cluster.
	Destroy() error
}
