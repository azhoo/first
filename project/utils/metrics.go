package utils

import (
	"fmt"

	"time"
)

func CalculateAndPrintMetrics(s *Scheduler, tasks []*Task) {

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
