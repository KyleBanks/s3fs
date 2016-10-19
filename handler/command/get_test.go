package command

import (
	"errors"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"testing"
	"time"

	"github.com/KyleBanks/s3fs/handler/command/context"
)

func TestGetCommand_Execute(t *testing.T) {
	// Function to validate the contents of a file match an expected string.
	validateFileContents := func(path, content string) {
		f, err := os.Open(path)
		if err != nil {
			t.Fatal(err)
		}
		defer f.Close()

		buf := make([]byte, len(content))
		f.Read(buf)
		if content != string(buf) {
			t.Fatalf("Unexpected file contents: %v", string(buf))
		}
	}

	// Mock the interfaces and data required to download.
	bucket := "bucket"
	folder := "folder"
	key := "file.txt"
	target := bucket + context.PathDelimiter + folder + context.PathDelimiter + key
	fileContents := "test file @ " + string(time.Now().UnixNano())

	var s3 mockS3Client
	var out mockOutputter
	var con context.Context

	s3.downloadObjectCallback = func(b, k string) (string, error) {
		if b != bucket || k != folder+context.PathDelimiter+key {
			t.Fatalf("Unexpected bucket/key provided to DownloadObject(%v, %v)", b, k)
		}

		f, _ := ioutil.TempFile("", "")
		f.WriteString(fileContents)
		return f.Name(), nil
	}

	// Positive: With destination as file
	{
		// Create a destination file to write the downloaded object to.
		dest, _ := ioutil.TempFile("", "")
		defer os.Remove(dest.Name())

		// Perform the download.
		get := NewGet(&s3, &con, []string{target, dest.Name()})
		if err := get.Execute(&out); err != nil {
			t.Fatal(err)
		}

		// Ensure the file exists and it's contents are correct.
		validateFileContents(dest.Name(), fileContents)
	}

	// Positive: With destination as folder
	{
		// Create a destination directory to write the downloaded object to.
		dest, _ := ioutil.TempDir("", "")

		// Perform the download.
		get := NewGet(&s3, &con, []string{target, dest})
		if err := get.Execute(&out); err != nil {
			t.Fatal(err)
		}

		// Ensure the file exists and it's contents are correct.
		validateFileContents(dest+string(os.PathSeparator)+key, fileContents)
	}

	// Positive: No destination
	{
		// Perform the download.
		get := NewGet(&s3, &con, []string{target})
		if err := get.Execute(&out); err != nil {
			t.Fatal(err)
		}

		// Ensure the file exists and it's contents are correct.
		expected, _ := filepath.Abs(key)
		defer os.Remove(expected)
		validateFileContents(expected, fileContents)
	}

	// Positive: With context
	{
		// Shadow the required interfaces.
		var con context.Context
		var s3 mockS3Client

		// Update context to point to a new path, and the downloadObjectCallback to
		// use the custom path.
		con.UpdatePath("bucket2/folder2/subfolder")
		s3.downloadObjectCallback = func(b, k string) (string, error) {
			if b != "bucket2" || k != "folder2/subfolder"+context.PathDelimiter+key {
				t.Fatalf("Unexpected bucket/key provided to DownloadObject(%v, %v)", b, k)
			}

			f, _ := ioutil.TempFile("", "")
			f.WriteString(fileContents)
			return f.Name(), nil
		}

		// Create a destination file to write the downloaded object to.
		dest, _ := ioutil.TempFile("", "")
		defer os.Remove(dest.Name())

		// Perform the download.
		get := NewGet(&s3, &con, []string{key, dest.Name()})
		if err := get.Execute(&out); err != nil {
			t.Fatal(err)
		}

		// Ensure the file exists and it's contents are correct.
		validateFileContents(dest.Name(), fileContents)
	}

	// Negative: No args
	{
		get := NewGet(&s3, nil, []string{})
		if err := get.Execute(&out); err == nil {
			t.Fatal("Expected Get to return an error when no arguments are provided")
		}
	}

	// Negative: S3 error
	{
		// Shadow the s3 client mock interface, and set the callback to return an error.
		var s3 mockS3Client
		mockErr := errors.New("Mock Err")
		s3.downloadObjectCallback = func(b, k string) (string, error) {
			return "", mockErr
		}

		// Perform the download.
		get := NewGet(&s3, &con, []string{target})
		if err := get.Execute(&out); err != mockErr {
			t.Fatalf("Expected the mock error to be bubbled up: %v", err)
		}
	}
}

func TestGetCommand_absDestination(t *testing.T) {
	// Get the home directory which will be required for some test cases.
	usr, err := user.Current()
	if err != nil {
		t.Fatal(err)
	}

	// Define test cases
	tests := []struct {
		args     []string
		s3Path   []string
		expected string
	}{
		{[]string{}, []string{"folder", "file1.txt"}, "file1.txt"},                               // No args
		{[]string{"", "/tmp"}, []string{"folder", "file1.txt"}, "/tmp/file1.txt"},                // Folder as destination
		{[]string{"", "/tmp/test.txt"}, []string{"folder", "file1.txt"}, "/tmp/test.txt"},        // File as destination
		{[]string{"", "~/test.txt"}, []string{"folder", "file1.txt"}, usr.HomeDir + "/test.txt"}, // Home directory
	}

	// For each test, run the destination function and validate the output.
	for _, test := range tests {
		get := NewGet(nil, nil, test.args)

		// Get the destination.
		res, err := get.absDestination(test.s3Path[len(test.s3Path)-1])
		if err != nil {
			t.Fatal(err)
		}

		// Determine the absolute path of the expected result.
		expected, err := filepath.Abs(test.expected)
		if err != nil {
			t.Fatal(err)
		}

		// Validate.
		if res != expected {
			t.Fatalf("Unexpected response for test [%v]: {Expected: %v, Actual: %v}", test, expected, res)
		}
	}
}

func TestGetCommand_IsLongRunning(t *testing.T) {
	get := NewGet(nil, nil, nil)

	if !get.IsLongRunning() {
		t.Fatal("Expected GetCommand to always be long running")
	}
}

func TestNewGet(t *testing.T) {
	var s3 mockS3Client
	var con context.Context
	args := []string{"file"}

	get := NewGet(&s3, &con, args)
	if get.s3 != &s3 {
		t.Fatalf("Unexpected S3 client stored on get command: %v", get.s3)
	} else if get.con != &con {
		t.Fatalf("Unexpected Context stored on get command: %v", get.con)
	} else if get.args[0] != args[0] {
		t.Fatalf("Unexpected args stored on get command: %v", get.args)
	}
}
