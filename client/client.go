// Package client provides AWS API wrapper functionality.
package client

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

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
	// Initialize the S3 request.
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

// ObjectExists returns a bool indicating if the specified object exists.
func (c Client) ObjectExists(bucket, key string) (bool, error) {
	// Perform a HEAD request to determine if the bucket exists.
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

// PathExists returns a bool indicating if the specified path exists in a given bucket.
func (c Client) PathExists(bucket, path string) (bool, error) {
	// Note: Using a HEAD request (HeadObject) won't work because S3 has no real concept of folders.
	// Since we simulate folders, we instead do a 'ListObjects' to validate that there are objects that start with
	// the path.

	// Initialize the S3 request.
	input := s3.ListObjectsInput{
		Bucket:  &bucket,
		Prefix:  &path,
		MaxKeys: aws.Int64(1),
	}

	// Get the object list from AWS.
	resp, err := c.s3.ListObjects(&input)
	if err != nil {
		return false, err
	}

	return len(resp.Contents) > 0, nil
}

// DownloadObject downloads the specified object from Amazon S3 and returns the name of a local temporary file
// containing the downloaded object.
//
// Note: It is the responsibility of the caller to clean up the temporary file as necessary.
func (c Client) DownloadObject(bucket, key string) (string, error) {
	// Construct the request.
	input := s3.GetObjectInput{
		Bucket: &bucket,
		Key:    &key,
	}

	// Perform the API request to get the object.
	output, err := c.s3.GetObject(&input)
	if err != nil {
		return "", err
	}
	defer output.Body.Close()

	// Create the temporary file.
	tmp, err := ioutil.TempFile(os.TempDir(), "")
	if err != nil {
		return "", err
	}
	defer tmp.Close()

	// Perform the write.
	if _, err := io.Copy(tmp, output.Body); err != nil {
		return "", err
	}

	return tmp.Name(), nil
}

// UploadObject uploads a file to the specified key in an Amazon S3 bucket.
//
// Note: If the key provided is a directory, the file will be stored in the directory with the
// same name as the original file. If the key exists, it will be overwritten.
func (c Client) UploadObject(bucket, key string, file *os.File) (string, error) {
	// Sanitize the key input if it's empty or is a directory.
	if len(key) == 0 {
		key = filepath.Base(file.Name())
	} else {
		// Determine if the key is a directory, and if so, append the file name.
		path := key
		if path != "/" {
			path = path + "/"
		}

		if isDir, err := c.PathExists(bucket, path); err != nil {
			return "", err
		} else if isDir {
			key = path + filepath.Base(file.Name())
		}
	}

	// Perform the upload.
	input := s3.PutObjectInput{
		Bucket: &bucket,
		Key:    &key,
		Body:   file,
	}
	if _, err := c.s3.PutObject(&input); err != nil {
		return "", err
	}

	return key, nil
}

// New returns an initialized Client.
func New(region string) Client {
	return Client{
		s3: s3.New(session.New(), &aws.Config{
			Region: aws.String(region),
		}),
	}
}
