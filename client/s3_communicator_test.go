package client

import (
	"github.com/aws/aws-sdk-go/service/s3"
)

// Mock s3Communicator

type mockS3Communicator struct {
	listBucketsCallback func(i *s3.ListBucketsInput) (*s3.ListBucketsOutput, error)
	listObjectsCallback func(i *s3.ListObjectsInput) (*s3.ListObjectsOutput, error)

	headBucketCallback func(i *s3.HeadBucketInput) (*s3.HeadBucketOutput, error)
	headObjectCallback func(i *s3.HeadObjectInput) (*s3.HeadObjectOutput, error)
}

func (m *mockS3Communicator) ListBuckets(i *s3.ListBucketsInput) (*s3.ListBucketsOutput, error) {
	return m.listBucketsCallback(i)
}

func (m *mockS3Communicator) ListObjects(i *s3.ListObjectsInput) (*s3.ListObjectsOutput, error) {
	return m.listObjectsCallback(i)
}

func (m *mockS3Communicator) HeadBucket(i *s3.HeadBucketInput) (*s3.HeadBucketOutput, error) {
	return m.headBucketCallback(i)
}

func (m *mockS3Communicator) HeadObject(i *s3.HeadObjectInput) (*s3.HeadObjectOutput, error) {
	return m.headObjectCallback(i)
}
