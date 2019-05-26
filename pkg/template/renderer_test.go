package template

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRender(t *testing.T) {
	r := NewRenderer()

	values := map[string]interface{}{
		"config": map[string]interface{}{
			"bar": "baz",
		},
	}

	expectedConfigMap := `---
apiVersion: v1
kind: ConfigMap
metadata:
  name: kcm-chart
  labels:
    app.kubernetes.io/name: chart
    helm.sh/chart: chart-0.1.0
data:
  bar: "baz"
  foo: "bar"
`

	expectedService := `---
apiVersion: v1
kind: Service
metadata:
  name: kcm-chart
  labels:
    app.kubernetes.io/name: chart
    helm.sh/chart: chart-0.1.0
spec:
  type: ClusterIP
  ports:
    - port: 80
      targetPort: http
      protocol: TCP
      name: http
  selector:
    app.kubernetes.io/name: chart
    app.kubernetes.io/instance: kcm
`

	renderedTemplates, err := r.Render("testdata/charts/chart", values)

	require.NoError(t, err)
	require.Len(t, renderedTemplates, 3)
	assert.Equal(t, expectedConfigMap, renderedTemplates["chart/templates/configmap.yaml"])
	assert.Equal(t, expectedService, renderedTemplates["chart/templates/service.yaml"])
}

func TestRenderError(t *testing.T) {
	r := NewRenderer()

	values := map[string]interface{}{}

	_, err := r.Render("testdata/charts/notahelmchart", values)

	require.Error(t, err)
}
