package listener

import (
	"strings"
)

// TextListener listens for incoming text commands.
type TextListener struct {
	input inputter

	ui indicator
}

// Listen prompts and waits for user input on Stdin.
func (t TextListener) Listen() (in []string, ok bool) {
	// Show the UI prompt.
	t.ui.ShowPrompt()

	// Listen for user input.
	for t.input.Scan() {
		// Tokenize and return the command input.
		return strings.Split(t.input.Text(), " "), true
	}

	return nil, false
}

// NewText initializes and returns a new TextListener type.
func NewText(ui indicator, input inputter) TextListener {
	return TextListener{
		input: input,
		ui:    ui,
	}
}
