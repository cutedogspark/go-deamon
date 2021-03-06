package gworker

import (
	"fmt"
	//"time"
)

type Worker struct {
	ID          int
	Group       string
	Name        string
	Work        chan Job
	WorkerQueue chan chan Job
	QuitChan    chan bool
}

func NewWorker(id int, jobQueue chan chan Job) Worker {
	worker := Worker{
		ID:          id,
		Work:        make(chan Job),
		WorkerQueue: jobQueue,
		QuitChan:    make(chan bool)}

	return worker
}

func (w *Worker) Start() {
	go func() {
		for {
			select {
			case w.WorkerQueue <- w.Work:
				work := <-w.Work
				fmt.Printf("[Worker-%d]: Running cmd : %s!\n", w.ID, work.context)
				work.Start()
			case <-w.QuitChan:
				return
			}
		}
	}()
}

func (w *Worker) Stop() {
	go func() {
		w.QuitChan <- true
	}()
}
