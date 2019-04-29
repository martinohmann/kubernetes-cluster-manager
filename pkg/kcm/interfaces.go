package kcm

// ClusterManager is the interface for a Kubernetes cluster manager that will
// orchestrate changes to the cluster infrastructure and the cluster itself.
type ClusterManager interface {
	// Provision performs all steps necessary to create and setup a cluster and
	// the required infrastructure. If a cluster already exists, it should
	// update it if there are pending changes to be rolled out. Depending on
	// the options it may or may not perform a dry run of the pending changes.
	Provision(*Options) error

	// Destroy deletes all applied manifests from a cluster and tears down the
	// cluster infrastructure. Depending on the options it may or may not
	// perform a dry run of the destruction process.
	Destroy(*Options) error

	// ApplyManifests renders and applies all manifests to the cluster. It also
	// takes care of pending resource deletions that should be performed before
	// and after applying.
	ApplyManifests(*Options) error

	// DeleteManifests renders and deletes all manifests from the cluster. It
	// also takes care of other resource deletions that should be performed
	// after the manifests have been deleted from the cluster.
	DeleteManifests(*Options) error
}

// Provisioner is the interface for an infrastructure provisioner.
type Provisioner interface {
	// Provision applies changes to the infrastructure. It should
	// automatically create or update a kubernetes cluster.
	Provision() error

	// Reconcile retrieves the current state of the infrastructure and
	// should log potential changes without actually applying them.
	Reconcile() error

	// Fetch obtains output values from the infrastructure provisioner.
	// These values are made available during kubernetes manifest
	// renderering.
	Fetch() (Values, error)

	// Destroy performs all actions needed to destroy the underlying
	// cluster infrastructure.
	Destroy() error
}

// Renderer is the interface for a Kubernetes manifest renderer.
type Renderer interface {
	// RenderManifest renders Kubernetes manifests.
	RenderManifests(Values) ([]*Manifest, error)
}

// CredentialSource provides credentials for a Kubernetes cluster.
type CredentialSource interface {
	// GetCredentials returns credentials for a Kubernetes cluster. Will return
	// an error if retrieving credentials fails.
	GetCredentials() (*Credentials, error)
}
