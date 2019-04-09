package git

import (
	"os"
	"path/filepath"

	"github.com/martinohmann/kubernetes-cluster-manager/pkg/fs"
)

// FileChanges is a container for a file and the changes that should be applied
// to it.
type FileChanges struct {
	filename string
	tmpf     *os.File
	changes  []byte
}

// NewFileChanges creates a new FileChanges value for given filename and the
// changes that should be written to it.
func NewFileChanges(filename string, changes []byte) (*FileChanges, error) {
	prefix := filepath.Base(filename)
	tmpf, err := fs.NewTempFile(prefix, changes)
	if err != nil {
		return nil, err
	}

	c := &FileChanges{
		tmpf:     tmpf,
		filename: filename,
		changes:  changes,
	}

	return c, nil
}

// Apply replaces the file content with the pending changes.
func (c *FileChanges) Apply() error {
	if err := fs.WriteFile(c.filename, c.changes); err != nil {
		return err
	}

	return nil
}

// Close implements io.Closer. Will remove any temporary files that were
// created during the lifecycle of the value.
func (c *FileChanges) Close() error {
	defer os.Remove(c.tmpf.Name())

	return c.tmpf.Close()
}

// Diff creates a diff for the file changes and returns it.
func (c *FileChanges) Diff() (string, error) {
	return Diff(c.filename, c.tmpf.Name())
}
