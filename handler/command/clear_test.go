package command

import (
	"runtime"
	"testing"
)

// Note: Cannot test ClearCommand Execute or it will clear the testing output

// func TestClearCommandExecute(t *testing.T) {
// 	var out mockOutputter

// 	clr := NewClear()
// 	clr.Execute(&out)

// 	// Nothing panicked, can't test anything else.
// }

func TestClearCommandCmdForSys(t *testing.T) {
	clr := NewClear()

	// Windows
	if c := clr.cmdForSys("windows"); c != "cls" {
		t.Fatalf("Unexpected clear command for windows: %v", c)
	}

	// Everything else
	archs := []string{"linux", "darwin", string(runtime.GOOS)}
	for _, arch := range archs {
		if c := clr.cmdForSys(arch); c != "clear" {
			t.Fatalf("Unexpected clear command for %v: %v", arch, c)
		}
	}
}

func TestClearCommandIsLongRunning(t *testing.T) {
	clr := NewClear()

	if clr.IsLongRunning() {
		t.Fatalf("Expected ClearCommand not to be long running")
	}
}

func TestNewClear(t *testing.T) {
	clr := NewClear()

	if clr != (ClearCommand{}) {
		t.Fatalf("Unexpected NewClear return: %v", clr)
	}
}
