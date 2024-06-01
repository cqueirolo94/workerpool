package workerpool

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
	// AddTask adds a task to the task queue to be processed later. Returns a boolean indicating if the addition was a success.
	AddTask(Task) bool
	// Close terminates the execution of the workers in an orderly manner. Returns a boolean indicating if the closing was a success.
	Close() bool
	// IsClosed returns a boolean indicating if the Workerpool is open to accept and proccess tasks.
	IsClosed() bool
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
	// closed indicates if the workerpool can process more tasks.
	closed bool
}

func (wp *workerpool) AddTask(task Task) bool {
	if wp.closed {
		return false
	}

	wp.tasks <- task
	return true
}

func (wp *workerpool) Close() bool {
	if wp.closed {
		return false
	}

	for range wp.workers {
		wp.stop <- true
	}
	for range wp.workers {
		<-wp.stopAck
	}

	wp.closed = true
	return true
}

func (wp *workerpool) IsClosed() bool {
	return wp.closed
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
		false,
	}
}

// worker is the entity that processes tasks.
type worker struct {
	// id is an identifier for a worker.
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
