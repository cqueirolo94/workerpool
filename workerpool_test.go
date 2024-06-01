package workerpool_test

import (
	"testing"

	"github.com/cqueirolo94/workerpool"
	"github.com/stretchr/testify/assert"
)

type TaskTest struct {
	Id      int
	Results chan int
}

func (tt *TaskTest) Process() {
	tt.Results <- tt.Id
}

func newTaskTestWith(results chan int, id int) *TaskTest {
	return &TaskTest{
		Results: results,
		Id:      id,
	}
}

// Given: a closed workerpool.
// When: closes the workerpool.
// Then: returns false.
func TestWorkerpool_Close_with_closed_workerpool(t *testing.T) {
	t.Run("buffered", func(t *testing.T) {
		// Arrange
		wp := workerpool.NewBufferedWorkerpool(1)
		wp.Close()

		// Act
		ok := wp.Close()

		// Assert
		assert.False(t, ok)
		assert.True(t, wp.IsClosed())
	})

	t.Run("unbuffered", func(t *testing.T) {
		// Arrange
		wp := workerpool.NewUnbufferedWorkerpool(1)
		wp.Close()

		// Act
		ok := wp.Close()

		// Assert
		assert.False(t, ok)
		assert.True(t, wp.IsClosed())
	})
}

// Given: a closed workerpool.
// When: adds a task.
// Then: returns false and no tasks are added.
func TestWorkerpool_AddTask_with_closed_workerpool(t *testing.T) {
	t.Run("buffered", func(t *testing.T) {
		// Arrange
		wp := workerpool.NewBufferedWorkerpool(1)
		wp.Close()

		// Act
		ok := wp.AddTask(newTaskTestWith(nil, 0))

		// Assert
		assert.False(t, ok)
	})

	t.Run("unbuffered", func(t *testing.T) {
		// Arrange
		wp := workerpool.NewUnbufferedWorkerpool(1)
		wp.Close()

		// Act
		ok := wp.AddTask(newTaskTestWith(nil, 0))

		// Assert
		assert.False(t, ok)
	})
}

// Given: a workerpool with one worker.
// When: process a list of tasks.
// Then: all tasks are processed in order.
func TestWorkerpool_with_all_tasks_finished_by_one_worker(t *testing.T) {
	// Prearrange
	tsize := 10
	results := make(chan int, tsize)
	tasks := make([]workerpool.Task, tsize)
	for i := range tsize {
		tasks[i] = newTaskTestWith(results, i)
	}

	t.Run("unbuffered", func(t *testing.T) {
		// Arrange
		wp := workerpool.NewUnbufferedWorkerpool(1)

		// Act
		for _, t := range tasks {
			wp.AddTask(t)
		}
		wp.Close()

		// Assert
		for i := range tsize {
			assert.True(t, i == <-results)
		}
	})

	t.Run("buffered", func(t *testing.T) {
		// Arrange
		wp := workerpool.NewBufferedWorkerpool(1)

		// Act
		for _, t := range tasks {
			wp.AddTask(t)
		}
		wp.Close()

		// Assert
		for i := range tsize {
			assert.True(t, i == <-results)
		}
	})
}

// Given: a workerpool with several workers.
// When: process a list of tasks.
// Then: all tasks are processed.
func TestWorkerpool_with_all_tasks_finished_by_multiple_workers(t *testing.T) {
	// Prearrange
	tsize := 10
	results := make(chan int, tsize)
	tasks := make([]*TaskTest, tsize)
	for i := range tsize {
		tasks[i] = newTaskTestWith(results, i)
	}

	t.Run("unbuffered", func(t *testing.T) {
		// Arrange
		wp := workerpool.NewUnbufferedWorkerpool(20)

		// Act
		for _, t := range tasks {
			wp.AddTask(t)
		}
		wp.Close()

		// Assert
		for range tsize {
			id := <-results
			assert.Equal(t, tasks[id].Id, id)
		}
	})

	t.Run("buffered", func(t *testing.T) {
		// Arrange
		wp := workerpool.NewBufferedWorkerpool(20)

		// Act
		for _, t := range tasks {
			wp.AddTask(t)
		}
		wp.Close()

		// Assert
		for range tsize {
			id := <-results
			assert.Equal(t, tasks[id].Id, id)
		}
	})
}
