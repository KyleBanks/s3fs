// Package context provides user state for a handler.
package context

import (
	"fmt"
	"strings"
)

const (
	// PathDelimiter is the delimiter to use between file path components.
	PathDelimiter = "/"
)

// Context represents the metadata of the current handler session.
type Context struct {
	path []string // Element zero is always the bucket name, the rest are path prefixes for folders in buckets
}

// UpdatePath modifies the context's path based on the string provided.
//
// For instance, provding "directory/another" will append "directory" and "another" to the path.
// Providing "../directory" will move up a level in the path and then append "directory".
func (c *Context) UpdatePath(p string) {
	// If the path is equal to a single path delimiter (ie. 'cd /') just reset the context path.
	if p == PathDelimiter {
		c.path = make([]string, 0)
		return
	}

	// Split path elements by the path delimiter.
	elems := strings.Split(p, PathDelimiter)

	// For each token in the path, update the context path accordingly
	for _, t := range elems {
		switch t {

		// Remove last element.
		case "..":
			if !c.IsRoot() {
				c.path = c.path[:len(c.path)-1]
			}

		// Ignore.
		case ".", "":
			// do nothing

		// Append to the path.
		default:
			c.path = append(c.path, t)
		}
	}
}

// IsRoot indicates if the context is at the root
func (c *Context) IsRoot() bool {
	return len(c.path) == 0
}

// Path returns the full current path as a string.
func (c *Context) Path() string {
	return strings.Join(c.path[:], PathDelimiter)
}

// PathWithoutBucket returns the current path without the Bucket, as a string.
func (c *Context) PathWithoutBucket() string {
	// Check if the path is empty or only contains a bucket.
	if len(c.path) <= 1 {
		return ""
	}

	// Return the path, without bucket, as a joined string delimited by PathDelimiter.
	return fmt.Sprintf("%v%v", strings.Join(c.path[1:], PathDelimiter), PathDelimiter)
}

// Bucket returns the bucket name as a string.
func (c *Context) Bucket() string {
	// Check if the path is empty.
	if len(c.path) == 0 {
		return ""
	}

	// Return the bucket only.
	return c.path[0]
}
