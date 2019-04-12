// +build integration

package manifest

import (
	"testing"

	"github.com/martinohmann/kubernetes-cluster-manager/pkg/api"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/command"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/config"
	"github.com/stretchr/testify/assert"
)

func TestRenderManifest(t *testing.T) {
	executor := command.NewExecutor()

	cfg := &config.HelmConfig{
		Chart: "helm/testdata/chart",
	}

	r := NewHelmRenderer(cfg, executor)

	values := api.Values{
		"config": map[string]interface{}{
			"bar": "baz",
		},
	}

	expected := `---
# Source: chart/templates/configmap.yaml
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: release-name-chart
  labels:
    app.kubernetes.io/name: chart
    helm.sh/chart: chart-0.1.0
data:
  bar: "baz"
  foo: "bar"

`

	manifest, err := r.RenderManifest(values)

	if assert.NoError(t, err) {
		assert.Equal(t, expected, string(manifest))
	}
}
