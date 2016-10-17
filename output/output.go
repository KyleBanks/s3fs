// Package output provides the ability to output system status and command results to the user.
package output

import (
	"fmt"
	"io"
)

// OutputWriter writes output to it's underlying io.Writer.
type OutputWriter struct {
	w io.Writer
}

// Write prints a string to the underlying Writer.
func (o OutputWriter) Write(out string) {
	fmt.Fprint(o.w, out)
}

// New intializes and returns an OutputWriter.
func New(w io.Writer) OutputWriter {
	return OutputWriter{
		w: w,
	}
}
