package store

import (
	"fmt"

	"github.com/peter-limawal/mdc/internal/domain"
)

type MemoryStore struct {
	jobs map[string]domain.Job
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		jobs: make(map[string]domain.Job),
	}
}

func (ms *MemoryStore) Save(job domain.Job) error {
	if _, exists := ms.jobs[job.ID]; exists {
		return fmt.Errorf("job %q already exists", job.ID)
	}

	ms.jobs[job.ID] = job
	return nil
}

func (ms *MemoryStore) Get(id string) (domain.Job, error) {
	job, exists := ms.jobs[id]

	if !exists {
		return domain.Job{}, fmt.Errorf("job %q not found", id)
	}

	return job, nil
}

func (ms *MemoryStore) Update(job domain.Job) error {
	if _, exists := ms.jobs[job.ID]; !exists {
		return fmt.Errorf("job %q not found", job.ID)
	}

	ms.jobs[job.ID] = job
	return nil
}
