package hook

import (
	"testing"

	"github.com/martinohmann/kubernetes-cluster-manager/pkg/resource"
	"github.com/stretchr/testify/assert"
)

func TestSliceMap_Bytes(t *testing.T) {
	m := SliceMap{
		PreDelete: Slice{
			{Resource: &resource.Resource{Content: []byte(`apiVersion: v1
kind: Job
metadata:
  name: delete-job
`)}},
		},
		PostCreate: Slice{
			{Resource: &resource.Resource{Content: []byte(`apiVersion: v1
kind: Job
metadata:
  name: qux
`)}},
			{Resource: &resource.Resource{Content: []byte(`apiVersion: v1
kind: Job
metadata:
  name: foo
`)}},
			{Resource: &resource.Resource{Content: []byte(`apiVersion: v1
kind: Job
metadata:
  name: bar
`)}},
		},
	}

	expected := []byte(`---
apiVersion: v1
kind: Job
metadata:
  name: qux

---
apiVersion: v1
kind: Job
metadata:
  name: foo

---
apiVersion: v1
kind: Job
metadata:
  name: bar

---
apiVersion: v1
kind: Job
metadata:
  name: delete-job

`)

	assert.Equal(t, string(expected), string(m.Bytes()))
}
