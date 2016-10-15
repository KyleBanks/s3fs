package indicator

import (
	"fmt"
	"sync"
	"time"
)

// CommandLineIndicator provides UI indications to the command line.
type CommandLineIndicator struct {
	mu      sync.Mutex
	loading bool
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
}

// ShowPrompt displays a command line prompt for input.
func (*CommandLineIndicator) ShowPrompt() {
	fmt.Printf("> ")
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
				fmt.Print(".")
				time.Sleep(time.Millisecond * 200)
			}
		}
	}()
}

// NewCommandLine initializes and returns a new CommandLineIndicator.
func NewCommandLine() *CommandLineIndicator {
	var c CommandLineIndicator
	c.init()

	return &c
}
