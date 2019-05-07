package provisioner

import "context"

// Null does not provision any infrastructure. All interface funcs
// will return successfully. Null can be used if you want to manage
// your cluster infrastructure by other means and just want to make use of the
// manifest rendering and applying part.
type Null struct{}

// NewNull creates a new Null provisioner.
func NewNull(o *Options) Provisioner {
	return &Null{}
}

// Provision implements Provision from the Provisioner interface.
func (*Null) Provision(ctx context.Context) error {
	return nil
}

// Destroy implements Destroy from the Provisioner interface.
func (*Null) Destroy(ctx context.Context) error {
	return nil
}
