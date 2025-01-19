package main

import (
	"fmt"
	"sort"
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
}

var res = 20

// Round-Robin Scheduling Algorithm
func RoundRobin(s *Scheduler) {
	if len(s.ReadyQueue) > 0 {
		for {
			task := <-s.ReadyQueue

			fmt.Print("\n.                                res = ", res, "\n")
			fmt.Print(task, "1\n")

			if task.IOTime > 0 {
				task.State = WaitingIO
				s.IOQueue <- task
			} else if task.Resources > 0 {
				task.State = WaitingRes
				fmt.Print("\n\n\nResQ_________________> ", task)
				s.ResQueue <- task
			} else if task.BurstTime > s.TimeSlice {

				res += task.ResourcesAllocated
				task.ResourcesAllocated = 0

				task.State = Running
				time.Sleep(time.Duration(s.TimeSlice) * time.Millisecond)
				task.BurstTime -= s.TimeSlice

				task.State = Ready

				s.ReadyQueue <- task

				fmt.Print(task, "2\n")

			} else {
				task.BurstTime = 0
				task.State = Completed

				res += task.ResourcesAllocated // retrieving resources back
				task.ResourcesAllocated = 0

				fmt.Print(task, "3\n")

				s.mu.Lock()
				s.Completed = append(s.Completed, task)
				s.mu.Unlock()
			}
		}
	}

}

// First-Come-First-Served (FCFS) Scheduling Algorithm
func FCFS(s *Scheduler) {

	for {
		if len(s.ReadyQueue) > 0 {
			task := <-s.ReadyQueue

			fmt.Print("\n.                                res = ", res, "\n")
			fmt.Print(task, "1\n")

			if task.IOTime > 0 {
				task.State = WaitingIO
				s.IOQueue <- task
			} else if task.Resources > 0 {
				task.State = WaitingRes
				fmt.Print("\n\n\nResQ_________________> ", task)
				s.ResQueue <- task
			} else {

				task.State = Running
				time.Sleep(time.Duration(task.BurstTime) * time.Millisecond)

				res += task.ResourcesAllocated // retrieving resources back
				task.ResourcesAllocated = 0

				task.BurstTime = 0
				task.State = Completed
				s.mu.Lock()
				s.Completed = append(s.Completed, task)
				s.mu.Unlock()

				fmt.Print(task, "3\n")
			}
		}
	}
}

func PriorityScheduling(s *Scheduler) {

	for {

		if len(s.ReadyQueue) > 0 {
			tasks := []*Task{}
			for len(s.ReadyQueue) > 0 {
				tasks = append(tasks, <-s.ReadyQueue)
			}

			// Sort tasks by priority (ascending)
			sort.Slice(tasks, func(i, j int) bool {
				return tasks[i].Priority < tasks[j].Priority
			})

			// Process sorted tasks
			for _, task := range tasks {

				if task.IOTime > 0 {
					task.State = WaitingIO
					s.IOQueue <- task
				} else if task.Resources > 0 {
					task.State = WaitingRes
					fmt.Print("\n\n\nResQ_________________> ", task)
					s.ResQueue <- task
				} else {

					task.State = Running
					time.Sleep(time.Duration(task.BurstTime) * time.Millisecond)

					res += task.ResourcesAllocated // retrieving resources back
					task.ResourcesAllocated = 0

					task.BurstTime = 0
					task.State = Completed
					s.mu.Lock()
					s.Completed = append(s.Completed, task)
					s.mu.Unlock()
				}
			}
		}
	}
}

func ShortestJobNext(s *Scheduler) {

	for {

		if len(s.ReadyQueue) > 0 {
			tasks := []*Task{}
			for len(s.ReadyQueue) > 0 {
				tasks = append(tasks, <-s.ReadyQueue)
			}

			// Sort tasks by BurstTime (ascending)
			sort.Slice(tasks, func(i, j int) bool {
				return tasks[i].BurstTime < tasks[j].BurstTime
			})

			// Process sorted tasks
			for _, task := range tasks {

				if task.IOTime > 0 {
					task.State = WaitingIO
					s.IOQueue <- task
				} else if task.Resources > 0 {
					task.State = WaitingRes
					fmt.Print("\n\n\nResQ_________________> ", task)
					s.ResQueue <- task
				} else {

					task.State = Running
					time.Sleep(time.Duration(task.BurstTime) * time.Millisecond)

					res += task.ResourcesAllocated // retrieving resources back
					task.ResourcesAllocated = 0

					task.BurstTime = 0
					task.State = Completed
					s.mu.Lock()
					s.Completed = append(s.Completed, task)
					s.mu.Unlock()
				}
			}
		}
	}
}

// Handle I/O Operations
func HandleIO(s *Scheduler) {

	for {
		if len(s.IOQueue) > 0 {
			task := <-s.IOQueue
			fmt.Print(task, "_________________IO\n")
			time.Sleep(time.Duration(task.IOTime) * time.Millisecond)
			task.IOTime = 0
			fmt.Print(task, "______________________________IO\n")
			task.State = Ready
			s.ReadyQueue <- task
		}
	}
}

