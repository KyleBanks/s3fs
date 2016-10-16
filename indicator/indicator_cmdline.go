package indicator

import (
	"sync"
	"time"
)

// CommandLineIndicator provides UI indications to the command line.
type CommandLineIndicator struct {
	mu      sync.Mutex
	loading bool

	out stringWriter
}

// stringWriter defines an interface that can recieve strings.
type stringWriter interface {
	Write(string)
}

// ShowLoader displays a command line loading indicator.
func (c *CommandLineIndicator) ShowLoader() {
	c.mu.Lock()
	c.loading = true
	c.mu.Unlock()
}

// HideLoader hides the command line loading indicator.
func (c *CommandLineIndicator) HideLoader() {
	c.mu.Lock()
	c.loading = false
	c.mu.Unlock()

	// Always write a blank line after loading finishes.
	c.out.Write("\n")
}

// ShowPrompt displays a command line prompt for input.
func (c *CommandLineIndicator) ShowPrompt() {
	c.out.Write("> ")
}

// init initializes the CommandLineIndicator.
//
// Note: This should only be called once!
func (c *CommandLineIndicator) init() {
	// Start the loading indicator goroutine.
	go func() {
		for {
			// Check if we're loading.
			c.mu.Lock()
			loading := c.loading
			c.mu.Unlock()

			// Update the loader if applicable.
			if loading {
				c.out.Write(".")
				time.Sleep(time.Millisecond * 200)
			}
		}
	}()
}

// NewCommandLine initializes and returns a new CommandLineIndicator.
func NewCommandLine(out stringWriter) *CommandLineIndicator {
	c := CommandLineIndicator{
		out: out,
	}
	c.init()

	return &c
}
