package helm

import (
	"io/ioutil"
	"os"
)

func writeFile(path string, content []byte) error {
	return ioutil.WriteFile(path, content, 0640)
}

func createFileIfNotExists(path string) error {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			f, err := os.Create(path)
			if err != nil {
				return err
			}
			f.Close()
		}
	}

	return nil
}

func next(filename string) string {
	return filename + ".next"
}
