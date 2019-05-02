package renderer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreate(t *testing.T) {
	r, err := Create("helm", &Options{})

	assert.NoError(t, err)
	assert.IsType(t, &Helm{}, r)
}

func TestCreateError(t *testing.T) {
	_, err := Create("", &Options{})

	assert.Error(t, err)
}
