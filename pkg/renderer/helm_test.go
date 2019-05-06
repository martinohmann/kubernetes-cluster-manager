package renderer

import (
	"testing"

	"github.com/martinohmann/kubernetes-cluster-manager/pkg/kcm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHelmRenderManifests(t *testing.T) {
	o := &Options{
		TemplatesDir: "testdata/helm",
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
  name: kcm-chart
  labels:
    app.kubernetes.io/name: chart
    helm.sh/chart: chart-0.1.0
data:
  bar: "baz"
  foo: "bar"

---
# Source: chart/templates/service.yaml
---
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

	manifests, err := r.RenderManifests(values)

	require.NoError(t, err)
	require.Len(t, manifests, 1)

	assert.Equal(t, "chart", manifests[0].Name)
	assert.Equal(t, expected, string(manifests[0].Content))
}
