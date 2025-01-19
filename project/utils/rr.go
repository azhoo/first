package utils

import (
	"fmt"

	"time"
)

func RoundRobin(s *Scheduler) {
	if len(s.ReadyQueue) > 0 {
		for {
			task := <-s.ReadyQueue

			fmt.Print("\n.                                res = ", s.Res, "\n")
			fmt.Print(task, "1\n")

			if task.IOTime > 0 {//if has io handle io first
				
				
				task.State = WaitingIO
				s.IOQueue <- task

				
			} else if task.Resources > 0 {//if need resources handle them first
				
				
				task.State = WaitingRes
				fmt.Print("\n\n\nResQ_________________> ", task)
				s.ResQueue <- task


			} else if task.BurstTime > s.TimeSlice {//if longer than time slice

				// s.res += task.ResourcesAllocated
				// task.ResourcesAllocated = 0

				task.State = Running
				time.Sleep(time.Duration(s.TimeSlice) * time.Millisecond)
				task.BurstTime -= s.TimeSlice

				task.State = Ready

				s.ReadyQueue <- task

				fmt.Print(task, "2\n")

			} else {  // if shorter than time slice

				task.State = Running
				time.Sleep(time.Duration(task.BurstTime) * time.Millisecond)

				task.BurstTime = 0
				task.State = Completed

				task.CompletionTime = time.Now()

				s.Res += task.ResourcesAllocated // retrieving resources back
				task.ResourcesAllocated = 0

				fmt.Print(task, "3\n")

				s.mu.Lock()
				s.Completed = append(s.Completed, task)
				s.mu.Unlock()
			}
		}
	}

}
