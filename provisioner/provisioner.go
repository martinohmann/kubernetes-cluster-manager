package provisioner

import (
	"os"

	"github.com/martinohmann/kubernetes-cluster-manager/infra"
	"github.com/martinohmann/kubernetes-cluster-manager/manifest"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/api"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/command"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/fs"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/git"
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

func (o *Options) ApplyDefaults() {
	if o.Manifest == "" {
		o.Manifest = "./manifest.yaml"
	}

	if o.Deletions == "" {
		o.Deletions = "./deletions.yaml"
	}

	if o.Values == "" {
		o.Values = "./values.yaml"
	}
}

type Provisioner struct {
	infraManager     infra.Manager
	manifestRenderer manifest.Renderer
	clusterOptions   *kubernetes.ClusterOptions
	executor         command.Executor
	values           api.Values
	deletions        *api.Deletions
}

func NewClusterProvisioner(
	clusterOptions *kubernetes.ClusterOptions,
	infraManager infra.Manager,
	manifestRenderer manifest.Renderer,
	executor command.Executor,
) *Provisioner {
	return &Provisioner{
		clusterOptions:   clusterOptions,
		infraManager:     infraManager,
		manifestRenderer: manifestRenderer,
		executor:         executor,
	}
}

func (p *Provisioner) prepare(o *Options) (err error) {
	p.values, err = loadValues(o.Values)
	if err != nil {
		return
	}

	p.deletions, err = loadDeletions(o.Deletions)

	return
}

func (p *Provisioner) Provision(o *Options) error {
	var err error

	if err = p.prepare(o); err != nil {
		return err
	}

	if o.DryRun {
		err = p.infraManager.Plan()
	} else if !o.OnlyManifest {
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

	err = p.finalizeChanges(o, o.Values, valueBytes)
	if err != nil {
		return err
	}

	manifest, err := p.manifestRenderer.RenderManifest(p.values)
	if err != nil {
		return err
	}

	p.clusterOptions.Update(p.values)

	kubectl := kubernetes.NewKubectl(p.clusterOptions, p.executor)

	if o.DryRun {
		log.Debug("Would wait for cluster to become available.")
	} else if err := kubectl.WaitForCluster(); err != nil {
		return err
	}

	err = p.finalizeChanges(o, o.Manifest, manifest)
	if err != nil {
		return err
	}

	defer p.finalizeDeletions(o, p.deletions)

	if err := processResourceDeletions(o, kubectl, p.deletions.PreApply); err != nil {
		return err
	}

	if o.DryRun {
		log.Warn("Would apply manifest")
		log.Debug(string(manifest))
	} else if err := kubectl.ApplyManifest(manifest); err != nil {
		return err
	}

	return processResourceDeletions(o, kubectl, p.deletions.PostApply)
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

	p.clusterOptions.Update(p.values)

	kubectl := kubernetes.NewKubectl(p.clusterOptions, p.executor)

	if o.DryRun {
		log.Warn("Would delete manifest")
		log.Debug(string(manifest))
	} else if err := kubectl.DeleteManifest(manifest); err != nil {
		return err
	}

	defer p.finalizeDeletions(o, p.deletions)

	if err := processResourceDeletions(o, kubectl, p.deletions.PreDestroy); err != nil {
		return err
	}

	if o.DryRun {
		log.Warn("Would destroy infrastructure")
	} else if !o.OnlyManifest {
		return p.infraManager.Destroy()
	}

	return nil
}

func (p *Provisioner) finalizeChanges(o *Options, filename string, content []byte) error {
	changes, err := git.NewFileChanges(filename, content)
	if err != nil {
		return err
	}

	defer changes.Close()

	if !fs.Exists(filename) {
		if err := fs.Touch(filename); err != nil {
			return err
		}

		if o.DryRun {
			defer os.Remove(filename)
		}
	}

	diff, err := changes.Diff()
	if err != nil {
		return err
	}

	if len(diff) > 0 {
		log.Infof("Changes to %s:\n%s", filename, diff)
	} else {
		log.Infof("No changes to %s", filename)
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
