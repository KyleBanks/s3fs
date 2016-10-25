package command

import (
	"errors"
	"os"
	"strings"

	"github.com/KyleBanks/s3fs/handler/command/context"
	"github.com/KyleBanks/s3fs/handler/command/util"
)

const (
	// putArgsIndexTarget indicates the expected argument index for the target file to put.
	putArgsIndexTarget = 0

	// putArgsIndexDestination indicates the expected argument index for the destination file location.
	putArgsIndexDestination = 1
)

// PutCommand uploads an object.
type PutCommand struct {
	s3  S3Client
	con *context.Context

	args []string
}

// Execute performs a 'put' command by uploading a file to S3.
func (p PutCommand) Execute(out Outputter) error {
	// Get the target to upload from the input arguments.
	if len(p.args) < putArgsIndexTarget+1 {
		return errors.New("Missing target file.")
	}
	target, err := util.AbsPath(p.args[putArgsIndexTarget])
	if err != nil {
		return err
	}

	// Open the target file.
	file, err := os.Open(target)
	if err != nil {
		return err
	}
	defer file.Close()

	// Get the (optional) destination.
	var destination string
	if len(p.args) >= putArgsIndexDestination+1 {
		destination = p.args[putArgsIndexDestination]
	}

	// Determine the S3 destination.
	path := p.con.CalculatePath(destination)
	if len(path) == 0 {
		return errors.New("Missing destination bucket.")
	}

	// Upload the object.
	uploadKey, err := p.s3.UploadObject(path[0], strings.Join(path[1:], context.PathDelimiter), file)
	if err != nil {
		return err
	}

	out.Write("File Uploaded: " + uploadKey)
	return nil
}

// IsLongRunning returns true because 'put' must always perform network requests.
func (PutCommand) IsLongRunning() bool {
	return true
}

// NewPut initializes and returns a PutCommand.
func NewPut(s3 S3Client, con *context.Context, args []string) PutCommand {
	return PutCommand{
		s3:   s3,
		con:  con,
		args: args,
	}
}
