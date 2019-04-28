package renderer

import (
	"reflect"

	"github.com/martinohmann/kubernetes-cluster-manager/pkg/kcm"
	"github.com/pkg/errors"
)

// Factory defines a factory func to create a manifest renderer.
type Factory func(*kcm.RendererOptions) (kcm.Renderer, error)

var (
	renderers = make(map[string]Factory)
)

// Register registers a factory for a manifest renderer with given
// name.
func Register(name string, factory Factory) {
	renderers[name] = factory
}

// Create creates a manifest renderer.
func Create(name string, o *kcm.RendererOptions) (kcm.Renderer, error) {
	if factory, ok := renderers[name]; ok {
		return factory(o)
	}

	return nil, errors.Errorf(
		"unsupported renderer %q. Available renderers: %s",
		name,
		reflect.ValueOf(renderers).MapKeys(),
	)
}
