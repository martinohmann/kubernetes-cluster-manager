package manifest

import (
	"github.com/martinohmann/cluster-manager/pkg/api"
)

type Renderer interface {
	RenderManifest(*api.InfraOutput) (*api.Manifest, error)
}
