package helm

import (
	"testing"

	"github.com/martinohmann/kubernetes-cluster-manager/pkg/command"
	"github.com/stretchr/testify/assert"
)

func TestChartRender(t *testing.T) {
	executor := command.NewMockExecutor()

	chart := NewChart("cluster", executor)

	_, err := chart.Render("values.yaml")

	assert.NoError(t, err)

	if assert.Len(t, executor.ExecutedCommands, 1) {
		assert.Equal(t, "helm template --values values.yaml cluster", executor.ExecutedCommands[0])
	}
}
