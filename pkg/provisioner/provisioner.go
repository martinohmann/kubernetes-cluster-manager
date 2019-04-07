package provisioner

import (
	"github.com/martinohmann/cluster-manager/pkg/command"
	"github.com/martinohmann/cluster-manager/pkg/config"
	"github.com/martinohmann/cluster-manager/pkg/infra"
	"github.com/martinohmann/cluster-manager/pkg/manifest"
)

type Provisioner struct {
	infraManager     infra.Manager
	manifestRenderer manifest.Renderer
	executor         command.Executor
}

func NewClusterProvisioner(infraManager infra.Manager, manifestRenderer manifest.Renderer, executor command.Executor) *Provisioner {
	return &Provisioner{
		infraManager:     infraManager,
		manifestRenderer: manifestRenderer,
		executor:         executor,
	}
}

func (p *Provisioner) Provision(cfg *config.Config) error {
	if !cfg.OnlyManifest {
		if err := p.infraManager.Apply(); err != nil {
			return err
		}
	}

	output, err := p.infraManager.GetOutput()
	if err != nil {
		return err
	}

	manifest, err := p.manifestRenderer.RenderManifest(output)
	if err != nil {
		return err
	}

	kubectl := NewKubectl(cfg, p.executor)

	deletions, err := loadDeletions(cfg.Deletions)
	if err != nil {
		return err
	}

	if !cfg.DryRun {
		defer saveDeletions(cfg.Deletions, deletions)
	}

	if err := processResourceDeletions(kubectl, deletions.PreApply); err != nil {
		return err
	}

	if err := kubectl.ApplyManifest(manifest); err != nil {
		return err
	}

	if err := processResourceDeletions(kubectl, deletions.PostApply); err != nil {
		return err
	}

	return nil
}

func (p *Provisioner) Destroy(cfg *config.Config) error {
	output, err := p.infraManager.GetOutput()
	if err != nil {
		return err
	}

	manifest, err := p.manifestRenderer.RenderManifest(output)
	if err != nil {
		return err
	}

	kubectl := NewKubectl(cfg, p.executor)

	if err := kubectl.DeleteManifest(manifest); err != nil {
		return err
	}

	deletions, err := loadDeletions(cfg.Deletions)
	if err != nil {
		return err
	}

	if !cfg.DryRun {
		defer saveDeletions(cfg.Deletions, deletions)
	}

	if err := processResourceDeletions(kubectl, deletions.PreDestroy); err != nil {
		return err
	}

	if !cfg.OnlyManifest {
		return p.infraManager.Destroy()
	}

	return nil
}
