package file

import (
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
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

// LoadYAML loads the contents of filename and unmarshals it into v.
func LoadYAML(filename string, v interface{}) error {
	buf, err := ioutil.ReadFile(filename)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	return yaml.Unmarshal(buf, v)
}
