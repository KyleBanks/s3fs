package command

import (
	"errors"
	"strings"
	"testing"

	"github.com/KyleBanks/s3fs/handler/command/context"
)

func TestLsCommand_Execute(t *testing.T) {
	// Root bucket list
	{
		var s3 mockS3Client
		var con context.Context
		var out mockOutputter

		// Define sample buckets.
		samples := []string{"bucket1", "bucket2"}

		// Override ls functions.
		s3.lsBucketsCallback = func() ([]string, error) {
			return samples, nil
		}
		s3.lsObjectsCallback = func(a, b string) ([]string, error) {
			t.Fatalf("LsObjects should not be called when context is at root")
			return nil, nil
		}

		// Execute the command.
		ls := NewLs(&s3, &con)
		if err := ls.Execute(&out); err != nil {
			t.Fatal(err)
		}

		// Validate output length.
		if len(out.output) != len(samples) {
			t.Fatalf("Unexpected output length for Bucket LS: %v", len(out.output))
		}

		// Validate the actual output.
		for i, sample := range samples {
			if !strings.Contains(out.output[i], sample) {
				t.Fatalf("Failed to find output [%v] in response: %v", sample, out.output)
			}
		}
	}

	// Bucket error
	{
		var s3 mockS3Client
		var con context.Context
		var out mockOutputter

		// Override ls commands to mock an error case.
		mockErr := errors.New("Mock Error")

		s3.lsBucketsCallback = func() ([]string, error) {
			return nil, mockErr
		}
		s3.lsObjectsCallback = func(a, b string) ([]string, error) {
			t.Fatalf("LsObjects should not be called when context is at root")
			return nil, nil
		}

		// Execute the command and validate the error is bubbled up.
		ls := NewLs(&s3, &con)
		if err := ls.Execute(&out); err != mockErr {
			t.Fatalf("Expected error to be passed up the stack: %v", err)
		}

		// Ensure no output was written.
		if len(out.output) != 0 {
			t.Fatalf("Error case should not produce output for Bucket LS: %v", len(out.output))
		}
	}

	// Object list
	{
		pathSets := []string{"bucket", "bucket/folder"}

		// For each sample path, test the LS command
		for _, path := range pathSets {
			var s3 mockS3Client
			var con context.Context
			var out mockOutputter

			// Update the path.
			con.UpdatePath(path)

			// Determine the bucket and prefix for the current test.
			sampleBucket := strings.Split(path, context.PathDelimiter)[0]
			samplePrefix := strings.Replace(path, sampleBucket, "", 1)

			// Sanitize the prefix.
			if len(samplePrefix) > 0 {
				// Remove starting /
				if string(samplePrefix[0]) == context.PathDelimiter {
					samplePrefix = string(samplePrefix[1:len(samplePrefix)])
				}

				// Append trailing /
				samplePrefix = samplePrefix + context.PathDelimiter
			}

			// Define the sample output.
			samples := []string{samplePrefix + "subfolder/", samplePrefix + "index.html"}

			// Override ls functions.
			s3.lsBucketsCallback = func() ([]string, error) {
				t.Fatalf("LsBuckets should not be called when context is not at root")
				return nil, nil
			}
			s3.lsObjectsCallback = func(bucket, prefix string) ([]string, error) {
				if bucket != sampleBucket || prefix != samplePrefix {
					t.Fatalf("Unexpected bucket/prefix provided to LsObjects: {Bucket: %v, ExpectedBucket: %v, Prefix: %v, ExpectedPrefix: %v}", bucket, sampleBucket, prefix, samplePrefix)
				}

				return samples, nil
			}

			// Execute the command.
			ls := NewLs(&s3, &con)
			if err := ls.Execute(&out); err != nil {
				t.Fatal(err)
			}

			// Validate output length.
			if len(out.output) != len(samples) {
				t.Fatalf("Unexpected output length for Object LS: %v", out.output)
			}

			// Validate the actual output.
			for i, sample := range samples {
				if !strings.Contains(out.output[i], strings.Replace(sample, samplePrefix, "", 1)) {
					t.Fatalf("Failed to find output: %v", sample)
				}
			}
		}
	}

	// Object list error
	{
		var s3 mockS3Client
		var con context.Context
		var out mockOutputter

		con.UpdatePath("bucket/prefix")
		mockErr := errors.New("Mock Error")

		s3.lsBucketsCallback = func() ([]string, error) {
			t.Fatalf("LsBuckets should not be called when context is not at root")
			return nil, nil
		}
		s3.lsObjectsCallback = func(bucket, prefix string) ([]string, error) {
			return nil, mockErr
		}

		ls := NewLs(&s3, &con)
		if err := ls.Execute(&out); err != mockErr {
			t.Fatalf("Expected error to be passed up the stack: %v", err)
		}

		if len(out.output) != 0 {
			t.Fatalf("Error case should not produce output for Bucket LS: %v", len(out.output))
		}
	}
}

func TestLsCommand_prefixOutput(t *testing.T) {
	var s3 mockS3Client
	var con context.Context
	ls := NewLs(&s3, &con)

	// Bucket
	{
		bucket := "bucket"

		if out := ls.prefixOutput(bucket, true); out != bucketPrefix+" "+bucket {
			t.Fatalf("Unexpected output for bucket: '%v'", out)
		}
	}

	// Folder
	{
		folder := "folder" + context.PathDelimiter

		if out := ls.prefixOutput(folder, false); out != folderPrefix+" "+folder {
			t.Fatalf("Unexpected output for folder: '%v'", out)
		}
	}

	// File
	{
		file := "file.txt"

		if out := ls.prefixOutput(file, false); out != filePrefix+" "+file {
			t.Fatalf("Unexpected output for file: '%v'", out)
		}
	}
}

func TestLsCommand_IsLongRunning(t *testing.T) {
	var s3 mockS3Client
	var con context.Context

	ls := NewLs(&s3, &con)
	if !ls.IsLongRunning() {
		t.Fatalf("Expected LsCommand to always be long running")
	}
}

func TestNewLs(t *testing.T) {
	var s3 mockS3Client
	var con context.Context

	ls := NewLs(&s3, &con)
	if ls.s3 != &s3 {
		t.Fatalf("Unexpected S3 client stored on ls command: %v", ls.s3)
	} else if ls.con != &con {
		t.Fatalf("Unexpected Context stored on ls command: %v", ls.con)
	}
}
