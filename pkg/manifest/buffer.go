package manifest

import "bytes"

type resourceBuffer struct {
	w bytes.Buffer
}

// Write implements io.Writer.
func (w *resourceBuffer) Write(p []byte) (n int, err error) {
	_, err = w.w.Write([]byte("---\n"))
	if err != nil {
		return
	}

	n, err = w.w.Write(p)
	if err != nil {
		return
	}

	_, err = w.w.Write([]byte("\n"))
	if err != nil {
		return
	}

	return n + 5, nil
}

// Bytes returns the content of the underlying bytes.Buffer.
func (w *resourceBuffer) Bytes() []byte {
	return w.w.Bytes()
}
