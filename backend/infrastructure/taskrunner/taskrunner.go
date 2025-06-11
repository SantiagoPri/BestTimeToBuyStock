package taskrunner

import (
	"fmt"
	"runtime/debug"
)

type TaskRunner struct {
	tasks chan func()
}

func New(bufferSize int) *TaskRunner {
	return &TaskRunner{
		tasks: make(chan func(), bufferSize),
	}
}

func (tr *TaskRunner) Start() {
	go func() {
		for task := range tr.tasks {
			go func(t func()) {
				defer func() {
					if r := recover(); r != nil {
						fmt.Printf("Task panicked: %v\nStack trace:\n%s\n", r, debug.Stack())
					}
				}()
				t()
			}(task)
		}
	}()
}

func (tr *TaskRunner) Dispatch(task func()) {
	tr.tasks <- task
}
