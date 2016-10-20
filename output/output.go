// Package output provides the ability to output system status and command results to the user.
package output

import (
	"fmt"
	"io"
)

// Output writes output to it's underlying io.Writer.
type Output struct {
	w io.Writer
}

// Write prints a string to the underlying Writer.
func (o Output) Write(out string) {
	fmt.Fprint(o.w, out)
}

// New intializes and returns an Output.
func New(w io.Writer) Output {
	return Output{
		w: w,
	}
}
