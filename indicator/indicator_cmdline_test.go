package indicator

import (
	"testing"
	"time"
)

func TestShowLoader(t *testing.T) {
	var out mockStringWriter

	ind := NewCommandLine(&out)
	ind.ShowLoader()

	if !ind.loading {
		t.Fatal("Expected CommandLineIndicator to be loading after calling ShowLoader()")
	}

	// Allow time for the loading prompt to display at least once.
	// Note: We may get more than one loading indicator displayed while sleeping, which is okay.
	time.Sleep(loaderSleepTime * 2)

	if len(out.output) == 0 || out.output[0] != loaderText {
		t.Fatal("Unexpected loader output: %v", out.output)
	}
}

func TestHideLoader(t *testing.T) {
	var out mockStringWriter

	ind := NewCommandLine(&out)

	ind.ShowLoader()
	ind.HideLoader()
	if ind.loading {
		t.Fatal("Expected CommandLineIndicator to stop loading after calling HideLoader()")
	}

	if len(out.output) == 0 || out.output[len(out.output)-1] != "\n" {
		t.Fatal("Expected CommandLineIndicator to write a newline when stopped loading")
	}
}

func TestShowPrompt(t *testing.T) {
	var out mockStringWriter

	ind := NewCommandLine(&out)
	ind.ShowPrompt()

	if len(out.output) != 1 || out.output[0] != promptText {
		t.Fatalf("Unexpected prompt: %v", out.output)
	}
}

func TestNewCommandLine(t *testing.T) {
	var out mockStringWriter

	ind := NewCommandLine(&out)
	if ind.out != &out {
		t.Fatal("CommandLineIndicator storing unknown stringWriter")
	}

	if ind.loading {
		t.Fatal("Expected CommandLineIndicator not to be loading on initialization")
	}
}
