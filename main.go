package main

import (
	"fmt"
)

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

	wp := NewUnbufferedWorkerpool(7)
	for _, t := range tasks {
		wp.AddTask(t)
	}
	defer wp.Close()
}
