package command

import (
	"testing"
)

// Note: Currently cannot test the exit Execute() function due to it's nature.

func TestExitCommand_IsLongRunning(t *testing.T) {
	e := NewExit()

	if e.IsLongRunning() {
		t.Fatalf("Expected ExitCommand to not be long running")
	}
}

func TestNewExit(t *testing.T) {
	e := NewExit()

	if e != (ExitCommand{}) {
		t.Fatalf("Unexpected NewExit return: %v", e)
	}
}
