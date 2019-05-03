package cluster

import (
	"testing"

	"github.com/martinohmann/kubernetes-cluster-manager/internal/commandtest"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/credentials"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/kubernetes"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestProcessResourceDeletions(t *testing.T) {
	commandtest.WithMockExecutor(func(executor *commandtest.MockExecutor) {
		kubectl := kubernetes.NewKubectl(&credentials.Credentials{})

		deletions := []*kubernetes.ResourceSelector{{Name: "foo", Kind: "pod"}}

		remaining, err := processResourceDeletions(&Options{}, log.New(), kubectl, deletions)

		assert.NoError(t, err)
		assert.Len(t, remaining, 0)
	})
}

func TestProcessResourceDeletionsDryRun(t *testing.T) {
	commandtest.WithMockExecutor(func(executor *commandtest.MockExecutor) {
		kubectl := kubernetes.NewKubectl(&credentials.Credentials{})

		deletions := []*kubernetes.ResourceSelector{{Name: "foo", Kind: "pod"}}

		remaining, err := processResourceDeletions(&Options{DryRun: true}, log.New(), kubectl, deletions)

		assert.NoError(t, err)
		assert.Len(t, remaining, 1)
	})
}
