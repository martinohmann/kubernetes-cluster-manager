package manifest

import (
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/api"
)

// Renderer is the interface for a kubernetes manifest renderer.
type Renderer interface {
	// RenderManifest renders a kubernetes manifest.
	RenderManifest(*api.InfraOutput) (*api.Manifest, error)
}
