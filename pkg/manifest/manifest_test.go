package manifest

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestManifest_Filename(t *testing.T) {
	m := &Manifest{Name: "manifest"}

	assert.Equal(t, "manifest.yaml", m.Filename())
}

func TestManifest_IsBlank(t *testing.T) {
	cases := []struct {
		name  string
		m     *Manifest
		blank bool
	}{
		{
			name:  "manifest nil",
			blank: true,
		},
		{
			name:  "manifest empty",
			m:     &Manifest{},
			blank: true,
		},
		{
			name:  "only whitespace",
			m:     &Manifest{content: []byte("\n  \n ")},
			blank: true,
		},
		{
			name:  "whitespace and comments",
			m:     &Manifest{content: []byte("# a comment \n  # another comment\n")},
			blank: true,
		},
		{
			name: "whitespace, comments and separators",
			m: &Manifest{content: []byte(`---
# a comment

  # another comment
--- # inline comment

`)},
			blank: true,
		},
		{
			name: "whitespace, comments, separators and key-value pair",
			m: &Manifest{content: []byte(`---
# a comment
somekey: somevalue
  # another comment
--- # inline comment

`)},
		},
		{
			name: "configmap",
			m: &Manifest{content: []byte(`---
apiVersion: v1
kind: ConfigMap
metadata:
  name: kcm-chart
  labels:
    app.kubernetes.io/name: chart
    helm.sh/chart: cluster-0.1.0
    app.kubernetes.io/instance: kcm
data:
  SOMEVAR: someval

`)},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.blank, tc.m.IsBlank())
		})
	}
}

func TestManifest_Content(t *testing.T) {
	renderedTemplates := map[string]string{
		"templates/NOTES.txt": "this note should not appear in the manifest content",
		"templates/job.yaml": `---
apiVersion: v1
kind: Job
metadata:
  name: another-job
  annotations:
    kcm/hooks: post-apply
  labels:
    app.kubernetes.io/name: chart
    helm.sh/chart: cluster-0.1.0
    app.kubernetes.io/instance: kcm
spec: {}
`,
		"templates/misc.yaml": `---
apiVersion: v1
kind: Deployment
metadata:
  name: kcm-chart
  labels:
    app.kubernetes.io/name: chart
    helm.sh/chart: cluster-0.1.0
    app.kubernetes.io/instance: kcm
spec: {}
---
apiVersion: v1
kind: Job
metadata:
  name: install-job
  annotations:
    kcm/hooks: post-apply
  labels:
    app.kubernetes.io/name: chart
    helm.sh/chart: cluster-0.1.0
    app.kubernetes.io/instance: kcm
spec: {}
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: kcm-chart
  labels:
    app.kubernetes.io/name: chart
    helm.sh/chart: cluster-0.1.0
    app.kubernetes.io/instance: kcm
data:
  SOMEVAR: someval
---
apiVersion: v1
kind: Job
metadata:
  name: deletion-job
  annotations:
    kcm/hooks: pre-delete
  labels:
    app.kubernetes.io/name: chart
    helm.sh/chart: cluster-0.1.0
    app.kubernetes.io/instance: kcm
spec: {}
`,
	}

	expected := `---
apiVersion: v1
data:
  SOMEVAR: someval
kind: ConfigMap
metadata:
  labels:
    app.kubernetes.io/instance: kcm
    app.kubernetes.io/name: chart
    helm.sh/chart: cluster-0.1.0
  name: kcm-chart

---
apiVersion: v1
kind: Deployment
metadata:
  labels:
    app.kubernetes.io/instance: kcm
    app.kubernetes.io/name: chart
    helm.sh/chart: cluster-0.1.0
  name: kcm-chart
spec: {}

---
apiVersion: v1
kind: Job
metadata:
  annotations:
    kcm/hooks: post-apply
  labels:
    app.kubernetes.io/instance: kcm
    app.kubernetes.io/name: chart
    helm.sh/chart: cluster-0.1.0
  name: another-job
spec: {}

---
apiVersion: v1
kind: Job
metadata:
  annotations:
    kcm/hooks: post-apply
  labels:
    app.kubernetes.io/instance: kcm
    app.kubernetes.io/name: chart
    helm.sh/chart: cluster-0.1.0
  name: install-job
spec: {}

---
apiVersion: v1
kind: Job
metadata:
  annotations:
    kcm/hooks: pre-delete
  labels:
    app.kubernetes.io/instance: kcm
    app.kubernetes.io/name: chart
    helm.sh/chart: cluster-0.1.0
  name: deletion-job
spec: {}

`

	m, err := NewManifest("foo", renderedTemplates)

	require.NoError(t, err)

	assert.Equal(t, expected, string(m.Content()))
}

func TestReadDir(t *testing.T) {
	expected := []byte(`---
kind: Pod
metadata:
  name: foo
  namespace: bar

`)

	manifests, err := ReadDir("testdata/manifests")

	require.NoError(t, err)
	require.Len(t, manifests, 1)
	assert.Equal(t, "foo", manifests[0].Name)
	assert.Equal(t, expected, manifests[0].Content())
}

type testRenderer struct{}

func (r *testRenderer) Render(dir string, v map[string]interface{}) (map[string]string, error) {
	tpl := filepath.Join(filepath.Base(dir), "template.yaml")
	return map[string]string{
		tpl: `
apiVersion: v1
kind: Pod
metadata:
  name: foo
  namespace: bar

`,
	}, nil
}

func TestRenderDir(t *testing.T) {
	r := &testRenderer{}

	manifests, err := RenderDir(r, "testdata/components", nil)

	require.NoError(t, err)
	require.Len(t, manifests, 2)
	assert.Equal(t, "one", manifests[0].Name)
	assert.Equal(t, "two", manifests[1].Name)
}

func TestFindMatching(t *testing.T) {
	manifests := []*Manifest{
		{Name: "foo"},
		{Name: "bar"},
	}

	m, ok := FindMatching(manifests, &Manifest{Name: "bar"})

	assert.True(t, ok)
	assert.Equal(t, "bar", m.Name)

	_, ok = FindMatching(manifests, &Manifest{Name: "baz"})

	assert.False(t, ok)
}
