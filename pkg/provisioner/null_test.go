package provisioner

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNull(t *testing.T) {
	p := NewNull(&Options{})

	require.NoError(t, p.Provision())
	require.NoError(t, p.Destroy())
}
