package store

import (
	"encoding/json"
	"path/filepath"
	"slices"
	"testing"

	"github.com/peter-limawal/mdc/internal/domain"
)

func TestNewSQLiteStoreOpensInMemoryDatabase(t *testing.T) {
	sqliteStore, err := NewSQLiteStore(":memory:")

	if err != nil {
		t.Fatalf("unexpected error opening in-memory database: %v", err)
	}

	defer sqliteStore.Close()

	if err := sqliteStore.db.Ping(); err != nil {
		t.Fatalf("unexpected error pinging in-memory database: %v", err)
	}
}

func TestNewSQLiteStoreCreatesJobsTable(t *testing.T) {
	const findJobsTable = `
	SELECT name
	FROM sqlite_master
	WHERE type = 'table' AND name = 'jobs'
	`
	sqliteStore, err := NewSQLiteStore(":memory:")

	if err != nil {
		t.Fatalf("unexpected error opening in-memory database: %v", err)
	}

	defer sqliteStore.Close()

	var tableName string

	err = sqliteStore.db.QueryRow(findJobsTable).Scan(&tableName)

	if err != nil {
		t.Fatalf("unexpected error finding jobs table: %v", err)
	}

	if tableName != "jobs" {
		t.Errorf("got table %q, want %q", tableName, "jobs")
	}
}

func TestSQLiteStoreSave(t *testing.T) {
	const getSavedJob = `
	SELECT id, command, state
	FROM jobs
	WHERE id = ?
	`

	sqliteStore, err := NewSQLiteStore(":memory:")

	if err != nil {
		t.Fatalf("unexpected error opening in-memory database: %v", err)
	}

	defer sqliteStore.Close()

	id := "job-001"
	command := []string{"echo", "hello"}

	job := domain.NewJob(id, command)

	if err := sqliteStore.Save(job); err != nil {
		t.Fatalf("unexpected error saving job: %v", err)
	}

	var (
		gotID          string
		gotCommandJSON string
		gotState       domain.JobState
	)

	err = sqliteStore.db.QueryRow(getSavedJob, job.ID).Scan(&gotID, &gotCommandJSON, &gotState)

	if err != nil {
		t.Fatalf("unexpected error getting saved job: %v", err)
	}

	if gotID != job.ID {
		t.Errorf("got ID %q, want %q", gotID, job.ID)
	}

	var gotCommand []string

	if err := json.Unmarshal([]byte(gotCommandJSON), &gotCommand); err != nil {
		t.Fatalf("unexpected error unmarshaling command: %v", err)
	}

	if !slices.Equal(gotCommand, job.Command) {
		t.Errorf("got command %v, want %v", gotCommand, job.Command)
	}

	if gotState != job.State {
		t.Errorf("got state %q, want %q", gotState, job.State)
	}
}

func TestSQLiteStoreSaveFailsWhenIDAlreadyExists(t *testing.T) {
	const getSavedCommand = `
	SELECT command
	FROM jobs
	WHERE id = ?
	`

	sqliteStore, err := NewSQLiteStore(":memory:")

	if err != nil {
		t.Fatalf("unexpected error opening in-memory database: %v", err)
	}

	defer sqliteStore.Close()

	id := "job-001"

	existingJob := domain.NewJob(id, []string{"echo", "existing"})
	incomingJob := domain.NewJob(id, []string{"echo", "incoming"})

	if err := sqliteStore.Save(existingJob); err != nil {
		t.Fatalf("unexpected error saving existing job: %v", err)
	}

	err = sqliteStore.Save(incomingJob)

	if err == nil {
		t.Fatal("expected error when saving a duplicate job ID")
	}

	var gotCommandJSON string

	err = sqliteStore.db.QueryRow(getSavedCommand, existingJob.ID).Scan(&gotCommandJSON)

	if err != nil {
		t.Fatalf("unexpected error getting saved command: %v", err)
	}

	var gotCommand []string

	if err := json.Unmarshal([]byte(gotCommandJSON), &gotCommand); err != nil {
		t.Fatalf("unexpected error unmarshaling command: %v", err)
	}

	if !slices.Equal(gotCommand, existingJob.Command) {
		t.Errorf("got command %v, want %v", gotCommand, existingJob.Command)
	}
}

func TestSQLiteStoreGet(t *testing.T) {
	sqliteStore, err := NewSQLiteStore(":memory:")

	if err != nil {
		t.Fatalf("unexpected error opening in-memory database: %v", err)
	}

	defer sqliteStore.Close()

	id := "job-001"
	command := []string{"echo", "hello"}

	job := domain.NewJob(id, command)

	if err := sqliteStore.Save(job); err != nil {
		t.Fatalf("unexpected error saving job: %v", err)
	}

	gotJob, err := sqliteStore.Get(job.ID)

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

func TestSQLiteStoreGetFailsWhenIDNotFound(t *testing.T) {
	sqliteStore, err := NewSQLiteStore(":memory:")

	if err != nil {
		t.Fatalf("unexpected error opening in-memory database: %v", err)
	}

	defer sqliteStore.Close()

	_, err = sqliteStore.Get("job-unknown")

	if err == nil {
		t.Fatal("expected error when getting an unknown job ID")
	}
}

func TestSQLiteStoreUpdate(t *testing.T) {
	sqliteStore, err := NewSQLiteStore(":memory:")

	if err != nil {
		t.Fatalf("unexpected error opening in-memory database: %v", err)
	}

	defer sqliteStore.Close()

	id := "job-001"
	command := []string{"echo", "hello"}

	job := domain.NewJob(id, command)

	if err := sqliteStore.Save(job); err != nil {
		t.Fatalf("unexpected error saving job: %v", err)
	}

	if err := job.Start(); err != nil {
		t.Fatalf("unexpected error starting job: %v", err)
	}

	if err := sqliteStore.Update(job); err != nil {
		t.Fatalf("unexpected error updating job: %v", err)
	}

	gotJob, err := sqliteStore.Get(job.ID)

	if err != nil {
		t.Fatalf("unexpected error getting updated job: %v", err)
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

func TestSQLiteStoreUpdateFailsWhenIDNotFound(t *testing.T) {
	sqliteStore, err := NewSQLiteStore(":memory:")

	if err != nil {
		t.Fatalf("unexpected error opening in-memory database: %v", err)
	}

	defer sqliteStore.Close()

	id := "job-unknown"
	command := []string{"echo", "hello"}

	job := domain.NewJob(id, command)

	err = sqliteStore.Update(job)

	if err == nil {
		t.Fatal("expected error when updating an unknown job ID")
	}
}

func TestSQLiteStorePersistsJobAcrossReopen(t *testing.T) {
	path := filepath.Join(t.TempDir(), "mdc.db")

	firstStore, err := NewSQLiteStore(path)

	if err != nil {
		t.Fatalf("unexpected error opening database: %v", err)
	}

	id := "job-001"
	command := []string{"echo", "hello"}

	job := domain.NewJob(id, command)

	if err := firstStore.Save(job); err != nil {
		t.Fatalf("unexpected error saving job: %v", err)
	}

	if err := firstStore.Close(); err != nil {
		t.Fatalf("unexpected error closing database: %v", err)
	}

	reopenedStore, err := NewSQLiteStore(path)

	if err != nil {
		t.Fatalf("unexpected error reopening database: %v", err)
	}

	defer reopenedStore.Close()

	gotJob, err := reopenedStore.Get(job.ID)

	if err != nil {
		t.Fatalf("unexpected error getting persisted job: %v", err)
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
