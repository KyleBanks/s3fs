package command

import (
	"testing"

	"github.com/KyleBanks/s3fs/handler/command/context"
)

func TestPwdCommand_Execute(t *testing.T) {
	var pwd PwdCommand
	var out mockOutputter
	var con context.Context

	validatePwd := func(expected string) {
		if err := pwd.Execute(&out); err != nil {
			t.Fatal(err)
		}

		if len(out.output) != 1 || out.output[0] != expected {
			t.Fatalf("Unexpected pwd output: {Actual: %v, Expected: %v}", out.output, expected)
		}
	}

	// No path
	out = mockOutputter{}
	con = context.Context{}
	pwd = NewPwd(&con)

	validatePwd(context.PathDelimiter + "\n")

	// With Bucket
	out = mockOutputter{}
	con = context.Context{}
	pwd = NewPwd(&con)

	con.UpdatePath("bucket")
	validatePwd("bucket" + context.PathDelimiter + "\n")

	// With bucket and directories
	out = mockOutputter{}
	con = context.Context{}
	pwd = NewPwd(&con)

	con.UpdatePath("bucket/folder/subfolder")
	validatePwd("bucket/folder/subfolder" + context.PathDelimiter + "\n")
}

func TestPwdCommand_IsLongRunning(t *testing.T) {
	pwd := NewPwd(&context.Context{})

	if pwd.IsLongRunning() {
		t.Fatalf("Expected pwd not to be long running: ", pwd.IsLongRunning())
	}
}

func TestNewPwd(t *testing.T) {
	var con context.Context

	pwd := NewPwd(&con)
	if pwd.con != &con {
		t.Fatal("pwd storing unexpected context")
	}
}
