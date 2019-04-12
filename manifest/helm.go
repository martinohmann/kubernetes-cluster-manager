package manifest

import (
	"github.com/martinohmann/kubernetes-cluster-manager/manifest/helm"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/api"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/command"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/config"
)

// HelmRenderer uses helm to render kubernetes manifests.
type HelmRenderer struct {
	cfg      *config.HelmConfig
	executor command.Executor
}

// NewHelmRenderer creates a new helm manifest renderer with given config
// and command executor.
func NewHelmRenderer(cfg *config.HelmConfig, executor command.Executor) *HelmRenderer {
	return &HelmRenderer{
		cfg:      cfg,
		executor: executor,
	}
}

// RenderManifest implements Renderer.
func (r *HelmRenderer) RenderManifest(v api.Values) (api.Manifest, error) {
	chart := helm.NewChart(r.cfg.Chart, r.executor)

	buf, err := chart.Render(v)
	if err != nil {
		return nil, err
	}

	return api.Manifest(buf), nil
}
