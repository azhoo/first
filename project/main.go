package main

import (
	"fmt"
	u "project/utils"
	"time"
)

func main() {
	var algoChoice int
	fmt.Println("Select Scheduling Algorithm:")
	fmt.Println("1. Round-Robin")
	fmt.Println("2. First-Come-First-Served (FCFS)")
	fmt.Println("3. Priority Scheduling")
	fmt.Println("4. Shortest Job Next (SJN)")
	fmt.Print("Enter choice: ")
	fmt.Scan(&algoChoice)

	scheduler := &u.Scheduler{
		ReadyQueue: make(chan *u.Task, 20),
		IOQueue:    make(chan *u.Task, 20),
		ResQueue:   make(chan *u.Task, 20),
		TimeSlice:  100,
		Res:        20,
	}

	// Add tasks to the ReadyQueue
	tasks := []*u.Task{
		{ID: 2, BurstTime: 2000, IOTime: 500, Resources: 5, ResourcesAllocated: 0, Priority: 1, State: u.Ready, StartTime: time.Now()},
		{ID: 3, BurstTime: 3000, IOTime: 200, Resources: 10, ResourcesAllocated: 0, Priority: 2, State: u.Ready, StartTime: time.Now()},
		{ID: 4, BurstTime: 500, IOTime: 5000, Resources: 15, ResourcesAllocated: 0, Priority: 3, State: u.Ready, StartTime: time.Now()},
		{ID: 5, BurstTime: 200, IOTime: 20000, Resources: 20, ResourcesAllocated: 0, Priority: 4, State: u.Ready, StartTime: time.Now()},
		{ID: 6, BurstTime: 100, IOTime: 0, Resources: 9, ResourcesAllocated: 0, Priority: 5, State: u.Ready, StartTime: time.Now()},
	}

	tasks2 := []*u.Task{
		{ID: 2, BurstTime: 2000, IOTime: 500, Resources: 5, ResourcesAllocated: 0, Priority: 1, State: u.Ready, StartTime: time.Now()},
		{ID: 3, BurstTime: 3000, IOTime: 200, Resources: 10, ResourcesAllocated: 0, Priority: 2, State: u.Ready, StartTime: time.Now()},
		{ID: 4, BurstTime: 500, IOTime: 5000, Resources: 15, ResourcesAllocated: 0, Priority: 3, State: u.Ready, StartTime: time.Now()},
		{ID: 5, BurstTime: 200, IOTime: 20000, Resources: 20, ResourcesAllocated: 0, Priority: 4, State: u.Ready, StartTime: time.Now()},
		{ID: 6, BurstTime: 100, IOTime: 0, Resources: 9, ResourcesAllocated: 0, Priority: 5, State: u.Ready, StartTime: time.Now()},
	}

	for i := 0; i < len(tasks); i++ {
		task := tasks[i]
		scheduler.ReadyQueue <- task

	}

	// Start the chosen scheduling algorithm
	switch algoChoice {
	case 1:

		go u.RoundRobin(scheduler)
	case 2:

		go u.FCFS(scheduler)
	case 3:
		go u.PriorityScheduling(scheduler)
	case 4:
		go u.ShortestJobNext(scheduler)
	default:
		fmt.Println("Invalid choice. Exiting.")
		return
	}

	// Start handling I/O
	go u.HandleIO(scheduler)

	// Start handling resources
	go u.HandleRes(scheduler)

	// Wait for goroutines to finish

	for i := 0; ; {

		if i > 600 { // some time near 15 - 20 sec
			fmt.Print("\nTimeout reached! Exiting.\n")
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
	u.CalculateAndPrintMetrics(scheduler, scheduler.Completed, tasks2)

}
