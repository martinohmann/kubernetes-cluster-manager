package cmd

import (
	"bytes"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestVersionOptionsValidate(t *testing.T) {
	o := &VersionOptions{Output: "xml"}

	err := o.Validate()

	assert.EqualError(t, errors.New("--output must be 'yaml' or 'json'"), err.Error())
}

func TestNewVersionCommand(t *testing.T) {
	buf := bytes.NewBuffer(nil)

	cmd := NewVersionCommand(buf)
	cmd.SetArgs([]string{"--short"})

	assert.NoError(t, cmd.Execute())

	assert.Equal(t, "v0.0.0-master\n", buf.String())
}
