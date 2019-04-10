package api

// InfraManager is the interface for a cloud infrastructure manager.
type InfraManager interface {
	// Apply will apply changes to the infrastructure. It will automatically
	// create or update a kubernetes cluster.
	Apply() error

	// Plan will plan changes to the infrastructure without acutally applying
	// them..
	Plan() error

	// GetValues obtains output values from the infrastructure manager.
	// These values are made available during kubernetes manifest
	// renderering.
	GetValues() (Values, error)

	// Destroy performs all actions needed to destroy a kubernetes cluster.
	Destroy() error
}
