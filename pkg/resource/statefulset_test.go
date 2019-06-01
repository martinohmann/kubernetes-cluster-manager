package resource

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSlice_PersistentVolumeClaimsForDeletion(t *testing.T) {
	s := Slice{
		{Kind: KindJob},
		{Kind: KindStatefulSet, Namespace: "foo", DeletePersistentVolumeClaims: true, Content: []byte(`---
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: web
  namespace: foo
  annotations:
    kcm/deletion-policy: delete-pvcs
spec:
  replicas: 2
  volumeClaimTemplates:
  - metadata:
      name: www
  - metadata:
      name: cache
`)},
		{Kind: KindStatefulSet, Namespace: "bar", DeletePersistentVolumeClaims: true, Content: []byte(`---
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: cache
  namespace: bar
  annotations:
    kcm/deletion-policy: delete-pvcs
spec:
  volumeClaimTemplates:
  - metadata:
      name: data
`)},
		{Kind: KindStatefulSet, Content: []byte(`---
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: db
spec:
  replicas: 3
  volumeClaimTemplates:
  - metadata:
      name: data
`)},
	}

	expected := []Head{
		{Kind: KindPersistentVolumeClaim, Metadata: Metadata{Name: "www-web-0", Namespace: "foo"}},
		{Kind: KindPersistentVolumeClaim, Metadata: Metadata{Name: "www-web-1", Namespace: "foo"}},
		{Kind: KindPersistentVolumeClaim, Metadata: Metadata{Name: "cache-web-0", Namespace: "foo"}},
		{Kind: KindPersistentVolumeClaim, Metadata: Metadata{Name: "cache-web-1", Namespace: "foo"}},
		{Kind: KindPersistentVolumeClaim, Metadata: Metadata{Name: "data-cache-0", Namespace: "bar"}},
	}

	pvcs := s.PersistentVolumeClaimsForDeletion()

	assert.Equal(t, expected, pvcs)
}