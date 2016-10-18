package command

import (
	"os"
	"os/exec"
	"runtime"
)

// ClearCommand simulates 'clear' functionality.
type ClearCommand struct{}

// Execute performs a 'clear' command by clearing the current program output.
// Note: Only applicable when running via command line.
func (clear ClearCommand) Execute(out Outputter) error {
	// Determine the clear command based on the architecture
	cmd := clear.cmdForSys(string(runtime.GOOS))

	// Execute the command.
	c := exec.Command(cmd)
	c.Stdout = os.Stdout
	c.Run()

	return nil
}

// cmdForSys returns the appropriate 'clear' command for the system provided.
func (ClearCommand) cmdForSys(sys string) string {
	switch sys {
	case "windows":
		return "cls"
	default:
		return "clear"
	}
}

// IsLongRunning returns false because 'clear' can execute without delay.
func (ClearCommand) IsLongRunning() bool {
	return false
}

// NewClear initializes and returns a CdCommand.
func NewClear() ClearCommand {
	return ClearCommand{}
}
