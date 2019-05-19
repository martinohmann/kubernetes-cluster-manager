package hook

import (
	"testing"

	"github.com/martinohmann/kubernetes-cluster-manager/pkg/resource"
	"github.com/stretchr/testify/assert"
)

func TestSliceMap_Bytes(t *testing.T) {
	m := SliceMap{
		TypePreDelete: Slice{
			{Resource: &resource.Resource{Content: []byte(`apiVersion: v1
kind: Job
metadata:
  name: delete-job
`)}},
		},
		TypePostApply: Slice{
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

func TestSliceMap_Get(t *testing.T) {
	s := Slice{{Resource: &resource.Resource{Kind: "Job", Name: "somejob"}}}
	m := SliceMap{TypePostDelete: s, TypePreApply: Slice{}}

	assert.True(t, m.Has(TypePostDelete))
	assert.False(t, m.Has(TypePreApply))
	assert.False(t, m.Has(TypePostApply))

	assert.Equal(t, s, m.Get(TypePostDelete))
}
