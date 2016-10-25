package command

import (
	"errors"
	"io/ioutil"
	"os"
	"testing"

	"github.com/KyleBanks/s3fs/handler/command/context"
)

func TestPutCommand_Execute(t *testing.T) {
	// Positive Case
	{
		var s3 mockS3Client
		var con context.Context
		con.UpdatePath("bucket/folder")
		var out mockOutputter
		file, _ := ioutil.TempFile("", "")
		defer os.Remove(file.Name())
		args := []string{file.Name()}

		s3.uploadObjectCallback = func(bucket, key string, f *os.File) (string, error) {
			if bucket != "bucket" || key != "folder" || f.Name() != file.Name() {
				t.Fatalf("Unexpected input to UploadObject(%v, %v, %v)", bucket, key, f.Name())
			}

			return "", nil
		}

		put := NewPut(&s3, &con, args)
		if err := put.Execute(&out); err != nil {
			t.Fatal(err)
		}
	}

	// Positive Case: With Destination
	{
		var s3 mockS3Client
		var con context.Context
		con.UpdatePath("bucket")
		var out mockOutputter
		file, _ := ioutil.TempFile("", "")
		defer os.Remove(file.Name())
		args := []string{file.Name(), "folder/subfolder"}

		s3.uploadObjectCallback = func(bucket, key string, f *os.File) (string, error) {
			if bucket != "bucket" || key != "folder/subfolder" || f.Name() != file.Name() {
				t.Fatalf("Unexpected input to UploadObject(%v, %v, %v)", bucket, key, f.Name())
			}

			return "", nil
		}

		put := NewPut(&s3, &con, args)
		if err := put.Execute(&out); err != nil {
			t.Fatal(err)
		}
	}

	// Negative: S3 error
	{
		var s3 mockS3Client
		var con context.Context
		con.UpdatePath("bucket/folder")
		var out mockOutputter
		file, _ := ioutil.TempFile("", "")
		defer os.Remove(file.Name())
		args := []string{file.Name()}
		mockErr := errors.New("Mock Err")

		s3.uploadObjectCallback = func(bucket, key string, f *os.File) (string, error) {
			return "", mockErr
		}

		put := NewPut(&s3, &con, args)
		if err := put.Execute(&out); err != mockErr {
			t.Fatalf("Expected mock error to be returned: %v", err)
		}
	}

	// Negative: No target
	{
		var s3 mockS3Client
		var con context.Context
		con.UpdatePath("bucket")
		var out mockOutputter
		args := []string{}

		put := NewPut(&s3, &con, args)
		if err := put.Execute(&out); err == nil {
			t.Fatal("Expected error for no target")
		}
	}

	// Negative: Invalid file
	{
		var s3 mockS3Client
		var con context.Context
		con.UpdatePath("bucket")
		var out mockOutputter
		args := []string{"/notarealfile.txt"}

		put := NewPut(&s3, &con, args)
		if err := put.Execute(&out); err == nil {
			t.Fatal("Expected error for invalid file")
		}
	}
}

func TestPutCommand_IsLongRunning(t *testing.T) {
	put := PutCommand{}

	if !put.IsLongRunning() {
		t.Fatalf("Expected put to be long running: %v", put.IsLongRunning())
	}
}

func TestNewPut(t *testing.T) {
	var s3 mockS3Client
	var con context.Context
	args := []string{"file"}

	put := NewPut(&s3, &con, args)
	if put.s3 != &s3 {
		t.Fatalf("Unexpected S3 client stored on put command: %v", put.s3)
	} else if put.con != &con {
		t.Fatalf("Unexpected Context stored on put command: %v", put.con)
	} else if put.args[0] != args[0] {
		t.Fatalf("Unexpected args stored on put command: %v", put.args)
	}
}
