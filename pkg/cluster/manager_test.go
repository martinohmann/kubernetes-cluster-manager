package cluster

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/martinohmann/kubernetes-cluster-manager/internal/commandtest"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/command"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/credentials"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/file"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/provisioner"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/template"
	"github.com/stretchr/testify/assert"
)

func createManager() *Manager {
	m := NewManager(
		credentials.NewStaticSource(&credentials.Credentials{Context: "test"}),
		provisioner.NewTerraform(&provisioner.Options{}),
		template.NewRenderer(),
	)

	return m
}

func TestProvision(t *testing.T) {
	commandtest.WithMockExecutor(func(executor commandtest.MockExecutor) {
		values, _ := file.NewTempFile("values.yaml", []byte(`baz: somevalue`))
		defer os.Remove(values.Name())
		manifest, _ := file.NewTempFile("manifest.yaml", []byte(``))
		defer os.Remove(manifest.Name())
		manifestsDir, _ := ioutil.TempDir("", "manifests")
		defer os.RemoveAll(manifestsDir)

		o := &Options{
			Values:       values.Name(),
			ManifestsDir: manifestsDir,
			TemplatesDir: "testdata/charts",
		}

		p := createManager()

		executor.ExpectCommand("terraform apply --auto-approve")
		executor.ExpectCommand("terraform output --json").WillReturn(`{"foo":{"value": "output-from-terraform"}}`)
		executor.ExpectCommand("kubectl cluster-info.*")
		executor.ExpectCommand("kubectl apply -f -")

		expectedManifest := `---
apiVersion: v1
data:
  bar: baz
  baz: somevalue
  foo: output-from-terraform
kind: ConfigMap
metadata:
  name: test
  namespace: kube-system

---
apiVersion: v1
kind: Service
metadata:
  name: my-service
spec:
  ports:
  - port: 80
    protocol: TCP
    targetPort: 9376
  selector:
    app: MyApp

---
apiVersion: batch/v1
kind: Job
metadata:
  labels:
    kcm/hook: post-create
  name: pi
spec:
  backoffLimit: 4
  template:
    spec:
      containers:
      - command:
        - perl
        - -Mbignum=bpi
        - -wle
        - print bpi(2000)
        image: perl
        name: pi
      restartPolicy: Never

`
		expectedValues := `baz: somevalue
foo: output-from-terraform
`

		assert.NoError(t, p.Provision(context.Background(), o))

		buf, _ := ioutil.ReadFile(filepath.Join(manifestsDir, "testchart.yaml"))

		assert.Equal(t, expectedManifest, string(buf))

		buf, _ = ioutil.ReadFile(values.Name())

		assert.Equal(t, expectedValues, string(buf))
		assert.NoError(t, executor.ExpectationsWereMet())
	}, command.NewExecutor(nil))
}

func TestDestroy(t *testing.T) {
	commandtest.WithMockExecutor(func(executor commandtest.MockExecutor) {
		values, _ := file.NewTempFile("values.yaml", []byte(`baz: somevalue`))
		defer os.Remove(values.Name())
		manifest, _ := file.NewTempFile("manifest.yaml", []byte(``))
		defer os.Remove(manifest.Name())
		manifestsDir, _ := ioutil.TempDir("", "manifests")
		defer os.RemoveAll(manifestsDir)

		o := &Options{
			Values:       values.Name(),
			ManifestsDir: manifestsDir,
			TemplatesDir: "testdata/charts",
		}

		p := createManager()

		executor.ExpectCommand("terraform output --json").WillReturn(`{}`)
		executor.ExpectCommand("kubectl cluster-info.*")
		executor.ExpectCommand("kubectl delete -f - --ignore-not-found --context test")
		executor.ExpectCommand("terraform destroy --auto-approve")

		assert.NoError(t, p.Destroy(context.Background(), o))

		assert.NoError(t, executor.ExpectationsWereMet())
	}, command.NewExecutor(nil))
}

func TestReadEmptyCredentials(t *testing.T) {
	m := &Manager{
		credentialSource: credentials.NewStaticSource(&credentials.Credentials{}),
	}

	_, err := m.readCredentials(context.Background(), &Options{})

	assert.Error(t, err)
}

func TestReadCredentials(t *testing.T) {
	expected := &credentials.Credentials{
		Server: "https://localhost:6443",
		Token:  "mytoken",
	}

	m := &Manager{
		credentialSource: credentials.NewStaticSource(expected),
	}

	creds, err := m.readCredentials(context.Background(), &Options{})

	assert.NoError(t, err)
	assert.Equal(t, expected, creds)
}
