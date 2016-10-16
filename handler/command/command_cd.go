package command

import (
	"errors"
	"strings"

	"github.com/KyleBanks/s3fs/handler/command/context"
)

// CdCommand simulates 'cd' functionality.
type CdCommand struct {
	s3  S3Client
	con *context.Context

	args []string
}

// Execute performs a 'cd' command by updating the underlying context path.
func (cd CdCommand) Execute(out Outputter) error {
	// Sanity, do nothing for cd without args.
	if len(cd.args) == 0 {
		return nil
	}
	target := cd.args[0]

	// Validate that we can 'cd' into the target.
	var ok bool
	var err error

	// Calculate the target path
	targetPath := cd.con.CalculatePath(target)

	// Perform the check based on the target path.
	switch len(targetPath) {

	// You can always cd back to root.
	case 0:
		ok = true

	// Targeting a bucket.
	case 1:
		ok, err = cd.s3.BucketExists(targetPath[0])

	// Targeting a specific directory.
	default:
		ok, err = cd.s3.ObjectExists(targetPath[0], strings.Join(targetPath[1:], context.PathDelimiter)+context.PathDelimiter)
	}

	// Ensure we can perform the command.
	if err != nil {
		return err
	} else if !ok {
		return errors.New("Cannot change into non-existent directory: " + target)
	}

	// Valid target, update the context path.
	cd.con.UpdatePath(target)

	return nil
}

// IsLongRunning returns true because 'cd' must validate with S3 that the target exists.
func (CdCommand) IsLongRunning() bool {
	return true
}

// NewCd initializes and returns a CdCommand.
func NewCd(s3 S3Client, con *context.Context, args []string) CdCommand {
	return CdCommand{
		s3:   s3,
		con:  con,
		args: args,
	}
}
