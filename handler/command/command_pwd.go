package command

import (
	"fmt"

	"github.com/KyleBanks/s3fs/handler/command/context"
)

// PwdCommand simulates 'pwd' functionality.
type PwdCommand struct {
	con *context.Context
}

// Execute performs a 'pwd' command by printing the current working path.
func (pwd PwdCommand) Execute() error {
	// Print the path followed by a PathDelimiter and newline.
	fmt.Printf("%v%v\n", pwd.con.Path(), context.PathDelimiter)

	return nil
}

// IsLongRunning returns false because 'pwd' can execute without delay.
func (PwdCommand) IsLongRunning() bool {
	return false
}

// NewPwd initializes and returns a PwdCommand.
func NewPwd(con *context.Context) PwdCommand {
	return PwdCommand{
		con: con,
	}
}
