package main

import (
	"fmt"
	"time"
)

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

func main() {
	fmt.Println("Hello World!")
	tasks := []Task{
		&EmailTask{
			Address: "address_1@email.com",
			Header:  "header 1",
			Body:    "body 1",
		},
		&ImageTask{
			Size:   50,
			Name:   "image_1",
			Format: "jpg",
		},
		&EmailTask{
			Address: "address_2@email.com",
			Header:  "header 2",
			Body:    "body 2",
		},
		&ImageTask{
			Size:   40,
			Name:   "image_2",
			Format: "png",
		},
		&EmailTask{
			Address: "address_3@email.com",
			Header:  "header 3",
			Body:    "body 3",
		},
		&ImageTask{
			Size:   20,
			Name:   "image_3",
			Format: "jpeg",
		},
		&EmailTask{
			Address: "address_4@email.com",
			Header:  "header 4",
			Body:    "body 4",
		},
	}

	wp := NewUnbufferedWorkerpool(9)
	defer wp.Close()
	for _, t := range tasks {
		wp.AddTask(t)
	}
}
