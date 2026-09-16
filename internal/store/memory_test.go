package store

import (
	"slices"
	"testing"

	"github.com/peter-limawal/mdc/internal/domain"
)

func TestMemoryStoreSave(t *testing.T) {
	id := "job-001"
	cmd := []string{"echo", "hello"}

	job := domain.NewJob(id, cmd)
	ms := NewMemoryStore()

	if err := ms.Save(job); err != nil {
		t.Fatalf("unexpected error saving job: %v", err)
	}

	savedJob, exists := ms.jobs[job.ID]

	if !exists {
		t.Fatalf("job %q was not saved", job.ID)
	}

	if savedJob.ID != job.ID {
		t.Errorf("got ID %q, want %q", savedJob.ID, job.ID)
	}
}

func TestMemoryStoreSaveFailsWhenIDAlreadyExists(t *testing.T) {
	id := "job-001"

	existingJob := domain.NewJob(id, []string{"echo", "existing"})
	incomingJob := domain.NewJob(id, []string{"echo", "incoming"})
	ms := NewMemoryStore()

	if err := ms.Save(existingJob); err != nil {
		t.Fatalf("unexpected error saving existing job: %v", err)
	}

	err := ms.Save(incomingJob)

	if err == nil {
		t.Fatal("expected error when saving a duplicate job ID")
	}

	savedJob, exists := ms.jobs[existingJob.ID]

	if !exists {
		t.Fatalf("original job %q was removed", existingJob.ID)
	}

	if !slices.Equal(savedJob.Command, existingJob.Command) {
		t.Errorf(
			"got command %v, want original command %v",
			savedJob.Command,
			existingJob.Command,
		)
	}
}

func TestMemoryStoreGet(t *testing.T) {
	id := "job-001"
	cmd := []string{"echo", "hello"}

	job := domain.NewJob(id, cmd)
	ms := NewMemoryStore()

	if err := ms.Save(job); err != nil {
		t.Fatalf("unexpected error saving job: %v", err)
	}

	gotJob, err := ms.Get(job.ID)

	if err != nil {
		t.Fatalf("unexpected error getting job: %v", err)
	}

	if gotJob.ID != job.ID {
		t.Errorf("got ID %q, want %q", gotJob.ID, job.ID)
	}

	if !slices.Equal(gotJob.Command, job.Command) {
		t.Errorf("got command %v, want %v", gotJob.Command, job.Command)
	}

	if gotJob.State != job.State {
		t.Errorf("got state %q, want %q", gotJob.State, job.State)
	}
}

func TestMemoryStoreGetReturnsErrorWhenIDNotFound(t *testing.T) {
	ms := NewMemoryStore()

	_, err := ms.Get("job-unknown")

	if err == nil {
		t.Fatal("expected error when getting an unknown job ID")
	}
}

func TestMemoryStoreUpdate(t *testing.T) {
	id := "job-001"
	cmd := []string{"echo", "hello"}

	job := domain.NewJob(id, cmd)
	ms := NewMemoryStore()

	if err := ms.Save(job); err != nil {
		t.Fatalf("unexpected error saving job: %v", err)
	}

	if err := job.Start(); err != nil {
		t.Fatalf("unexpected error starting job: %v", err)
	}

	if err := ms.Update(job); err != nil {
		t.Fatalf("unexpected error updating job: %v", err)
	}

	gotJob, err := ms.Get(job.ID)

	if err != nil {
		t.Fatalf("unexpected error getting updated job: %v", err)
	}

	if gotJob.State != domain.StateRunning {
		t.Errorf("got state %q, want %q", gotJob.State, domain.StateRunning)
	}
}

func TestMemoryStoreUpdateReturnsErrorWhenIDNotFound(t *testing.T) {
	id := "job-unknown"
	cmd := []string{"echo", "hello"}

	job := domain.NewJob(id, cmd)
	ms := NewMemoryStore()

	err := ms.Update(job)

	if err == nil {
		t.Fatal("expected error when updating an unknown job ID")
	}
}
