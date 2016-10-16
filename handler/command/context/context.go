// Package context provides user state for a handler.
package context

import (
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
// UpdatePath differs from CalculatePath in that the underlying path of the context is updated.
func (c *Context) UpdatePath(p string) {
	c.path = c.CalculatePath(p)
}

// CalculatePath determines the context's path based on the string provided.
//
// For instance, provding "directory/another" will append "directory" and "another" to the path.
// Providing "../directory" will move up a level in the path and then append "directory".
//
// CalculatePath differs from UpdatePath in that the underlying path of the context is not actually updated,
// the result is simply returned.
func (c *Context) CalculatePath(p string) []string {
	// Get a reference to the current path.
	path := c.path

	// Sanity.
	if len(p) == 0 {
		return path
	}

	// If the path is equal to, or starts with, a single path delimiter (ie. 'cd /') then reset the context path.
	if p == PathDelimiter || p[:1] == PathDelimiter {
		path = make([]string, 0)
	}

	// Split path elements by the path delimiter.
	elems := strings.Split(p, PathDelimiter)

	// For each token in the path, update the context path accordingly
	for _, t := range elems {
		switch t {

		// Remove up a level by removing the last element of the path.
		case "..":
			if !c.IsRoot() {
				path = path[:len(path)-1]
			}

		// Ignore.
		case ".", "":
			// do nothing

		// Append to the path.
		default:
			path = append(path, t)
		}
	}

	return path
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
	return strings.Join(c.path[1:], PathDelimiter)
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
