package command

import (
	"os"
	"os/exec"
	"runtime"
)

// ClearCommand simulates 'clear' functionality.
type ClearCommand struct {
}

// Execute performs a 'clear' command by clearing the current program output.
// Note: Only applicable when running via command line.
func (clear ClearCommand) Execute(out Outputter) error {
	var cmd string

	// Determine the clear command based on the architecture
	switch runtime.GOOS {
	case "windows":
		cmd = "cls"
	default:
		cmd = "clear"
	}

	// Execute the command.
	c := exec.Command(cmd)
	c.Stdout = os.Stdout
	c.Run()

	return nil
}

// IsLongRunning returns false because 'clear' can execute without delay.
func (ClearCommand) IsLongRunning() bool {
	return false
}

// NewClear initializes and returns a CdCommand.
func NewClear() ClearCommand {
	return ClearCommand{}
}
