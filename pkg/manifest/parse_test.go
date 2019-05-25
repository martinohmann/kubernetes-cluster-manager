package manifest

import (
	"testing"

	"github.com/martinohmann/kubernetes-cluster-manager/pkg/hook"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/resource"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	cases := []struct {
		description       string
		buf               []byte
		expectError       bool
		expectedResources resource.Slice
		expectedHooks     hook.SliceMap
	}{
		{
			description:       "empty",
			expectedResources: resource.Slice{},
			expectedHooks:     hook.SliceMap{},
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
			expectedResources: resource.Slice{},
			expectedHooks:     hook.SliceMap{},
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
			expectedResources: resource.Slice{},
			expectedHooks:     hook.SliceMap{},
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
			expectedResources: resource.Slice{
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
			expectedHooks: hook.SliceMap{},
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
			expectedResources: resource.Slice{
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
			expectedHooks: hook.SliceMap{},
		},
		{
			description: "unsupported hook resource kind",
			buf: []byte(`
apiVersion: v1
kind: StatefulSet
metadata:
  annotations:
    kcm/hook: pre-delete
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
    kcm/hook: pre-delete
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
    kcm/hook: pre-delete
  labels:
    app.kubernetes.io/instance: kcm
    app.kubernetes.io/name: chart
    helm.sh/chart: cluster-0.1.0
  name: deletion-job
spec: {}
`),
			expectedResources: resource.Slice{},
			expectedHooks: hook.SliceMap{
				hook.TypePreDelete: hook.Slice{
					{
						Resource: &resource.Resource{
							Name: "deletion-job",
							Kind: "Job",
							Content: []byte(`apiVersion: v1
kind: Job
metadata:
  annotations:
    kcm/hook: pre-delete
  labels:
    app.kubernetes.io/instance: kcm
    app.kubernetes.io/name: chart
    helm.sh/chart: cluster-0.1.0
  name: deletion-job
spec: {}
`),
						},
						Type: hook.TypePreDelete,
					},
					{
						Resource: &resource.Resource{
							Name: "deletion-job2",
							Kind: "Job",
							Content: []byte(`apiVersion: v1
kind: Job
metadata:
  annotations:
    kcm/hook: pre-delete
  labels:
    app.kubernetes.io/instance: kcm
    app.kubernetes.io/name: chart
    helm.sh/chart: cluster-0.1.0
  name: deletion-job2
spec: {}
`),
						},
						Type: hook.TypePreDelete,
					},
				},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.description, func(t *testing.T) {
			r, h, err := Parse(tc.buf)

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
