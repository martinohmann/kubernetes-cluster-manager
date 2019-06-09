package diff

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrinter(t *testing.T) {
	cases := []struct {
		description string
		o           Options
		expected    string
	}{
		{
			description: "no diff",
			o:           Options{A: []byte("foo"), B: []byte("foo")},
			expected:    ``,
		},
		{
			description: "plain diff",
			o:           Options{A: []byte("foo"), B: []byte("bar")},
			expected: `
  @@ -1 +1 @@
  -foo
  +bar
`,
		},
		{
			description: "no diff with filename",
			o:           Options{A: []byte("foo"), B: []byte("foo"), Filename: "baz.yaml"},
			expected:    ``,
		},
		{
			description: "diff with filename",
			o:           Options{A: []byte("foo"), B: []byte("bar"), Filename: "baz.yaml"},
			expected: `changes to baz.yaml:

  --- baz.yaml
  +++ baz.yaml
  @@ -1 +1 @@
  -foo
  +bar
`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.description, func(t *testing.T) {
			var buf bytes.Buffer

			p := NewPrinter(&buf)

			p.Print(tc.o)

			assert.Equal(t, tc.expected, buf.String())
		})
	}
}
