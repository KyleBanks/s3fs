package command

import (
	"os"
)

// ExitCommand simulates 'exit' functionality.
type ExitCommand struct {
}

// Execute performs an 'exit' command by quitting the program.
func (exit ExitCommand) Execute() error {
	os.Exit(0)

	return nil
}

// IsLongRunning returns false because 'exit' can execute without delay.
func (ExitCommand) IsLongRunning() bool {
	return false
}

// NewExit initializes and returns a CdCommand.
func NewExit() ExitCommand {
	return ExitCommand{}
}
