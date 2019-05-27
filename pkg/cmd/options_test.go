package cmd

import (
	"os"
	"testing"

	"github.com/martinohmann/kubernetes-cluster-manager/pkg/credentials"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/file"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestOptionsComplete(t *testing.T) {
	o := &Options{}

	config := `---
workingDir: ~/foo
provisioner: minikube
managerOptions:
  values: /values.yaml
  deletions: /deletions.yaml
credentials:
  kubeconfig: /tmp/kubeconfig
`

	f, err := file.NewTempFile("config.yaml", []byte(config))
	defer os.Remove(f.Name())

	assert.NoError(t, err)

	cmd := &cobra.Command{}
	o.AddFlags(cmd)

	flags := []string{
		"--config", f.Name(),
		"--values", "/tmp/values.yaml",
		"--dry-run",
		"--provisioner", "foo",
	}

	assert.NoError(t, cmd.ParseFlags(flags))
	assert.NoError(t, o.Complete(cmd))

	home, _ := homedir.Dir()

	assert.Equal(t, home+"/foo", o.WorkingDir)
	assert.Equal(t, "minikube", o.Provisioner)
	assert.Equal(t, "/tmp/kubeconfig", o.Credentials.Kubeconfig)
	assert.Equal(t, true, o.ManagerOptions.DryRun)
	assert.Equal(t, "/values.yaml", o.ManagerOptions.Values)
}

func TestOptionsCreateManager(t *testing.T) {
	cases := []struct {
		name        string
		o           *Options
		expectError bool
	}{
		{
			name:        "invalid provisioner",
			o:           &Options{Provisioner: "foo"},
			expectError: true,
		},
		{
			name:        "missing cluster options",
			o:           &Options{Provisioner: "null"},
			expectError: true,
		},
		{
			name: "valid cluster options",
			o: &Options{
				Provisioner: "null",
				Credentials: credentials.Credentials{
					Kubeconfig: "/tmp/kubeconfig",
				},
			},
			expectError: false,
		},
		{
			name: "value fetcher credential source",
			o: &Options{
				Provisioner: "terraform",
			},
			expectError: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			m, err := tc.o.createManager()
			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, m)
			}
		})
	}
}
