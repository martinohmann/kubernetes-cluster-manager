package hook

import (
	"testing"

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
		Annotation:       "pre-apply",
		PolicyAnnotation: "foo",
	}

	hook, err := New(r, annotations)

	require.NoError(t, err)
	assert.Equal(t, TypePreApply, hook.Type)
	assert.Equal(t, "foo", hook.Policy)
}

func TestNewError(t *testing.T) {
	r := &resource.Resource{
		Name: "foo",
		Kind: "StatefulSet",
	}

	_, err := New(r, map[string]string{})

	assert.Error(t, err)
}
