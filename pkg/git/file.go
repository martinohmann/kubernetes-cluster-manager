package git

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

const (
	fileMode os.FileMode = 0660
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
	f, err := openFile(filename)
	if err != nil {
		return nil, err
	}

	prefix := filepath.Base(f.Name())
	tmpf, err := createTempFile(prefix, changes)
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
	if err := ioutil.WriteFile(c.filename, c.changes, fileMode); err != nil {
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

// openFile opens given file if it exists or creates it otherwise.
func openFile(path string) (*os.File, error) {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, fileMode)
		}
	}

	return os.OpenFile(path, os.O_RDWR, fileMode)
}

// createTemplFile creates a temporary for with given prefix and content.
func createTempFile(prefix string, content []byte) (*os.File, error) {
	f, err := ioutil.TempFile("", prefix)
	if err != nil {
		return nil, err
	}

	if _, err := f.Write(content); err != nil {
		return nil, err
	}

	return f, nil
}
