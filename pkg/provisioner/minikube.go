package provisioner

import (
	"context"
	"fmt"
	"os/exec"

	"github.com/martinohmann/kubernetes-cluster-manager/pkg/command"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
)

// Minikube uses minikube instead of an actual infrastructure provisioner.
// This is useful for local testing.
type Minikube struct{}

// NewMinikube creates a new Minikube provisioner.
func NewMinikube(o *Options) Provisioner {
	return &Minikube{}
}

func (m *Minikube) status() error {
	cmd := exec.Command("minikube", "status")

	_, err := command.RunSilently(cmd)
	if err != nil {
		err = errors.New("minikube not running")
	}

	return err
}

// Provision implements Provision from the Provisioner interface.
func (m *Minikube) Provision(ctx context.Context) error {
	if err := m.status(); err == nil {
		return nil
	}

	cmd := exec.Command("minikube", "start", "--keep-context")

	_, err := command.Run(cmd)

	return err
}

// Output implements Outputter.
func (m *Minikube) Output(ctx context.Context) (map[string]interface{}, error) {
	home, _ := homedir.Dir()

	v := map[string]interface{}{
		"kubeconfig": fmt.Sprintf("%s/.kube/config", home),
		"context":    "minikube",
	}

	return v, nil
}

// Destroy implements Destroy from the Provisioner interface.
func (m *Minikube) Destroy(ctx context.Context) error {
	if err := m.status(); err != nil {
		return err
	}

	cmd := exec.Command("minikube", "delete")

	_, err := command.Run(cmd)

	return err
}
