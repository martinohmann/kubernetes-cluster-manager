package provisioner

import (
	"errors"
	"testing"

	"github.com/martinohmann/kubernetes-cluster-manager/pkg/api"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/command"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/kubernetes"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestProcessResourceDeletions(t *testing.T) {
	o := &Options{}
	l := log.New()
	kubectl := kubernetes.NewKubectl(&kubernetes.ClusterOptions{}, command.NewMockExecutor(nil))

	deletions := []*api.Deletion{{Name: "foo", Kind: "pod"}}

	err := processResourceDeletions(o, l, kubectl, deletions)

	if assert.NoError(t, err) {
		assert.True(t, deletions[0].Deleted())
	}
}

func TestProcessResourceDeletionsDryRun(t *testing.T) {
	o := &Options{DryRun: true}
	l := log.New()
	kubectl := kubernetes.NewKubectl(&kubernetes.ClusterOptions{}, command.NewMockExecutor(nil))

	deletions := []*api.Deletion{{Name: "foo", Kind: "pod"}}

	err := processResourceDeletions(o, l, kubectl, deletions)

	if assert.NoError(t, err) {
		assert.False(t, deletions[0].Deleted())
	}
}

func TestProcessResourceDeletionsFailed(t *testing.T) {
	o := &Options{}
	l := log.New()
	executor := command.NewMockExecutor(nil)
	kubectl := kubernetes.NewKubectl(&kubernetes.ClusterOptions{}, executor)

	deletions := []*api.Deletion{{Name: "foo", Kind: "pod"}}
	expectedError := errors.New("deletion failed")

	executor.NextCommand().WillReturnError(expectedError)

	err := processResourceDeletions(o, l, kubectl, deletions)

	if assert.Equal(t, expectedError, err) {
		assert.False(t, deletions[0].Deleted())
	}
}
