package store

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/peter-limawal/mdc/internal/domain"
	_ "modernc.org/sqlite"
)

const createJobsTable = `
CREATE TABLE IF NOT EXISTS jobs (
    id TEXT PRIMARY KEY,
    command TEXT NOT NULL,
    state TEXT NOT NULL
);
`

const insertJob = `
INSERT INTO jobs (id, command, state)
VALUES (?, ?, ?);
`

const selectJobByID = `
SELECT id, command, state
FROM jobs
WHERE id = ?;
`

const updateJob = `
UPDATE jobs
SET command = ?, state = ?
WHERE id = ?;
`

type SQLiteStore struct {
	db *sql.DB
}

func NewSQLiteStore(path string) (*SQLiteStore, error) {
	db, err := sql.Open("sqlite", path)

	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, err
	}

	if _, err := db.Exec(createJobsTable); err != nil {
		_ = db.Close()
		return nil, err
	}

	return &SQLiteStore{
		db: db,
	}, nil
}

func (sqliteStore *SQLiteStore) Save(job domain.Job) error {
	commandJSON, err := json.Marshal(job.Command)

	if err != nil {
		return fmt.Errorf("marshal command for job %q: %w", job.ID, err)
	}

	_, err = sqliteStore.db.Exec(insertJob, job.ID, string(commandJSON), string(job.State))

	if err != nil {
		return fmt.Errorf("save job %q: %w", job.ID, err)
	}

	return nil
}

func (sqliteStore *SQLiteStore) Get(id string) (domain.Job, error) {
	var (
		job         domain.Job
		commandJSON string
	)

	err := sqliteStore.db.QueryRow(selectJobByID, id).Scan(&job.ID, &commandJSON, &job.State)

	if errors.Is(err, sql.ErrNoRows) {
		return domain.Job{}, fmt.Errorf("job %q not found", id)
	}

	if err != nil {
		return domain.Job{}, fmt.Errorf("get job %q: %w", id, err)
	}

	if err := json.Unmarshal([]byte(commandJSON), &job.Command); err != nil {
		return domain.Job{}, fmt.Errorf("unmarshal command for job %q: %w", id, err)
	}

	return job, nil
}

func (sqliteStore *SQLiteStore) Update(job domain.Job) error {
	commandJSON, err := json.Marshal(job.Command)

	if err != nil {
		return fmt.Errorf("marshal command for job %q: %w", job.ID, err)
	}

	result, err := sqliteStore.db.Exec(updateJob, string(commandJSON), string(job.State), job.ID)

	if err != nil {
		return fmt.Errorf("update job %q: %w", job.ID, err)
	}

	rowsAffected, err := result.RowsAffected()

	if err != nil {
		return fmt.Errorf("get updated rows for job %q: %w", job.ID, err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("job %q not found", job.ID)
	}

	return nil
}

func (sqliteStore *SQLiteStore) Close() error {
	return sqliteStore.db.Close()
}
