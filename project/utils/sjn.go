package utils

import (
	"fmt"
	"sort"

	"time"
)

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
					fmt.Println("Task ", task.ID, "-----> IO queue")
					s.IOQueue <- task
				} else if task.Resources > 0 {
					task.State = WaitingRes
					// fmt.Print("\n\n\nResQ_________________> ", task)
					fmt.Println("Task ", task.ID, "-----> Resource queue")
					s.ResQueue <- task
				} else {

					task.State = Running
					fmt.Println("Task ", task.ID, "-----> running")
					time.Sleep(time.Duration(task.BurstTime) * time.Millisecond)

					s.Res += task.ResourcesAllocated // retrieving resources back
					task.ResourcesAllocated = 0

					task.BurstTime = 0
					task.State = Completed

					task.CompletionTime = time.Now()

					fmt.Println("Task ", task.ID, "-----> completed")

					s.mu.Lock()
					s.Completed = append(s.Completed, task)
					s.mu.Unlock()
				}
			}
		}
	}
}
