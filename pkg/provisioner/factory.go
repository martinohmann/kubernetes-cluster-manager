package provisioner

import (
	"reflect"

	"github.com/pkg/errors"
)

// Factory defines a factory func to create an infrastructure provisioner.
type Factory func(*Options) (Provisioner, error)

var (
	provisioners = make(map[string]Factory)
)

// Register registers a factory for an infrastructure provisioner with given
// name.
func Register(name string, factory Factory) {
	provisioners[name] = factory
}

// Create creates an infrastructure provisioner.
func Create(name string, o *Options) (Provisioner, error) {
	if factory, ok := provisioners[name]; ok {
		return factory(o)
	}

	return nil, errors.Errorf(
		"unsupported provisioner %q. Available provisioners: %s",
		name,
		reflect.ValueOf(provisioners).MapKeys(),
	)
}
