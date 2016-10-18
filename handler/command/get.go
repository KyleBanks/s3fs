package command

import (
	"github.com/KyleBanks/s3fs/handler/command/context"
)

// GetCommand downloads a remote file.
type GetCommand struct {
	s3  S3Client
	con *context.Context

	args []string
}

// Execute performs a 'get' by downloading a remote file.
func (get GetCommand) Execute(out Outputter) error {
	return nil
}

// IsLongRunning returns true because a 'get' is always long running.
func (GetCommand) IsLongRunning() bool {
	return true
}

// NewGet initializes and returns a GetCommand.
func NewGet(s3 S3Client, con *context.Context, args []string) GetCommand {
	return GetCommand{
		s3:   s3,
		con:  con,
		args: args,
	}
}
