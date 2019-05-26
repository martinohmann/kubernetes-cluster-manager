package provisioner

import (
	"context"
)

// Provisioner is the interface for an infrastructure provisioner.
type Provisioner interface {
	// Provision applies changes to the infrastructure. It should
	// automatically create or update a kubernetes cluster.
	Provision(context.Context) error

	// Destroy performs all actions needed to destroy the underlying
	// cluster infrastructure.
	Destroy(context.Context) error
}

// Reconciler can reconcile infrastucture status.
type Reconciler interface {
	// Reconcile retrieves the current state of the infrastructure and
	// should log potential changes without actually applying them.
	Reconcile(context.Context) error
}

// Outputter can output values that are made available while rendering
// templates and may also contain kubernetes credentials.
type Outputter interface {
	// Output obtains output values from the infrastructure provisioner. These
	// values are made available during kubernetes manifest renderering.
	Output(context.Context) (map[string]interface{}, error)
}

// Options are made available to infrastructure provisioners.
type Options struct {
	Parallelism int `json:"parallelism,omitempty" yaml:"parallelism,omitempty"`
}
