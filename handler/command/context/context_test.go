package context

import (
	"testing"
)

func TestContext_UpdatePath(t *testing.T) {
	// Helper function to check the length of a path.
	validatePathLength := func(p []string, l int) {
		if len(p) != l {
			t.Fatalf("Unexpected path length: {Actual: %v, Expected: %v}, %v", len(p), l, p)
		}
	}

	var c Context

	// Most simple use case
	c.UpdatePath("bucket")
	validatePathLength(c.path, 1)
	if c.Bucket() != "bucket" {
		t.Fatalf("Unexpected Path: %v", c.path)
	}

	// Go back in the path
	c.UpdatePath("../")
	validatePathLength(c.path, 0)

	// Two elements in the path
	c.UpdatePath("bucket/key")
	validatePathLength(c.path, 2)
	if c.Bucket() != "bucket" || c.PathWithoutBucket() != "key" {
		t.Fatalf("Unexpected path: %v", c.path)
	}

	// Just the path delimiter (ie. 'cd /')
	c.UpdatePath(PathDelimiter)
	validatePathLength(c.path, 0)

	// Empty string should be ignored
	c.UpdatePath("bucket")
	c.UpdatePath("")
	validatePathLength(c.path, 1)
	if c.Bucket() != "bucket" || len(c.PathWithoutBucket()) > 0 {
		t.Fatalf("Unexpected path after cd-ing with empty string: %v", c.Path())
	}
}

func TestContext_CalculatePath(t *testing.T) {
	var c Context
	var p []string

	// Most simple use case
	p = c.CalculatePath("bucket")
	if len(p) != 1 || p[0] != "bucket" {
		t.Fatalf("Unexpected Path: %v", p)
	}
}

func TestContext_IsRoot(t *testing.T) {
	var c Context

	if !c.IsRoot() {
		t.Fatalf("Expected context to be at root in it's zero-value: %v", c)
	}

	c.UpdatePath("bucket")
	if c.IsRoot() {
		t.Fatalf("Expected context NOT to be at root after cd-ing to a bucket: %v", c)
	}

	c.UpdatePath("../")
	if !c.IsRoot() {
		t.Fatalf("Expected context to be at root after cd-ing ../: %v", c)
	}

	c.UpdatePath("bucket/folder")
	if c.IsRoot() {
		t.Fatalf("Expected context NOT to be at root after cd-ing to a bucket+folder: %v", c)
	}

	c.UpdatePath("/")
	if !c.IsRoot() {
		t.Fatalf("Expected context to be at root after cd-ing /: %v", c)
	}
}

func TestContext_Path(t *testing.T) {
	var c Context

	if len(c.Path()) > 0 {
		t.Fatalf("Expected Path() len to be zero in the context's zero value: %v", c.Path())
	}

	c.UpdatePath("bucket")
	if c.Path() != "bucket" {
		t.Fatalf("Unexpected Path() after moving into a bucket: %v", c.Path())
	}

	c.UpdatePath("folder/subfolder")
	if c.Path() != "bucket/folder/subfolder" {
		t.Fatalf("Unexpected Path() after moving into subdirectories: %v", c.Path())
	}
}

func TestContext_PathWithoutBucket(t *testing.T) {
	var c Context

	if len(c.PathWithoutBucket()) > 0 {
		t.Fatalf("Expected PathWithoutBucket() len to be zero in the context's zero value: %v", c.PathWithoutBucket())
	}

	c.UpdatePath("bucket")
	if len(c.PathWithoutBucket()) > 0 {
		t.Fatalf("Expected PathWithoutBucket() len to be zero after moving into a bucket: %v", c.PathWithoutBucket())
	}

	c.UpdatePath("folder/subfolder")
	if c.PathWithoutBucket() != "folder/subfolder" {
		t.Fatalf("Unexpected PathWithoutBucket() after moving into subdirectories: %v", c.PathWithoutBucket())
	}
}

func TestContext_Bucket(t *testing.T) {
	var c Context

	if len(c.Bucket()) > 0 {
		t.Fatalf("Expected Bucket() len to be zero in the context's zero value: %v", c.Bucket())
	}

	c.UpdatePath("bucket")
	if c.Bucket() != "bucket" {
		t.Fatalf("Unexpected Bucket() after moving into a bucket: %v", c.Bucket())
	}

	c.UpdatePath("folder/subfolder")
	if c.Bucket() != "bucket" {
		t.Fatalf("Unexpected Bucket() after moving into a subdirectories: %v", c.Bucket())
	}
}
