package manifest

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/martinohmann/kubernetes-cluster-manager/pkg/hook"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/resource"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/template"
	"github.com/pkg/errors"
)

// Manifest contains a kubernetes manifest split into resources and hooks.
type Manifest struct {
	Name      string
	Resources resource.Slice
	Hooks     hook.SliceMap

	content []byte
}

// New creates a new manifest with name and given content. Will error if
// content parsing fails.
func New(name string, content []byte) (*Manifest, error) {
	resources, hooks, err := Parse(content)
	if err != nil {
		return nil, err
	}

	m := &Manifest{
		Name:      name,
		Resources: resources,
		Hooks:     hooks,
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

// Content returns the rendered manifest as raw bytes. Resources and hooks are
// sorted to make the output of this stable.
func (m *Manifest) Content() []byte {
	if m.content == nil {
		var buf bytes.Buffer

		buf.Write(m.Resources.Sort(resource.ApplyOrder).Bytes())
		buf.Write(m.Hooks.SortSlices().Bytes())

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

		name := strings.TrimSuffix(f.Name(), ext)
		filename := filepath.Join(dir, f.Name())

		buf, err := ioutil.ReadFile(filename)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		manifest, err := New(name, buf)
		if err != nil {
			return nil, err
		}

		manifests = append(manifests, manifest)
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

		var buf resource.Buffer

		for path, content := range renderedTemplates {
			ext := filepath.Ext(path)

			if ext != ".yaml" && ext != ".yml" {
				// We are only interested in rendered yaml files and will just
				// discard the rest.
				continue
			}

			buf.Write([]byte(content))
		}

		manifest, err := New(name, buf.Bytes())
		if err != nil {
			return nil, err
		}

		manifests = append(manifests, manifest)
	}

	return manifests, nil
}

// FindMatching finds a manifest in a haystack. Matching is done by name.
func FindMatching(haystack []*Manifest, needle *Manifest) (*Manifest, bool) {
	for _, m := range haystack {
		if m.Name == needle.Name {
			return m, true
		}
	}

	return nil, false
}
