package main

import (
	"fmt"
	"time"
)

type Task interface {
	Process()
}

type EmailTask struct {
	Address string
	Header  string
	Body    string
}

func (e *EmailTask) Process() {
	time.Sleep(2 * time.Second)
	fmt.Printf("Email process - Address: %s, Header: %s, Body: %s\n", e.Address, e.Header, e.Body)
}

type ImageTask struct {
	Size   int
	Name   string
	Format string
}

func (i *ImageTask) Process() {
	time.Sleep(4 * time.Second)
	fmt.Printf("Image process - Name: %s, Format: %s, Size: %d\n", i.Name, i.Format, i.Size)
}

type Workerpool struct {
	workers []*worker
	Tasks   chan Task
	Stop    chan bool
	StopAck chan bool
}

func (wp *Workerpool) AddTask(task Task) {
	wp.Tasks <- task
}

func (wp *Workerpool) Close() {
	for range wp.workers {
		wp.Stop <- true
	}
	for range wp.workers {
		<-wp.StopAck
	}
}

func NewBufferedWorkerpool(n uint) *Workerpool {
	tasks := make(chan Task, n)
	return newWorkerpool(n, tasks)
}

func NewUnbufferedWorkerpool(n uint) *Workerpool {
	tasks := make(chan Task)
	return newWorkerpool(n, tasks)
}

func newWorkerpool(n uint, tasks chan Task) *Workerpool {
	workers := make([]*worker, n)
	stop := make(chan bool)
	stopAck := make(chan bool)

	for i := range n {
		workers[i] = newWorker(i)
		go workers[i].execute(tasks, stop, stopAck)
	}

	return &Workerpool{
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
