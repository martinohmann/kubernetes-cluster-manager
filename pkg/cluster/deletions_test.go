package cluster

import (
	"errors"
	"testing"

	"github.com/martinohmann/kubernetes-cluster-manager/internal/commandtest"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/command"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/kcm"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/kubernetes"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestProcessResourceDeletions(t *testing.T) {
	commandtest.WithMockExecutor(func(executor *command.MockExecutor) {
		kubectl := kubernetes.NewKubectl(&kcm.Credentials{})

		deletions := []*kcm.Deletion{{Name: "foo", Kind: "pod"}}

		err := processResourceDeletions(&kcm.Options{}, log.New(), kubectl, deletions)

		if assert.NoError(t, err) {
			assert.True(t, deletions[0].Deleted())
		}
	})
}

func TestProcessResourceDeletionsDryRun(t *testing.T) {
	commandtest.WithMockExecutor(func(executor *command.MockExecutor) {
		kubectl := kubernetes.NewKubectl(&kcm.Credentials{})

		deletions := []*kcm.Deletion{{Name: "foo", Kind: "pod"}}

		err := processResourceDeletions(&kcm.Options{DryRun: true}, log.New(), kubectl, deletions)

		if assert.NoError(t, err) {
			assert.False(t, deletions[0].Deleted())
		}
	})
}

func TestProcessResourceDeletionsFailed(t *testing.T) {
	commandtest.WithMockExecutor(func(executor *command.MockExecutor) {
		kubectl := kubernetes.NewKubectl(&kcm.Credentials{})

		deletions := []*kcm.Deletion{{Name: "foo", Kind: "pod"}}
		expectedError := errors.New("deletion failed")

		executor.NextCommand().WillReturnError(expectedError)

		err := processResourceDeletions(&kcm.Options{}, log.New(), kubectl, deletions)

		if assert.Equal(t, expectedError, err) {
			assert.False(t, deletions[0].Deleted())
		}
	})
}
