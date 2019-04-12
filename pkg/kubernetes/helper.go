package kubernetes

import (
	"time"

	"github.com/cenkalti/backoff"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
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
func (k *Kubectl) WaitForCluster() error {
	log.Info("Waiting for cluster to become available...")

	err := backoff.Retry(
		func() error {
			out, err := k.ClusterInfo()
			return errors.Wrapf(err, "failed to connect to cluster due to:\n%s", out)
		},
		pollingStrategy,
	)

	return err
}
