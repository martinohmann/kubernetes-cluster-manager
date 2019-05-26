package resource

import "bytes"

// Buffer wraps a bytes.Buffer and delimits the bytes of every write so that
// the resulting byte slice is valid multi-document yaml.
type Buffer struct {
	w bytes.Buffer
}

// Write implements io.Writer.
func (w *Buffer) Write(p []byte) (n int, err error) {
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
func (w *Buffer) Bytes() []byte {
	return w.w.Bytes()
}
