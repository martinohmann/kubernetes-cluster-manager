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
	f    string
	a, b []byte
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
		f: filename,
		a: content,
		b: changes,
	}

	return c, nil
}

// Filename returns the filename for the change set.
func (c *ChangeSet) Filename() string {
	return c.f
}

// HasChanges returns true if the ChangeSet would produce file changes.
func (c *ChangeSet) HasChanges() bool {
	return bytes.Compare(c.a, c.b) != 0
}

// Apply applies the ChangeSet by replacing the content of the underlying file
// with the changes.
func (c *ChangeSet) Apply() error {
	return ioutil.WriteFile(c.f, c.b, 0660)
}

// Diff returns the diff for the ChangeSet.
func (c *ChangeSet) Diff() string {
	diff := difflib.UnifiedDiff{
		A:        difflib.SplitLines(string(c.a)),
		B:        difflib.SplitLines(string(c.b)),
		FromFile: c.f,
		ToFile:   c.f,
		Context:  5,
		Color:    true,
	}

	out, _ := difflib.GetUnifiedDiffString(diff)

	return out
}
