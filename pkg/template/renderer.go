package template

import (
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/kubernetes"
	"gopkg.in/yaml.v2"
	"k8s.io/helm/pkg/chartutil"
	"k8s.io/helm/pkg/proto/hapi/chart"
	"k8s.io/helm/pkg/renderutil"
	"k8s.io/helm/pkg/timeconv"
)

// Renderer defines a template renderer
type Renderer interface {
	// Render renders all templates in dir with values v and returns a map of
	// template-file-path => rendered-template-content. Should return an error
	// if rendering of any of the templates fails.
	Render(dir string, v map[string]interface{}) (map[string]string, error)
}

type renderer struct {
	Name      string
	Namespace string
}

// NewRenderer creates a new template renderer.
func NewRenderer() Renderer {
	return &renderer{
		Name:      "kcm",
		Namespace: kubernetes.DefaultNamespace,
	}
}

// Render implements Renderer.
func (r *renderer) Render(dir string, v map[string]interface{}) (map[string]string, error) {
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

	opts := renderutil.Options{
		ReleaseOptions: chartutil.ReleaseOptions{
			Name:      r.Name,
			Namespace: r.Namespace,
			Time:      timeconv.Now(),
		},
	}

	return renderutil.Render(c, config, opts)
}
