package provisioner

import (
	"reflect"

	"github.com/martinohmann/kubernetes-cluster-manager/pkg/command"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/kcm"
	"github.com/pkg/errors"
)

// Factory defines a factory func to create an infrastructure provisioner.
type Factory func(*kcm.ProvisionerOptions, command.Executor) (kcm.Provisioner, error)

var (
	provisioners = make(map[string]Factory)
)

// Register registers a factory for an infrastructure provisioner with given
// name.
func Register(name string, factory Factory) {
	provisioners[name] = factory
}

// Create creates an infrastructure provisioner.
func Create(name string, o *kcm.ProvisionerOptions, executor command.Executor) (kcm.Provisioner, error) {
	if factory, ok := provisioners[name]; ok {
		return factory(o, executor)
	}

	return nil, errors.Errorf(
		"unsupported provisioner %q. Available provisioners: %s",
		name,
		reflect.ValueOf(provisioners).MapKeys(),
	)
}
