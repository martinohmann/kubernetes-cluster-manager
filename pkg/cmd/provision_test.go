package cmd

import (
	"os"
	"testing"

	"github.com/martinohmann/kubernetes-cluster-manager/pkg/file"
	homedir "github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestOptionsComplete(t *testing.T) {
	o := &Options{logger: log.StandardLogger()}

	config := `---
workingDir: ~/foo
provisioner: minikube
managerOptions:
  values: /values.yaml
  deletions: /deletions.yaml
clusterOptions:
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
	assert.Equal(t, "/tmp/kubeconfig", o.ClusterOptions.Kubeconfig)
	assert.Equal(t, true, o.ManagerOptions.DryRun)
	assert.Equal(t, "/values.yaml", o.ManagerOptions.Values)
	assert.Equal(t, "/deletions.yaml", o.ManagerOptions.Deletions)
}
