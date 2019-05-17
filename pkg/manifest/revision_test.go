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

func TestRevisionSlice_Reverse(t *testing.T) {
	s := RevisionSlice{
		{Current: &Manifest{Name: "foo"}},
		{Current: &Manifest{Name: "bar"}},
		{Current: &Manifest{Name: "baz"}},
	}

	expected := RevisionSlice{
		{Current: &Manifest{Name: "baz"}},
		{Current: &Manifest{Name: "bar"}},
		{Current: &Manifest{Name: "foo"}},
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
			revision:    &Revision{Next: &Manifest{}},
			initial:     true,
		},
		{
			description: "upgrade",
			revision:    &Revision{Current: &Manifest{}, Next: &Manifest{}},
			upgrade:     true,
		},
		{
			description: "removal",
			revision:    &Revision{Current: &Manifest{}},
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
		added       ResourceSlice
		changed     ResourceSlice
		unchanged   ResourceSlice
		removed     ResourceSlice
		hooks       HookSliceMap
	}{
		{
			description: "removal",
			revision: &Revision{
				Current: &Manifest{
					resources: testResourceNameSlice("bar"),
					hooks: HookSliceMap{
						HookTypePreApply: testHookNameSlice("baz"),
					},
				},
			},
			removed: testResourceNameSlice("bar"),
			hooks: HookSliceMap{
				HookTypePreApply: testHookNameSlice("baz"),
			},
		},
		{
			description: "initial",
			revision: &Revision{
				Next: &Manifest{
					resources: testResourceNameSlice("bar"),
					hooks: HookSliceMap{
						HookTypePreApply: testHookNameSlice("baz"),
					},
				},
			},
			added: testResourceNameSlice("bar"),
			hooks: HookSliceMap{
				HookTypePreApply: testHookNameSlice("baz"),
			},
		},
		{
			description: "upgrade",
			revision: &Revision{
				Current: &Manifest{
					resources: testResourceNameSlice("foo", "bar", "qux"),
					hooks: HookSliceMap{
						HookTypePreApply: testHookNameSlice("bar"),
					},
				},
				Next: &Manifest{
					resources: append(testResourceNameSlice("bar", "baz"), &Resource{Name: "qux", Content: []byte("---\nchanges")}),
					hooks: HookSliceMap{
						HookTypePreApply: testHookNameSlice("baz"),
					},
				},
			},
			added:     testResourceNameSlice("baz"),
			changed:   ResourceSlice{{Name: "qux", Content: []byte("---\nchanges")}},
			unchanged: testResourceNameSlice("bar"),
			removed:   testResourceNameSlice("foo"),
			hooks: HookSliceMap{
				HookTypePreApply: testHookNameSlice("baz"),
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

func testResourceNameSlice(names ...string) ResourceSlice {
	s := make(ResourceSlice, len(names))

	for i, name := range names {
		s[i] = &Resource{Name: name}
	}

	return s
}

func testHookNameSlice(names ...string) HookSlice {
	s := make(HookSlice, len(names))

	for i, name := range names {
		s[i] = &Hook{Resource: &Resource{Name: name}}
	}

	return s
}
