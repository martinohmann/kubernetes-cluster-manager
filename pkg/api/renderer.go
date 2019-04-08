package api

// ManifestRenderer is the interface for a kubernetes manifest renderer.
type ManifestRenderer interface {
	// RenderManifest renders a kubernetes manifest.
	RenderManifest(Values) (Manifest, error)
}
