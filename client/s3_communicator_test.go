package client

import (
	"io"

	"github.com/aws/aws-sdk-go/service/s3"
)

// Mock s3Communicator

type mockS3Communicator struct {
	listBucketsCallback func(i *s3.ListBucketsInput) (*s3.ListBucketsOutput, error)
	listObjectsCallback func(i *s3.ListObjectsInput) (*s3.ListObjectsOutput, error)

	headBucketCallback func(i *s3.HeadBucketInput) (*s3.HeadBucketOutput, error)
	headObjectCallback func(i *s3.HeadObjectInput) (*s3.HeadObjectOutput, error)

	getObjectCallback func(i *s3.GetObjectInput) (*s3.GetObjectOutput, error)
	putObjectCallback func(i *s3.PutObjectInput) (*s3.PutObjectOutput, error)
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

func (m *mockS3Communicator) GetObject(i *s3.GetObjectInput) (*s3.GetObjectOutput, error) {
	return m.getObjectCallback(i)
}

func (m *mockS3Communicator) PutObject(i *s3.PutObjectInput) (*s3.PutObjectOutput, error) {
	return m.putObjectCallback(i)
}

// Mock ReadCloser

type mockReadCloser struct {
	index int
	data  []byte
}

func (m *mockReadCloser) Read(b []byte) (int, error) {
	var i int
	for i = 0; i < len(b); i++ {
		if i+m.index >= len(m.data) {
			return i, io.EOF
		}

		b[i] = m.data[i+m.index]
	}

	m.index += i
	return i, nil
}

func (m *mockReadCloser) Close() error {
	return nil
}
