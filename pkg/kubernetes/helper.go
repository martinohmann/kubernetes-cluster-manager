package kubernetes

import (
	"context"
	"time"

	"github.com/cenkalti/backoff"
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
