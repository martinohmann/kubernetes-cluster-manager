package hook

import (
	"testing"
	"time"

	"github.com/martinohmann/kubernetes-cluster-manager/pkg/resource"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	cases := []struct {
		description string
		resource    *resource.Resource
		annotations map[string]string
		expected    *Hook
		expectError bool
	}{
		{
			description: "unsupported resource kind",
			resource:    &resource.Resource{Name: "foo", Kind: resource.StatefulSet},
			expectError: true,
		},
		{
			description: "missing hook annotation",
			resource:    &resource.Resource{Name: "foo", Kind: resource.Job},
			expectError: true,
		},
		{
			description: "invalid hook annotation",
			resource:    &resource.Resource{Name: "foo", Kind: resource.Job},
			annotations: map[string]string{Annotation: "nonexistent-hook-type"},
			expectError: true,
		},
		{
			description: "invalid wait timeout",
			resource:    &resource.Resource{Name: "foo", Kind: resource.Job},
			annotations: map[string]string{
				Annotation:            PreCreate,
				WaitTimeoutAnnotation: "bar",
			},
			expectError: true,
		},
		{
			description: "valid hook with wait condition",
			resource:    &resource.Resource{Name: "foo", Kind: resource.Job},
			annotations: map[string]string{
				Annotation:        PreCreate,
				WaitForAnnotation: "condition=complete",
			},
			expected: &Hook{
				Resource: &resource.Resource{Name: "foo", Kind: resource.Job},
				Type:     PreCreate,
				WaitFor:  "condition=complete",
			},
		},
		{
			description: "valid hook with wait condition and timeout",
			resource:    &resource.Resource{Name: "foo", Kind: resource.Job},
			annotations: map[string]string{
				Annotation:            PreCreate,
				WaitForAnnotation:     "condition=complete",
				WaitTimeoutAnnotation: "100s",
			},
			expected: &Hook{
				Resource:    &resource.Resource{Name: "foo", Kind: resource.Job},
				Type:        PreCreate,
				WaitFor:     "condition=complete",
				WaitTimeout: 100 * time.Second,
			},
		},
		{
			description: "valid hook with wait condition, timeout and delete-after-completion policy",
			resource:    &resource.Resource{Name: "foo", Kind: resource.Job},
			annotations: map[string]string{
				Annotation:            PreCreate,
				WaitForAnnotation:     "condition=complete",
				WaitTimeoutAnnotation: "100s",
				PolicyAnnotation:      DeleteAfterCompletionPolicy,
			},
			expected: &Hook{
				Resource:              &resource.Resource{Name: "foo", Kind: resource.Job},
				Type:                  PreCreate,
				WaitFor:               "condition=complete",
				WaitTimeout:           100 * time.Second,
				DeleteAfterCompletion: true,
			},
		},
		{
			description: "missing wait-for condition when delete-after-completion policy defined",
			resource:    &resource.Resource{Name: "foo", Kind: resource.Job},
			annotations: map[string]string{
				Annotation:       PreCreate,
				PolicyAnnotation: DeleteAfterCompletionPolicy,
			},
			expectError: true,
		},
		{
			description: "invalid hook policy",
			resource:    &resource.Resource{Name: "foo", Kind: resource.Job},
			annotations: map[string]string{
				Annotation:            PreCreate,
				WaitForAnnotation:     "condition=complete",
				WaitTimeoutAnnotation: "100s",
				PolicyAnnotation:      "foo,bar",
			},
			expectError: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.description, func(t *testing.T) {
			h, err := New(tc.resource, tc.annotations)

			if tc.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expected, h)
			}
		})
	}
}
