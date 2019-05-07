package provisioner

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNull(t *testing.T) {
	p := NewNull(&Options{})

	require.NoError(t, p.Provision(context.Background()))
	require.NoError(t, p.Destroy(context.Background()))
}
