package manifest

import (
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/api"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/command"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/config"
	"github.com/pkg/errors"
)

// Renderer is the interface for a kubernetes manifest renderer.
type Renderer interface {
	// RenderManifest renders a kubernetes manifest.
	RenderManifest(api.Values) (api.Manifest, error)
}

// CreateRenderer creates a manifest renderer based on the config.
func CreateRenderer(cfg *config.Config, executor command.Executor) (Renderer, error) {
	switch cfg.ManifestRenderer {
	case "helm":
		return NewHelmRenderer(&cfg.Helm, executor), nil
	default:
		return nil, errors.Errorf("unsupported manifest renderer: %s", cfg.ManifestRenderer)
	}
}
