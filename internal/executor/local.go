package executor

import (
	"errors"
	"os/exec"
)

type LocalExecutor struct{}

func (LocalExecutor) Run(command []string) (string, error) {
	if len(command) == 0 {
		return "", errors.New("command is required")
	}

	cmd := exec.Command(command[0], command[1:]...)

	output, err := cmd.CombinedOutput()

	return string(output), err
}
