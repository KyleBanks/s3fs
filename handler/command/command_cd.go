package command

import (
	"github.com/KyleBanks/s3fs/handler/command/context"
)

// CdCommand simulates 'cd' functionality.
type CdCommand struct {
	con *context.Context

	args []string
}

// Execute performs a 'cd' command by updating the underlying context path.
func (cd CdCommand) Execute() error {
	cd.con.UpdatePath(cd.args[0])

	return nil
}

// IsLongRunning returns false because 'cd' can execute without delay.
func (CdCommand) IsLongRunning() bool {
	return false
}

// NewCd initializes and returns a CdCommand.
func NewCd(con *context.Context, args []string) CdCommand {
	return CdCommand{
		con:  con,
		args: args,
	}
}
