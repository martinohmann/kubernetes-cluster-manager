// +build integration

package provisioner

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/martinohmann/kubernetes-cluster-manager/pkg/command"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/config"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/fs"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/kubernetes/helm"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/terraform"
	"github.com/stretchr/testify/assert"
)

func createProvisioner(cfg *config.Config) (*Provisioner, *command.MockExecutor) {
	e := command.NewMockExecutor(command.NewExecutor())
	p := NewClusterProvisioner(
		terraform.NewInfraManager(&cfg.Terraform, e),
		helm.NewManifestRenderer(&cfg.Helm, e),
		e,
	)

	return p, e
}

func TestProvision(t *testing.T) {
	deletions, _ := fs.NewTempFile("deletions.yaml", []byte(`
preApply:
- kind: Pod
  name: foo
  namespace: kube-system
postApply:
- kind: Deployment
  name: bar`))
	defer os.Remove(deletions.Name())
	values, _ := fs.NewTempFile("values.yaml", []byte(``))
	defer os.Remove(values.Name())
	manifest, _ := fs.NewTempFile("manifest.yaml", []byte(``))
	defer os.Remove(manifest.Name())

	cfg := &config.Config{
		Helm: config.HelmConfig{
			Chart: "testdata/testchart",
		},
		Values:    values.Name(),
		Deletions: deletions.Name(),
		Manifest:  manifest.Name(),
	}

	p, executor := createProvisioner(cfg)

	executor.Command("terraform apply --auto-approve").WillSucceed()
	executor.Command("terraform output --json").WillReturn(`{"foo":{"value": "output-from-terraform"}}`)
	executor.Pattern("helm template --values .*").WillExecute()
	executor.Command("kubectl cluster-info").WillSucceed()
	executor.Command("kubectl delete pod --ignore-not-found --namespace kube-system foo").WillSucceed()
	executor.Pattern("kubectl apply -f -").WillSucceed()
	executor.Command("kubectl delete deployment --ignore-not-found --namespace default bar").WillSucceed()

	expectedManifest := `---
# Source: testchart/templates/configmap.yaml
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: test
  namespace: kube-system
data:
  foo: output-from-terraform
  bar: baz

`
	expectedDeletions := `preApply: []
postApply: []
preDestroy: []
`
	expectedValues := `foo: output-from-terraform
`

	err := p.Provision(cfg)

	assert.NoError(t, err)

	buf, _ := ioutil.ReadFile(manifest.Name())

	assert.Equal(t, expectedManifest, string(buf))

	buf, _ = ioutil.ReadFile(deletions.Name())

	assert.Equal(t, expectedDeletions, string(buf))

	buf, _ = ioutil.ReadFile(values.Name())

	assert.Equal(t, expectedValues, string(buf))
}
