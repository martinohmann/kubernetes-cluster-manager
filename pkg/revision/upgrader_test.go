package revision

import (
	"context"
	"fmt"
	"testing"

	"github.com/martinohmann/kubernetes-cluster-manager/pkg/manifest"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/resource"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockClient struct {
	applyCalled, deleteCalled int
}

func (c *mockClient) ApplyManifest(ctx context.Context, buf []byte) error {
	fmt.Println(string(buf))
	c.applyCalled++
	return nil
}

func (c *mockClient) DeleteManifest(ctx context.Context, buf []byte) error {
	c.deleteCalled++
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

	assert.Equal(t, 1, client.deleteCalled)

	err = u.Upgrade(context.Background(), rev2)

	require.NoError(t, err)

	assert.Equal(t, 1, client.applyCalled)

	err = u.Upgrade(context.Background(), rev3)

	require.NoError(t, err)

	assert.Equal(t, 2, client.applyCalled)
	assert.Equal(t, 2, client.deleteCalled)
}
