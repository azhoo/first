package utils

import (
	"fmt"
	"sort"

	"time"
)

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
			if task.Resources <= s.Res {
				s.Res -= task.Resources
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
