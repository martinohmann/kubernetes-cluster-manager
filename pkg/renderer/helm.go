package renderer

import (
	"bytes"
	"path/filepath"
	"sort"
	"strings"

	"github.com/martinohmann/kubernetes-cluster-manager/pkg/kcm"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/kubernetes"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/manifest"
	yaml "gopkg.in/yaml.v2"
	"k8s.io/helm/pkg/chartutil"
	"k8s.io/helm/pkg/proto/hapi/chart"
	"k8s.io/helm/pkg/renderutil"
	"k8s.io/helm/pkg/timeconv"
)

// Helm uses helm to render kubernetes manifests.
type Helm struct {
	TemplatesDir string
}

// NewHelm create a new helm template renderer.
func NewHelm(o *Options) Renderer {
	return &Helm{
		TemplatesDir: o.TemplatesDir,
	}
}

// RenderManifests implements Renderer.
func (r *Helm) RenderManifests(v kcm.Values) ([]*Manifest, error) {
	return renderManifests(r.TemplatesDir, v, renderChart)
}

// renderChart renders a helm chart and satisfies the signature for
// renderManifestFunc.
func renderChart(dir string, v kcm.Values) (*Manifest, error) {
	if ok, err := chartutil.IsChartDir(dir); err != nil || !ok {
		return nil, skipError{err}
	}

	c, err := chartutil.Load(dir)
	if err != nil {
		return nil, err
	}

	rawVals, err := yaml.Marshal(v)
	if err != nil {
		return nil, err
	}

	config := &chart.Config{
		Raw:    string(rawVals),
		Values: map[string]*chart.Value{},
	}

	renderOpts := renderutil.Options{
		ReleaseOptions: chartutil.ReleaseOptions{
			Name:      "kcm",
			Time:      timeconv.Now(),
			Namespace: kubernetes.DefaultNamespace,
		},
	}

	renderedTemplates, err := renderutil.Render(c, config, renderOpts)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer

	for _, manifest := range sortedManifests(renderedTemplates) {
		b := filepath.Base(manifest.Name)
		if strings.HasPrefix(b, "_") || b == "NOTES.txt" {
			continue
		}

		writeSourceHeader(&buf, manifest.Name)

		buf.Write(manifest.Content)
		buf.WriteString("\n")
	}

	m := &Manifest{
		Name:    filepath.Base(dir),
		Content: buf.Bytes(),
	}

	return m, nil
}

// sortedManifests transforms a map of rendered templates into a sorted slice
// of manifests.
func sortedManifests(m map[string]string) []*Manifest {
	manifests := make([]*Manifest, 0, len(m))

	for source, data := range m {
		manifests = append(manifests, &Manifest{
			Name:    source,
			Content: []byte(data),
		})
	}

	sort.Sort(manifest.ByName(manifests))

	return manifests
}
