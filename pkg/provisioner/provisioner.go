package provisioner

import (
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/kcm"
)

// Provisioner is the interface for an infrastructure provisioner.
type Provisioner interface {
	// Provision applies changes to the infrastructure. It should
	// automatically create or update a kubernetes cluster.
	Provision() error

	// Destroy performs all actions needed to destroy the underlying
	// cluster infrastructure.
	Destroy() error
}

// Reconciler can reconcile infrastucture status.
type Reconciler interface {
	// Reconcile retrieves the current state of the infrastructure and
	// should log potential changes without actually applying them.
	Reconcile() error
}

// Outputter can output kcm.Values.
type Outputter interface {
	// Output obtains output values from the infrastructure provisioner. These
	// values are made available during kubernetes manifest renderering.
	Output() (kcm.Values, error)
}

// Options are made available to infrastructure provisioners.
type Options struct {
	Terraform TerraformOptions `json:"terraform" yaml:"terraform"`
}

// TerraformOptions configure the terraform provisioner.
type TerraformOptions struct {
	Parallelism int `json:"parallelism" yaml:"parallelism"`
}
