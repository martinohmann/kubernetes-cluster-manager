package helm

import (
	"bytes"
	"fmt"
	"log"

	"github.com/martinohmann/cluster-manager/pkg/api"
	"github.com/martinohmann/cluster-manager/pkg/config"
	"github.com/martinohmann/cluster-manager/pkg/executor"
	"github.com/martinohmann/cluster-manager/pkg/git"
	"github.com/martinohmann/cluster-manager/pkg/manifest"
	"gopkg.in/yaml.v2"
)

var _ manifest.Renderer = &Renderer{}

type Renderer struct {
	cfg *config.Config
}

func NewManifestRenderer(cfg *config.Config) *Renderer {
	return &Renderer{cfg: cfg}
}

func (r *Renderer) RenderManifest(out *api.InfraOutput) (*api.Manifest, error) {
	if err := r.renderValues(out); err != nil {
		return nil, err
	}

	return r.renderManifest()
}

func (r *Renderer) renderValues(out *api.InfraOutput) error {
	valuesFile := r.cfg.Helm.Values
	content, err := yaml.Marshal(out.Values)
	if err != nil {
		return err
	}

	return r.writeAndDiff(valuesFile, next(valuesFile), content)
}

func (r *Renderer) generateManifest(valuesFile string) (*api.Manifest, error) {
	args := []string{
		"helm",
		"template",
		"--values",
		valuesFile,
		r.cfg.Helm.Chart,
	}

	var buf bytes.Buffer

	err := executor.Execute(&buf, args...)
	if err != nil {
		return nil, err
	}

	return &api.Manifest{Content: buf.Bytes()}, nil
}

func (r *Renderer) renderManifest() (*api.Manifest, error) {
	manifestFile := r.cfg.Manifest
	valuesFile := r.cfg.Helm.Values

	if r.cfg.DryRun {
		valuesFile = next(valuesFile)
	}

	manifest, err := r.generateManifest(valuesFile)
	if err != nil {
		return nil, err
	}

	if err := r.writeAndDiff(manifestFile, next(manifestFile), manifest.Content); err != nil {
		return nil, err
	}

	return manifest, nil
}

func (r *Renderer) writeAndDiff(current, next string, content []byte) error {
	if err := createFileIfNotExists(current); err != nil {
		return err
	}

	if err := writeFile(next, content); err != nil {
		return err
	}

	if diff, err := git.Diff(current, next); err != nil {
		return err
	} else {
		fmt.Println("changes:")
		fmt.Println(diff)
	}

	if !r.cfg.DryRun {
		log.Printf("updating %s\n", current)
		if err := writeFile(current, content); err != nil {
			return err
		}
	}

	return nil
}
