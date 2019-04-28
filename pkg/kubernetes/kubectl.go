package kubernetes

import (
	"bytes"
	"fmt"
	"os/exec"
	"sort"
	"strings"

	"github.com/cenkalti/backoff"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/command"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/kcm"
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
	credentials *kcm.Credentials
}

// NewKubectl create a new kubectl interactor.
func NewKubectl(c *kcm.Credentials) *Kubectl {
	return &Kubectl{
		credentials: c,
	}
}

// ApplyManifest applies the manifest via kubectl.
func (k *Kubectl) ApplyManifest(manifest kcm.Manifest) error {
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
			_, err := command.Run(cmd)
			return err
		},
		backoffStrategy,
	)

	return err
}

// DeleteManifest deletes the manifest via kubectl.
func (k *Kubectl) DeleteManifest(manifest kcm.Manifest) error {
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
			_, err := command.Run(cmd)
			return err
		},
		backoffStrategy,
	)

	return err
}

// DeleteResource deletes a resource via kubectl.
func (k *Kubectl) DeleteResource(deletion *kcm.Deletion) error {
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

	args = append(args, k.buildCredentialArgs()...)

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

	_, err := command.Run(cmd)

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

	_, err := command.RunSilently(cmd)

	return err
}

// ClusterInfo fetches the kubernetes cluster info.
func (k *Kubectl) ClusterInfo() (string, error) {
	args := []string{
		"kubectl",
		"cluster-info",
	}

	args = append(args, k.buildCredentialArgs()...)

	cmd := exec.Command(args[0], args[1:]...)

	return command.RunSilently(cmd)
}

// buildCredentialArgs builds kubectl args from credentials.
func (k *Kubectl) buildCredentialArgs() (args []string) {
	if k.credentials.Kubeconfig != "" {
		args = append(args, "--kubeconfig", k.credentials.Kubeconfig)

		if k.credentials.Context != "" {
			args = append(args, "--context", k.credentials.Context)
		}
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
