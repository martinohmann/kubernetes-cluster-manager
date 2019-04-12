// +build skip

package provisioner

import (
	"errors"
	"testing"

	"github.com/martinohmann/kubernetes-cluster-manager/pkg/api"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/command"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/config"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/kubernetes/helm"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/terraform"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

type mockManager struct{}

func (mockManager) Apply() error {
	return nil
}

func (mockManager) Plan() error {
	return nil
}

func (mockManager) Destroy() error {
	return nil
}

func (mockManager) GetValues() (api.Values, error) {
	return api.Values{}, nil
}

type mockRenderer struct{}

func (mockRenderer) RenderManifest(v api.Values) (api.Manifest, error) {
	return api.Manifest{}, nil
}

func createProvisioner(cfg *config.Config) (*Provisioner, *command.MockExecutor) {
	e := command.NewMockExecutor()
	p := NewClusterProvisioner(
		terraform.NewInfraManager(&cfg.Terraform, e),
		helm.NewManifestRenderer(&cfg.Helm, e),
		e,
	)

	return p, e
}

func TestProvisioner(t *testing.T) {
	cfg := &config.Config{
		DryRun:    true,
		Values:    "testdata/values.yaml",
		Deletions: "testdata/deletions.yaml",
		Manifest:  "testdata/manifest.yaml",
	}

	log.SetLevel(log.DebugLevel)

	p, executor := createProvisioner(cfg)

	executor.Pattern("terraform .*").WillReturnError(errors.New("foo"))

	err := p.Provision(cfg)

	assert.NoError(t, err)
}
