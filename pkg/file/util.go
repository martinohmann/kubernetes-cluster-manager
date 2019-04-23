package file

import (
	"io/ioutil"
	"os"
)

// NewTempFile creates a temporary file with given prefix and content.
func NewTempFile(prefix string, content []byte) (*os.File, error) {
	f, err := ioutil.TempFile("", prefix)
	if err != nil {
		return nil, err
	}

	if _, err := f.Write(content); err != nil {
		return nil, os.Remove(f.Name())
	}

	return f, nil
}

// Exists returns true if path exists.
func Exists(path string) bool {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}

	return true
}
