package domain

import (
	"slices"
	"testing"
)

func TestNewJobInitializesQueuedState(t *testing.T) {
	id := "job-001"
	cmd := []string{"echo", "hello"}

	job := NewJob(id, cmd)

	wantID := "job-001"

	if job.ID != wantID {
		t.Errorf("got ID %q, want %q", job.ID, wantID)
	}

	if job.State != StateQueued {
		t.Errorf("got state %q, want %q", job.State, StateQueued)
	}

	wantCommand := []string{"echo", "hello"}

	if !slices.Equal(job.Command, wantCommand) {
		t.Errorf("got command %v, want %v", job.Command, wantCommand)
	}
}

func TestJobStart(t *testing.T) {
	id := "job-001"
	cmd := []string{"echo", "hello"}

	job := NewJob(id, cmd)

	if err := job.Start(); err != nil {
		t.Fatalf("unexpected error starting job: %v", err)
	}

	if job.State != StateRunning {
		t.Errorf("got state %q, want %q", job.State, StateRunning)
	}
}

func TestJobStartFailsWhenAlreadyRunning(t *testing.T) {
	id := "job-001"
	cmd := []string{"echo", "hello"}

	job := NewJob(id, cmd)

	if err := job.Start(); err != nil {
		t.Fatalf("unexpected error on first start: %v", err)
	}

	err := job.Start()

	if err == nil {
		t.Fatal("expected error when starting a running job")
	}

	if job.State != StateRunning {
		t.Errorf("got state %q, want %q", job.State, StateRunning)
	}
}

func TestJobStartFailsWhenAlreadySucceeded(t *testing.T) {
	id := "job-001"
	cmd := []string{"echo", "hello"}

	job := NewJob(id, cmd)

	if err := job.Start(); err != nil {
		t.Fatalf("unexpected error starting job: %v", err)
	}

	if err := job.Succeed(); err != nil {
		t.Fatalf("unexpected error succeeding job: %v", err)
	}

	err := job.Start()

	if err == nil {
		t.Fatal("expected error when starting a succeeded job")
	}

	if job.State != StateSucceeded {
		t.Errorf("got state %q, want %q", job.State, StateSucceeded)
	}
}

func TestJobStartFailsWhenAlreadyFailed(t *testing.T) {
	id := "job-001"
	cmd := []string{"echo", "hello"}

	job := NewJob(id, cmd)

	if err := job.Start(); err != nil {
		t.Fatalf("unexpected error starting job: %v", err)
	}

	if err := job.Fail(); err != nil {
		t.Fatalf("unexpected error failing job: %v", err)
	}

	err := job.Start()

	if err == nil {
		t.Fatal("expected error when starting a failed job")
	}

	if job.State != StateFailed {
		t.Errorf("got state %q, want %q", job.State, StateFailed)
	}
}

func TestJobStartFailsWhenAlreadyCancelled(t *testing.T) {
	id := "job-001"
	cmd := []string{"echo", "hello"}

	job := NewJob(id, cmd)

	if err := job.Start(); err != nil {
		t.Fatalf("unexpected error starting job: %v", err)
	}

	if err := job.Cancel(); err != nil {
		t.Fatalf("unexpected error cancelling job: %v", err)
	}

	err := job.Start()

	if err == nil {
		t.Fatal("expected error when starting a cancelled job")
	}

	if job.State != StateCancelled {
		t.Errorf("got state %q, want %q", job.State, StateCancelled)
	}
}

func TestJobSucceed(t *testing.T) {
	id := "job-001"
	cmd := []string{"echo", "hello"}

	job := NewJob(id, cmd)

	if err := job.Start(); err != nil {
		t.Fatalf("unexpected error starting job: %v", err)
	}

	if err := job.Succeed(); err != nil {
		t.Fatalf("unexpected error succeeding job: %v", err)
	}

	if job.State != StateSucceeded {
		t.Errorf("got state %q, want %q", job.State, StateSucceeded)
	}
}

