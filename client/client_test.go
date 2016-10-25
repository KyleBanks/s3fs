package client

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

func TestClient_LsBuckets(t *testing.T) {
	// Positive case
	{
		sample := s3.ListBucketsOutput{
			Buckets: []*s3.Bucket{
				{Name: aws.String("test")},
			},
		}

		var mockS3 mockS3Communicator
		mockS3.listBucketsCallback = func(i *s3.ListBucketsInput) (*s3.ListBucketsOutput, error) {
			return &sample, nil
		}

		c := Client{&mockS3}

		buckets, err := c.LsBuckets()
		if err != nil {
			t.Fatal(err)
		} else if len(buckets) != len(sample.Buckets) {
			t.Fatalf("Unexpected number of buckets returned: %v", buckets)
		}

		for i, bucket := range buckets {
			if bucket != *sample.Buckets[i].Name {
				t.Fatalf("Unexpected response from LsBuckets: {Expected: %v, Actual: %v}", *sample.Buckets[i].Name, bucket)
			}
		}
	}

	// Negative case
	{
		mockErr := errors.New("Mock Error")

		var mockS3 mockS3Communicator
		mockS3.listBucketsCallback = func(i *s3.ListBucketsInput) (*s3.ListBucketsOutput, error) {
			return nil, mockErr
		}

		c := Client{&mockS3}

		if _, err := c.LsBuckets(); err != mockErr {
			t.Fatalf("Unexpected error returned: %v", err)
		}
	}
}

func TestClient_LsObjects(t *testing.T) {
	// Positive case
	{
		bucket := "bucket"
		prefix := "prefix"
		sample := s3.ListObjectsOutput{
			Contents: []*s3.Object{
				{Key: aws.String("test")},
			},
		}

		var mockS3 mockS3Communicator
		mockS3.listObjectsCallback = func(i *s3.ListObjectsInput) (*s3.ListObjectsOutput, error) {
			if *i.Bucket != bucket || *i.Prefix != prefix {
				t.Fatalf("Unexpected ListObjectsInput: %v", i)
			}

			return &sample, nil
		}

		c := Client{&mockS3}

		objects, err := c.LsObjects(bucket, prefix)
		if err != nil {
			t.Fatal(err)
		} else if len(objects) != len(sample.Contents) {
			t.Fatalf("Unexpected number of objects returned: %v", objects)
		}

		for i, obj := range objects {
			if obj != *sample.Contents[i].Key {
				t.Fatalf("Unexpected response from LsObjects: {Expected: %v, Actual: %v}", *sample.Contents[i].Key, obj)
			}
		}
	}

	// Negative case
	{
		bucket := "bucket"
		prefix := "prefix"
		mockErr := errors.New("Mock Error")

		var mockS3 mockS3Communicator
		mockS3.listObjectsCallback = func(i *s3.ListObjectsInput) (*s3.ListObjectsOutput, error) {
			return nil, mockErr
		}

		c := Client{&mockS3}

		if _, err := c.LsObjects(bucket, prefix); err != mockErr {
			t.Fatalf("Unexpected error returned: %v", err)
		}
	}
}

func TestClient_BucketExists(t *testing.T) {
	// Positive case
	{
		bucket := "bucket"

		var mockS3 mockS3Communicator
		mockS3.headBucketCallback = func(i *s3.HeadBucketInput) (*s3.HeadBucketOutput, error) {
			if *i.Bucket != bucket {
				t.Fatalf("Unexpected HeadBucketInput: %v", i)
			}

			return nil, nil
		}

		c := Client{&mockS3}
		if exists, err := c.BucketExists(bucket); err != nil {
			t.Fatal(err)
		} else if !exists {
			t.Fatal("Exists should be true in positive case")
		}
	}

	// Negative case
	{
		bucket := "bucket"

		var mockS3 mockS3Communicator
		mockS3.headBucketCallback = func(i *s3.HeadBucketInput) (*s3.HeadBucketOutput, error) {
			return nil, errors.New("Fake error")
		}

		c := Client{&mockS3}
		if exists, err := c.BucketExists(bucket); err != nil {
			t.Fatal(err)
		} else if exists {
			t.Fatal("Exists should be false in negative case")
		}
	}
}

