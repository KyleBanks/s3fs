package listener

import (
	"bufio"
	"os"
	"strings"
)

// StdinListener listens for incoming commands on Stdin.
type StdinListener struct {
	scanner *bufio.Scanner

	ui indicator
}

// indicator defines a UI interface to display status updates to the user.
type indicator interface {
	ShowPrompt()
}

// Listen prompts and waits for user input on Stdin.
func (s StdinListener) Listen() (in []string, ok bool) {
	// Show the UI prompt.
	s.ui.ShowPrompt()

	// Listen for user input.
	for s.scanner.Scan() {
		// Tokenize and return the command input.
		return strings.Split(s.scanner.Text(), " "), true
	}

	return nil, false
}

// NewStdin initializes and returns a new StdinListener type.
func NewStdin(ui indicator) StdinListener {
	return StdinListener{
		scanner: bufio.NewScanner(os.Stdin),
		ui:      ui,
	}
}
