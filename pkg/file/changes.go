package file

import (
	"io/ioutil"
	"os"

	"github.com/martinohmann/go-difflib/difflib"
)

// Changes is a container for a file and the changes that should be applied
// to it.
type Changes struct {
	filename string
	changes  []byte
}

// NewChanges creates a new Changes value for given filename and the
// changes that should be written to it.
func NewChanges(filename string, changes []byte) *Changes {
	return &Changes{
		filename: filename,
		changes:  changes,
	}
}

// Apply replaces the file content with the pending changes.
func (c *Changes) Apply() error {
	return ioutil.WriteFile(c.filename, c.changes, 0660)
}

// Diff creates a diff for the file changes and returns it.
func (c *Changes) Diff() (string, error) {
	buf, err := ioutil.ReadFile(c.filename)
	if err != nil && !os.IsNotExist(err) {
		return "", err
	}

	diff := difflib.UnifiedDiff{
		A:        difflib.SplitLines(string(buf)),
		B:        difflib.SplitLines(string(c.changes)),
		FromFile: c.filename,
		ToFile:   c.filename,
		Context:  5,
		Color:    true,
	}

	return difflib.GetUnifiedDiffString(diff)
}
