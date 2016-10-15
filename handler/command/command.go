// Package command provides the ability to execute commands.
package command

const (
	// CmdLs lists directory contents based on the current context.
	CmdLs = "ls"

	// CmdCd changes directory.
	CmdCd = "cd"

	// CmdPwd prints the present working directory.
	CmdPwd = "pwd"

	// CmdClear clears the current output.
	CmdClear = "clear"

	// CmdExit exit the program.
	CmdExit = "exit"
)

// Command defines an interface for executable instructions.
type Command interface {
	// Execute runs the command.
	Execute() error

	// IsLongRunning indicates if the command may require a waiting period (ie. remote commands).
	IsLongRunning() bool
}

// s3client defines an interface that communicates with Amazon S3.
type s3client interface {
	LsBuckets() ([]string, error)
	LsObjects(bucket, prefix string) ([]string, error)
}