func HandleRes(s *Scheduler) {
	tasks := []*Task{}
	threshold := 1 * time.Second

	for {
		if len(s.ResQueue) > 0 {
			task1 := <-s.ResQueue
			fmt.Print("    \nthis task is going in resqueue :    ", task1)
			task1.EnqueueTime = time.Now()
			tasks = append(tasks, task1)

			for i := 0; i < len(tasks); i++ {
				task := tasks[i]

				// Check if the task's waiting time exceeds the threshold
				if time.Since(task.EnqueueTime) > threshold {
					tasks[i].Priority-- // Increase priority (lower number means higher priority)
					fmt.Printf("Escalating priority for task: %d (new priority: %d)\n", task.ID, task.Priority)
				}
			}

			sort.Slice(tasks, func(i, j int) bool {
				return tasks[i].Priority < tasks[j].Priority
			})

			task := tasks[0]

			fmt.Print("\n", task, "_________________res\n")
			if task.Resources <= res {
				res -= task.Resources
				task.ResourcesAllocated = task.Resources
				task.Resources = 0
				fmt.Print("\n", task, "_____________________________Res\n")
				task.State = Ready
				s.ReadyQueue <- task

				tasks = tasks[1:]

			}
		}
	}
}

func calculateAndPrintMetrics(s *Scheduler, tasks []*Task) {

	var totalWaitTime, totalTurnaroundTime time.Duration
	var totalCPUTime int
	startTime := tasks[0].EnqueueTime
	endTime := time.Now()

	for _, task := range s.Completed {
		// Calculate task-specific metrics
		waitTime := task.StartTime.Sub(task.EnqueueTime)
		turnaroundTime := task.CompletionTime.Sub(task.EnqueueTime)
		totalCPUTime += task.BurstTime // CPU time is the sum of burst times for all completed tasks

		// Accumulate metrics
		totalWaitTime += waitTime
		totalTurnaroundTime += turnaroundTime

		// Print individual task metrics
		fmt.Printf("Task %d: Wait Time: %.2f ms, Turnaround Time: %.2f ms\n",
			task.ID, waitTime.Seconds()*1000, turnaroundTime.Seconds()*1000)
	}

	// Calculate averages
	numTasks := len(s.Completed)
	avgWaitTime := totalWaitTime / time.Duration(numTasks)
	avgTurnaroundTime := totalTurnaroundTime / time.Duration(numTasks)
	totalTime := endTime.Sub(startTime)
	cpuUtilization := (float64(totalCPUTime) / float64(totalTime.Milliseconds())) * 100

	// Print summary metrics
	fmt.Println("\n--- Simulation Metrics ---")
	fmt.Printf("Average Wait Time: %.2f ms\n", avgWaitTime.Seconds()*1000)
	fmt.Printf("Average Turnaround Time: %.2f ms\n", avgTurnaroundTime.Seconds()*1000)
	fmt.Printf("CPU Utilization: %.2f%%\n", cpuUtilization)
	fmt.Printf("Total Simulation Time: %.2f ms\n", totalTime.Seconds()*1000)
	fmt.Println("--------------------------")
}

func main() {
	var algoChoice int
	fmt.Println("Select Scheduling Algorithm:")
	fmt.Println("1. Round-Robin")
	fmt.Println("2. First-Come-First-Served (FCFS)")
	fmt.Println("3. Priority Scheduling")
	fmt.Println("4. Shortest Job Next (SJN)")
	fmt.Print("Enter choice: ")
	fmt.Scan(&algoChoice)

	scheduler := &Scheduler{
		ReadyQueue: make(chan *Task, 20),
		IOQueue:    make(chan *Task, 20),
		ResQueue:   make(chan *Task, 20),
		TimeSlice:  100,
	}

	// Add tasks to the ReadyQueue
	tasks := []*Task{
		{ID: 2, BurstTime: 2000, IOTime: 500, Resources: 5, ResourcesAllocated: 0, Priority: 1, State: Ready, StartTime: time.Now()},
		{ID: 3, BurstTime: 3000, IOTime: 200, Resources: 10, ResourcesAllocated: 0, Priority: 2, State: Ready, StartTime: time.Now()},
		{ID: 4, BurstTime: 500, IOTime: 5000, Resources: 15, ResourcesAllocated: 0, Priority: 3, State: Ready, StartTime: time.Now()},
		{ID: 5, BurstTime: 200, IOTime: 20000, Resources: 20, ResourcesAllocated: 0, Priority: 4, State: Ready, StartTime: time.Now()},
		{ID: 6, BurstTime: 100, IOTime: 0, Resources: 9, ResourcesAllocated: 0, Priority: 5, State: Ready, StartTime: time.Now()},
	}

	for i := 0; i < len(tasks); i++ {
		task := tasks[i]
		scheduler.ReadyQueue <- task

	}

	// Close ReadyQueue after adding all tasks

	// Add goroutines to WaitGroup

	// Start the chosen scheduling algorithm
	switch algoChoice {
	case 1:

		go RoundRobin(scheduler)
	case 2:

		go FCFS(scheduler)
	case 3:
		go PriorityScheduling(scheduler)
	case 4:
		go ShortestJobNext(scheduler)
	default:
		fmt.Println("Invalid choice. Exiting.")
		return
	}

	// Start handling I/O
	go HandleIO(scheduler)
	// Start handling resources
	go HandleRes(scheduler)

	// Wait for goroutines to finish

	for i := 0; ; {

		if i > 600 { // some time near 15 - 20 sec
			fmt.Println("\nTimeout reached! Exiting.\n")
			return
		} else {

			if len(scheduler.Completed) == len(tasks) {
				break
			}
		}

		time.Sleep(100 * time.Millisecond) // Check periodically
		i++
	}
	close(scheduler.ReadyQueue)
	close(scheduler.IOQueue)
	close(scheduler.ResQueue)

	// Print completed tasks
	for _, task := range scheduler.Completed {
		fmt.Printf("Task %d completed with state: %s\n", task.ID, task.State)
	}

	// After all tasks are completed and channels are closed
	calculateAndPrintMetrics(scheduler, scheduler.Completed)

}
