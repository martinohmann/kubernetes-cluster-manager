package resource

import (
	"bytes"
	"io"
)

// Printer can print resources in a formatted way.
type Printer struct {
	w io.Writer
}

// Printer creates a new Printer with w as the backing writer for the formatted
// output.
func NewPrinter(w io.Writer) *Printer {
	return &Printer{w: w}
}

// Print prints the formatted resource.
func (p *Printer) Print(r *Resource) error {
	return p.PrintSlice(Slice{r})
}

// Print prints a formatted resource slice.
func (p *Printer) PrintSlice(s Slice) error {
	if p == nil || len(s) == 0 {
		return nil
	}

	buf := bytes.NewBufferString(FormatSlice(s))

	_, err := buf.WriteTo(p.w)

	return err
}
