// Package client provides AWS API wrapper functionality.
package client

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// Client defines a wrapper for the Amazon S3 API.
type Client struct {
	s3 *s3.S3
}

// LsBuckets performs a request to retrieve all buckets, and returns their names.
func (c Client) LsBuckets() ([]string, error) {
	// Get the bucket list from AWS.
	resp, err := c.s3.ListBuckets(&s3.ListBucketsInput{})
	if err != nil {
		return nil, err
	}

	// Create a slice of bucket names to return.
	buckets := make([]string, 0, len(resp.Buckets))
	for _, b := range resp.Buckets {
		buckets = append(buckets, *b.Name)
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
	objects := make([]string, 0, len(resp.Contents))
	for _, o := range resp.Contents {
		objects = append(objects, *o.Key)
	}

	return objects, nil
}

// New returns an initialized Client.
func New(region string) Client {
	return Client{
		s3: s3.New(session.New(), &aws.Config{
			Region: aws.String(region),
		}),
	}
}
