package main

import (
	"os/exec"
	"strings"
	"testing"
)

func TestCLIRunCommandSucceeds(t *testing.T) {
	execCmd := exec.Command("go", "run", ".", "run", "echo", "hello")

	output, err := execCmd.CombinedOutput()

	if err != nil {
		t.Fatalf("unexpected CLI error: %v\noutput: %q", err, output)
	}

	gotOutput := string(output)

	if !strings.Contains(gotOutput, "hello\n") {
		t.Errorf("output %q does not contain command output", gotOutput)
	}

	if !strings.Contains(gotOutput, ": succeeded\n") {
		t.Errorf("output %q does not contain succeeded state", gotOutput)
	}
}

func TestCLIRunCommandFailsWhenCommandIsUnknown(t *testing.T) {
	execCmd := exec.Command("go", "run", ".", "run", "__mdc_unknown_command__")

	output, err := execCmd.CombinedOutput()

	if err == nil {
		t.Fatal("expected CLI error for unknown command")
	}

	gotOutput := string(output)

	if !strings.Contains(gotOutput, ": failed\n") {
		t.Errorf("output %q does not contain failed state", gotOutput)
	}
}

func TestCLIShowsUsageWhenArgumentsAreMissing(t *testing.T) {
	execCmd := exec.Command("go", "run", ".")

	output, err := execCmd.CombinedOutput()

	if err == nil {
		t.Fatal("expected CLI error when arguments are missing")
	}

	gotOutput := string(output)
	wantOutput := "usage: mdc run <command> [args...]"

	if !strings.Contains(gotOutput, wantOutput) {
		t.Errorf("output %q does not contain usage message", gotOutput)
	}
}

func TestCLIRunCommandSeparatesOutputFromState(t *testing.T) {
	execCmd := exec.Command("go", "run", ".", "run", "printf", "hello")

	output, err := execCmd.CombinedOutput()

	if err != nil {
		t.Fatalf("unexpected CLI error: %v\noutput: %q", err, output)
	}

	gotOutput := string(output)

	if !strings.Contains(gotOutput, "hello\njob") {
		t.Errorf("output %q does not separate command output from job state", gotOutput)
	}
}
