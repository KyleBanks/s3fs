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

// Handle takes a cmd as input and performs the required processing.
func (s S3Handler) Handle(cmd []string, out command.Outputter) error {
	if len(cmd) == 0 {
		return nil
	}

	// Determine the action to take based on the cmd.
	e, err := s.commandFromArgs(cmd)
	if err != nil {
		return err
	}

	// Show the loading indicator if applicable.
	if e.IsLongRunning() {
		s.ui.ShowLoader()
	}

	// Execute the command.
	err = e.Execute(out)

	// Notify the UI channel that we're done.
	if e.IsLongRunning() {
		s.ui.HideLoader()
	}

	return err
}

// commandFromArgs takes an arg slice and returns the appropriate command executor.
func (s S3Handler) commandFromArgs(args []string) (ex command.Executor, err error) {
	switch args[0] {

	case command.CmdLs:
		ex = command.NewLs(s.s3, s.con)
	case command.CmdCd:
		ex = command.NewCd(s.s3, s.con, args[1:])
	case command.CmdGet:
		ex = command.NewGet(s.s3, s.con, args[1:])
	case command.CmdPut:
		ex = command.NewPut(s.s3, s.con, args[1:])
	case command.CmdPwd:
		ex = command.NewPwd(s.con)
	case command.CmdClear:
		ex = command.NewClear()
	case command.CmdExit:
		ex = command.NewExit()

	default:
		err = errors.New("Unknown Command: " + args[0])
	}

	return ex, err
}

// NewS3 initializes and returns an S3Handler.
func NewS3(s3 command.S3Client, ui indicator) S3Handler {
	return S3Handler{
		s3:  s3,
		ui:  ui,
		con: &context.Context{},
	}
}
