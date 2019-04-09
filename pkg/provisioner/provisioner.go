package provisioner

import (
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/api"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/command"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/config"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/git"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/kubernetes"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type Provisioner struct {
	infraManager     api.InfraManager
	manifestRenderer api.ManifestRenderer
	executor         command.Executor
}

func NewClusterProvisioner(infraManager api.InfraManager, manifestRenderer api.ManifestRenderer, executor command.Executor) *Provisioner {
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

	newValues, err := p.infraManager.GetValues()
	if err != nil {
		return err
	}

	values, err := loadValues(cfg.Values)
	if err != nil {
		return err
	}

	if err := values.Merge(newValues); err != nil {
		return err
	}

	valueBytes, err := yaml.Marshal(values)
	if err != nil {
		return err
	}

	err = p.finalizeChanges(cfg, cfg.Values, valueBytes)
	if err != nil {
		return err
	}

	manifest, err := p.manifestRenderer.RenderManifest(values)
	if err != nil {
		return err
	}

	updateCredentialsFromValues(cfg, values)

	kubectl := kubernetes.NewKubectl(cfg, p.executor)

	deletions, err := loadDeletions(cfg.Deletions)
	if err != nil {
		return err
	}

	err = p.finalizeChanges(cfg, cfg.Manifest, manifest)
	if err != nil {
		return err
	}

	defer p.finalizeDeletions(cfg, deletions)

	if err := processResourceDeletions(kubectl, deletions.PreApply); err != nil {
		return err
	}

	if err := kubectl.ApplyManifest(manifest); err != nil {
		return err
	}

	return processResourceDeletions(kubectl, deletions.PostApply)
}

func (p *Provisioner) Destroy(cfg *config.Config) error {
	values, err := loadValues(cfg.Values)
	if err != nil {
		return err
	}

	manifest, err := p.manifestRenderer.RenderManifest(values)
	if err != nil {
		return err
	}

	updateCredentialsFromValues(cfg, values)

	kubectl := kubernetes.NewKubectl(cfg, p.executor)

	if err := kubectl.DeleteManifest(manifest); err != nil {
		return err
	}

	deletions, err := loadDeletions(cfg.Deletions)
	if err != nil {
		return err
	}

	defer p.finalizeDeletions(cfg, deletions)

	if err := processResourceDeletions(kubectl, deletions.PreDestroy); err != nil {
		return err
	}

	if !cfg.OnlyManifest {
		return p.infraManager.Destroy()
	}

	return nil
}

func (p *Provisioner) finalizeChanges(cfg *config.Config, filename string, content []byte) error {
	changes, err := git.NewFileChanges(filename, content)
	if err != nil {
		return err
	}

	defer changes.Close()

	diff, err := changes.Diff()
	if err != nil {
		return err
	}

	if len(diff) > 0 {
		log.Infof("Changes to %s:\n%s", filename, diff)
	} else {
		log.Infof("No changes to %s", filename)
	}

	if cfg.DryRun {
		return nil
	}

	return changes.Apply()
}

func (p *Provisioner) finalizeDeletions(cfg *config.Config, deletions *api.Deletions) error {
	buf, err := yaml.Marshal(deletions.FilterPending())
	if err != nil {
		return err
	}

	return p.finalizeChanges(cfg, cfg.Deletions, buf)
}
