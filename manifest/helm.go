package manifest

import (
	"github.com/martinohmann/kubernetes-cluster-manager/manifest/helm"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/api"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/command"
)

func init() {
	RegisterRenderer("helm", func(o *RendererOptions, e command.Executor) (Renderer, error) {
		return NewHelmRenderer(&o.Helm, e), nil
	})
}

type HelmOptions struct {
	Chart string `json:"chart" yaml:"chart"`
}

// HelmRenderer uses helm to render kubernetes manifests.
type HelmRenderer struct {
	options  *HelmOptions
	executor command.Executor
}

// NewHelmRenderer creates a new helm manifest renderer with given options
// and command executor.
func NewHelmRenderer(o *HelmOptions, executor command.Executor) *HelmRenderer {
	return &HelmRenderer{
		options:  o,
		executor: executor,
	}
}

// RenderManifest implements Renderer.
func (r *HelmRenderer) RenderManifest(v api.Values) (api.Manifest, error) {
	chart := helm.NewChart(r.options.Chart, r.executor)

	buf, err := chart.Render(v)
	if err != nil {
		return nil, err
	}

	return api.Manifest(buf), nil
}
