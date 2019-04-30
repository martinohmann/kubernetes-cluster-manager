package file

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChangeSetApply(t *testing.T) {
	f, err := NewTempFile("foo.yaml", []byte(`bar`))
	if !assert.NoError(t, err) {
		return
	}

	defer os.Remove(f.Name())

	c, err := NewChangeSet(f.Name(), []byte(`baz`))
	assert.NoError(t, err)

	if !assert.NoError(t, c.Apply()) {
		return
	}

	buf, err := ioutil.ReadFile(f.Name())
	if assert.NoError(t, err) {
		assert.Equal(t, []byte(`baz`), buf)
	}
}

func TestChangeSetDiff(t *testing.T) {
	c, err := NewChangeSet("foo.yaml", []byte(`foo`))
	assert.NoError(t, err)

	expected := `--- foo.yaml
+++ foo.yaml
@@ -1 +1 @@
-
+foo
`

	assert.Equal(t, expected, c.Diff())
}
