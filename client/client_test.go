package client

import (
	"errors"
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
			t.Fatal("Unexpected error returned: %v", err)
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
			if *i.Bucket != bucket || *i.Prefix != prefix {
				t.Fatalf("Unexpected ListObjectsInput: %v", i)
			}

			return nil, mockErr
		}

		c := Client{&mockS3}

		if _, err := c.LsObjects(bucket, prefix); err != mockErr {
			t.Fatal("Unexpected error returned: %v", err)
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
			if *i.Bucket != bucket {
				t.Fatalf("Unexpected HeadBucketInput: %v", i)
			}

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
			if *i.Bucket != bucket || *i.Key != key {
				t.Fatalf("Unexpected HeadObjectInput: %v", i)
			}

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

func TestNew(t *testing.T) {
	c := New("region")

	if c.s3 == nil {
		t.Fatal("Expected client to be initialized with an s3communicator")
	}
}
