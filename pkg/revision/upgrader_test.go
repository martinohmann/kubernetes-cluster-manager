package revision

import (
	"context"
	"sync/atomic"
	"testing"

	"github.com/martinohmann/kubernetes-cluster-manager/pkg/hook"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/kubernetes"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/manifest"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/resource"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockClient struct {
	applyCalled          uint64
	deleteCalled         uint64
	waitCalled           uint64
	deleteResourceCalled uint64
}

func (c *mockClient) ApplyManifest(ctx context.Context, buf []byte) error {
	atomic.AddUint64(&c.applyCalled, 1)
	return nil
}

func (c *mockClient) DeleteManifest(ctx context.Context, buf []byte) error {
	atomic.AddUint64(&c.deleteCalled, 1)
	return nil
}

func (c *mockClient) DeleteResource(ctx context.Context, selector resource.Head) error {
	atomic.AddUint64(&c.deleteResourceCalled, 1)
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
			Kind:      resource.KindStatefulSet,
			Name:      "bar",
			Namespace: "baz",
			Content: []byte(`apiVersion: v1
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: bar
  namespace: baz
  annotations:
    kcm/deletion-policy: delete-pvcs
spec:
  replicas: 2
  volumeClaimTemplates:
  - metadata:
      name: data
`),
			DeletePersistentVolumeClaims: true,
		},
	}

	resources2 := resource.Slice{
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
	assert.Equal(t, uint64(2), client.deleteResourceCalled)

	err = u.Upgrade(context.Background(), rev2)

	require.NoError(t, err)

	assert.Equal(t, uint64(1), client.applyCalled)

	err = u.Upgrade(context.Background(), rev3)

	require.NoError(t, err)

	assert.Equal(t, uint64(2), client.applyCalled)
	assert.Equal(t, uint64(2), client.deleteCalled)
	assert.Equal(t, uint64(4), client.deleteResourceCalled)
}

func TestUpgrader_execHooks(t *testing.T) {
	client := &mockClient{}

	hooks := hook.Slice{
		{
			Type: hook.TypePreCreate,
			Resource: &resource.Resource{
				Name: "foo",
				Kind: resource.KindJob,
			},
			WaitFor:               "condition=complete",
			DeleteAfterCompletion: true,
		},
		{
			Type: hook.TypePreCreate,
			Resource: &resource.Resource{
				Name: "bar",
				Kind: resource.KindJob,
			},
		},
		{
			Type: hook.TypePreCreate,
			Resource: &resource.Resource{
				Name: "baz",
				Kind: resource.KindJob,
			},
			WaitFor: "condition=complete",
		},
	}

	u := &upgrader{client: client, logger: log.NewEntry(log.StandardLogger()), options: &UpgraderOptions{NoSave: true}}

	err := u.execHooks(context.Background(), hooks)

	require.NoError(t, err)

	assert.Equal(t, uint64(2), client.deleteCalled)
	assert.Equal(t, uint64(2), client.waitCalled)
	assert.Equal(t, uint64(1), client.applyCalled)
}

func TestUpgrader_execHooksDryRun(t *testing.T) {
	client := &mockClient{}

	hooks := hook.Slice{
		{
			Type: hook.TypePreCreate,
			Resource: &resource.Resource{
				Name: "foo",
				Kind: resource.KindJob,
			},
			WaitFor:               "condition=complete",
			DeleteAfterCompletion: true,
		},
		{
			Type: hook.TypePreCreate,
			Resource: &resource.Resource{
				Name: "bar",
				Kind: resource.KindJob,
			},
		},
		{
			Type: hook.TypePreCreate,
			Resource: &resource.Resource{
				Name: "baz",
				Kind: resource.KindJob,
			},
			WaitFor: "condition=complete",
		},
	}

	u := &upgrader{client: client, logger: log.NewEntry(log.StandardLogger()), options: &UpgraderOptions{DryRun: true, NoSave: true}}

	err := u.execHooks(context.Background(), hooks)

	require.NoError(t, err)

	assert.Equal(t, uint64(0), client.deleteCalled)
	assert.Equal(t, uint64(0), client.waitCalled)
	assert.Equal(t, uint64(0), client.applyCalled)
}