func TestJobSucceedFailsWhenQueued(t *testing.T) {
	id := "job-001"
	cmd := []string{"echo", "hello"}

	job := NewJob(id, cmd)

	err := job.Succeed()

	if err == nil {
		t.Fatal("expected error when succeeding a queued job")
	}

	if job.State != StateQueued {
		t.Errorf("got state %q, want %q", job.State, StateQueued)
	}
}

func TestJobSucceedFailsWhenAlreadySucceeded(t *testing.T) {
	id := "job-001"
	cmd := []string{"echo", "hello"}

	job := NewJob(id, cmd)

	if err := job.Start(); err != nil {
		t.Fatalf("unexpected error starting job: %v", err)
	}

	if err := job.Succeed(); err != nil {
		t.Fatalf("unexpected error succeeding job: %v", err)
	}

	err := job.Succeed()

	if err == nil {
		t.Fatal("expected error when succeeding a succeeded job")
	}

	if job.State != StateSucceeded {
		t.Errorf("got state %q, want %q", job.State, StateSucceeded)
	}
}

func TestJobSucceedFailsWhenAlreadyFailed(t *testing.T) {
	id := "job-001"
	cmd := []string{"echo", "hello"}

	job := NewJob(id, cmd)

	if err := job.Start(); err != nil {
		t.Fatalf("unexpected error starting job: %v", err)
	}

	if err := job.Fail(); err != nil {
		t.Fatalf("unexpected error failing job: %v", err)
	}

	err := job.Succeed()

	if err == nil {
		t.Fatal("expected error when succeeding a failed job")
	}

	if job.State != StateFailed {
		t.Errorf("got state %q, want %q", job.State, StateFailed)
	}
}

func TestJobSucceedFailsWhenAlreadyCancelled(t *testing.T) {
	id := "job-001"
	cmd := []string{"echo", "hello"}

	job := NewJob(id, cmd)

	if err := job.Start(); err != nil {
		t.Fatalf("unexpected error starting job: %v", err)
	}

	if err := job.Cancel(); err != nil {
		t.Fatalf("unexpected error cancelling job: %v", err)
	}

	err := job.Succeed()

	if err == nil {
		t.Fatal("expected error when succeeding a cancelled job")
	}

	if job.State != StateCancelled {
		t.Errorf("got state %q, want %q", job.State, StateCancelled)
	}
}

func TestJobFail(t *testing.T) {
	id := "job-001"
	cmd := []string{"echo", "hello"}

	job := NewJob(id, cmd)

	if err := job.Start(); err != nil {
		t.Fatalf("unexpected error starting job: %v", err)
	}

	if err := job.Fail(); err != nil {
		t.Fatalf("unexpected error failing job: %v", err)
	}

	if job.State != StateFailed {
		t.Errorf("got state %q, want %q", job.State, StateFailed)
	}
}

func TestJobFailFailsWhenQueued(t *testing.T) {
	id := "job-001"
	cmd := []string{"echo", "hello"}

	job := NewJob(id, cmd)

	err := job.Fail()

	if err == nil {
		t.Fatal("expected error when failing a queued job")
	}

	if job.State != StateQueued {
		t.Errorf("got state %q, want %q", job.State, StateQueued)
	}
}

func TestJobFailFailsWhenAlreadyFailed(t *testing.T) {
	id := "job-001"
	cmd := []string{"echo", "hello"}

	job := NewJob(id, cmd)

	if err := job.Start(); err != nil {
		t.Fatalf("unexpected error starting job: %v", err)
	}

	if err := job.Fail(); err != nil {
		t.Fatalf("unexpected error failing job: %v", err)
	}

	err := job.Fail()

	if err == nil {
		t.Fatal("expected error when failing a failed job")
	}

	if job.State != StateFailed {
		t.Errorf("got state %q, want %q", job.State, StateFailed)
	}
}

