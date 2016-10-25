package client

import (
	"github.com/aws/aws-sdk-go/service/s3"
)

// s3Communicator defines an interface that performs actual AWS API operations.
type s3Communicator interface {
	ListBuckets(*s3.ListBucketsInput) (*s3.ListBucketsOutput, error)
	ListObjects(*s3.ListObjectsInput) (*s3.ListObjectsOutput, error)

	HeadBucket(*s3.HeadBucketInput) (*s3.HeadBucketOutput, error)
	HeadObject(*s3.HeadObjectInput) (*s3.HeadObjectOutput, error)

	GetObject(*s3.GetObjectInput) (*s3.GetObjectOutput, error)
	PutObject(*s3.PutObjectInput) (*s3.PutObjectOutput, error)
}