func TestClient_ObjectExists(t *testing.T) {
	// Positive case
	{
		bucket := "bucket"
		key := "key"

		var mockS3 mockS3Communicator
		mockS3.headObjectCallback = func(i *s3.HeadObjectInput) (*s3.HeadObjectOutput, error) {
			if *i.Bucket != bucket || *i.Key != key {
				t.Fatalf("Unexpected HeadObjectInput: %v", i)
			}

			return nil, nil
		}

		c := Client{&mockS3}
		if exists, err := c.ObjectExists(bucket, key); err != nil {
			t.Fatal(err)
		} else if !exists {
			t.Fatal("Exists should be true in positive case")
		}
	}

	// Negative case
	{
		bucket := "bucket"
		key := "key"

		var mockS3 mockS3Communicator
		mockS3.headObjectCallback = func(i *s3.HeadObjectInput) (*s3.HeadObjectOutput, error) {
			return nil, errors.New("Fake Error")
		}

		c := Client{&mockS3}
		if exists, err := c.ObjectExists(bucket, key); err != nil {
			t.Fatal(err)
		} else if exists {
			t.Fatal("Exists should be false in negative case")
		}
	}
}

func TestClient_PathExists(t *testing.T) {
	// Positive case
	{
		bucket := "bucket"
		key := "key"

		var mockS3 mockS3Communicator
		mockS3.listObjectsCallback = func(i *s3.ListObjectsInput) (*s3.ListObjectsOutput, error) {
			if *i.Bucket != bucket || *i.Prefix != key || *i.MaxKeys != int64(1) {
				t.Fatalf("Unexpected ListObjectsInput: %v", i)
			}

			return &s3.ListObjectsOutput{
				Contents: make([]*s3.Object, 1),
			}, nil
		}

		c := Client{&mockS3}
		if exists, err := c.PathExists(bucket, key); err != nil {
			t.Fatal(err)
		} else if !exists {
			t.Fatal("Exists should be true in positive case")
		}
	}

	// Negative case
	{
		bucket := "bucket"
		key := "key"
		mockErr := errors.New("Mock err")

		var mockS3 mockS3Communicator
		mockS3.listObjectsCallback = func(i *s3.ListObjectsInput) (*s3.ListObjectsOutput, error) {
			return nil, mockErr
		}

		c := Client{&mockS3}
		if exists, err := c.PathExists(bucket, key); err != mockErr {
			t.Fatalf("Unexpected error returned: %v", err)
		} else if exists {
			t.Fatal("Exists should be false in negative case")
		}
	}
}

func TestClient_DownloadObject(t *testing.T) {
	// Positive Case
	{
		bucket := "bucket"
		key := "key"

		// Create a mock file to download.
		data := make([]byte, 10000)
		for i := 0; i < len(data); i++ {
			data[i] = byte(i)
		}

		// Override the getObjectCallback to return the mock file.
		var mockS3 mockS3Communicator
		mockS3.getObjectCallback = func(i *s3.GetObjectInput) (*s3.GetObjectOutput, error) {
			if *i.Bucket != bucket || *i.Key != key {
				t.Fatalf("Unexpected GetObjectInput: %v", i)
			}

			return &s3.GetObjectOutput{
				Body: &mockReadCloser{
					data: data,
				},
			}, nil
		}

		// Download the sample object.
		c := Client{&mockS3}
		f, err := c.DownloadObject(bucket, key)
		if err != nil {
			t.Fatal(err)
		}
		defer os.Remove(f)

		// Ensure the correct file contents were written.
		contents, err := ioutil.ReadFile(f)
		if err != nil {
			t.Fatal(err)
		}

		if len(contents) != len(data) {
			t.Fatalf("Downloaded file length mismatch: {Expected: %v, Actual: %v}", len(data), len(contents))
		}

		for i := 0; i < len(contents); i++ {
			if contents[i] != data[i] {
				t.Fatal("Downloaded file contents incorrect!")
			}
		}
	}

	// S3 Error
	{
		bucket := "bucket"
		key := "key"
		mockErr := errors.New("Mock Error")

		// Override the getObjectCallback to return the mock error.
		var mockS3 mockS3Communicator
		mockS3.getObjectCallback = func(i *s3.GetObjectInput) (*s3.GetObjectOutput, error) {
			return nil, mockErr
		}

		c := Client{&mockS3}
		if f, err := c.DownloadObject(bucket, key); err != mockErr {
			t.Fatalf("Expected mock error to be returned: %v", err)
		} else if len(f) > 0 {
			t.Fatalf("Unexpected file returned: %v", f)
		}
	}
}

