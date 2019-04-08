package manifest

import (
	"io/ioutil"
	"os"

	"github.com/imdario/mergo"
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

	values, err := loadValues(valuesFile)
	if err != nil {
		return nil, err
	}

	if err := mergo.Merge(&values, out.Values, mergo.WithOverride); err != nil {
		return nil, err
	}

	valueBytes, err := yaml.Marshal(values)
	if err != nil {
		return nil, err
	}

	valueChanges, err := r.processChanges(valuesFile, valueBytes)
	if err != nil {
		return nil, err
	}

	defer valueChanges.Close()

	chart := helm.NewChart(r.cfg.Helm.Chart, r.executor)

	manifest, err := chart.Render(valueChanges.Filename())
	if err != nil {
		return nil, err
	}

	manifestChanges, err := r.processChanges(manifestFile, manifest.Content)
	if err != nil {
		return nil, err
	}

	defer manifestChanges.Close()

	return manifest, nil
}

func (r *HelmRenderer) processChanges(filename string, content []byte) (*git.FileChanges, error) {
	changes, err := git.NewFileChanges(filename, content)
	if err != nil {
		return nil, err
	}

	diff, err := git.DiffFileChanges(changes)
	if err != nil {
		return nil, err
	}

	if len(diff) > 0 {
		log.Infof("Changes to %s:\n%s", filename, diff)
	} else {
		log.Infof("No changes to %s", filename)
	}

	if r.cfg.DryRun {
		return changes, nil
	}

	return changes, changes.Apply()
}

func loadValues(valuesFile string) (map[string]interface{}, error) {
	values := make(map[string]interface{})

	content, err := ioutil.ReadFile(valuesFile)
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	err = yaml.Unmarshal(content, &values)

	return values, err
}
