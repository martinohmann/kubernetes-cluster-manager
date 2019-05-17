package manifest

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResourceBuffer_Write(t *testing.T) {
	expected := []byte(`---
foo
---
barbaz
`)

	var buf resourceBuffer

	n, err := buf.Write([]byte("foo"))

	require.NoError(t, err)
	require.Equal(t, 8, n)

	n, err = buf.Write([]byte("barbaz"))

	require.NoError(t, err)
	require.Equal(t, 11, n)

	assert.Equal(t, expected, buf.Bytes())
}
