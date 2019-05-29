package hook

import (
	"testing"
	"time"

	"github.com/martinohmann/kubernetes-cluster-manager/pkg/resource"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	r := &resource.Resource{
		Name: "foo",
		Kind: "Job",
	}

	annotations := map[string]string{
		Annotation:            TypePreCreate,
		WaitForAnnotation:     "condition=complete",
		WaitTimeoutAnnotation: "100s",
	}

	hook, err := New(r, annotations)

	require.NoError(t, err)
	assert.Equal(t, TypePreCreate, hook.Type)
	assert.Equal(t, "condition=complete", hook.WaitFor)
	assert.Equal(t, 100*time.Second, hook.WaitTimeout)
}

func TestNewError(t *testing.T) {
	r := &resource.Resource{
		Name: "foo",
		Kind: "StatefulSet",
	}

	_, err := New(r, map[string]string{})

	assert.Error(t, err)
}
