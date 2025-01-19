package utils

import (
	"fmt"

	"time"
)

func FCFS(s *Scheduler) {

	for {
		if len(s.ReadyQueue) > 0 {
			task := <-s.ReadyQueue

			fmt.Print("\n.                                res = ", s.Res, "\n")
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

				s.Res += task.ResourcesAllocated // retrieving resources back
				task.ResourcesAllocated = 0

				task.BurstTime = 0
				task.State = Completed

				task.CompletionTime = time.Now()

				s.mu.Lock()
				s.Completed = append(s.Completed, task)
				s.mu.Unlock()

				fmt.Print(task, "3\n")
			}
		}
	}
}
