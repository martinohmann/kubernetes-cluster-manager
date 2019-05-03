package renderer

import "github.com/martinohmann/kubernetes-cluster-manager/pkg/kcm"

// Renderer is the interface for a Kubernetes manifest renderer.
type Renderer interface {
	// RenderManifest renders Kubernetes manifests.
	RenderManifests(kcm.Values) ([]*ManifestInfo, error)
}

// Options are made available to manifest renderers.
type Options struct {
	Helm HelmOptions `json:"helm" yaml:"helm"`
}

// HelmOptions configure the helm manifest renderer.
type HelmOptions struct {
	ChartsDir string `json:"chartsDir" yaml:"chartsDir"`
}

// ManifestInfo contains a kubernetes manifest as raw bytes and the filename.
type ManifestInfo struct {
	Filename string
	Content  []byte
}
