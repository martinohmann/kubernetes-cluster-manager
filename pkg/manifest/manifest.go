package manifest

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/martinohmann/kubernetes-cluster-manager/pkg/template"
	"github.com/pkg/errors"
)

// Manifest contains a kubernetes manifest split into resources and hooks.
type Manifest struct {
	Name string

	content   []byte
	resources ResourceSlice
	hooks     HookSliceMap
}

// NewManifest creates a new manifest with name from given rendered template
// files.
func NewManifest(name string, renderedTemplates map[string]string) (*Manifest, error) {
	resources := make(ResourceSlice, 0)
	hooks := make(HookSliceMap)

	for _, content := range renderedTemplates {
		r, h, err := parseResources([]byte(content))
		if err != nil {
			return nil, err
		}

		resources = append(resources, r...)

		for k, v := range h {
			if _, ok := hooks[k]; ok {
				hooks[k] = append(hooks[k], v...)
			} else {
				hooks[k] = v
			}
		}
	}

	m := &Manifest{
		Name:      name,
		resources: resources.Sort(ApplyOrder),
		hooks:     hooks.Sort(ApplyOrder),
	}

	return m, nil
}

// Filename returns the filename for the manifest.
func (m *Manifest) Filename() string {
	return fmt.Sprintf("%s.yaml", m.Name)
}

// IsBlank returns true if a manifest does contain nothing but whitespace,
// comments and document separators (`---`). In this case it is semantically
// equivalent to an empty manifest. A nil manifest is considered blank.
func (m *Manifest) IsBlank() bool {
	if m == nil || len(m.Content()) == 0 {
		return true
	}

	buf := bytes.NewBuffer(m.Content())
	s := bufio.NewScanner(buf)

	for s.Scan() {
		line := bytes.TrimSpace(s.Bytes())

		if len(line) == 0 || line[0] == '#' || bytes.HasPrefix(line, []byte(`---`)) {
			continue
		}

		return false
	}
	return true
}

// Content returns the rendered manifest as raw bytes.
func (m *Manifest) Content() []byte {
	if m.content == nil {
		var buf bytes.Buffer

		buf.Write(m.resources.Bytes())
		buf.Write(m.hooks.Bytes())

		m.content = buf.Bytes()
	}

	return m.content
}

// ReadDir reads dir and returns all found manifests. It will ignore
// subdirectories.
func ReadDir(dir string) ([]*Manifest, error) {
	files, err := ioutil.ReadDir(dir)
	if os.IsNotExist(err) {
		return []*Manifest{}, nil
	}

	if err != nil {
		return nil, errors.WithStack(err)
	}

	manifests := make([]*Manifest, 0, len(files))

	for _, f := range files {
		if f.IsDir() {
			continue
		}

		ext := filepath.Ext(f.Name())
		if ext != ".yaml" && ext != ".yml" {
			continue
		}

		filename := filepath.Join(dir, f.Name())

		buf, err := ioutil.ReadFile(filename)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		resources, hooks, err := parseResources(buf)
		if err != nil {
			return nil, err
		}

		manifests = append(manifests, &Manifest{
			Name:      strings.TrimSuffix(f.Name(), ext),
			resources: resources,
			hooks:     hooks,
		})
	}

	return manifests, nil
}

// RenderDir renders manifests for all subdirectories of dir.
func RenderDir(r template.Renderer, dir string, v map[string]interface{}) ([]*Manifest, error) {
	dirs, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open component dir")
	}

	manifests := make([]*Manifest, 0)

	for _, d := range dirs {
		if !d.IsDir() {
			continue
		}

		name := d.Name()
		dirPath := filepath.Join(dir, name)

		renderedTemplates, err := r.Render(dirPath, v)
		if err != nil {
			return nil, err
		}

		manifest, err := NewManifest(name, renderedTemplates)
		if err != nil {
			return nil, err
		}

		manifests = append(manifests, manifest)
	}

	return manifests, nil
}
