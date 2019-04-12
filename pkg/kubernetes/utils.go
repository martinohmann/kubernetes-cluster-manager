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
