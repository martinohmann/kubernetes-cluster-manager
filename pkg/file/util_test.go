package file

import (
	"os"
	"testing"

	"github.com/martinohmann/kubernetes-cluster-manager/pkg/kcm"
	"github.com/stretchr/testify/assert"
)

func TestLoadYAML(t *testing.T) {
	content := []byte("---\npreApply:\n- kind: pod\n  name: foo")
	f, err := NewTempFile("deletions.yaml", content)
	if !assert.NoError(t, err) {
		return
	}

	defer os.Remove(f.Name())

	deletions := &kcm.Deletions{}

	if !assert.NoError(t, LoadYAML(f.Name(), deletions)) {
		return
	}

	if assert.Len(t, deletions.PreApply, 1) {
		assert.Equal(t, "pod", deletions.PreApply[0].Kind)
		assert.Equal(t, "foo", deletions.PreApply[0].Name)
	}
}
