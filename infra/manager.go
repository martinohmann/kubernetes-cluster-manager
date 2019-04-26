package infra

import (
	"reflect"

	"github.com/martinohmann/kubernetes-cluster-manager/pkg/command"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/kcm"
	"github.com/pkg/errors"
)

// Manager is the interface for a cloud infrastructure manager.
type Manager interface {
	// Apply will apply changes to the infrastructure. It will automatically
	// create or update a kubernetes cluster.
	Apply() error

	// Plan will plan changes to the infrastructure without actually applying
	// them.
	Plan() error

	// GetValues obtains output values from the infrastructure manager.
	// These values are made available during kubernetes manifest
	// renderering.
	GetValues() (kcm.Values, error)

	// Destroy performs all actions needed to destroy a kubernetes cluster.
	Destroy() error
}

type ManagerOptions struct {
	Terraform TerraformOptions `json:"terraform" yaml:"terraform"`
}

// ManagerFactory defines a factory func to create an infrastructure manager.
type ManagerFactory func(*ManagerOptions, command.Executor) (Manager, error)

var managers = make(map[string]ManagerFactory)

// RegisterManager registers a factory for an infrastructure manager with given
// name.
func RegisterManager(name string, factory ManagerFactory) {
	managers[name] = factory
}

// CreateManager creates an infrastructure manager.
func CreateManager(name string, o *ManagerOptions, executor command.Executor) (Manager, error) {
	if factory, ok := managers[name]; ok {
		return factory(o, executor)
	}

	return nil, errors.Errorf(
		"unsupported infrastructure manager %q. Available managers: %s",
		name,
		reflect.ValueOf(managers).MapKeys(),
	)
}
