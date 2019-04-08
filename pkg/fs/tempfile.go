package fs

import (
	"io/ioutil"
	"os"
)

// FileMode is used for all new files.
const FileMode os.FileMode = 0660

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

// OpenFile opens given file if it exists or creates it otherwise.
func OpenFile(path string) (*os.File, error) {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, FileMode)
		}
	}

	return os.OpenFile(path, os.O_RDWR, FileMode)
}

// WriteFile writes content to a file.
func WriteFile(path string, content []byte) error {
	return ioutil.WriteFile(path, content, FileMode)
}
