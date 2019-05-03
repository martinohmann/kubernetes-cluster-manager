package provisioner

import (
	"fmt"
	"os/exec"

	"github.com/martinohmann/kubernetes-cluster-manager/pkg/command"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/kcm"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
)

// Minikube uses minikube instead of an actual infrastructure provisioner.
// This is useful for local testing.
type Minikube struct{}

func (m *Minikube) status() error {
	cmd := exec.Command("minikube", "status")

	_, err := command.RunSilently(cmd)
	if err != nil {
		err = errors.New("minikube not running")
	}

	return err
}

// Provision implements Provision from the Provisioner interface.
func (m *Minikube) Provision() error {
	if err := m.status(); err == nil {
		return nil
	}

	cmd := exec.Command("minikube", "start")

	_, err := command.Run(cmd)

	return err
}

// Output implements Outputter.
func (m *Minikube) Output() (kcm.Values, error) {
	home, _ := homedir.Dir()

	v := kcm.Values{
		"kubeconfig": fmt.Sprintf("%s/.kube/config", home),
		"context":    "minikube",
	}

	return v, nil
}

// Destroy implements Destroy from the Provisioner interface.
func (m *Minikube) Destroy() error {
	if err := m.status(); err != nil {
		return err
	}

	cmd := exec.Command("minikube", "delete")

	_, err := command.Run(cmd)

	return err
}
