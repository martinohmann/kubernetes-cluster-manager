package kubernetes

import (
	"bytes"
	"fmt"
	"os/exec"
	"sort"
	"strings"

	"github.com/cenkalti/backoff"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/command"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/credentials"
	"github.com/pkg/errors"
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

// ResourceSelector is used to select kubernetes resources.
type ResourceSelector struct {
	Kind      string            `json:"kind,omitempty" yaml:"kind,omitempty"`
	Name      string            `json:"name,omitempty" yaml:"name,omitempty"`
	Namespace string            `json:"namespace,omitempty" yaml:"namespace,omitempty"`
	Labels    map[string]string `json:"labels,omitempty" yaml:"labels,omitempty"`
}

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
func (k *Kubectl) ApplyManifest(manifest []byte) error {
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
func (k *Kubectl) DeleteManifest(manifest []byte) error {
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
func (k *Kubectl) DeleteResource(selector *ResourceSelector) error {
	namespace := selector.Namespace
	if namespace == "" {
		namespace = DefaultNamespace
	}

	args := []string{
		"kubectl",
		"delete",
		strings.ToLower(selector.Kind),
		"--ignore-not-found",
		"--namespace",
		namespace,
	}

	args = append(args, k.buildCredentialArgs()...)

	if selector.Name != "" {
		args = append(args, selector.Name)
	} else if len(selector.Labels) > 0 {
		keys := make([]string, 0, len(selector.Labels))
		for k := range selector.Labels {
			keys = append(keys, k)
		}

		sort.Strings(keys)

		pairs := make([]string, 0, len(selector.Labels))
		for _, k := range keys {
			pairs = append(pairs, fmt.Sprintf("%s=%s", k, selector.Labels[k]))
		}

		args = append(args, "--selector", strings.Join(pairs, ","))
	} else {
		return errors.Errorf(
			"either a name or labels must be specified in the resource selector (kind=%s,namespace=%s)",
			selector.Kind,
			namespace,
		)
	}

	cmd := exec.Command(args[0], args[1:]...)

	_, err := command.Run(cmd)

	return err
}

// DeleteResources deletes resources via kubectl. Returns a slice containing
// the resources that were not deleted to to an error.
func (k *Kubectl) DeleteResources(resources []*ResourceSelector) ([]*ResourceSelector, error) {
	for i, selector := range resources {
		if err := k.DeleteResource(selector); err != nil {
			return resources[i:], err
		}
	}

	return []*ResourceSelector{}, nil
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
