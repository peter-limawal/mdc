package service

import (
	"github.com/peter-limawal/mdc/internal/domain"
	"github.com/peter-limawal/mdc/internal/executor"
)

type JobStore interface {
	Save(job domain.Job) error
	Get(id string) (domain.Job, error)
	Update(job domain.Job) error
}

type Service struct {
	jobStore JobStore
	runner   executor.LocalExecutor
}

func New(jobStore JobStore, runner executor.LocalExecutor) *Service {
	return &Service{
		jobStore: jobStore,
		runner:   runner,
	}
}

func (s *Service) Run(job domain.Job) (domain.Job, string, error) {
	queuedJob := job

	if err := job.Start(); err != nil {
		return job, "", err
	}

	if err := s.jobStore.Save(queuedJob); err != nil {
		return job, "", err
	}

	if err := s.jobStore.Update(job); err != nil {
		return job, "", err
	}

	output, execErr := s.runner.Run(job.Command)

	if execErr != nil {
		if err := job.Fail(); err != nil {
			return job, output, err
		}

		if err := s.jobStore.Update(job); err != nil {
			return job, output, err
		}

		return job, output, execErr
	}

	if err := job.Succeed(); err != nil {
		return job, output, err
	}

	if err := s.jobStore.Update(job); err != nil {
		return job, output, err
	}

	return job, output, nil
}
