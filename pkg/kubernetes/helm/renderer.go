package helm

import (
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/api"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/command"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/config"
)

type ManifestRenderer struct {
	cfg      *config.Config
	executor command.Executor
}

func NewManifestRenderer(cfg *config.Config, executor command.Executor) *ManifestRenderer {
	return &ManifestRenderer{
		cfg:      cfg,
		executor: executor,
	}
}

// RenderManifest implements api.ManifestRenderer.
func (r *ManifestRenderer) RenderManifest(v api.Values) (api.Manifest, error) {
	chart := NewChart(r.cfg.Helm.Chart, r.executor)

	buf, err := chart.Render(v)
	if err != nil {
		return nil, err
	}

	return api.Manifest(buf), nil
}
