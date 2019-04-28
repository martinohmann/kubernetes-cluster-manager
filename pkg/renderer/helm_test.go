// +build integration

package renderer

import (
	"testing"

	"github.com/martinohmann/kubernetes-cluster-manager/pkg/command"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/kcm"
	"github.com/stretchr/testify/assert"
)

func TestHelmRenderManifest(t *testing.T) {
	executor := command.NewExecutor(nil)

	o := &kcm.HelmOptions{
		Chart: "testdata/helm/chart",
	}

	r := NewHelm(o, executor)

	values := kcm.Values{
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
