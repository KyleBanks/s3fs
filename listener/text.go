package listener

import (
	"strings"
)

const (
	// cmdDelimiter is used to separate multiple commands in the same input.
	cmdDelimiter = "&&"

	// argDelimiter is used to separate arguments for a command.
	argDelimiter = " "
)

// TextListener listens for incoming text commands.
type TextListener struct {
	input inputter

	ui indicator
}

// Listen prompts and waits for user input on Stdin.
func (t TextListener) Listen() ([]InputCommand, bool) {
	// Show the UI prompt.
	t.ui.ShowPrompt()

	// Listen for user input.
	for t.input.Scan() {
		// Split input command(s) by the cmdDelimiter.
		cmdStrs := strings.Split(t.input.Text(), cmdDelimiter)
		cmds := make([]InputCommand, len(cmdStrs))

		// For each command, construct the InputCommand.
		for i, str := range cmdStrs {
			str = strings.TrimSpace(str)

			// Split the command arguments by the argDelimiter and construct the InputCommand.
			cmds[i] = InputCommand{
				Args: strings.Split(str, argDelimiter),
			}
		}

		return cmds, true
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
