package provisioner

import (
	"errors"
	"os"
	"testing"

	"github.com/martinohmann/kubernetes-cluster-manager/pkg/api"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/command"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/config"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/fs"
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

func TestLoadDeletions(t *testing.T) {
	content := []byte("---\npreApply:\n- kind: pod\n  name: foo")
	f, err := fs.NewTempFile("deletions.yaml", content)
	if !assert.NoError(t, err) {
		return
	}

	defer os.Remove(f.Name())

	deletions, err := loadDeletions(f.Name())

	if !assert.NoError(t, err) {
		return
	}

	if assert.Len(t, deletions.PreApply, 1) {
		assert.Equal(t, "pod", deletions.PreApply[0].Kind)
		assert.Equal(t, "foo", deletions.PreApply[0].Name)
	}
}
