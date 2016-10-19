package util

import (
	"os/user"
	"path/filepath"
	"testing"
)

func TestAbsPath(t *testing.T) {
	// Get the home directory which will be required for some test cases.
	usr, err := user.Current()
	if err != nil {
		t.Fatal(err)
	}

	// Get the current directory.
	dir, err := filepath.Abs("")
	if err != nil {
		t.Fatal(err)
	}

	// Define the test cases
	tests := []struct {
		input  string
		output string
	}{
		{"~", usr.HomeDir},
		{"~/", usr.HomeDir},
		{"~/file.txt", usr.HomeDir + "/file.txt"},
		{"/tmp", "/tmp"},
		{"/tmp/file.txt", "/tmp/file.txt"},
		{".", dir},
		{"test", dir + "/test"},
		{"./file.txt", dir + "/file.txt"},
		{"test/file.txt", dir + "/test/file.txt"},
	}

	for _, test := range tests {
		if out, err := AbsPath(test.input); err != nil {
			t.Fatal(err)
		} else if out != test.output {
			t.Fatalf("Unexpected output: {Expected: %v, Actual: %v}", test.output, out)
		}
	}
}
