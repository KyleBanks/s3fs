package output

import (
	"testing"
)

// Mock Writer

type mockWriter struct {
	written []string
}

func (m *mockWriter) Write(b []byte) (int, error) {
	m.written = append(m.written, string(b))

	return len(b), nil
}

func TestOutputter_Write(t *testing.T) {
	var w mockWriter

	out := New(&w)
	out.Write("Hello Output")

	if len(w.written) != 1 || w.written[0] != "Hello Output" {
		t.Fatalf("Unexpected output written: %v", w.written)
	}
}

func TestNew(t *testing.T) {
	var w mockWriter

	out := New(&w)
	if out.w != &w {
		t.Fatalf("OutputWriter storing unexpected Writer: %v", out.w)
	}
}
