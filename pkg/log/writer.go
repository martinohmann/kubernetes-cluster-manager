package log

import (
	"bufio"
	"bytes"
)

// LineWriter is a printf-style func which satisfies the io.Writer interface.
type LineWriter func(args ...interface{})

// Write implements io.Writer.
func (w LineWriter) Write(p []byte) (n int, err error) {
	s := bufio.NewScanner(bytes.NewReader(p))
	s.Split(bufio.ScanLines)

	for s.Scan() {
		w(s.Text())
	}

	return len(p), nil
}
