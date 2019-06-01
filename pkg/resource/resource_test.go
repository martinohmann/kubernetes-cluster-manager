package resource

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	cases := []struct {
		description string
		head        Head
		expected    *Resource
		expectError bool
	}{
		{
			description: "empty head",
			expected:    &Resource{},
		},
		{
			description: "head with metadata",
			head:        Head{Kind: KindJob, Metadata: Metadata{Name: "foo", Namespace: "bar"}},
			expected:    &Resource{Kind: KindJob, Name: "foo", Namespace: "bar"},
		},
		{
			description: "stateful set with valid deletion policy",
			head: Head{
				Kind: KindStatefulSet,
				Metadata: Metadata{
					Name:      "foo",
					Namespace: "bar",
					Annotations: map[string]string{
						AnnotationDeletionPolicy: PolicyDeletePersistentVolumeClaims,
					},
				},
			},
			expected: &Resource{Kind: KindStatefulSet, Name: "foo", Namespace: "bar", DeletePersistentVolumeClaims: true},
		},
		{
			description: "stateful set with invalid deletion policy",
			head: Head{
				Kind: KindStatefulSet,
				Metadata: Metadata{
					Name:      "foo",
					Namespace: "bar",
					Annotations: map[string]string{
						AnnotationDeletionPolicy: "baz",
					},
				},
			},
			expectError: true,
		},
		{
			description: "resource that does not support deletion policy annotation",
			head: Head{
				Kind: KindJob,
				Metadata: Metadata{
					Name:      "foo",
					Namespace: "bar",
					Annotations: map[string]string{
						AnnotationDeletionPolicy: PolicyDeletePersistentVolumeClaims,
					},
				},
			},
			expectError: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.description, func(t *testing.T) {
			r, err := New(nil, tc.head)

			if tc.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expected, r)
			}
		})
	}
}
