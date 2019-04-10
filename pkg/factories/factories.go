package factories

import (
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/api"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/command"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/config"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/kubernetes/helm"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/terraform"
	"github.com/pkg/errors"
)

// CreateManifestRenderer creates a manifest renderer based on the config.
func CreateManifestRenderer(cfg *config.Config, executor command.Executor) (api.ManifestRenderer, error) {
	switch cfg.ManifestRenderer {
	case "helm":
		return helm.NewManifestRenderer(cfg, executor), nil
	default:
		return nil, errors.Errorf("unsupported manifest renderer: %s", cfg.ManifestRenderer)
	}
}

// CreateInfraManager creates an infrastructure manager based on the config.
func CreateInfraManager(cfg *config.Config, executor command.Executor) (api.InfraManager, error) {
	switch cfg.InfraManager {
	case "terraform":
		return terraform.NewInfraManager(cfg, executor), nil
	default:
		return nil, errors.Errorf("unsupported infrastructure manager: %s", cfg.InfraManager)
	}
}
