package file

import (
	"bytes"
	"io/ioutil"
	"os"

	"github.com/martinohmann/go-difflib/difflib"
)

// ChangeSet is a container for a file and the changes that should be applied
// to it.
type ChangeSet struct {
	Filename string
	Changes  []byte

	content []byte
}

// NewChangeSet creates a new ChangeSet for given filename and the changes that
// should be written to it. Will return an error if reading the file fails with
// and error other than os.ErrNotExist.
func NewChangeSet(filename string, changes []byte) (*ChangeSet, error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	c := &ChangeSet{
		Filename: filename,
		Changes:  changes,
		content:  content,
	}

	return c, nil
}

// Content returns the content of the underlying file.
func (c *ChangeSet) Content() []byte {
	return c.content
}

// HasChanges returns true if the ChangeSet would produce file changes.
func (c *ChangeSet) HasChanges() bool {
	return bytes.Compare(c.Changes, c.content) != 0
}

// Apply applies the ChangeSet by replacing the content of the underlying file
// with the changes.
func (c *ChangeSet) Apply() error {
	return ioutil.WriteFile(c.Filename, c.Changes, 0660)
}

// Diff returns the diff for the ChangeSet.
func (c *ChangeSet) Diff() string {
	diff := difflib.UnifiedDiff{
		A:        difflib.SplitLines(string(c.content)),
		B:        difflib.SplitLines(string(c.Changes)),
		FromFile: c.Filename,
		ToFile:   c.Filename,
		Context:  5,
		Color:    true,
	}

	out, _ := difflib.GetUnifiedDiffString(diff)

	return out
}
