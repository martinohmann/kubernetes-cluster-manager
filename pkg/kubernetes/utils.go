package kubernetes

import (
	"time"

	"github.com/cenkalti/backoff"
)

const (
	maxPollingRetries = 60
	pollingTimeout    = 2 * time.Second
)

var (
	pollingStrategy = backoff.WithMaxRetries(
		backoff.NewConstantBackOff(pollingTimeout),
		maxPollingRetries,
	)
)

func WaitForCluster(kubectl *Kubectl) error {
	err := backoff.Retry(
		func() error {
			_, err := kubectl.ClusterInfo()
			return err
		},
		pollingStrategy,
	)

	return err
}
