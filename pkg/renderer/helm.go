package renderer

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"sort"

	"github.com/martinohmann/kubernetes-cluster-manager/pkg/kcm"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/kubernetes/helm"
	"github.com/pkg/errors"
)

func init() {
	Register("helm", func(o *kcm.RendererOptions) (kcm.Renderer, error) {
		return NewHelm(&o.Helm), nil
	})
}

// Helm uses helm to render kubernetes manifests.
type Helm struct {
	chartsDir string
}

// NewHelm creates a new helm manifest renderer with given options.
func NewHelm(o *kcm.HelmOptions) *Helm {
	return &Helm{
		chartsDir: o.ChartsDir,
	}
}

// RenderManifests implements RenderManifests from the kcm.Renderer interface.
func (r *Helm) RenderManifests(v kcm.Values) ([]*kcm.Manifest, error) {
	files, err := ioutil.ReadDir(r.chartsDir)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	dirs := make([]string, 0, len(files))

	for _, f := range files {
		if !f.IsDir() && !helm.IsChartDir(f.Name()) {
			continue
		}

		dirs = append(dirs, f.Name())
	}

	sort.Strings(dirs)

	manifests := make([]*kcm.Manifest, len(dirs))

	for i, d := range dirs {
		manifest, err := r.renderChart(d, v)
		if err != nil {
			return nil, err
		}

		manifests[i] = manifest
	}

	return manifests, nil
}

// renderChart renders a helm chart.
func (r *Helm) renderChart(chartName string, v kcm.Values) (*kcm.Manifest, error) {
	chart := helm.NewChart(filepath.Join(r.chartsDir, chartName))
	buf, err := chart.Render(v)
	if err != nil {
		return nil, err
	}

	m := &kcm.Manifest{
		Filename: fmt.Sprintf("%s.yaml", chartName),
		Content:  buf,
	}

	return m, nil
}
