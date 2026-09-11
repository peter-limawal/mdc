package domain

import "fmt"

type JobState string

const (
	StateQueued    JobState = "queued"
	StateRunning   JobState = "running"
	StateSucceeded JobState = "succeeded"
	StateFailed    JobState = "failed"
	StateCancelled JobState = "cancelled"
)

type Job struct {
	ID      string
	Command []string
	State   JobState
}

func NewJob(id string, command []string) Job {
	return Job{
		ID:      id,
		Command: command,
		State:   StateQueued,
	}
}

func (j *Job) Start() error {
	if j.State != StateQueued {
		return fmt.Errorf(
			"job %q cannot start from state %q",
			j.ID,
			j.State,
		)
	}

	j.State = StateRunning
	return nil
}

func (j *Job) Succeed() error {
	if j.State != StateRunning {
		return fmt.Errorf(
			"job %q cannot succeed from state %q",
			j.ID,
			j.State,
		)
	}

	j.State = StateSucceeded
	return nil
}

func (j *Job) Fail() error {
	if j.State != StateRunning {
		return fmt.Errorf(
			"job %q cannot fail from state %q",
			j.ID,
			j.State,
		)
	}

	j.State = StateFailed
	return nil
}

func (j *Job) Cancel() error {
	switch j.State {
	case StateQueued, StateRunning:
		j.State = StateCancelled
		return nil

	default:
		return fmt.Errorf(
			"job %q cannot cancel from state %q",
			j.ID,
			j.State,
		)
	}
}
