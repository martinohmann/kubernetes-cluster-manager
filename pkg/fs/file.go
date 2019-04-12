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
	if Exists(path) {
		return os.OpenFile(path, os.O_RDWR, FileMode)
	}

	return os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, FileMode)
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

// Touch touches the file a path.
func Touch(path string) error {
	f, err := OpenFile(path)
	defer f.Close()

	return err
}

// WriteFile writes content to a file.
func WriteFile(path string, content []byte) error {
	return ioutil.WriteFile(path, content, FileMode)
}
