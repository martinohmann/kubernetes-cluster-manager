package cmd

import (
	"bytes"
	"os"
	"testing"

	"github.com/martinohmann/kubernetes-cluster-manager/pkg/fs"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestDumpConfigOptionsValidate(t *testing.T) {
	o := &DumpConfigOptions{Output: "xml"}

	err := o.Validate()

	assert.EqualError(t, errors.New("--output must be 'yaml' or 'json'"), err.Error())
}

func TestNewDumpConfigCommand(t *testing.T) {
	buf := bytes.NewBuffer(nil)

	config := `---
workingDir: /tmp/cluster
infraManager:
  terraform:
    parallelism: 1
`

	expected := `workingDir: /tmp/cluster
infraManager:
  terraform:
    parallelism: 1

`

	f, err := fs.NewTempFile("config.yaml", []byte(config))
	defer os.Remove(f.Name())

	assert.NoError(t, err)

	cmd := NewDumpConfigCommand(buf)
	cmd.SetArgs([]string{"--output", "yaml", "--config", f.Name()})

	assert.NoError(t, cmd.Execute())

	assert.Equal(t, expected, buf.String())
}
