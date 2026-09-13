package executor

import "testing"

func TestLocalExecutorRun(t *testing.T) {
	runner := LocalExecutor{}

	cmd := []string{"echo", "hello"}
	output, err := runner.Run(cmd)

	if err != nil {
		t.Fatalf("unexpected error running command: %v", err)
	}

	if wantOutput := "hello\n"; output != wantOutput {
		t.Errorf("got output %q, want %q", output, wantOutput)
	}
}

func TestLocalExecutorRunFailsWhenCommandIsEmpty(t *testing.T) {
	runner := LocalExecutor{}

	_, err := runner.Run([]string{})

	if err == nil {
		t.Fatal("expected error when running an empty command")
	}
}

func TestLocalExecutorRunFailsWhenCommandIsUnknown(t *testing.T) {
	runner := LocalExecutor{}

	_, err := runner.Run([]string{"__mdc_unknown_command__"})

	if err == nil {
		t.Fatal("expected error when running an unknown command")
	}
}
