package renderer

import (
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/kcm"
)

// Null is a renderer which renders nothing.
type Null struct{}

// RenderManifests implements Renderer.
func (*Null) RenderManifests(v kcm.Values) ([]*ManifestInfo, error) {
	return []*ManifestInfo{}, nil
}
