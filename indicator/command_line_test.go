package indicator

import (
	"testing"
	"time"
)

func TestCommandLineIndicator_ShowLoader(t *testing.T) {
	var out mockStringWriter

	ind := NewCommandLine(&out)
	ind.ShowLoader()

	// Allow time for the loading prompt to display at least once.
	// Note: We may get more than one loading indicator displayed while sleeping, which is okay.
	time.Sleep(loaderSleepTime * 2)

	if len(out.output) == 0 || out.output[0] != loaderText {
		t.Fatalf("Unexpected loader output: %v", out.output)
	}
}

func TestCommandLineIndicator_HideLoader(t *testing.T) {
	var out mockStringWriter

	ind := NewCommandLine(&out)

	ind.ShowLoader()
	ind.HideLoader()

	if len(out.output) == 0 || out.output[len(out.output)-1] != "\n" {
		t.Fatal("Expected CommandLineIndicator to write a newline when stopped loading")
	}
}

func TestCommandLineIndicator_ShowPrompt(t *testing.T) {
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
}
