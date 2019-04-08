package manifest

import (
	"os/exec"

	"github.com/martinohmann/kubernetes-cluster-manager/pkg/api"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/command"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/config"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/git"
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
	diffTool := &git.DiffTool{DiffOnly: r.cfg.DryRun}
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

	defer valueChanges.Close()

	diff, err := diffTool.Apply(valueChanges)
	if err != nil {
		return nil, err
	}

	if len(diff) > 0 {
		log.Infof("Changes to values:\n%s", diff)
	}

	manifest, err := r.generateManifest(valueChanges.Filename(), r.cfg.Helm.Chart)
	if err != nil {
		return nil, err
	}

	manifestChanges, err := git.NewFileChanges(manifestFile, manifest.Content)
	if err != nil {
		return nil, err
	}

	defer manifestChanges.Close()

	diff, err = diffTool.Apply(manifestChanges)
	if err != nil {
		return nil, err
	}

	if len(diff) > 0 {
		log.Infof("Changes to manifest:\n%s", diff)
	}

	return manifest, nil
}

func (r *HelmRenderer) generateManifest(values string, chart string) (*api.Manifest, error) {
	args := []string{
		"helm",
		"template",
		"--values",
		values,
		chart,
	}

	cmd := exec.Command(args[0], args[1:]...)

	out, err := r.executor.RunSilently(cmd)
	if err != nil {
		return nil, err
	}

	return &api.Manifest{Content: []byte(out)}, nil
}
