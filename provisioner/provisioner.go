package provisioner

import (
	"github.com/martinohmann/kubernetes-cluster-manager/infra"
	"github.com/martinohmann/kubernetes-cluster-manager/manifest"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/api"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/command"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/credentials"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/file"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/kubernetes"
	log "github.com/sirupsen/logrus"
	yaml "gopkg.in/yaml.v2"
)

type Options struct {
	DryRun       bool   `json:"dryRun" yaml:"dryRun"`
	Manifest     string `json:"manifest" yaml:"manifest"`
	Values       string `json:"values" yaml:"values"`
	Deletions    string `json:"deletions" yaml:"deletions"`
	OnlyManifest bool   `json:"onlyManifest" yaml:"onlyManifest"`
}

type Provisioner struct {
	credentialProvider credentials.Provider
	infraManager       infra.Manager
	manifestRenderer   manifest.Renderer
	executor           command.Executor
	values             api.Values
	deletions          *api.Deletions
	logger             *log.Logger
}

func NewClusterProvisioner(
	credentialProvider credentials.Provider,
	infraManager infra.Manager,
	manifestRenderer manifest.Renderer,
	executor command.Executor,
	logger *log.Logger,
) *Provisioner {
	return &Provisioner{
		credentialProvider: credentialProvider,
		infraManager:       infraManager,
		manifestRenderer:   manifestRenderer,
		executor:           executor,
		logger:             logger,
		deletions:          &api.Deletions{},
		values:             make(api.Values),
	}
}

func (p *Provisioner) prepare(o *Options) error {
	if err := file.LoadYAML(o.Values, &p.values); err != nil {
		return err
	}

	return file.LoadYAML(o.Deletions, &p.deletions)
}

func (p *Provisioner) Provision(o *Options) error {
	var err error

	if err = p.prepare(o); err != nil {
		return err
	}

	if !o.OnlyManifest {
		if o.DryRun {
			err = p.infraManager.Plan()
		} else {
			err = p.infraManager.Apply()
		}
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

	err = p.finalizeChanges(o, o.Values, valueBytes)
	if err != nil {
		return err
	}

	manifest, err := p.manifestRenderer.RenderManifest(p.values)
	if err != nil {
		return err
	}

	creds, err := p.credentialProvider.GetCredentials()
	if err != nil {
		return err
	}

	kubectl := kubernetes.NewKubectl(creds, p.executor)

	if !o.DryRun {
		p.logger.Info("Waiting for cluster to become available...")

		if err := kubectl.WaitForCluster(); err != nil {
			return err
		}
	}

	err = p.finalizeChanges(o, o.Manifest, manifest)
	if err != nil {
		return err
	}

	defer p.finalizeDeletions(o, p.deletions)

	if err := processResourceDeletions(o, p.logger, kubectl, p.deletions.PreApply); err != nil {
		return err
	}

	if o.DryRun {
		p.logger.Warn("Would apply manifest")
		p.logger.Debug(string(manifest))
	} else if err := kubectl.ApplyManifest(manifest); err != nil {
		return err
	}

	return processResourceDeletions(o, p.logger, kubectl, p.deletions.PostApply)
}

func (p *Provisioner) Destroy(o *Options) error {
	if err := p.prepare(o); err != nil {
		return err
	}

	currentValues, err := p.infraManager.GetValues()
	if err != nil {
		return err
	}

	if err := p.values.Merge(currentValues); err != nil {
		return err
	}

	manifest, err := p.manifestRenderer.RenderManifest(p.values)
	if err != nil {
		return err
	}

	creds, err := p.credentialProvider.GetCredentials()
	if err != nil {
		return err
	}

	kubectl := kubernetes.NewKubectl(creds, p.executor)

	if o.DryRun {
		p.logger.Warn("Would delete manifest")
		p.logger.Debug(string(manifest))
	} else if err := kubectl.DeleteManifest(manifest); err != nil {
		return err
	}

	defer p.finalizeDeletions(o, p.deletions)

	if err := processResourceDeletions(o, p.logger, kubectl, p.deletions.PreDestroy); err != nil {
		return err
	}

	if !o.OnlyManifest {
		if o.DryRun {
			p.logger.Warn("Would destroy infrastructure")
		} else {
			return p.infraManager.Destroy()
		}
	}

	return nil
}

func (p *Provisioner) finalizeChanges(o *Options, filename string, content []byte) error {
	changes := file.NewChanges(filename, content)

	diff, err := changes.Diff()
	if err != nil {
		return err
	}

	if len(diff) > 0 {
		p.logger.Infof("Changes to %s:\n%s", filename, diff)
	} else {
		p.logger.Infof("No changes to %s", filename)
	}

	if o.DryRun {
		return nil
	}

	return changes.Apply()
}

func (p *Provisioner) finalizeDeletions(o *Options, deletions *api.Deletions) error {
	buf, err := yaml.Marshal(deletions.FilterPending())
	if err != nil {
		return err
	}

	return p.finalizeChanges(o, o.Deletions, buf)
}
