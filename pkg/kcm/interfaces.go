package kcm

// Provisioner is the interface for a cloud infrastructure provisioner.
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

// Renderer is the interface for a kubernetes manifest renderer.
type Renderer interface {
	// RenderManifest renders a kubernetes manifest.
	RenderManifest(Values) (Manifest, error)
}

// CredentialSource provides credentials for a Kubernetes cluster.
type CredentialSource interface {
	// GetCredentials returns credentials for a Kubernetes cluster. Will return
	// an error if retrieving credentials fails.
	GetCredentials() (*Credentials, error)
}
