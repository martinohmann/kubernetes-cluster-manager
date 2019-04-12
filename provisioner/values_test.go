package provisioner

import (
	"os"
	"testing"

	"github.com/martinohmann/kubernetes-cluster-manager/pkg/fs"
	"github.com/stretchr/testify/assert"
)

func TestLoadValues(t *testing.T) {
	content := []byte("---\nfoo: bar")
	f, err := fs.NewTempFile("values.yaml", content)
	if !assert.NoError(t, err) {
		return
	}

	defer os.Remove(f.Name())

	values, err := loadValues(f.Name())

	if assert.NoError(t, err) {
		assert.Equal(t, "bar", values["foo"])
	}
}