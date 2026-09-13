package main

import (
	"fmt"
	"os"
	"time"

	"github.com/peter-limawal/mdc/internal/domain"
	"github.com/peter-limawal/mdc/internal/executor"
	"github.com/peter-limawal/mdc/internal/service"
	"github.com/peter-limawal/mdc/internal/store"
)

func main() {
	if len(os.Args) < 3 || os.Args[1] != "run" {
		fmt.Fprintln(os.Stderr, "usage: mdc run <command> [args...]")
		os.Exit(1)
	}

	id := fmt.Sprintf("job-%d", time.Now().UnixNano())
	command := os.Args[2:]

	job := domain.NewJob(id, command)

	ms := store.NewMemoryStore()
	svc := service.New(ms, executor.LocalExecutor{})

	finalJob, output, err := svc.Run(job)

	if output != "" {
		fmt.Print(output)

		if output[len(output)-1] != '\n' {
			fmt.Println()
		}
	}

	fmt.Printf("job %s: %s\n", finalJob.ID, finalJob.State)

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
