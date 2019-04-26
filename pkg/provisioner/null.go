package provisioner

import (
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/command"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/kcm"
)

func init() {
	Register("null", func(_ *kcm.ProvisionerOptions, _ command.Executor) (kcm.Provisioner, error) {
		return &Null{}, nil
	})
}

// Null does not provision any infrastructure. All interface funcs
// will return successfully. Null can be used if you want to manage
// your cluster infrastructure by other means and just want to make use of the
// manifest rendering and applying part.
type Null struct{}

// Provision implements Provision from the kcm.Provisioner interface.
func (*Null) Provision() error {
	return nil
}

// Reconcile implements Reconcile from the kcm.Provisioner interface.
func (*Null) Reconcile() error {
	return nil
}

// Fetch implements Fetch from the kcm.Provisioner interface.
func (*Null) Fetch() (kcm.Values, error) {
	return kcm.Values{}, nil
}

// Destroy implements Destroy from the Provisioner interface.
func (*Null) Destroy() error {
	return nil
}
