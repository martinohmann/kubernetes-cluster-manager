package provisioner

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/cenkalti/backoff"
	"github.com/martinohmann/cluster-manager/pkg/api"
	"github.com/martinohmann/cluster-manager/pkg/config"
	"github.com/martinohmann/cluster-manager/pkg/executor"
	"github.com/martinohmann/cluster-manager/pkg/infra"
	"github.com/martinohmann/cluster-manager/pkg/manifest"
	"gopkg.in/yaml.v2"
)

const (
	maxRetries = 3
)

type Provisioner struct {
	manifestRenderer manifest.Renderer
	infraManager     infra.Manager
	w                io.Writer
}

func NewClusterProvisioner(im infra.Manager, mr manifest.Renderer, w io.Writer) *Provisioner {
	if w == nil {
		w = os.Stdout
	}

	return &Provisioner{
		infraManager:     im,
		manifestRenderer: mr,
		w:                w,
	}
}

func (p *Provisioner) Provision(cfg *config.Config, deletions *api.Deletions) error {
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

	if deletions == nil {
		deletions = &api.Deletions{}
	}

	if err := p.processDeletions(cfg, deletions.PreApply); err != nil {
		return err
	}

	if err := p.applyManifest(cfg, manifest); err != nil {
		return err
	}

	if err := p.processDeletions(cfg, deletions.PostApply); err != nil {
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

	if err := p.deleteManifest(cfg, manifest); err != nil {
		return err
	}

	if !cfg.OnlyManifest && !cfg.DryRun {
		if err := p.infraManager.Destroy(); err != nil {
			return err
		}
	} else if cfg.DryRun {
		fmt.Println("would destroy infrastructure")
	}

	return nil
}

func (p *Provisioner) applyManifest(cfg *config.Config, manifest *api.Manifest) error {
	return p.submitManifest("apply", cfg, manifest)
}

func (p *Provisioner) deleteManifest(cfg *config.Config, manifest *api.Manifest) error {
	if cfg.DryRun {
		fmt.Println("would delete manifest:")
		fmt.Println(manifest)

		return nil
	}

	return p.submitManifest("delete", cfg, manifest)
}

func (p *Provisioner) submitManifest(action string, cfg *config.Config, manifest *api.Manifest) error {
	args := []string{
		"kubectl",
		action,
		"-f",
		"-",
	}

	if cfg.Kubeconfig != "" {
		args = append(args, "--kubeconfig", cfg.Kubeconfig)
	}

	if action != "delete" && cfg.DryRun {
		args = append(args, "--dry-run")
	}

	if action == "delete" {
		args = append(args, "--ignore-not-found")
	}

	err := backoff.Retry(
		func() error {
			in := bytes.NewBuffer(manifest.Content)
			return executor.Pipe(in, p.w, args...)
		},
		backoff.WithMaxRetries(backoff.NewExponentialBackOff(), maxRetries))

	return err
}

func (p *Provisioner) processDeletions(cfg *config.Config, deletions []api.Deletion) error {
	if cfg.DryRun {
		buf, err := yaml.Marshal(deletions)
		if err != nil {
			return err
		}

		fmt.Println("Would delete the following resources:")
		fmt.Println(string(buf))

		return nil
	}

	for _, deletion := range deletions {
		namespace := deletion.Namespace
		if namespace == "" {
			namespace = "default"
		}

		args := []string{
			"kubectl",
			"delete",
			fmt.Sprintf("--namespace=%s", namespace),
			deletion.Kind,
		}

		if deletion.Name != "" {
			args = append(args, deletion.Name)
		} else if len(deletion.Labels) > 0 {
			args = append(args, fmt.Sprintf("--selector=%s", deletion.Labels))
		} else {
			return fmt.Errorf(
				"either a name or labels must be specified for a deletion (kind=%s,namespace=%s)",
				deletion.Kind,
				deletion.Namespace,
			)
		}

		if err := executor.Execute(p.w, args...); err != nil {
			return err
		}
	}

	return nil
}
