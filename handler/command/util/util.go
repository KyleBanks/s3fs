// Package util provides shared utility and convenience functions for commands.
package util

import (
	"os/user"
	"path/filepath"
	"strings"
)

const (
	// homeSymbol defines the symbol representing the HOME directory.
	homeSymbol = "~"
)

// AbsPath converts a file/directory path into an absolute path, including support for handling the home directory symbol
// represented by 'homeSymbol'.
func AbsPath(path string) (string, error) {
	// Replace the home symbol with the actual home directory if necessary.
	if strings.Contains(path, homeSymbol) {
		usr, err := user.Current()
		if err != nil {
			return "", err
		}

		path = strings.Replace(path, homeSymbol, usr.HomeDir, 1)
	}

	// Convert the path to an absolute path.
	return filepath.Abs(path)
}
