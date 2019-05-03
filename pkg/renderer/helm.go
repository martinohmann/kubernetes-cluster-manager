package renderer

import (
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/kcm"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/kubernetes/helm"
)

// Helm uses helm to render kubernetes manifests.
type Helm struct {
	TemplatesDir string
}

func NewHelm(o *Options) Renderer {
	return &Helm{
		TemplatesDir: o.TemplatesDir,
	}
}

// RenderManifests implements Renderer.
func (r *Helm) RenderManifests(v kcm.Values) ([]*ManifestInfo, error) {
	return renderManifests(r.TemplatesDir, v, renderChart)
}

// renderChart renders a helm chart and satisfies the signature for
// renderManifestFunc.
func renderChart(dir string, v kcm.Values) (*ManifestInfo, error) {
	if !helm.IsChartDir(dir) {
		return nil, skipError{dir}
	}

	buf, err := helm.NewChart(dir).Render(v)
	if err != nil {
		return nil, err
	}

	m := &ManifestInfo{
		Filename: manifestFilename(dir),
		Content:  buf,
	}

	return m, nil
}
