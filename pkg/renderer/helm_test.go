// +build integration

package renderer

import (
	"testing"

	"github.com/martinohmann/kubernetes-cluster-manager/pkg/kcm"
	"github.com/stretchr/testify/assert"
)

func TestHelmRenderManifests(t *testing.T) {
	o := &kcm.HelmOptions{
		ChartsDir: "testdata/helm",
	}

	r := NewHelm(o)

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

	manifests, err := r.RenderManifests(values)

	if !assert.NoError(t, err) {
		return
	}

	if assert.Len(t, manifests, 1) {
		assert.Equal(t, "chart.yaml", string(manifests[0].Filename))
		assert.Equal(t, expected, string(manifests[0].Content))
	}
}
