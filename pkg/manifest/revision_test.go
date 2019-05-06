package manifest

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateRevisions(t *testing.T) {
	cases := []struct {
		name       string
		prev, next []*Manifest
		expected   []Revision
		hasNext    bool
	}{
		{
			name:     "empty",
			expected: []Revision{},
		},
		{
			name: "one removed",
			prev: []*Manifest{{Name: "one"}},
			expected: []Revision{
				{
					Prev: &Manifest{Name: "one"},
				},
			},
		},
		{
			name: "present in both",
			prev: []*Manifest{{Name: "one"}},
			next: []*Manifest{{Name: "one"}},
			expected: []Revision{
				{
					Prev: &Manifest{Name: "one"},
					Next: &Manifest{Name: "one"},
				},
			},
		},
		{
			name: "one added",
			next: []*Manifest{{Name: "one"}},
			expected: []Revision{
				{
					Next: &Manifest{Name: "one"},
				},
			},
		},
		{
			name: "one added, one removed",
			prev: []*Manifest{{Name: "one"}},
			next: []*Manifest{{Name: "two"}},
			expected: []Revision{
				{
					Prev: &Manifest{Name: "one"},
				},
				{
					Next: &Manifest{Name: "two"},
				},
			},
		},
		{
			name: "one added, one removed, one in both",
			prev: []*Manifest{{Name: "three"}, {Name: "one"}},
			next: []*Manifest{{Name: "two"}, {Name: "three"}},
			expected: []Revision{
				{
					Prev: &Manifest{Name: "three"},
					Next: &Manifest{Name: "three"},
				},
				{
					Prev: &Manifest{Name: "one"},
				},
				{
					Next: &Manifest{Name: "two"},
				},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			actual := CreateRevisions(tc.prev, tc.next)

			assert.Equal(t, tc.expected, actual)
		})
	}
}
