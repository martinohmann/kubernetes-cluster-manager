package infra

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/martinohmann/kubernetes-cluster-manager/pkg/command"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/kcm"
	homedir "github.com/mitchellh/go-homedir"
)

func init() {
	RegisterManager("minikube", func(_ *ManagerOptions, e command.Executor) (Manager, error) {
		return NewMinikubeManager(e), nil
	})
}

// MinikubeManager uses minikube instead of an actual infrastructure manager.
// This is useful for local testing.
type MinikubeManager struct {
	executor command.Executor
}

// NewMinikubeManager creates a new minikube manager.
func NewMinikubeManager(executor command.Executor) *MinikubeManager {
	return &MinikubeManager{
		executor: executor,
	}
}

func (m *MinikubeManager) status() error {
	args := []string{
		"minikube",
		"status",
	}

	cmd := exec.Command(args[0], args[1:]...)

	_, err := m.executor.Run(cmd)

	return err
}

func (m *MinikubeManager) start() error {
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

// Apply implements Apply from the Manager interface.
func (m *MinikubeManager) Apply() error {
	return m.start()
}

// Plan implements Plan from the Manager interface.
func (m *MinikubeManager) Plan() error {
	return m.start()
}

// GetValues implements GetValues from the Manager interface.
func (m *MinikubeManager) GetValues() (kcm.Values, error) {
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

// Destroy implements Destroy from the Manager interface.
func (m *MinikubeManager) Destroy() error {
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
