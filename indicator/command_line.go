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
	promptText = "\n> "
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
	var didPrint bool

	for {
		// Sleep a while before and between each loading indicator print.
		time.Sleep(loaderSleepTime)

		select {

		// Check if we need to stop loading.
		case <-c.stopLoading:
			if didPrint {
				c.out.Write("\n")
			}
			return

		// Update the loader indicator as required.
		default:
			didPrint = true
			c.out.Write(loaderText)
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
