package renderer

import (
	"testing"

	"github.com/martinohmann/kubernetes-cluster-manager/pkg/kcm"
	"github.com/stretchr/testify/require"
)

func TestNull(t *testing.T) {
	p := NewNull(&Options{})

	manifests, err := p.RenderManifests(kcm.Values{})

	require.NoError(t, err)
	require.Len(t, manifests, 0)
}
