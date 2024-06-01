package main

import (
	"fmt"
)

// Task is the interface that has to be implemented by the entities that will be processed by the Workerpool.
type Task interface {
	// Process is the method that the workers will execute.
	Process()
}

// Workerpool is an entity that holds a list of workers. These workers will run in the background processing tasks from a channel.
type Workerpool interface {
	// AddTask adds a task to the task queue to be processed later.
	AddTask(Task)
	// Close terminates the execution of the workers in an orderly manner.
	Close()
}

// workerpool implements Workerpool.
type workerpool struct {
	// workers contains the full list of workers that will process tasks.
	workers []*worker
	// tasks is the channel that will hold the Task's to process.
	tasks chan Task
	// stop is the channel that will signal the workers to stop listening for tasks and terminating the execution.
	stop chan bool
	// stopAck is the channel that the workers will use to acknowledge that the stop signal was received.
	stopAck chan bool
}

func (wp *workerpool) AddTask(task Task) {
	wp.tasks <- task
}

func (wp *workerpool) Close() {
	for range wp.workers {
		wp.stop <- true
	}
	for range wp.workers {
		<-wp.stopAck
	}
}

// NewBufferedWorkerpool returns a new instance of workerpool, making the Task channel buffered with size n.
func NewBufferedWorkerpool(n uint) Workerpool {
	tasks := make(chan Task, n)
	return newWorkerpool(n, tasks)
}

// NewUnbufferedWorkerpool returns a new instance of workerpool, making the Task channel unbuffered.
func NewUnbufferedWorkerpool(n uint) Workerpool {
	tasks := make(chan Task)
	return newWorkerpool(n, tasks)
}

func newWorkerpool(n uint, tasks chan Task) Workerpool {
	workers := make([]*worker, n)
	stop := make(chan bool)
	stopAck := make(chan bool)

	for i := range n {
		workers[i] = newWorker(i)
		go workers[i].execute(tasks, stop, stopAck)
	}

	return &workerpool{
		workers,
		tasks,
		stop,
		stopAck,
	}
}

type worker struct {
	id uint
}

func (w *worker) execute(tasks chan Task, stop chan bool, stopAck chan bool) {
	for {
		select {
		case task := <-tasks:
			fmt.Printf("worker %d received task, starting...\n", w.id)
			task.Process()
			fmt.Printf("worker %d finished task, waiting for more tasks...\n", w.id)
		case <-stop:
			fmt.Printf("close signal detected, worker %d going home\n", w.id)
			stopAck <- true
			return
		}
	}
}

func newWorker(id uint) *worker {
	return &worker{id: id}
}
