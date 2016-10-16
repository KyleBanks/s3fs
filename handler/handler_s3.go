package handler

import (
	"errors"

	"github.com/KyleBanks/s3fs/handler/command"
	"github.com/KyleBanks/s3fs/handler/command/context"
)

// S3Handler defines a struct that handles commands and dispatches them through the Amazon S3 API.
type S3Handler struct {
	s3 command.S3Client
	ui indicator

	con *context.Context
}

// indicator defines a UI interface to display status updates to the user.
type indicator interface {
	ShowLoader()
	HideLoader()
}

// Handle takes a cmd as input and performs the required processing.
func (s S3Handler) Handle(cmd []string, out command.Outputter) error {
	if len(cmd) == 0 {
		return nil
	}

	// Determine the action to take based on the cmd.
	var e command.Executer

	switch cmd[0] {

	// API operations:
	case command.CmdLs:
		e = command.NewLs(s.s3, s.con)

	// Context operations:
	case command.CmdCd:
		e = command.NewCd(s.s3, s.con, cmd[1:])
	case command.CmdPwd:
		e = command.NewPwd(s.con)
	case command.CmdClear:
		e = command.NewClear()
	case command.CmdExit:
		e = command.NewExit()

	default:
		return errors.New("Unknown Command: " + cmd[0])
	}

	// Show the loading indicator if applicable.
	if e.IsLongRunning() {
		s.ui.ShowLoader()
	}

	// Execute the command.
	err := e.Execute(out)

	// Notify the UI channel that we're done.
	if e.IsLongRunning() {
		s.ui.HideLoader()
	}

	return err
}

// NewS3 initializes and returns an S3Handler.
func NewS3(s command.S3Client, ui indicator) S3Handler {
	return S3Handler{
		s3:  s,
		ui:  ui,
		con: &context.Context{},
	}
}
