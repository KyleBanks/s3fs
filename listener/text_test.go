package listener

import (
	"testing"
)

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
