package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeletionsFilterPending(t *testing.T) {
	d1 := &Deletion{
		Name:      "some-pod",
		Kind:      "pod",
		Namespace: "kube-system",
	}

	d2 := &Deletion{
		Name:      "some-pvc",
		Kind:      "pvc",
		Namespace: "kube-system",
	}

	d := &Deletions{
		PreApply: []*Deletion{d1, d2},
	}

	actual := d.FilterPending()

	assert.Len(t, actual.PreApply, 2)

	d1.MarkDeleted()

	actual = d.FilterPending()

	if assert.Len(t, actual.PreApply, 1) {
		assert.Equal(t, "some-pvc", actual.PreApply[0].Name)
	}
}

func TestDeletionString(t *testing.T) {
	d := &Deletion{
		Name:      "foo",
		Kind:      "pod",
		Namespace: "default",
		Labels:    map[string]string{"app.kubernetes.io/version": "v0.0.1"},
	}

	expected := `kind: pod
name: foo
namespace: default
labels:
  app.kubernetes.io/version: v0.0.1
`

	assert.Equal(t, expected, d.String())
}
