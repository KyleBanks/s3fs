// Package handler provides command processing functionality.
package handler

const (
	// cmdLs lists directory contents based on the current context.
	cmdLs = "ls"

	// cmdCd changes directory.
	cmdCd = "cd"

	// cmdPwd prints the present working directory.
	cmdPwd = "pwd"

	// cmdClear clears the current output.
	cmdClear = "clear"

	// cmdExit exit the program.
	cmdExit = "exit"
)

// Handler defines an interface that can accept and process commands.
type Handler interface {
	Handle(cmd []string) error
}
