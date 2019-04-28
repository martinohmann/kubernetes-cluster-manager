package provisioner

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/martinohmann/kubernetes-cluster-manager/pkg/command"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/kcm"
	homedir "github.com/mitchellh/go-homedir"
)

func init() {
	Register("minikube", func(_ *kcm.ProvisionerOptions, e command.Executor) (kcm.Provisioner, error) {
		return NewMinikube(e), nil
	})
}

// Minikube uses minikube instead of an actual infrastructure provisioner.
// This is useful for local testing.
type Minikube struct {
	executor command.Executor
}

// NewMinikube creates a new minikube manager.
func NewMinikube(executor command.Executor) *Minikube {
	return &Minikube{
		executor: executor,
	}
}

func (m *Minikube) status() error {
	args := []string{
		"minikube",
		"status",
	}

	cmd := exec.Command(args[0], args[1:]...)

	_, err := m.executor.Run(cmd)

	return err
}

func (m *Minikube) start() error {
	if err := m.status(); err == nil {
		return nil
	}

	args := []string{
		"minikube",
		"start",
	}

	cmd := exec.Command(args[0], args[1:]...)

	_, err := m.executor.Run(cmd)

	return err
}

// Provision implements Provision from the kcm.Provisioner interface.
func (m *Minikube) Provision() error {
	return m.start()
}

// Reconcile implements Reconcile from the kcm.Provisioner interface.
func (m *Minikube) Reconcile() error {
	return m.start()
}

// Fetch implements Fetch from the kcm.Provisioner interface.
func (m *Minikube) Fetch() (kcm.Values, error) {
	args := []string{
		"minikube",
		"ip",
	}

	cmd := exec.Command(args[0], args[1:]...)

	out, err := m.executor.RunSilently(cmd)
	if err != nil {
		return nil, err
	}

	home, _ := homedir.Dir()

	v := kcm.Values{
		"server":     fmt.Sprintf("https://%s:8443", strings.Trim(out, "\n")),
		"kubeconfig": fmt.Sprintf("%s/.kube/config", home),
		// this will force the correct kubectl kubeconfig context
		"context": "minikube",
	}

	return v, nil
}

// Destroy implements Destroy from the kcm.Provisioner interface.
func (m *Minikube) Destroy() error {
	if err := m.status(); err != nil {
		return err
	}

	args := []string{
		"minikube",
		"delete",
	}

	cmd := exec.Command(args[0], args[1:]...)

	_, err := m.executor.Run(cmd)

	return err
}
