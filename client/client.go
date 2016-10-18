// Package client provides AWS API wrapper functionality.
package client

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// Client defines a wrapper for the Amazon S3 API.
type Client struct {
	s3 s3Communicator
}

// LsBuckets performs a request to retrieve all buckets, and returns their names.
func (c Client) LsBuckets() ([]string, error) {
	// Get the bucket list from AWS.
	resp, err := c.s3.ListBuckets(&s3.ListBucketsInput{})
	if err != nil {
		return nil, err
	}

	// Create a slice of bucket names to return.
	buckets := make([]string, len(resp.Buckets))
	for i, b := range resp.Buckets {
		buckets[i] = *b.Name
	}

	return buckets, nil
}

// LsObjects performs a request to retrieve all objects, and returns their keys.
func (c Client) LsObjects(bucket, prefix string) ([]string, error) {
	// Use the current context to initialize the S3 request.
	input := s3.ListObjectsInput{
		Bucket: &bucket,
		Prefix: &prefix,
	}

	// Get the object list from AWS.
	resp, err := c.s3.ListObjects(&input)
	if err != nil {
		return nil, err
	}

	// Create a slice of object keys to return.
	objects := make([]string, len(resp.Contents))
	for i, o := range resp.Contents {
		objects[i] = *o.Key
	}

	return objects, nil
}

// BucketExists returns a bool indicating if the specified bucket exists.
func (c Client) BucketExists(bucket string) (bool, error) {
	// Perform a HEAD request to determine if the bucket exists.
	input := s3.HeadBucketInput{
		Bucket: &bucket,
	}

	// Assume an error means that the bucket doesn't exist.
	// TODO: Not a great assumption, check the actual error.
	if _, err := c.s3.HeadBucket(&input); err != nil {
		return false, nil
	}

	return true, nil
}

// ObjectExists returns a bool indicating if the specified object exists in a given bucket.
func (c Client) ObjectExists(bucket, key string) (bool, error) {
	// Perform a HEAD request to determine if the object exists.
	input := s3.HeadObjectInput{
		Bucket: &bucket,
		Key:    &key,
	}

	// Assume an error means that the object doesn't exist.
	// TODO: Not a great assumption, check the actual error.
	if _, err := c.s3.HeadObject(&input); err != nil {
		return false, nil
	}

	return true, nil
}

// New returns an initialized Client.
func New(region string) Client {
	return Client{
		s3: s3.New(session.New(), &aws.Config{
			Region: aws.String(region),
		}),
	}
}