func TestJobFailFailsWhenAlreadySucceeded(t *testing.T) {
	id := "job-001"
	cmd := []string{"echo", "hello"}

	job := NewJob(id, cmd)

	if err := job.Start(); err != nil {
		t.Fatalf("unexpected error starting job: %v", err)
	}

	if err := job.Succeed(); err != nil {
		t.Fatalf("unexpected error succeeding job: %v", err)
	}

	err := job.Fail()

	if err == nil {
		t.Fatal("expected error when failing a succeeded job")
	}

	if job.State != StateSucceeded {
		t.Errorf("got state %q, want %q", job.State, StateSucceeded)
	}
}

func TestJobFailFailsWhenAlreadyCancelled(t *testing.T) {
	id := "job-001"
	cmd := []string{"echo", "hello"}

	job := NewJob(id, cmd)

	if err := job.Start(); err != nil {
		t.Fatalf("unexpected error starting job: %v", err)
	}

	if err := job.Cancel(); err != nil {
		t.Fatalf("unexpected error cancelling job: %v", err)
	}

	err := job.Fail()

	if err == nil {
		t.Fatal("expected error when failing a cancelled job")
	}

	if job.State != StateCancelled {
		t.Errorf("got state %q, want %q", job.State, StateCancelled)
	}
}

func TestJobCancelWhenQueued(t *testing.T) {
	id := "job-001"
	cmd := []string{"echo", "hello"}

	job := NewJob(id, cmd)

	if err := job.Cancel(); err != nil {
		t.Fatalf("unexpected error cancelling queued job: %v", err)
	}

	if job.State != StateCancelled {
		t.Errorf("got state %q, want %q", job.State, StateCancelled)
	}
}

func TestJobCancelWhenRunning(t *testing.T) {
	id := "job-001"
	cmd := []string{"echo", "hello"}

	job := NewJob(id, cmd)

	if err := job.Start(); err != nil {
		t.Fatalf("unexpected error starting job: %v", err)
	}

	if err := job.Cancel(); err != nil {
		t.Fatalf("unexpected error cancelling running job: %v", err)
	}

	if job.State != StateCancelled {
		t.Errorf("got state %q, want %q", job.State, StateCancelled)
	}
}

func TestJobCancelFailsWhenAlreadySucceeded(t *testing.T) {
	id := "job-001"
	cmd := []string{"echo", "hello"}

	job := NewJob(id, cmd)

	if err := job.Start(); err != nil {
		t.Fatalf("unexpected error starting job: %v", err)
	}

	if err := job.Succeed(); err != nil {
		t.Fatalf("unexpected error succeeding job: %v", err)
	}

	err := job.Cancel()

	if err == nil {
		t.Fatal("expected error when cancelling a succeeded job")
	}

	if job.State != StateSucceeded {
		t.Errorf("got state %q, want %q", job.State, StateSucceeded)
	}
}

func TestJobCancelFailsWhenAlreadyFailed(t *testing.T) {
	id := "job-001"
	cmd := []string{"echo", "hello"}

	job := NewJob(id, cmd)

	if err := job.Start(); err != nil {
		t.Fatalf("unexpected error starting job: %v", err)
	}

	if err := job.Fail(); err != nil {
		t.Fatalf("unexpected error failing job: %v", err)
	}

	err := job.Cancel()

	if err == nil {
		t.Fatal("expected error when cancelling a failed job")
	}

	if job.State != StateFailed {
		t.Errorf("got state %q, want %q", job.State, StateFailed)
	}
}

func TestJobCancelFailsWhenAlreadyCancelled(t *testing.T) {
	id := "job-001"
	cmd := []string{"echo", "hello"}

	job := NewJob(id, cmd)

	if err := job.Start(); err != nil {
		t.Fatalf("unexpected error starting job: %v", err)
	}

	if err := job.Cancel(); err != nil {
		t.Fatalf("unexpected error cancelling running job: %v", err)
	}

	err := job.Cancel()

	if err == nil {
		t.Fatal("expected error when cancelling a cancelled job")
	}

	if job.State != StateCancelled {
		t.Errorf("got state %q, want %q", job.State, StateCancelled)
	}
}
