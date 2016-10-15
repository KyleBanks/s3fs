package handler

import (
	"errors"

	"github.com/KyleBanks/s3fs/handler/command"
	"github.com/KyleBanks/s3fs/handler/command/context"
)

// S3Handler defines a struct that handles commands and dispatches them through the Amazon S3 API.
type S3Handler struct {
	s3 s3client
	ui indicator

	con *context.Context
}

// s3client defines an interface that communicates with Amazon S3.
type s3client interface {
	LsBuckets() ([]string, error)
	LsObjects(bucket, prefix string) ([]string, error)
}

// indicator defines a UI interface to display status updates to the user.
type indicator interface {
	ShowLoader()
	HideLoader()
}

// Handle takes a cmd as input and performs the required processing.
func (s S3Handler) Handle(cmd []string) error {
	if len(cmd) == 0 {
		return nil
	}

	// Determine the action to take based on the cmd.
	var c command.Command

	switch cmd[0] {

	// API operations:
	case command.CmdLs:
		c = command.NewLs(s.s3, s.con)

	// Context operations:
	case command.CmdCd:
		c = command.NewCd(s.con, cmd[1:])
	case command.CmdPwd:
		c = command.NewPwd(s.con)
	case command.CmdClear:
		c = command.NewClear()
	case command.CmdExit:
		c = command.NewExit()

	default:
		return errors.New("Unknown Command: " + cmd[0])
	}

	// Show the loading indicator if applicable.
	if c.IsLongRunning() {
		s.ui.ShowLoader()
	}

	// Execute the command.
	err := c.Execute()

	// Notify the UI channel that we're done.
	if c.IsLongRunning() {
		s.ui.HideLoader()
	}

	return err
}

// NewS3 initializes and returns an S3Handler.
func NewS3(s s3client, ui indicator) S3Handler {
	return S3Handler{
		s3:  s,
		ui:  ui,
		con: &context.Context{},
	}
}
