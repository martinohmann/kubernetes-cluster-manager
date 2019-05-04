package renderer

import (
	"bytes"
	"path/filepath"
	"strings"

	"github.com/martinohmann/kubernetes-cluster-manager/pkg/kcm"
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
func (r *Helm) RenderManifests(v kcm.Values) ([]*ManifestInfo, error) {
	return renderManifests(r.TemplatesDir, v, renderChart)
}

// renderChart renders a helm chart and satisfies the signature for
// renderManifestFunc.
func renderChart(dir string, v kcm.Values) (*ManifestInfo, error) {
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
			IsInstall: true,
			IsUpgrade: false,
			Time:      timeconv.Now(),
			Namespace: "default",
		},
	}

	renderedTemplates, err := renderutil.Render(c, config, renderOpts)

	var buf bytes.Buffer

	for name, data := range renderedTemplates {
		b := filepath.Base(name)
		if strings.HasPrefix(b, "_") {
			continue
		}

		buf.WriteString("---\n# Source: ")
		buf.WriteString(name)
		buf.WriteString("\n")
		buf.WriteString(data)
		buf.WriteString("\n")
	}

	m := &ManifestInfo{
		Filename: manifestFilename(dir),
		Content:  buf.Bytes(),
	}

	return m, nil
}
