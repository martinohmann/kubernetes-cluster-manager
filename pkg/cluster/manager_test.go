// build +integration

package cluster

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/martinohmann/kubernetes-cluster-manager/internal/commandtest"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/command"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/credentials"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/file"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/provisioner"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/renderer"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func createManager() *Manager {
	p := provisioner.NewTerraform(&provisioner.TerraformOptions{})
	m := NewManager(
		credentials.NewStaticSource(&credentials.Credentials{Context: "test"}),
		p,
		renderer.NewHelm(&renderer.HelmOptions{ChartsDir: "testdata/charts"}),
		log.StandardLogger(),
	)

	return m
}

func TestProvision(t *testing.T) {
	commandtest.WithMockExecutor(func(executor *commandtest.MockExecutor) {
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
		manifestsDir, _ := ioutil.TempDir("", "manifests")
		defer os.RemoveAll(manifestsDir)

		o := &Options{
			Values:       values.Name(),
			Deletions:    deletions.Name(),
			ManifestsDir: manifestsDir,
		}

		p := createManager()

		executor.Command("terraform apply --auto-approve").WillSucceed()
		executor.Command("terraform output --json").WillReturn(`{"foo":{"value": "output-from-terraform"}}`)
		executor.Pattern("helm template --values .*").WillExecute()
		executor.Pattern("kubectl cluster-info.*").WillSucceed()
		executor.Pattern("kubectl delete pod --ignore-not-found --namespace kube-system --context test foo").WillSucceed()
		executor.Pattern("kubectl apply -f -").WillSucceed()
		executor.Pattern("kubectl delete deployment --ignore-not-found --namespace default --context test bar").WillSucceed()

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
`

		assert.NoError(t, p.Provision(o))

		buf, _ := ioutil.ReadFile(filepath.Join(manifestsDir, "testchart.yaml"))

		assert.Equal(t, expectedManifest, string(buf))

		buf, _ = ioutil.ReadFile(deletions.Name())

		assert.Equal(t, expectedDeletions, string(buf))

		buf, _ = ioutil.ReadFile(values.Name())

		assert.Equal(t, expectedValues, string(buf))
	}, command.NewExecutor(nil))
}

func TestDestroy(t *testing.T) {
	commandtest.WithMockExecutor(func(executor *commandtest.MockExecutor) {
		deletions, _ := file.NewTempFile("deletions.yaml", []byte(`
preDestroy:
- kind: PersistentVolumeClaim
  name: bar`))
		defer os.Remove(deletions.Name())
		values, _ := file.NewTempFile("values.yaml", []byte(`baz: somevalue`))
		defer os.Remove(values.Name())
		manifest, _ := file.NewTempFile("manifest.yaml", []byte(``))
		defer os.Remove(manifest.Name())
		manifestsDir, _ := ioutil.TempDir("", "manifests")
		defer os.RemoveAll(manifestsDir)

		o := &Options{
			Values:       values.Name(),
			Deletions:    deletions.Name(),
			ManifestsDir: manifestsDir,
		}

		p := createManager()

		executor.Command("terraform output --json").WillReturn(`{}`)
		executor.Pattern("helm template --values .*").WillExecute()
		executor.Pattern("kubectl cluster-info.*").WillSucceed()
		executor.Pattern("kubectl delete -f - --ignore-not-found --context test").WillSucceed()
		executor.Pattern("kubectl delete persistentvolumeclaim --ignore-not-found --namespace default --context test bar").WillSucceed()
		executor.Command("terraform destroy --auto-approve").WillSucceed()

		expectedDeletions := `preApply: []
postApply: []
preDestroy: []
`

		assert.NoError(t, p.Destroy(o))

		buf, _ := ioutil.ReadFile(deletions.Name())

		assert.Equal(t, expectedDeletions, string(buf))
	}, command.NewExecutor(nil))
}

func TestReadEmptyCredentials(t *testing.T) {
	m := &Manager{
		credentialSource: credentials.NewStaticSource(&credentials.Credentials{}),
	}

	_, err := m.readCredentials(&Options{})

	assert.Error(t, err)
}

func TestReadCredentials(t *testing.T) {
	expected := &credentials.Credentials{
		Server: "https://localhost:6443",
		Token:  "mytoken",
	}

	m := &Manager{
		credentialSource: credentials.NewStaticSource(expected),
		logger:           log.StandardLogger(),
	}

	creds, err := m.readCredentials(&Options{})

	assert.NoError(t, err)
	assert.Equal(t, expected, creds)
}
