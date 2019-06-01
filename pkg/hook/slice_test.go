package hook

import (
	"testing"
	"time"

	"github.com/martinohmann/kubernetes-cluster-manager/pkg/resource"
	"github.com/stretchr/testify/assert"
)

func TestSlice_String(t *testing.T) {
	s := Slice{
		{Resource: &resource.Resource{Name: "foo", Kind: resource.KindJob}, Type: TypePreCreate},
		{Resource: &resource.Resource{Name: "bar", Kind: resource.KindJob}, Type: TypePostDelete, WaitFor: "condition=complete"},
		{Resource: &resource.Resource{Name: "baz", Kind: resource.KindJob}, Type: TypePreUpgrade, WaitFor: "condition=complete", WaitTimeout: 10 * time.Second},
	}

	expected := `pre-create/job/foo
post-delete/job/bar (wait-for=condition=complete)
pre-upgrade/job/baz (wait-for=condition=complete,wait-timeout=10s)`

	assert.Equal(t, expected, s.String())
}
