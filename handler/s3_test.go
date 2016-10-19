package handler

import (
	"reflect"
	"strings"
	"testing"

	"github.com/KyleBanks/s3fs/handler/command"
)

func TestS3Handler_Handle(t *testing.T) {
	// Empty command
	{
		var ui mockIndicator
		s3 := NewS3(nil, &ui)

		if err := s3.Handle([]string{}, nil); err != nil {
			t.Fatal(err)
		}

		if ui.showLoaderCalled || ui.hideLoaderCalled {
			t.Fatalf("Loader methods should not have been called for blank cmd: %v", ui)
		}
	}

	// Invalid command
	{
		var ui mockIndicator
		s3 := NewS3(nil, &ui)

		if err := s3.Handle([]string{"fake"}, nil); err == nil {
			t.Fatal("Expected error for unknown command")
		} else if !strings.Contains(err.Error(), "fake") {
			t.Fatalf("Expected error to contain the command name: %v", err.Error())
		}

		if ui.showLoaderCalled || ui.hideLoaderCalled {
			t.Fatalf("Loader methods should not have been called for unknown cmd: %v", ui)
		}
	}

	// Valid command (short running)
	{
		var ui mockIndicator
		var out mockOutputter
		s3 := NewS3(nil, &ui)

		if err := s3.Handle([]string{command.CmdPwd}, &out); err != nil {
			t.Fatal(err)
		}

		if len(out.out) == 0 {
			t.Fatal("Expected output to be written to the Outputter interface")
		}

		if ui.showLoaderCalled || ui.hideLoaderCalled {
			t.Fatalf("Loader methods should not have been called for short running cmd: %v", ui)
		}
	}

	// Valid command (long running)
	{
		var ui mockIndicator
		var out mockOutputter
		var mockS3 mockS3Client

		mockS3.lsBucketsCallback = func() ([]string, error) {
			return []string{"bucket", "bucket2"}, nil
		}

		s3 := NewS3(&mockS3, &ui)

		if err := s3.Handle([]string{command.CmdLs}, &out); err != nil {
			t.Fatal(err)
		}

		if len(out.out) == 0 {
			t.Fatal("Expected output to be written to the Outputter interface")
		}

		if !ui.showLoaderCalled || !ui.hideLoaderCalled {
			t.Fatalf("Loader methods should have been called for long running cmd: %v", ui)
		}
	}
}

func TestS3Handler_commandFromArgs(t *testing.T) {
	// Known Commands
	{
		// Define the command types, and the expected Executor for each.
		cmds := []struct {
			name     string
			expected command.Executor
		}{
			{command.CmdLs, command.LsCommand{}},
			{command.CmdCd, command.CdCommand{}},
			{command.CmdGet, command.GetCommand{}},
			{command.CmdPwd, command.PwdCommand{}},
			{command.CmdClear, command.ClearCommand{}},
			{command.CmdExit, command.ExitCommand{}},
		}

		s3 := NewS3(nil, nil)

		// For each command, ensure the proper Executor is returned.
		for _, cmd := range cmds {
			c, err := s3.commandFromArgs([]string{cmd.name})
			if err != nil {
				t.Fatal(err)
			}

			// Check the returned interface type and ensure equality with what's expected.
			if reflect.TypeOf(c) != reflect.TypeOf(cmd.expected) {
				t.Fatalf("Unexpected Executor returned for command[%v]: %v", cmd.name, reflect.TypeOf(c))
			}
		}
	}

	// Unknown Command
	{
		s3 := NewS3(nil, nil)
		unknown := []string{"fake command"}

		if _, err := s3.commandFromArgs(unknown); err == nil {
			t.Fatal("Expected an error for unknown commands.")
		} else if !strings.Contains(err.Error(), unknown[0]) {
			t.Fatalf("Expected error to contain the command name: %v", err.Error())
		}
	}
}

func TestNewS3(t *testing.T) {
	var ui mockIndicator
	var mockS3 mockS3Client

	s3 := NewS3(&mockS3, &ui)

	if s3.con == nil {
		t.Fatalf("Expected S3Handler to be initialized with a Context: %v", s3.con)
	} else if s3.ui != &ui {
		t.Fatalf("S3Handler storing unknown indicator: %v", s3.ui)
	} else if s3.s3 != &mockS3 {
		t.Fatalf("S3Handler storing unknown s3client: %v", s3.s3)
	}
}
