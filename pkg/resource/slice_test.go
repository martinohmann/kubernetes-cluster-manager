package resource

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSlice_Bytes(t *testing.T) {
	s := Slice{
		{Content: []byte(`apiVersion: v1
kind: ConfigMap
metadata:
  name: bar
  namespace: baz
`)},
		{Content: []byte(`apiVersion: v1
kind: Pod
metadata:
  name: foo
  namespace: bar
`)},
	}

	expected := []byte(`---
apiVersion: v1
kind: ConfigMap
metadata:
  name: bar
  namespace: baz

---
apiVersion: v1
kind: Pod
metadata:
  name: foo
  namespace: bar

`)

	assert.Equal(t, string(expected), string(s.Bytes()))
}
