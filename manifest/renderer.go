package manifest

import (
	"reflect"

	"github.com/martinohmann/kubernetes-cluster-manager/pkg/api"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/command"
	"github.com/pkg/errors"
)

// Renderer is the interface for a kubernetes manifest renderer.
type Renderer interface {
	// RenderManifest renders a kubernetes manifest.
	RenderManifest(api.Values) (api.Manifest, error)
}

type RendererOptions struct {
	Helm HelmOptions `json:"helm" yaml:"helm"`
}

// RendererFactory defines a factory func to create a manifest renderer.
type RendererFactory func(*RendererOptions, command.Executor) (Renderer, error)

var renderers = make(map[string]RendererFactory)

// RegisterRenderer registers a factory for a manifest renderer with given
// name.
func RegisterRenderer(name string, factory RendererFactory) {
	renderers[name] = factory
}

// CreateRenderer creates a manifest renderer.
func CreateRenderer(name string, o *RendererOptions, executor command.Executor) (Renderer, error) {
	if factory, ok := renderers[name]; ok {
		return factory(o, executor)
	}

	return nil, errors.Errorf(
		"unsupported manifest renderer %q. Available renderers: %s",
		name,
		reflect.ValueOf(renderers).MapKeys(),
	)
}
