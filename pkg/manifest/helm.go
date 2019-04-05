package manifest

import (
	"bytes"
	"os"

	"github.com/martinohmann/cluster-manager/pkg/api"
	"github.com/martinohmann/cluster-manager/pkg/config"
	"github.com/martinohmann/cluster-manager/pkg/executor"
	"github.com/martinohmann/cluster-manager/pkg/git"
	"gopkg.in/yaml.v2"
)

var _ Renderer = &HelmRenderer{}

type HelmRenderer struct {
	cfg *config.Config
}

func NewHelmRenderer(cfg *config.Config) *HelmRenderer {
	return &HelmRenderer{cfg: cfg}
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

	defer valueChanges.Close()

	if err = git.DiffAndApply(os.Stdout, valueChanges, !r.cfg.DryRun); err != nil {
		return nil, err
	}

	manifest, err := generateManifest(valueChanges.Filename(), r.cfg.Helm.Chart)
	if err != nil {
		return nil, err
	}

	manifestChanges, err := git.NewFileChanges(manifestFile, manifest.Content)
	if err != nil {
		return nil, err
	}

	defer manifestChanges.Close()

	if err = git.DiffAndApply(os.Stdout, manifestChanges, !r.cfg.DryRun); err != nil {
		return nil, err
	}

	return manifest, nil
}

func generateManifest(values string, chart string) (*api.Manifest, error) {
	args := []string{
		"helm",
		"template",
		"--values",
		values,
		chart,
	}

	var buf bytes.Buffer

	err := executor.Execute(&buf, args...)
	if err != nil {
		return nil, err
	}

	return &api.Manifest{Content: buf.Bytes()}, nil
}
