package utils

import (
	"fmt"
	"time"
)

func CalculateAndPrintMetrics(s *Scheduler, completedTasks []*Task, firstTasks []*Task) {
	var totalWaitTime, totalTurnaroundTime int64
	var totalCPUTime int64

	// Ensure there are completedTasks to process
	if len(completedTasks) == 0 {
		fmt.Println("No completedTasks to process.")
		return
	}

	startTime := completedTasks[0].StartTime
	endTime := time.Now()

	// Check if there are completed completedTasks to avoid division by zero
	numcompletedTasks := len(s.Completed)
	if numcompletedTasks == 0 {
		fmt.Println("No completed completedTasks to process metrics.")
		return
	}

	// Map to store turnaround times by task ID
	turnaroundTimes := make(map[int]int64)

	// First loop: Calculate Turnaround Time for each task in s.Completed
	for _, task := range s.Completed {
		turnaroundTime := task.CompletionTime.Sub(task.StartTime).Milliseconds()
		turnaroundTimes[task.ID] = turnaroundTime

		fmt.Printf("Task %d: Turnaround Time: %d ms\n", task.ID, turnaroundTime)
		totalTurnaroundTime += turnaroundTime
	}

	// Second loop: Calculate Waiting Time for each task using firstTasks
	for _, task := range firstTasks {
		turnaroundTime, exists := turnaroundTimes[task.ID]
		if !exists {
			fmt.Printf("Task %d: No turnaround time found. Skipping.\n", task.ID)
			continue
		}

		waitTime := turnaroundTime - int64(task.BurstTime)

		fmt.Printf("Task %d: Waiting Time: %d ms \n", task.ID, waitTime)

		totalWaitTime += waitTime

		totalCPUTime += int64(task.BurstTime)
	}

	// Calculate averages
	avgWaitTime := float64(totalWaitTime) / float64(numcompletedTasks)
	avgTurnaroundTime := float64(totalTurnaroundTime) / float64(numcompletedTasks)

	// Calculate total simulation time and CPU utilization
	totalTime := endTime.Sub(startTime).Milliseconds()
	cpuUtilization := (float64(totalCPUTime) / float64(totalTime)) * 100

	// Print summary metrics
	fmt.Println("\n--- Simulation Metrics ---")
	fmt.Printf("Average Wait Time: %.2f ms\n", avgWaitTime)
	fmt.Printf("Average Turnaround Time: %.2f ms\n", avgTurnaroundTime)
	fmt.Printf("CPU Utilization: %.2f%%\n", cpuUtilization)
	fmt.Printf("Total Simulation Time: %d ms\n", totalTime)
	fmt.Println("--------------------------")
}
