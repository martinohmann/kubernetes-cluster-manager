package helm

import (
	"testing"

	"github.com/martinohmann/kubernetes-cluster-manager/pkg/command"
	"github.com/stretchr/testify/assert"
)

func TestChartRender(t *testing.T) {
	executor := command.NewMockExecutor(nil)
	restoreExecutor := command.SetExecutorWithRestore(executor)
	defer restoreExecutor()

	chart := NewChart("cluster")

	_, err := chart.Render(make(map[string]interface{}))

	assert.NoError(t, err)

	if assert.Len(t, executor.ExecutedCommands, 1) {
		assert.Regexp(t, "helm template --values .*values.yaml.* cluster", executor.ExecutedCommands[0])
	}
}
