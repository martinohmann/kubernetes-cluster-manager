package revision

import (
	"context"
	"sync/atomic"
	"testing"

	"github.com/martinohmann/kubernetes-cluster-manager/pkg/hook"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/kubernetes"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/manifest"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/resource"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockClient struct {
	applyCalled, deleteCalled, waitCalled uint64
}

func (c *mockClient) ApplyManifest(ctx context.Context, buf []byte) error {
	atomic.AddUint64(&c.applyCalled, 1)
	return nil
}

func (c *mockClient) DeleteManifest(ctx context.Context, buf []byte) error {
	atomic.AddUint64(&c.deleteCalled, 1)
	return nil
}

func (c *mockClient) Wait(ctx context.Context, o kubernetes.WaitOptions) error {
	atomic.AddUint64(&c.waitCalled, 1)
	return nil
}

func TestUpgrader_Upgrade(t *testing.T) {
	client := &mockClient{}

	resources1 := resource.Slice{
		{
			Kind:      "ConfigMap",
			Name:      "bar",
			Namespace: "baz",
			Content: []byte(`apiVersion: v1
kind: ConfigMap
metadata:
  name: bar
  namespace: baz
`),
		},
	}

	resources2 := resource.Slice{
		{
			Kind:      "Pod",
			Name:      "bar",
			Namespace: "baz",
			Content: []byte(`apiVersion: v1
kind: Pod
metadata:
  name: bar
  namespace: baz
`),
		},
	}

	current1 := &manifest.Manifest{
		Name:      "foo",
		Resources: resources1,
	}

	next1 := &manifest.Manifest{
		Name:      "foo",
		Resources: resources1,
	}

	current2 := &manifest.Manifest{
		Name:      "foo",
		Resources: resources1,
	}

	next2 := &manifest.Manifest{
		Name:      "foo",
		Resources: resources2,
	}

	rev1 := &Revision{
		Current: current1,
	}

	rev2 := &Revision{
		Next: next1,
	}

	rev3 := &Revision{
		Current: current2,
		Next:    next2,
	}

	u := NewUpgrader(client, &UpgraderOptions{NoSave: true})

	err := u.Upgrade(context.Background(), rev1)

	require.NoError(t, err)

	assert.Equal(t, uint64(1), client.deleteCalled)

	err = u.Upgrade(context.Background(), rev2)

	require.NoError(t, err)

	assert.Equal(t, uint64(1), client.applyCalled)

	err = u.Upgrade(context.Background(), rev3)

	require.NoError(t, err)

	assert.Equal(t, uint64(2), client.applyCalled)
	assert.Equal(t, uint64(2), client.deleteCalled)
}

func TestUpgrader_execHooks(t *testing.T) {
	client := &mockClient{}

	hooks := hook.Slice{
		{
			Type: hook.TypePreCreate,
			Resource: &resource.Resource{
				Name: "foo",
				Kind: "Job",
			},
			WaitFor:               "condition=complete",
			DeleteAfterCompletion: true,
		},
		{
			Type: hook.TypePreCreate,
			Resource: &resource.Resource{
				Name: "bar",
				Kind: "Job",
			},
		},
		{
			Type: hook.TypePreCreate,
			Resource: &resource.Resource{
				Name: "baz",
				Kind: "Job",
			},
			WaitFor: "condition=complete",
		},
	}

	u := &upgrader{client: client, noSave: true}

	err := u.execHooks(context.Background(), hooks)

	require.NoError(t, err)

	assert.Equal(t, uint64(2), client.deleteCalled)
	assert.Equal(t, uint64(2), client.waitCalled)
	assert.Equal(t, uint64(1), client.applyCalled)
}
