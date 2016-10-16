package command

import (
	"fmt"
	"strings"

	"github.com/KyleBanks/s3fs/handler/command/context"
)

const (
	// bucketPrefix is the prefix used when outputting bucket names.
	bucketPrefix = "[B]"

	// folderPrefix is the prefix used when outputting folder names.
	folderPrefix = "[F]"

	// filePrefix is the prefix used when outputting file names.
	filePrefix = "   "
)

// LsCommand simulates 'ls' functionality.
type LsCommand struct {
	s3  S3Client
	con *context.Context
}

// Execute performs a 'ls' command by printing the buckets/objects in the pwd based on the underlying context.
func (ls LsCommand) Execute(out Outputter) error {
	var res []string
	var err error
	var prefix string
	var isBucketList bool

	// Determine which type of 'ls' to perform based on the context.
	if ls.con.IsRoot() {
		isBucketList = true
		res, err = ls.s3.LsBuckets()
	} else {
		// If we have a prefix, store it and provide it to the LsObject command.
		prefix = ls.con.PathWithoutBucket()

		res, err = ls.s3.LsObjects(ls.con.Bucket(), prefix)
	}

	// Sanity.
	if err != nil {
		return err
	}

	// Add a blank line prior to printing to ensure we don't mix up the first object/bucket name with
	// previous output (ie. the loading indicator).
	out.Write("\n")

	// Print the 'ls' results, grouping folders together.
	cache := make(map[string]bool)
	for _, f := range res {
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
			out.Write(prefixOutput(f, isBucketList) + "\n")
			cache[f] = true
		}
	}

	return nil
}

// prefixOutput returns a modified version of a bucket/folder/filename by prepending the appropriate prefix.
func prefixOutput(out string, isBucketList bool) string {
	var prefix string

	if isBucketList { // Check if it's a bucket...
		prefix = bucketPrefix
	} else if string(out[len(out)-1]) == context.PathDelimiter { // ... or folder...
		prefix = folderPrefix
	} else { // ... must be a file.
		prefix = filePrefix
	}

	return fmt.Sprintf("%v %v", prefix, out)
}

// IsLongRunning returns true because 'ls' requires a network operation.
func (LsCommand) IsLongRunning() bool {
	return true
}

// NewLs initializes and returns an LsCommand.
func NewLs(s3 S3Client, con *context.Context) LsCommand {
	return LsCommand{
		s3:  s3,
		con: con,
	}
}
