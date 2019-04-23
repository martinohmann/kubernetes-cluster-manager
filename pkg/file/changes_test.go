package file

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChangesApply(t *testing.T) {
	f, err := NewTempFile("foo.yaml", []byte(`bar`))
	if !assert.NoError(t, err) {
		return
	}

	defer os.Remove(f.Name())

	c := NewChanges(f.Name(), []byte(`baz`))

	if !assert.NoError(t, c.Apply()) {
		return
	}

	buf, err := ioutil.ReadFile(f.Name())
	if assert.NoError(t, err) {
		assert.Equal(t, []byte(`baz`), buf)
	}
}

func TestChangesDiff(t *testing.T) {
	c := NewChanges("foo.yaml", []byte(`foo`))

	diff, err := c.Diff()

	assert.NoError(t, err)

	expected := `--- foo.yaml
+++ foo.yaml
@@ -1 +1 @@
-
+foo
`

	assert.Equal(t, expected, diff)
}
