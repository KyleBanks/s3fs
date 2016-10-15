package handler

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/KyleBanks/s3fs/handler/context"
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
	var err error
	switch cmd[0] {

	// API operations:
	case cmdLs:
		s.ui.ShowLoader()
		err = s.handleLs()

	// Context operations:
	case cmdCd:
		err = s.handleCd(cmd)
	case cmdPwd:
		err = s.handlePwd()
	case cmdClear:
		err = s.handleClear()
	case cmdExit:
		os.Exit(0)

	default:
		fmt.Println("Unknown Command:", cmd[0])
	}

	// Notify the UI channel that we're done.
	s.ui.HideLoader()

	return err
}

// handleLs handles an 'ls' command by printing the buckets/objects in the pwd based on the underlying context.
func (s S3Handler) handleLs() error {
	var ls []string
	var err error
	var prefix string

	// Determine which type of 'ls' to perform based on the context.
	if s.con.IsRoot() {
		ls, err = s.s3.LsBuckets()
	} else {
		// If we have a prefix, store it and provide it to the LsObject command.
		prefix = s.con.PathWithoutBucket()

		ls, err = s.s3.LsObjects(s.con.Bucket(), prefix)
	}

	// Sanity.
	if err != nil {
		return err
	}

	// Add a blank line prior to printing to ensure we don't mix up the first object/bucket name with
	// previous output (ie. the loading indicator).
	fmt.Println("")

	// Print the 'ls' results, grouping folders together.
	cache := make(map[string]bool)
	for _, f := range ls {
		// Remove the prefix if applicable.
		if len(prefix) > 0 && strings.Contains(f, prefix) {
			f = strings.Replace(f, prefix, "", 1)
		}

		// Only display the folder name if present.
		if strings.Contains(f, context.PathDelimiter) {
			f = fmt.Sprintf("%v%v", strings.Split(f, context.PathDelimiter)[0], context.PathDelimiter)
		}

		// Check if we've already printed this key.
		if _, ok := cache[f]; !ok {
			fmt.Println(f)
			cache[f] = true
		}
	}

	return nil
}

// handleCd handles the 'cd' command by updating the underlying context path.
func (s S3Handler) handleCd(cmd []string) error {
	s.con.UpdatePath(cmd[1])

	return nil
}

// handlePwd handles a 'pwd' command by printing the present working directory of the underlying context.
func (s S3Handler) handlePwd() error {
	// Print the path followed by a PathDelimiter and newline.
	fmt.Printf("%v%v\n", s.con.Path(), context.PathDelimiter)

	return nil
}

// handleClear handles the 'clear' command by creating the current output.
func (S3Handler) handleClear() error {
	var cmd string

	// Determine the clear command based on the architecture
	switch runtime.GOOS {
	case "windows":
		cmd = "cls"
	default:
		cmd = "clear"
	}

	// Execute the command.
	c := exec.Command(cmd)
	c.Stdout = os.Stdout
	c.Run()

	return nil
}

// NewS3 initializes and returns an S3Handler.
func NewS3(s s3client, ui indicator) S3Handler {
	return S3Handler{
		s3:  s,
		ui:  ui,
		con: &context.Context{},
	}
}
