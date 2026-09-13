package service

import (
	"testing"

	"github.com/peter-limawal/mdc/internal/domain"
	"github.com/peter-limawal/mdc/internal/executor"
	"github.com/peter-limawal/mdc/internal/store"
)

func TestServiceRunSucceeds(t *testing.T) {
	id := "job-001"
	command := []string{"echo", "hello"}

	job := domain.NewJob(id, command)

	ms := store.NewMemoryStore()
	runner := executor.LocalExecutor{}
	svc := New(ms, runner)

	gotJob, output, err := svc.Run(job)

	if err != nil {
		t.Fatalf("unexpected error running job: %v", err)
	}

	if wantOutput := "hello\n"; output != wantOutput {
		t.Errorf("got output %q, want %q", output, wantOutput)
	}

	if gotJob.State != domain.StateSucceeded {
		t.Errorf("got final state %q, want %q", gotJob.State, domain.StateSucceeded)
	}

	storedJob, err := ms.Get(job.ID)

	if err != nil {
		t.Fatalf("unexpected error getting stored job: %v", err)
	}

	if storedJob.State != domain.StateSucceeded {
		t.Errorf("got stored state %q, want %q", storedJob.State, domain.StateSucceeded)
	}
}

func TestServiceRunFailsWhenCommandUnknown(t *testing.T) {
	id := "job-002"
	command := []string{"__mdc_unknown_command__"}

	job := domain.NewJob(id, command)

	ms := store.NewMemoryStore()
	runner := executor.LocalExecutor{}
	svc := New(ms, runner)

	gotJob, _, err := svc.Run(job)

	if err == nil {
		t.Fatal("expected error when command is unknown")
	}

	if gotJob.State != domain.StateFailed {
		t.Errorf("got final state %q, want %q", gotJob.State, domain.StateFailed)
	}

	storedJob, err := ms.Get(job.ID)

	if err != nil {
		t.Fatalf("unexpected error getting stored job: %v", err)
	}

	if storedJob.State != domain.StateFailed {
		t.Errorf("got stored state %q, want %q", storedJob.State, domain.StateFailed)
	}
}

func TestServiceRunRejectsDuplicateJobID(t *testing.T) {
	id := "job-001"
	command := []string{"echo", "hello"}

	job := domain.NewJob(id, command)

	ms := store.NewMemoryStore()
	runner := executor.LocalExecutor{}
	svc := New(ms, runner)

	if _, _, err := svc.Run(job); err != nil {
		t.Fatalf("unexpected error running first job: %v", err)
	}

	_, _, err := svc.Run(job)

	if err == nil {
		t.Fatal("expected error when running a duplicate job ID")
	}

	storedJob, err := ms.Get(job.ID)

	if err != nil {
		t.Fatalf("unexpected error getting stored job: %v", err)
	}

	if storedJob.State != domain.StateSucceeded {
		t.Errorf("got stored state %q, want %q", storedJob.State, domain.StateSucceeded)
	}
}

func TestServiceRunRejectsNonQueuedJob(t *testing.T) {
	id := "job-001"
	command := []string{"echo", "hello"}

	job := domain.NewJob(id, command)

	if err := job.Start(); err != nil {
		t.Fatalf("unexpected error starting job: %v", err)
	}

	ms := store.NewMemoryStore()
	runner := executor.LocalExecutor{}
	svc := New(ms, runner)

	_, _, err := svc.Run(job)

	if err == nil {
		t.Fatal("expected error when running a non-queued job")
	}

	_, err = ms.Get(job.ID)

	if err == nil {
		t.Fatalf("expected non-queued job not to be stored")
	}
}
