package infra

import (
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/command"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/kcm"
)

func init() {
	RegisterManager("null", func(_ *ManagerOptions, _ command.Executor) (Manager, error) {
		return &NullManager{}, nil
	})
}

// NullManager does not manage any infrastructure. All interface funcs will
// return successfully. NullManager can be used if you want to manage your
// cluster infrastructure by other means and just want to make use of the
// manifest rendering and applying part.
type NullManager struct{}

// Apply implements Apply from the Manager interface.
func (*NullManager) Apply() error {
	return nil
}

// Plan implements Plan from the Manager interface.
func (*NullManager) Plan() error {
	return nil
}

// GetValues implements GetValues from the Manager interface.
func (*NullManager) GetValues() (kcm.Values, error) {
	return kcm.Values{}, nil
}

// Destroy implements Destroy from the Manager interface.
func (*NullManager) Destroy() error {
	return nil
}
