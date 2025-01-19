package utils

import (
	"sync"
	"time"
)

type State string

const (
	Ready      State = "Ready"
	Running    State = "Running"
	WaitingIO  State = "Waiting I/O"
	WaitingRes State = "Waiting Resource"
	Completed  State = "Completed"
)

type Task struct {
	ID                 int
	BurstTime          int
	IOTime             int
	Resources          int
	ResourcesAllocated int
	Priority           int // Lower value = higher priority
	State              State
	EnqueueTime        time.Time
	StartTime          time.Time
	CompletionTime     time.Time
}

type Scheduler struct {
	ReadyQueue chan *Task
	IOQueue    chan *Task
	ResQueue   chan *Task
	Completed  []*Task
	TimeSlice  int
	mu         sync.Mutex
	Res        int
}
