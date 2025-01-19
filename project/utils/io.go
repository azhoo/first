package utils

import (
	"fmt"

	"time"
)

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
