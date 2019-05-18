package manifest

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResource_matches(t *testing.T) {
	cases := []struct {
		description string
		a, b        *Resource
		matches     bool
	}{
		{
			description: "both nil",
			matches:     true,
		},
		{
			description: "a nil",
			b:           &Resource{Name: "foo"},
		},
		{
			description: "b nil",
			a:           &Resource{Name: "foo"},
		},
		{
			description: "different kind",
			a:           &Resource{Name: "foo", Kind: "Deployment"},
			b:           &Resource{Name: "foo", Kind: "Pod"},
		},
		{
			description: "different namespace",
			a:           &Resource{Name: "foo", Kind: "Pod", Namespace: "default"},
			b:           &Resource{Name: "foo", Kind: "Pod", Namespace: "kube-system"},
		},
		{
			description: "same name, kind and namespace",
			a:           &Resource{Name: "foo", Kind: "Pod", Namespace: "kube-system"},
			b:           &Resource{Name: "foo", Kind: "Pod", Namespace: "kube-system"},
			matches:     true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.description, func(t *testing.T) {
			assert.Equal(t, tc.matches, tc.a.matches(tc.b))
		})
	}
}

func TestParseResources(t *testing.T) {
	cases := []struct {
		description       string
		buf               []byte
		expectError       bool
		expectedResources ResourceSlice
		expectedHooks     HookSliceMap
	}{
		{
			description:       "empty",
			expectedResources: ResourceSlice{},
			expectedHooks:     HookSliceMap{},
		},
		{
			description: "missing resource name should be skipped",
			buf: []byte(`
apiVersion: v1
kind: StatefulSet
metadata:
  labels:
    app.kubernetes.io/instance: kcm
    app.kubernetes.io/name: chart
    helm.sh/chart: cluster-0.1.0
spec: {}
`),
			expectedResources: ResourceSlice{},
			expectedHooks:     HookSliceMap{},
		},
		{
			description: "missing resource kind should be skipped",
			buf: []byte(`
apiVersion: v1
metadata:
  labels:
    app.kubernetes.io/instance: kcm
    app.kubernetes.io/name: chart
    helm.sh/chart: cluster-0.1.0
  name: some-statefulset
spec: {}
`),
			expectedResources: ResourceSlice{},
			expectedHooks:     HookSliceMap{},
		},
		{
			description: "resources",
			buf: []byte(`---
apiVersion: v1
kind: Prometheus
metadata:
  name: prometheus
spec: {}
---
apiVersion: v1
kind: Alertmanager
metadata:
  name: alertmanager
spec: {}
---
apiVersion: v1
kind: Pod
metadata:
  name: pod
spec: {}
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: cm
data: {}
---
apiVersion: v1
kind: Prometheus
metadata:
  name: another-prometheus
spec: {}
`),
			expectedResources: ResourceSlice{
				{
					Name: "cm",
					Kind: "ConfigMap",
					Content: []byte(`apiVersion: v1
data: {}
kind: ConfigMap
metadata:
  name: cm
`),
				},
				{
					Name: "pod",
					Kind: "Pod",
					Content: []byte(`apiVersion: v1
kind: Pod
metadata:
  name: pod
spec: {}
`),
				},
				{
					Name: "alertmanager",
					Kind: "Alertmanager",
					Content: []byte(`apiVersion: v1
kind: Alertmanager
metadata:
  name: alertmanager
spec: {}
`),
				},
				{
					Name: "another-prometheus",
					Kind: "Prometheus",
					Content: []byte(`apiVersion: v1
kind: Prometheus
metadata:
  name: another-prometheus
spec: {}
`),
				},
				{
					Name: "prometheus",
					Kind: "Prometheus",
					Content: []byte(`apiVersion: v1
kind: Prometheus
metadata:
  name: prometheus
spec: {}
`),
				},
			},
			expectedHooks: HookSliceMap{},
		},
		{
			description: "unparsable yaml should be ignored",
			buf: []byte(`
just some test
---
apiVersion: v1
kind: StatefulSet
metadata:
  labels:
    app.kubernetes.io/instance: kcm
    app.kubernetes.io/name: chart
    helm.sh/chart: cluster-0.1.0
  name: some-statefulset
spec: {}
`),
			expectedResources: ResourceSlice{
				{
					Name: "some-statefulset",
					Kind: "StatefulSet",
					Content: []byte(`apiVersion: v1
kind: StatefulSet
metadata:
  labels:
    app.kubernetes.io/instance: kcm
    app.kubernetes.io/name: chart
    helm.sh/chart: cluster-0.1.0
  name: some-statefulset
spec: {}
`),
				},
			},
			expectedHooks: HookSliceMap{},
		},
		{
			description: "unsupported hook resource kind",
			buf: []byte(`
apiVersion: v1
kind: StatefulSet
metadata:
  annotations:
    kcm/hooks: pre-delete
  labels:
    app.kubernetes.io/instance: kcm
    app.kubernetes.io/name: chart
    helm.sh/chart: cluster-0.1.0
  name: deletion-job
spec: {}
`),
			expectError: true,
		},
		{
			description: "hooks",
			buf: []byte(`
apiVersion: v1
kind: Job
metadata:
  annotations:
    kcm/hooks: pre-delete
  labels:
    app.kubernetes.io/instance: kcm
    app.kubernetes.io/name: chart
    helm.sh/chart: cluster-0.1.0
  name: deletion-job2
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
`),
			expectedResources: ResourceSlice{},
			expectedHooks: HookSliceMap{
				HookTypePreDelete: HookSlice{
					{
						Resource: &Resource{
							Name: "deletion-job",
							Kind: "Job",
							Content: []byte(`apiVersion: v1
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
`),
						},
						types: []HookType{HookTypePreDelete},
					},
					{
						Resource: &Resource{
							Name: "deletion-job2",
							Kind: "Job",
							Content: []byte(`apiVersion: v1
kind: Job
metadata:
  annotations:
    kcm/hooks: pre-delete
  labels:
    app.kubernetes.io/instance: kcm
    app.kubernetes.io/name: chart
    helm.sh/chart: cluster-0.1.0
  name: deletion-job2
spec: {}
`),
						},
						types: []HookType{HookTypePreDelete},
					},
				},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.description, func(t *testing.T) {
			r, h, err := parseResources(tc.buf)

			if tc.expectError {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)

			assert.Equal(t, tc.expectedResources, r)
			assert.Equal(t, tc.expectedHooks, h)
		})
	}
}
