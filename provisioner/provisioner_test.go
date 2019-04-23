// +build integration

package provisioner

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/martinohmann/kubernetes-cluster-manager/infra"
	"github.com/martinohmann/kubernetes-cluster-manager/manifest"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/command"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/file"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/kubernetes"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func createProvisioner() (*Provisioner, *command.MockExecutor) {
	e := command.NewMockExecutor(command.NewExecutor())
	p := NewClusterProvisioner(
		&kubernetes.ClusterOptions{},
		infra.NewTerraformManager(&infra.TerraformOptions{}, e),
		manifest.NewHelmRenderer(&manifest.HelmOptions{Chart: "testdata/testchart"}, e),
		e,
		log.StandardLogger(),
	)

	return p, e
}

func TestProvision(t *testing.T) {
	deletions, _ := file.NewTempFile("deletions.yaml", []byte(`
preApply:
- kind: Pod
  name: foo
  namespace: kube-system
postApply:
- kind: Deployment
  name: bar`))
	defer os.Remove(deletions.Name())
	values, _ := file.NewTempFile("values.yaml", []byte(`baz: somevalue`))
	defer os.Remove(values.Name())
	manifest, _ := file.NewTempFile("manifest.yaml", []byte(``))
	defer os.Remove(manifest.Name())

	o := &Options{
		Values:    values.Name(),
		Deletions: deletions.Name(),
		Manifest:  manifest.Name(),
	}

	p, executor := createProvisioner()

	executor.Command("terraform apply --auto-approve").WillSucceed()
	executor.Command("terraform output --json").WillReturn(`{"foo":{"value": "output-from-terraform"},"kubeconfig":{"value":"/tmp/kubeconfig"}}`)
	executor.Pattern("helm template --values .*").WillExecute()
	executor.Pattern("kubectl cluster-info.*").WillSucceed()
	executor.Pattern("kubectl delete pod --ignore-not-found --namespace kube-system --kubeconfig /tmp/kubeconfig foo").WillSucceed()
	executor.Pattern("kubectl apply -f -").WillSucceed()
	executor.Pattern("kubectl delete deployment --ignore-not-found --namespace default --kubeconfig /tmp/kubeconfig bar").WillSucceed()

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
  baz: somevalue

`
	expectedDeletions := `preApply: []
postApply: []
preDestroy: []
`
	expectedValues := `baz: somevalue
foo: output-from-terraform
kubeconfig: /tmp/kubeconfig
`

	err := p.Provision(o)

	assert.NoError(t, err)

	buf, _ := ioutil.ReadFile(manifest.Name())

	assert.Equal(t, expectedManifest, string(buf))

	buf, _ = ioutil.ReadFile(deletions.Name())

	assert.Equal(t, expectedDeletions, string(buf))

	buf, _ = ioutil.ReadFile(values.Name())

	assert.Equal(t, expectedValues, string(buf))
}
