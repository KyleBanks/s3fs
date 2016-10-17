package listener

import (
	"testing"
)

func TestTextListenerListen(t *testing.T) {
	// Positive case, scan successful
	{
		var ui mockIndicator
		var input mockInputter
		var text TextListener

		text = NewText(&ui, &input)
		input.scanCallback = func() bool {
			if !ui.promptShown {
				t.Fatalf("Expected TextListener to show prompt before scanning")
			}

			return true
		}
		input.textCallback = func() string {
			return "one two"
		}

		res, ok := text.Listen()
		if !ok {
			t.Fatal("Expected positive response from Listen() when Scan returns true")
		}

		if len(res) != 2 || res[0] != "one" || res[1] != "two" {
			t.Fatalf("Unexpected response from Listen(): %v", res)
		}
	}

	// Negative case, scan failed
	{
		var ui mockIndicator
		var input mockInputter
		var text TextListener

		text = NewText(&ui, &input)
		input.scanCallback = func() bool {
			if !ui.promptShown {
				t.Fatalf("Expected TextListener to show prompt before scanning")
			}

			return false
		}
		input.textCallback = func() string {
			t.Fatalf("Expected Text() not to be called when Scan() returns false")
			return ""
		}

		res, ok := text.Listen()
		if ok {
			t.Fatal("Expected negative response from Listen() when Scan returns false")
		}

		if res != nil {
			t.Fatalf("Unexpected response from Listen(): %v", res)
		}
	}

}

func TestNewText(t *testing.T) {
	var ui indicator
	var input inputter

	l := NewText(ui, input)
	if l.ui != ui {
		t.Fatalf("Unexpected UI indicator stored on TextListener: %v", l.ui)
	} else if l.input != input {
		t.Fatalf("Unexpected inputter stored on TextListener: %v", l.input)
	}
}
