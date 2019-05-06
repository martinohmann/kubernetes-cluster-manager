package renderer

import (
	"io"
	"io/ioutil"
	"path/filepath"

	"github.com/martinohmann/kubernetes-cluster-manager/pkg/kcm"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/manifest"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// Renderer is the interface for a Kubernetes manifest renderer.
type Renderer interface {
	// RenderManifest renders Kubernetes manifests.
	RenderManifests(kcm.Values) ([]*Manifest, error)
}

// Options are made available to manifest renderers.
type Options struct {
	TemplatesDir string `json:"templatesDir,omitempty" yaml:"templatesDir,omitempty"`
}

// Manifest is a type alias for manifest.Manifest.
type Manifest = manifest.Manifest

// skipError can be returned while iterating directories to indicate that the
// directory should be skipped.
type skipError struct {
	err error
}

// Error implements error.
func (e skipError) Error() string {
	return e.err.Error()
}

// renderManifestFunc is a function that renders a manifest with the values v
// for the contents of dir.
type renderManifestFunc func(dir string, v kcm.Values) (*Manifest, error)

// renderManifests interates dir and renders manifests using render.
func renderManifests(dir string, v kcm.Values, render renderManifestFunc) ([]*Manifest, error) {
	dirs, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open templates dir")
	}

	manifests := make([]*Manifest, 0, len(dirs))

	for _, d := range dirs {
		if !d.IsDir() {
			continue
		}

		fullPath := filepath.Join(dir, d.Name())

		manifest, err := render(fullPath, v)
		if e, ok := err.(skipError); ok {
			log.Warnf("rendering of component %s was skipped due to: %s", d.Name(), e.Error())
			continue
		}

		if err != nil {
			return nil, err
		}

		manifests = append(manifests, manifest)
	}

	return manifests, nil
}

// writeSourceHeader writes the manifest source file header to w.
func writeSourceHeader(w io.StringWriter, source string) {
	w.WriteString("---\n# Source: ")
	w.WriteString(source)
	w.WriteString("\n")
}
