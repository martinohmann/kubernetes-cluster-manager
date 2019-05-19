package revision

import (
	"testing"

	"github.com/martinohmann/kubernetes-cluster-manager/pkg/hook"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/manifest"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/resource"
	"github.com/stretchr/testify/assert"
)

func TestNewSlice(t *testing.T) {
	cases := []struct {
		name          string
		current, next []*manifest.Manifest
		expected      Slice
		hasNext       bool
	}{
		{
			name:     "empty",
			expected: Slice{},
		},
		{
			name:    "one removed",
			current: []*manifest.Manifest{{Name: "one"}},
			expected: Slice{
				{
					Current: &manifest.Manifest{Name: "one"},
				},
			},
		},
		{
			name:    "present in both",
			current: []*manifest.Manifest{{Name: "one"}},
			next:    []*manifest.Manifest{{Name: "one"}},
			expected: Slice{
				{
					Current: &manifest.Manifest{Name: "one"},
					Next:    &manifest.Manifest{Name: "one"},
				},
			},
		},
		{
			name: "one added",
			next: []*manifest.Manifest{{Name: "one"}},
			expected: Slice{
				{
					Next: &manifest.Manifest{Name: "one"},
				},
			},
		},
		{
			name:    "one added, one removed",
			current: []*manifest.Manifest{{Name: "one"}},
			next:    []*manifest.Manifest{{Name: "two"}},
			expected: Slice{
				{
					Current: &manifest.Manifest{Name: "one"},
				},
				{
					Next: &manifest.Manifest{Name: "two"},
				},
			},
		},
		{
			name:    "one added, one removed, one in both",
			current: []*manifest.Manifest{{Name: "three"}, {Name: "one"}},
			next:    []*manifest.Manifest{{Name: "two"}, {Name: "three"}},
			expected: Slice{
				{
					Current: &manifest.Manifest{Name: "three"},
					Next:    &manifest.Manifest{Name: "three"},
				},
				{
					Current: &manifest.Manifest{Name: "one"},
				},
				{
					Next: &manifest.Manifest{Name: "two"},
				},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			actual := NewSlice(tc.current, tc.next)

			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestSlice_Reverse(t *testing.T) {
	s := Slice{
		{Current: &manifest.Manifest{Name: "foo"}},
		{Current: &manifest.Manifest{Name: "bar"}},
		{Current: &manifest.Manifest{Name: "baz"}},
	}

	expected := Slice{
		{Current: &manifest.Manifest{Name: "baz"}},
		{Current: &manifest.Manifest{Name: "bar"}},
		{Current: &manifest.Manifest{Name: "foo"}},
	}

	assert.Equal(t, expected, s.Reverse())
}

func TestRevision_Types(t *testing.T) {
	cases := []struct {
		description               string
		revision                  *Revision
		initial, upgrade, removal bool
	}{
		{
			description: "empty",
			revision:    &Revision{},
		},
		{
			description: "initial",
			revision:    &Revision{Next: &manifest.Manifest{}},
			initial:     true,
		},
		{
			description: "upgrade",
			revision:    &Revision{Current: &manifest.Manifest{}, Next: &manifest.Manifest{}},
			upgrade:     true,
		},
		{
			description: "removal",
			revision:    &Revision{Current: &manifest.Manifest{}},
			removal:     true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.description, func(t *testing.T) {
			assert.Equal(t, tc.initial, tc.revision.IsInitial())
			assert.Equal(t, tc.upgrade, tc.revision.IsUpgrade())
			assert.Equal(t, tc.removal, tc.revision.IsRemoval())
		})
	}
}

func TestRevision_ChangeSet(t *testing.T) {
	cases := []struct {
		description string
		revision    *Revision
		added       resource.Slice
		changed     resource.Slice
		unchanged   resource.Slice
		removed     resource.Slice
		hooks       hook.SliceMap
	}{
		{
			description: "removal",
			revision: &Revision{
				Current: &manifest.Manifest{
					Resources: testResourceNameSlice("bar"),
					Hooks: hook.SliceMap{
						hook.TypePreApply: testHookNameSlice("baz"),
					},
				},
			},
			removed: testResourceNameSlice("bar"),
			hooks: hook.SliceMap{
				hook.TypePreApply: testHookNameSlice("baz"),
			},
		},
		{
			description: "initial",
			revision: &Revision{
				Next: &manifest.Manifest{
					Resources: testResourceNameSlice("bar"),
					Hooks: hook.SliceMap{
						hook.TypePreApply: testHookNameSlice("baz"),
					},
				},
			},
			added: testResourceNameSlice("bar"),
			hooks: hook.SliceMap{
				hook.TypePreApply: testHookNameSlice("baz"),
			},
		},
		{
			description: "upgrade",
			revision: &Revision{
				Current: &manifest.Manifest{
					Resources: testResourceNameSlice("foo", "bar", "qux"),
					Hooks: hook.SliceMap{
						hook.TypePreApply: testHookNameSlice("bar"),
					},
				},
				Next: &manifest.Manifest{
					Resources: append(testResourceNameSlice("bar", "baz"), &resource.Resource{Name: "qux", Content: []byte("---\nchanges")}),
					Hooks: hook.SliceMap{
						hook.TypePreApply: testHookNameSlice("baz"),
					},
				},
			},
			added:     testResourceNameSlice("baz"),
			changed:   resource.Slice{{Name: "qux", Content: []byte("---\nchanges")}},
			unchanged: testResourceNameSlice("bar"),
			removed:   testResourceNameSlice("foo"),
			hooks: hook.SliceMap{
				hook.TypePreApply: testHookNameSlice("baz"),
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.description, func(t *testing.T) {
			c := tc.revision.ChangeSet()

			assert.Equal(t, tc.added, c.AddedResources)
			assert.Equal(t, tc.changed, c.ChangedResources)
			assert.Equal(t, tc.unchanged, c.UnchangedResources)
			assert.Equal(t, tc.removed, c.RemovedResources)
			assert.Equal(t, tc.hooks, c.Hooks)
		})
	}
}

func testResourceNameSlice(names ...string) resource.Slice {
	s := make(resource.Slice, len(names))

	for i, name := range names {
		s[i] = &resource.Resource{Name: name}
	}

	return s
}

func testHookNameSlice(names ...string) hook.Slice {
	s := make(hook.Slice, len(names))

	for i, name := range names {
		s[i] = &hook.Hook{Resource: &resource.Resource{Name: name}}
	}

	return s
}
