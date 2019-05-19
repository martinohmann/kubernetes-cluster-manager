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
		HooksAnnotation:  "pre-apply, post-delete ",
		PolicyAnnotation: "foo",
	}

	hook, err := New(r, annotations)

	require.NoError(t, err)
	assert.Equal(t, []Type{TypePreApply, TypePostDelete}, hook.Types)
	assert.Equal(t, Policy("foo"), hook.policy)
}

func TestNewError(t *testing.T) {
	r := &resource.Resource{
		Name: "foo",
		Kind: "StatefulSet",
	}

	_, err := New(r, map[string]string{})

	assert.Error(t, err)
}
