package kubernetes

import (
	"bytes"
	"context"
	"os/exec"
	"strings"

	"github.com/cenkalti/backoff"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/command"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/credentials"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/resource"
)

const (
	// maxRetries defines the number of retries for kubectl commands.
	maxRetries = 10

	// DefaultNamespace is the namespace that should be used where namespace is
	// omitted.
	DefaultNamespace = "default"
)

var (
	// backoffStrategy is the retries strategy used for failed kubectl commands.
	backoffStrategy = backoff.WithMaxRetries(backoff.NewExponentialBackOff(), maxRetries)
)

// Kubectl defines a type for interacting with kubectl.
type Kubectl struct {
	credentials *credentials.Credentials
}

// NewKubectl create a new kubectl interactor.
func NewKubectl(c *credentials.Credentials) *Kubectl {
	return &Kubectl{
		credentials: c,
	}
}

// ApplyManifest applies the manifest via kubectl.
func (k *Kubectl) ApplyManifest(ctx context.Context, manifest []byte) error {
	args := []string{
		"kubectl",
		"apply",
		"-f",
		"-",
	}

	args = append(args, k.buildCredentialArgs()...)

	err := backoff.Retry(
		func() error {
			cmd := exec.Command(args[0], args[1:]...)
			cmd.Stdin = bytes.NewBuffer(manifest)
			_, err := command.RunWithContext(ctx, cmd)
			return err
		},
		backoffStrategy,
	)

	return err
}

// DeleteManifest deletes the manifest via kubectl.
func (k *Kubectl) DeleteManifest(ctx context.Context, manifest []byte) error {
	args := []string{
		"kubectl",
		"delete",
		"-f",
		"-",
		"--ignore-not-found",
	}

	args = append(args, k.buildCredentialArgs()...)

	err := backoff.Retry(
		func() error {
			cmd := exec.Command(args[0], args[1:]...)
			cmd.Stdin = bytes.NewBuffer(manifest)
			_, err := command.RunWithContext(ctx, cmd)
			return err
		},
		backoffStrategy,
	)

	return err
}

// DeleteResource deletes a resource by its kind, name and namespace.
func (k *Kubectl) DeleteResource(ctx context.Context, selector resource.Head) error {
	namespace := selector.Metadata.Namespace
	if namespace == "" {
		namespace = DefaultNamespace
	}

	args := []string{
		"kubectl",
		"delete",
		strings.ToLower(selector.Kind),
		selector.Metadata.Name,
		"--namespace",
		namespace,
		"--ignore-not-found",
	}

	args = append(args, k.buildCredentialArgs()...)

	err := backoff.Retry(
		func() error {
			cmd := exec.Command(args[0], args[1:]...)
			_, err := command.RunWithContext(ctx, cmd)
			return err
		},
		backoffStrategy,
	)

	return err
}

// ClusterInfo fetches the kubernetes cluster info.
func (k *Kubectl) ClusterInfo(ctx context.Context) (string, error) {
	args := []string{
		"kubectl",
		"cluster-info",
	}

	args = append(args, k.buildCredentialArgs()...)

	cmd := exec.Command(args[0], args[1:]...)

	return command.RunSilentlyWithContext(ctx, cmd)
}

// buildCredentialArgs builds kubectl args from credentials.
func (k *Kubectl) buildCredentialArgs() (args []string) {
	if k.credentials.Context != "" {
		args = append(args, "--context", k.credentials.Context)
	}

	if k.credentials.Kubeconfig != "" {
		args = append(args, "--kubeconfig", k.credentials.Kubeconfig)
	} else {
		if k.credentials.Server != "" {
			args = append(args, "--server", k.credentials.Server)
		}

		if k.credentials.Token != "" {
			args = append(args, "--token", k.credentials.Token)
		}
	}

	return args
}
