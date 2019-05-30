package kubernetes

import (
	"context"
	"fmt"
	"os/exec"
	"time"

	"github.com/cenkalti/backoff"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/command"
	"github.com/pkg/errors"
)

const (
	maxPollingRetries = 30
	pollingTimeout    = 2 * time.Second
)

var (
	pollingStrategy = backoff.WithMaxRetries(
		backoff.NewConstantBackOff(pollingTimeout),
		maxPollingRetries,
	)
)

// WaitOptions are passed to kubectl when waiting for a wait condition to be
// met.
type WaitOptions struct {
	Kind      string
	Name      string
	Namespace string
	For       string
	Timeout   time.Duration
}

// Wait waits until the condition in the WaitOptions is met.
func (k *Kubectl) Wait(ctx context.Context, o WaitOptions) error {
	namespace := o.Namespace
	if namespace == "" {
		namespace = DefaultNamespace
	}

	args := []string{
		"kubectl",
		"wait",
		"--for",
		o.For,
		"--namespace",
		namespace,
		fmt.Sprintf("%s/%s", o.Kind, o.Name),
	}

	if o.Timeout > 0 {
		args = append(args, "--timeout", fmt.Sprintf("%ds", int64(o.Timeout.Seconds())))
	}

	args = append(args, k.buildCredentialArgs()...)

	cmd := exec.Command(args[0], args[1:]...)

	_, err := command.RunWithContext(ctx, cmd)

	return err
}

// WaitForCluster waits until the api-server is reachable. Will retry every 2
// seconds in case of error. After 30 failed attempts it will give up and
// return the last error.
func (k *Kubectl) WaitForCluster(ctx context.Context) error {
	err := backoff.Retry(
		func() error {
			out, err := k.ClusterInfo(ctx)
			return errors.Wrapf(err, "failed to connect to cluster due to:\n%s", out)
		},
		pollingStrategy,
	)

	return err
}
