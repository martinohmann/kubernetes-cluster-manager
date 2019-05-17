package manifest

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateRevisions(t *testing.T) {
	cases := []struct {
		name          string
		current, next []*Manifest
		expected      RevisionSlice
		hasNext       bool
	}{
		{
			name:     "empty",
			expected: RevisionSlice{},
		},
		{
			name:    "one removed",
			current: []*Manifest{{Name: "one"}},
			expected: RevisionSlice{
				{
					Current: &Manifest{Name: "one"},
				},
			},
		},
		{
			name:    "present in both",
			current: []*Manifest{{Name: "one"}},
			next:    []*Manifest{{Name: "one"}},
			expected: RevisionSlice{
				{
					Current: &Manifest{Name: "one"},
					Next:    &Manifest{Name: "one"},
				},
			},
		},
		{
			name: "one added",
			next: []*Manifest{{Name: "one"}},
			expected: RevisionSlice{
				{
					Next: &Manifest{Name: "one"},
				},
			},
		},
		{
			name:    "one added, one removed",
			current: []*Manifest{{Name: "one"}},
			next:    []*Manifest{{Name: "two"}},
			expected: RevisionSlice{
				{
					Current: &Manifest{Name: "one"},
				},
				{
					Next: &Manifest{Name: "two"},
				},
			},
		},
		{
			name:    "one added, one removed, one in both",
			current: []*Manifest{{Name: "three"}, {Name: "one"}},
			next:    []*Manifest{{Name: "two"}, {Name: "three"}},
			expected: RevisionSlice{
				{
					Current: &Manifest{Name: "three"},
					Next:    &Manifest{Name: "three"},
				},
				{
					Current: &Manifest{Name: "one"},
				},
				{
					Next: &Manifest{Name: "two"},
				},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			actual := CreateRevisions(tc.current, tc.next)

			assert.Equal(t, tc.expected, actual)
		})
	}
}
