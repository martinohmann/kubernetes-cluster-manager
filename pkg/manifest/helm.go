package manifest

import (
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/api"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/command"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/config"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/git"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/kubernetes/helm"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type HelmRenderer struct {
	cfg      *config.Config
	executor command.Executor
}

func NewHelmRenderer(cfg *config.Config, executor command.Executor) *HelmRenderer {
	return &HelmRenderer{
		cfg:      cfg,
		executor: executor,
	}
}

func (r *HelmRenderer) RenderManifest(out *api.InfraOutput) (*api.Manifest, error) {
	valuesFile := r.cfg.Helm.Values
	manifestFile := r.cfg.Manifest

	values, err := yaml.Marshal(out.Values)
	if err != nil {
		return nil, err
	}

	valueChanges, err := git.NewFileChanges(valuesFile, values)
	if err != nil {
		return nil, err
	}

	defer func() {
		if !r.cfg.DryRun {
			valueChanges.Apply()
			valueChanges.Close()
		}
	}()

	diff, err := git.DiffFileChanges(valueChanges)
	if err != nil {
		return nil, err
	}

	if len(diff) > 0 {
		log.Infof("Changes to values:\n%s", diff)
	}

	chart := helm.NewChart(r.cfg.Helm.Chart, r.executor)

	manifest, err := chart.Render(valueChanges.Filename())
	if err != nil {
		return nil, err
	}

	manifestChanges, err := git.NewFileChanges(manifestFile, manifest.Content)
	if err != nil {
		return nil, err
	}

	defer func() {
		if !r.cfg.DryRun {
			manifestChanges.Apply()
			manifestChanges.Close()
		}
	}()

	diff, err = git.DiffFileChanges(manifestChanges)
	if err != nil {
		return nil, err
	}

	if len(diff) > 0 {
		log.Infof("Changes to manifest:\n%s", diff)
	}

	return manifest, nil
}
