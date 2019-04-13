package kubernetes

import (
	"bytes"
	"fmt"
	"os/exec"
	"sort"
	"strings"

	"github.com/cenkalti/backoff"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/api"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/command"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/config"
	"github.com/pkg/errors"
)

const (
	// maxRetries defines the number of retries for kubectl commands.
	maxRetries = 10

	// defaultNamespace is the namespace that should be used where namespace is
	// omitted.
	defaultNamespace = "default"
)

var (
	// backoffStrategy is the retries strategy used for failed kubectl commands.
	backoffStrategy = backoff.WithMaxRetries(backoff.NewExponentialBackOff(), maxRetries)
)

// Kubectl defines a type for interacting with kubectl.
type Kubectl struct {
	cfg        *config.ClusterConfig
	executor   command.Executor
	globalArgs []string
}

// NewKubectl create a new kubectl interactor.
func NewKubectl(cfg *config.ClusterConfig, executor command.Executor) *Kubectl {
	return &Kubectl{
		cfg:        cfg,
		executor:   executor,
		globalArgs: buildGlobalKubectlArgs(cfg),
	}
}

// ApplyManifest applies the manifest via kubectl.
func (k *Kubectl) ApplyManifest(manifest api.Manifest) error {
	args := []string{
		"kubectl",
		"apply",
		"-f",
		"-",
	}

	args = append(args, k.globalArgs...)

	err := backoff.Retry(
		func() error {
			cmd := exec.Command(args[0], args[1:]...)
			cmd.Stdin = bytes.NewBuffer(manifest)
			_, err := k.executor.Run(cmd)
			return err
		},
		backoffStrategy,
	)

	return err
}

// DeleteManifest deletes the manifest via kubectl.
func (k *Kubectl) DeleteManifest(manifest api.Manifest) error {
	args := []string{
		"kubectl",
		"delete",
		"-f",
		"-",
		"--ignore-not-found",
	}

	args = append(args, k.globalArgs...)

	err := backoff.Retry(
		func() error {
			cmd := exec.Command(args[0], args[1:]...)
			cmd.Stdin = bytes.NewBuffer(manifest)
			_, err := k.executor.Run(cmd)
			return err
		},
		backoffStrategy,
	)

	return err
}

// DeleteResource deletes a resource via kubectl.
func (k *Kubectl) DeleteResource(deletion *api.Deletion) error {
	namespace := deletion.Namespace
	if namespace == "" {
		namespace = defaultNamespace
	}

	args := []string{
		"kubectl",
		"delete",
		strings.ToLower(deletion.Kind),
		"--ignore-not-found",
		"--namespace",
		namespace,
	}

	args = append(args, k.globalArgs...)

	if deletion.Name != "" {
		args = append(args, deletion.Name)
	} else if len(deletion.Labels) > 0 {
		keys := make([]string, 0, len(deletion.Labels))
		for k := range deletion.Labels {
			keys = append(keys, k)
		}

		sort.Strings(keys)

		pairs := make([]string, 0, len(deletion.Labels))
		for _, k := range keys {
			pairs = append(pairs, fmt.Sprintf("%s=%s", k, deletion.Labels[k]))
		}

		args = append(args, "--selector", strings.Join(pairs, ","))
	} else {
		return errors.Errorf(
			"either a name or labels must be specified for a deletion (kind=%s,namespace=%s)",
			deletion.Kind,
			namespace,
		)
	}

	cmd := exec.Command(args[0], args[1:]...)

	_, err := k.executor.Run(cmd)

	return err
}

// UseContext sets the active kubernetes context
func (k *Kubectl) UseContext(context string) error {
	args := []string{
		"kubectl",
		"config",
		"use-context",
		context,
	}

	cmd := exec.Command(args[0], args[1:]...)

	out, err := k.executor.RunSilently(cmd)
	if err != nil {
		return errors.Wrapf(err, "failed to set context: %s", out)
	}

	return nil
}

// ClusterInfo fetches the kubernetes cluster info.
func (k *Kubectl) ClusterInfo() (string, error) {
	args := []string{
		"kubectl",
		"cluster-info",
	}

	args = append(args, k.globalArgs...)

	cmd := exec.Command(args[0], args[1:]...)

	return k.executor.RunSilently(cmd)
}

// buildGlobalKubectlArgs builds global kubectl args from the config.
func buildGlobalKubectlArgs(cfg *config.ClusterConfig) (args []string) {
	if cfg.Kubeconfig != "" {
		args = append(args, "--kubeconfig", cfg.Kubeconfig)

		if cfg.Context != "" {
			args = append(args, "--context", cfg.Context)
		}
	} else {
		if cfg.Server != "" {
			args = append(args, "--server", cfg.Server)
		}

		if cfg.Token != "" {
			args = append(args, "--token", cfg.Token)
		}
	}

	return args
}