func TestClient_UploadObject(t *testing.T) {
	// Positive case, not directory
	{
		file, _ := ioutil.TempFile("", "")
		defer os.Remove(file.Name())
		bucket := "bucket"
		key := "folder/file.txt"

		var mockS3 mockS3Communicator
		mockS3.putObjectCallback = func(i *s3.PutObjectInput) (*s3.PutObjectOutput, error) {
			if *i.Bucket != bucket || *i.Key != key {
				t.Fatalf("Unexpected PutObjectInput: %v", i)
			}

			return nil, nil
		}
		mockS3.listObjectsCallback = func(i *s3.ListObjectsInput) (*s3.ListObjectsOutput, error) {
			if *i.Bucket != bucket || *i.Prefix != key+"/" {
				t.Fatalf("Unexpected ListObjectsOutput: %v", i)
			}

			return &s3.ListObjectsOutput{
				Contents: make([]*s3.Object, 0),
			}, nil
		}

		c := Client{&mockS3}
		if path, err := c.UploadObject(bucket, key, file); err != nil {
			t.Fatal(err)
		} else if path != key {
			t.Fatalf("Unexpected path returned: %v", path)
		}
	}

	// Positive case, directory
	{
		file, _ := ioutil.TempFile("", "")
		defer os.Remove(file.Name())
		bucket := "bucket"
		key := "folder"
		expectedKey := key + "/" + filepath.Base(file.Name())

		var mockS3 mockS3Communicator
		mockS3.putObjectCallback = func(i *s3.PutObjectInput) (*s3.PutObjectOutput, error) {
			if *i.Bucket != bucket || *i.Key != expectedKey {
				t.Fatalf("Unexpected PutObjectInput: %v", i)
			}

			return nil, nil
		}
		mockS3.listObjectsCallback = func(i *s3.ListObjectsInput) (*s3.ListObjectsOutput, error) {
			return &s3.ListObjectsOutput{
				Contents: make([]*s3.Object, 1),
			}, nil
		}

		c := Client{&mockS3}
		if path, err := c.UploadObject(bucket, key, file); err != nil {
			t.Fatal(err)
		} else if path != expectedKey {
			t.Fatalf("Unexpected path returned: %v", path)
		}
	}

	// Negative case
	{
		file, _ := ioutil.TempFile("", "")
		defer os.Remove(file.Name())
		bucket := "bucket"
		key := "folder"
		mockErr := errors.New("Mock error")

		var mockS3 mockS3Communicator
		mockS3.putObjectCallback = func(i *s3.PutObjectInput) (*s3.PutObjectOutput, error) {
			return nil, mockErr
		}
		mockS3.listObjectsCallback = func(i *s3.ListObjectsInput) (*s3.ListObjectsOutput, error) {
			return &s3.ListObjectsOutput{
				Contents: make([]*s3.Object, 0),
			}, nil
		}

		c := Client{&mockS3}
		if _, err := c.UploadObject(bucket, key, file); err != mockErr {
			t.Fatalf("Expected mock error to be returned: %v", err)
		}
	}
}

func TestNew(t *testing.T) {
	c := New("region")

	if c.s3 == nil {
		t.Fatal("Expected client to be initialized with an s3communicator")
	}
}
