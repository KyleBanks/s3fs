package context

import "testing"

func TestUpdatePath(t *testing.T) {
	c := Context{}

	// Most simple use case
	c.UpdatePath("bucket")
	validatePathLength(t, c.path, 1)
	if c.path[0] != "bucket" {
		t.Fatalf("Unexpected Path: %v", c.path)
	}

	// Go back in the path
	c.UpdatePath("../")
	validatePathLength(t, c.path, 0)

	// Two elements in the path
	c.UpdatePath("bucket/key")
	validatePathLength(t, c.path, 2)
	if c.path[0] != "bucket" || c.path[1] != "key" {
		t.Fatalf("Unexpected path: %v", c.path)
	}

	// Just the path delimiter (ie. 'cd /')
	c.UpdatePath(PathDelimiter)
	validatePathLength(t, c.path, 0)
}

func validatePathLength(t *testing.T, p []string, l int) {
	if len(p) != l {
		t.Fatalf("Unexpected path length: {Actual: %v, Expected: %v}, %v", len(p), l, p)
	}
}

func TestIsRoot(t *testing.T) {
	c := Context{}

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
}
