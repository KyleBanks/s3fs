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

// Executor defines an interface for executable instructions.
type Executor interface {
	// Execute runs the command.
	Execute(Outputter) error

	// IsLongRunning indicates if the command may require a waiting period (ie. remote commands).
	IsLongRunning() bool
}

// Outputter defines a type that can receive command output in the form of strings.
type Outputter interface {
	Write(string)
}

// S3Client defines an interface that communicates with Amazon S3.
type S3Client interface {
	LsBuckets() ([]string, error)
	LsObjects(bucket, prefix string) ([]string, error)

	BucketExists(string) (bool, error)
	ObjectExists(string, string) (bool, error)
}
