package indicator

import (
	"time"
)

const (
	// loaderSleepTime is the time between loading indicator updates.
	loaderSleepTime = time.Millisecond * 200

	// loaderText is the text displayed when the loading indicator is enabled.
	loaderText = "."
	// promptText is the text displayed when ShowPrompt() is called.
	promptText = "> "
)

// CommandLine provides UI indications to the command line.
type CommandLine struct {
	stopLoading chan bool

	out stringWriter
}

// ShowLoader displays a command line loading indicator.
func (c *CommandLine) ShowLoader() {
	go c.startLoading()
}

// HideLoader hides the command line loading indicator.
func (c *CommandLine) HideLoader() {
	c.stopLoading <- true
}

// ShowPrompt displays a command line prompt for input.
func (c *CommandLine) ShowPrompt() {
	c.out.Write(promptText)
}

// startLoading initializes prints the loading indicator until the stop signal is received.
func (c *CommandLine) startLoading() {
	for {
		select {

		// Check if we need to stop loading.
		case <-c.stopLoading:
			// Always write a blank line after loading finishes.
			c.out.Write("\n")
			return

		// Update the loader indicator as required.
		default:
			c.out.Write(loaderText)

			// Sleep a while before displaying the loading indicator again.
			time.Sleep(loaderSleepTime)
		}
	}
}

// NewCommandLine initializes and returns a new CommandLine.
func NewCommandLine(out stringWriter) *CommandLine {
	return &CommandLine{
		out:         out,
		stopLoading: make(chan bool),
	}
}
