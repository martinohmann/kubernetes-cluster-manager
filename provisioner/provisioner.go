package provisioner

import (
	"fmt"

	"github.com/martinohmann/kubernetes-cluster-manager/infra"
	"github.com/martinohmann/kubernetes-cluster-manager/manifest"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/api"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/command"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/config"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/git"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/kubernetes"
	log "github.com/sirupsen/logrus"
	yaml "gopkg.in/yaml.v2"
)

type Provisioner struct {
	infraManager     infra.Manager
	manifestRenderer manifest.Renderer
	executor         command.Executor
	values           api.Values
	deletions        *api.Deletions
}

func NewClusterProvisioner(infraManager infra.Manager, manifestRenderer manifest.Renderer, executor command.Executor) *Provisioner {
	return &Provisioner{
		infraManager:     infraManager,
		manifestRenderer: manifestRenderer,
		executor:         executor,
	}
}

func (p *Provisioner) prepare(cfg *config.Config) (err error) {
	p.values, err = loadValues(cfg.Values)
	if err != nil {
		return
	}

	p.deletions, err = loadDeletions(cfg.Deletions)

	return
}

func (p *Provisioner) Provision(cfg *config.Config) error {
	var err error

	if err = p.prepare(cfg); err != nil {
		return err
	}

	if cfg.DryRun {
		err = p.infraManager.Plan()
	} else if !cfg.OnlyManifest {
		err = p.infraManager.Apply()
	}

	if err != nil {
		return err
	}

	newValues, err := p.infraManager.GetValues()
	if err != nil {
		return err
	}

	if err := p.values.Merge(newValues); err != nil {
		return err
	}

	valueBytes, err := yaml.Marshal(p.values)
	if err != nil {
		return err
	}

	err = p.finalizeChanges(cfg, cfg.Values, valueBytes)
	if err != nil {
		return err
	}

	manifest, err := p.manifestRenderer.RenderManifest(p.values)
	if err != nil {
		return err
	}

	cfg.Cluster.Update(p.values)

	kubectl := kubernetes.NewKubectl(&cfg.Cluster, p.executor)

	if cfg.DryRun {
		log.Debug("Would wait for cluster to become available.")
	} else if err := kubectl.WaitForCluster(); err != nil {
		return err
	}

	err = p.finalizeChanges(cfg, cfg.Manifest, manifest)
	if err != nil {
		return err
	}

	defer p.finalizeDeletions(cfg, p.deletions)

	if err := processResourceDeletions(cfg, kubectl, p.deletions.PreApply); err != nil {
		return err
	}

	if cfg.DryRun {
		log.Warnf("Would apply manifest:\n%s", manifest)
	} else if err := kubectl.ApplyManifest(manifest); err != nil {
		return err
	}

	return processResourceDeletions(cfg, kubectl, p.deletions.PostApply)
}

func (p *Provisioner) Destroy(cfg *config.Config) error {
	if err := p.prepare(cfg); err != nil {
		return err
	}

	manifest, err := p.manifestRenderer.RenderManifest(p.values)
	if err != nil {
		return err
	}

	cfg.Cluster.Update(p.values)

	kubectl := kubernetes.NewKubectl(&cfg.Cluster, p.executor)

	if cfg.DryRun {
		log.Debug("Would wait for cluster to become available.")
	} else if err := kubectl.WaitForCluster(); err != nil {
		return err
	}

	if cfg.DryRun {
		log.Warnf("Would delete manifest:\n%s", manifest)
	} else if err := kubectl.DeleteManifest(manifest); err != nil {
		return err
	}

	defer p.finalizeDeletions(cfg, p.deletions)

	if err := processResourceDeletions(cfg, kubectl, p.deletions.PreDestroy); err != nil {
		return err
	}

	if cfg.DryRun {
		log.Warn("Would destroy infrastructure")
	} else if !cfg.OnlyManifest {
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

	fmt.Println(cfg.Deletions)

	return p.finalizeChanges(cfg, cfg.Deletions, buf)
}
