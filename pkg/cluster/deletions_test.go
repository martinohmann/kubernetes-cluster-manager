package cluster

import (
	"context"
	"testing"

	"github.com/martinohmann/kubernetes-cluster-manager/internal/commandtest"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/credentials"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/kubernetes"
	"github.com/stretchr/testify/assert"
)

func TestProcessResourceDeletions(t *testing.T) {
	commandtest.WithMockExecutor(func(executor commandtest.MockExecutor) {
		kubectl := kubernetes.NewKubectl(&credentials.Credentials{})

		deletions := []ResourceSelector{{Name: "foo", Kind: "pod"}}

		executor.ExpectCommand("kubectl .*")

		remaining, err := processResourceDeletions(context.Background(), &Options{}, kubectl, deletions)

		assert.NoError(t, err)
		assert.Len(t, remaining, 0)

		assert.NoError(t, executor.ExpectationsWereMet())
	})
}

func TestProcessResourceDeletionsDryRun(t *testing.T) {
	commandtest.WithMockExecutor(func(executor commandtest.MockExecutor) {
		kubectl := kubernetes.NewKubectl(&credentials.Credentials{})

		deletions := []ResourceSelector{{Name: "foo", Kind: "pod"}}

		remaining, err := processResourceDeletions(context.Background(), &Options{DryRun: true}, kubectl, deletions)

		assert.NoError(t, err)
		assert.Len(t, remaining, 1)

		assert.NoError(t, executor.ExpectationsWereMet())
	})
}
