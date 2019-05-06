package manifest

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

// Manifest contains a kubernetes manifest as raw bytes and its name.
type Manifest struct {
	Name    string
	Content []byte
}

// Filename returns the filename for the manifest.
func (m *Manifest) Filename() string {
	return fmt.Sprintf("%s.yaml", m.Name)
}

// Matches returns true if other matches m.
func (m *Manifest) Matches(other *Manifest) bool {
	if m == other {
		return true
	}

	if m == nil || other == nil {
		return false
	}

	return m.Name == other.Name
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

		m := &Manifest{
			Name:    strings.TrimSuffix(f.Name(), ext),
			Content: buf,
		}

		manifests = append(manifests, m)
	}

	return manifests, nil
}
