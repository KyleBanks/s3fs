package output

import (
	"fmt"
)

// Stdout is an output type that writes to stdout.
type Stdout struct {
}

// Write takes a string and writes it to stdout.
func (Stdout) Write(s string) {
	fmt.Printf(s)
}
