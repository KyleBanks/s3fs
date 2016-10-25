package command

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/KyleBanks/s3fs/handler/command/context"
	"github.com/KyleBanks/s3fs/handler/command/util"
)

const (
	// getArgsIndexTarget indicates the expected argument index for the target file to get.
	getArgsIndexTarget = 0

	// getArgsIndexDestination indicates the expected argument index for the destination file location.
	getArgsIndexDestination = 1
)

// GetCommand downloads a remote file.
type GetCommand struct {
	s3  S3Client
	con *context.Context

	args []string
}

// Execute performs a 'get' by downloading a remote file to a local destination.
func (get GetCommand) Execute(out Outputter) error {
	// Get the target to download from the input arguments.
	if len(get.args) < getArgsIndexTarget+1 {
		return errors.New("Missing target file.")
	}
	target := get.args[getArgsIndexTarget]

	// Calculate the S3 object path.
	path := get.con.CalculatePath(target)
	if len(path) <= 1 {
		return fmt.Errorf("Target is not a file: %v", strings.Join(path, context.PathDelimiter))
	}

	// Download the object.
	f, err := get.s3.DownloadObject(path[0], strings.Join(path[1:], context.PathDelimiter))
	if err != nil {
		return err
	}

	// Get the destination to put the downloaded file.
	dst, err := get.absDestination(path[len(path)-1])
	if err != nil {
		return err
	}

	// Move the file to the proper destination.
	if err := os.Rename(f, dst); err != nil {
		return err
	}

	return nil
}

// absDestination returns the absolute path of the destination argument where the downloaded object should be placed.
//
// If the destination argument is not set, the current working directory will be used.
// If the destination argument points to a directory, the defaultName provided will be used as the file name.
func (get GetCommand) absDestination(defaultName string) (string, error) {
	// Use the default name if the destination argument was not provided.
	dst := defaultName
	if len(get.args) >= getArgsIndexDestination+1 && len(get.args[getArgsIndexDestination]) > 0 {
		dst = get.args[getArgsIndexDestination]
	}

	// Convert to the absolute path.
	dst, err := util.AbsPath(dst)
	if err != nil {
		return "", err
	}

	// Check if the destination is a directory, and append the default name as necessary.
	if dstInfo, _ := os.Stat(dst); dstInfo != nil && dstInfo.IsDir() {
		dst = filepath.Join(dst, defaultName)
	}

	return dst, nil
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
