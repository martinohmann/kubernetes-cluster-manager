package diff

import (
	"bytes"
	"io"

	"github.com/kr/text"
)

// Printer prints diffs.
type Printer struct {
	w io.Writer
}

// NewPrinter creates a new Printer with w as the backing writer for the formatted
// output.
func NewPrinter(w io.Writer) *Printer {
	return &Printer{w: w}
}

// Print prints the formatted resource.
func (p *Printer) Print(o Options) error {
	if p == nil {
		return nil
	}

	diff := Diff(o)

	if diff == "" {
		return nil
	}

	var buf bytes.Buffer

	if o.Filename != "" {
		buf.WriteString("changes to ")
		buf.WriteString(o.Filename)
		buf.WriteString(":\n")
	}

	buf.WriteByte('\n')

	buf.WriteString(text.Indent(diff, "  "))

	_, err := buf.WriteTo(p.w)

	return err
}
