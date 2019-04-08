package git

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/martinohmann/kubernetes-cluster-manager/pkg/fs"
)

// FileChanges is a container for the current content of a file and the changes
// that should be applied to it.
type FileChanges struct {
	filename string
	f        *os.File
	content  []byte
	tmpf     *os.File
	changes  []byte
}

func NewFileChanges(filename string, changes []byte) (*FileChanges, error) {
	f, err := fs.OpenFile(filename)
	if err != nil {
		return nil, err
	}

	prefix := filepath.Base(f.Name())
	tmpf, err := fs.NewTempFile(prefix, changes)
	if err != nil {
		return nil, err
	}

	c := &FileChanges{
		f:        f,
		tmpf:     tmpf,
		filename: filename,
		changes:  changes,
	}

	return c, nil
}

func (c *FileChanges) Content() []byte {
	if c.content == nil {
		content, _ := ioutil.ReadAll(c.f)
		c.content = content
	}

	return c.content
}

func (c *FileChanges) Changes() []byte {
	return c.changes
}

func (c *FileChanges) Filename() string {
	return c.tmpf.Name()
}

func (c *FileChanges) Apply() error {
	if err := fs.WriteFile(c.filename, c.changes); err != nil {
		return err
	}

	c.content = c.changes

	return nil
}

func (c *FileChanges) Close() error {
	defer os.Remove(c.tmpf.Name())

	err := c.tmpf.Close()

	if err1 := c.f.Close(); err == nil {
		err = err1
	}

	return err
}

// GetDiff creates a diff for the file changes and returns it.
func (c *FileChanges) Diff() (string, error) {
	return Diff(c.filename, c.tmpf.Name())
}
