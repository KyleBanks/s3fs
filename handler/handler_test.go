package handler

// Mock indicator

type mockIndicator struct {
	showLoaderCalled bool
	hideLoaderCalled bool
}

func (m *mockIndicator) ShowLoader() {
	m.showLoaderCalled = true
}

func (m *mockIndicator) HideLoader() {
	m.hideLoaderCalled = true
}

// Mock command.Outputter

type mockOutputter struct {
	out []string
}

func (m *mockOutputter) Write(s string) {
	m.out = append(m.out, s)
}

// Mock command.S3Client
// TODO: This is duplicated in command_test.go, would be nice to find a way to share it.

type mockS3Client struct {
	lsBucketsCallback func() ([]string, error)
	lsObjectsCallback func(string, string) ([]string, error)

	bucketExistsCallback func(string) (bool, error)
	objectExistsCallback func(string, string) (bool, error)
	pathExistsCallback   func(string, string) (bool, error)

	downloadObjectCallback func(string, string) (string, error)
}

func (m mockS3Client) LsBuckets() ([]string, error) {
	return m.lsBucketsCallback()
}

func (m mockS3Client) LsObjects(bucket, prefix string) ([]string, error) {
	return m.lsObjectsCallback(bucket, prefix)
}

func (m mockS3Client) BucketExists(bucket string) (bool, error) {
	return m.bucketExistsCallback(bucket)
}

func (m mockS3Client) ObjectExists(bucket, key string) (bool, error) {
	return m.objectExistsCallback(bucket, key)
}

func (m mockS3Client) PathExists(bucket, path string) (bool, error) {
	return m.pathExistsCallback(bucket, path)
}

func (m mockS3Client) DownloadObject(bucket, key string) (string, error) {
	return m.downloadObjectCallback(bucket, key)
}
