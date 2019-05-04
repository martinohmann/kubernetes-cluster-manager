package renderer

import (
	"bytes"
	"path/filepath"
	"strings"

	"github.com/martinohmann/kubernetes-cluster-manager/pkg/kcm"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/kubernetes"
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
		return nil, skipError{dir}
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

	var buf bytes.Buffer

	for source, data := range renderedTemplates {
		b := filepath.Base(source)
		if strings.HasPrefix(b, "_") {
			continue
		}

		writeSourceHeader(&buf, source)

		buf.WriteString(data)
		buf.WriteString("\n")
	}

	m := &Manifest{
		Name:    filepath.Base(dir),
		Content: buf.Bytes(),
	}

	return m, nil
}
