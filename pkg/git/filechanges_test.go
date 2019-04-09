package git

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/martinohmann/kubernetes-cluster-manager/pkg/fs"
	"github.com/stretchr/testify/assert"
)

func TestFileChangesApply(t *testing.T) {
	f, err := fs.NewTempFile("foo.yaml", []byte(`bar`))
	if !assert.NoError(t, err) {
		return
	}

	defer os.Remove(f.Name())

	c, err := NewFileChanges(f.Name(), []byte(`baz`))
	if !assert.NoError(t, err) {
		return
	}

	defer c.Close()

	err = c.Apply()
	if !assert.NoError(t, err) {
		return
	}

	buf, err := ioutil.ReadFile(f.Name())
	if assert.NoError(t, err) {
		assert.Equal(t, []byte(`baz`), buf)
	}
}
