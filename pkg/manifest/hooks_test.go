package manifest

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewHook(t *testing.T) {
	r := &Resource{
		Name: "foo",
		Kind: "Job",
	}

	annotations := map[string]string{
		HooksAnnotation:      "pre-apply, post-delete ",
		HookPolicyAnnotation: "foo",
	}

	hook, err := newHook(r, annotations)

	require.NoError(t, err)
	assert.Equal(t, []HookType{HookTypePreApply, HookTypePostDelete}, hook.types)
	assert.Equal(t, HookPolicy("foo"), hook.policy)
}

func TestNewHookError(t *testing.T) {
	r := &Resource{
		Name: "foo",
		Kind: "StatefulSet",
	}

	_, err := newHook(r, map[string]string{})

	assert.Error(t, err)
}
