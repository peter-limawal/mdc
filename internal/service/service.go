package service

import (
	"github.com/peter-limawal/mdc/internal/domain"
	"github.com/peter-limawal/mdc/internal/executor"
	"github.com/peter-limawal/mdc/internal/store"
)

type Service struct {
	ms     *store.MemoryStore
	runner executor.LocalExecutor
}

func New(ms *store.MemoryStore, runner executor.LocalExecutor) *Service {
	return &Service{ms: ms, runner: runner}
}

func (s *Service) Run(job domain.Job) (domain.Job, string, error) {
	initialJob := job

	if err := initialJob.Start(); err != nil {
		return job, "", err
	}

	if err := s.ms.Save(job); err != nil {
		return job, "", err
	}

	if err := s.ms.Update(initialJob); err != nil {
		return job, "", err
	}

	job = initialJob

	output, execErr := s.runner.Run(job.Command)

	if execErr != nil {
		if err := job.Fail(); err != nil {
			return job, output, err
		}

		if err := s.ms.Update(job); err != nil {
			return job, output, err
		}

		return job, output, execErr
	}

	if err := job.Succeed(); err != nil {
		return job, output, err
	}

	if err := s.ms.Update(job); err != nil {
		return job, output, err
	}

	return job, output, nil
}
