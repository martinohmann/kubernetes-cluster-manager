package renderer

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/martinohmann/kubernetes-cluster-manager/pkg/kcm"
	"github.com/pkg/errors"
)

// Renderer is the interface for a Kubernetes manifest renderer.
type Renderer interface {
	// RenderManifest renders Kubernetes manifests.
	RenderManifests(kcm.Values) ([]*ManifestInfo, error)
}

// Options are made available to manifest renderers.
type Options struct {
	TemplatesDir string `json:"templatesDir" yaml:"templatesDir"`
}

// ManifestInfo contains a kubernetes manifest as raw bytes and the filename.
type ManifestInfo struct {
	Filename string
	Content  []byte
}

// skipError can be returned while iterating directories to indicate that the
// directory should be skipped.
type skipError struct {
	dir string
}

// Error implements error.
func (e skipError) Error() string {
	return fmt.Sprintf("%s skipped", e.dir)
}

// renderManifestFunc is a function that renders a manifest with the values v
// for the contents of dir.
type renderManifestFunc func(dir string, v kcm.Values) (*ManifestInfo, error)

// renderManifests interates dir and renders manifests using render.
func renderManifests(dir string, v kcm.Values, render renderManifestFunc) ([]*ManifestInfo, error) {
	dirs, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	manifests := make([]*ManifestInfo, 0, len(dirs))

	for _, d := range dirs {
		if !d.IsDir() {
			continue
		}

		fullPath := filepath.Join(dir, d.Name())

		manifest, err := render(fullPath, v)
		if _, ok := err.(skipError); ok {
			continue
		}

		if err != nil {
			return nil, err
		}

		manifests = append(manifests, manifest)
	}

	return manifests, nil
}

// manifestFilename creates the manifest filename for dir.
func manifestFilename(dir string) string {
	return fmt.Sprintf("%s.yaml", filepath.Base(dir))
}
