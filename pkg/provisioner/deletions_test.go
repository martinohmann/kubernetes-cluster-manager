package provisioner

import (
	"errors"
	"testing"

	"github.com/martinohmann/kubernetes-cluster-manager/pkg/api"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/command"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/config"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/kubernetes"
	"github.com/stretchr/testify/assert"
)

func TestProcessResourceDeletions(t *testing.T) {
	kubectl := kubernetes.NewKubectl(&config.Config{}, command.NewMockExecutor())

	deletions := []*api.Deletion{{Name: "foo", Kind: "pod"}}

	err := processResourceDeletions(kubectl, deletions)

	if assert.NoError(t, err) {
		assert.True(t, deletions[0].Deleted())
	}
}

func TestProcessResourceDeletionsFailed(t *testing.T) {
	executor := command.NewMockExecutor()
	kubectl := kubernetes.NewKubectl(&config.Config{}, executor)

	deletions := []*api.Deletion{{Name: "foo", Kind: "pod"}}
	expectedError := errors.New("deletion failed")

	executor.WillErrorWith(expectedError)

	err := processResourceDeletions(kubectl, deletions)

	if assert.Equal(t, err, expectedError) {
		assert.False(t, deletions[0].Deleted())
	}
}
