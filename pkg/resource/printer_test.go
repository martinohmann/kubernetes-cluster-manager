package resource

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPrinter_PrintSlice(t *testing.T) {
	cases := []struct {
		description string
		s           Slice
		expected    string
	}{
		{
			description: "empty",
		},
		{
			description: "one resource",
			s:           Slice{{Name: "foo", Namespace: "bar", Kind: StatefulSet, hint: Addition}},
			expected:    "1 resource (+ addition: 1)\n\n  + bar/statefulset/foo\n\n",
		},
		{
			description: "multiple resources",
			s: Slice{
				{Name: "foo", Namespace: "bar", Kind: StatefulSet, hint: Addition},
				{Name: "bar", Namespace: "baz", Kind: PersistentVolumeClaim},
				{Name: "baz", Namespace: "qux", Kind: Job, hint: Removal},
				{Name: "qux", Kind: StatefulSet, hint: Update, contentHint: []byte("old"), Content: []byte("new")},
			},
			expected: `4 resources (* no change: 1, + addition: 1, ~ update: 1, - removal: 1)

  + bar/statefulset/foo

  * baz/persistentvolumeclaim/bar

  - qux/job/baz

  ~ statefulset/qux

  @@ -1 +1 @@
  -old
  +new

`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.description, func(t *testing.T) {
			buf := bytes.NewBuffer(nil)
			p := NewPrinter(buf)

			require.NoError(t, p.PrintSlice(tc.s))
			assert.Equal(t, tc.expected, buf.String())
		})
	}
}
