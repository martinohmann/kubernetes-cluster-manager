package renderer

import (
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/kcm"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/kubernetes/helm"
)

func init() {
	Register("helm", func(o *kcm.RendererOptions) (kcm.Renderer, error) {
		return NewHelm(&o.Helm), nil
	})
}

// Helm uses helm to render kubernetes manifests.
type Helm struct {
	chart *helm.Chart
}

// NewHelm creates a new helm manifest renderer with given options.
func NewHelm(o *kcm.HelmOptions) *Helm {
	return &Helm{
		chart: helm.NewChart(o.Chart),
	}
}

// RenderManifest implements Renderer.
func (r *Helm) RenderManifest(v kcm.Values) (kcm.Manifest, error) {
	buf, err := r.chart.Render(v)
	if err != nil {
		return nil, err
	}

	return kcm.Manifest(buf), nil
}
