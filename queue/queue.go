package queue

import (
	"sync"
)

// TaskType represents the type of a task.
type TaskType int

const (
	TaskTypeDownload TaskType = iota
	TaskTypeConvert
)

// Task represents a unit of work.
type Task struct {
	URL      string
	Type     TaskType
	FilePath string // Used for PDF conversion tasks
}

// TaskQueue is a concurrent queue for tasks.
type TaskQueue struct {
	tasks []Task
	mu    sync.Mutex
	cond  *sync.Cond
}

// NewTaskQueue creates a new TaskQueue.
func NewTaskQueue() *TaskQueue {
	q := &TaskQueue{}
	q.cond = sync.NewCond(&q.mu)
	return q
}

// Enqueue adds a task to the queue.
func (q *TaskQueue) Enqueue(task Task) {
	q.mu.Lock()
	q.tasks = append(q.tasks, task)
	q.cond.Signal() // Signal a waiting worker
	q.mu.Unlock()
}

// Dequeue removes and returns a task from the queue.
func (q *TaskQueue) Dequeue() (Task, bool) {
	q.mu.Lock()
	defer q.mu.Unlock()
	for len(q.tasks) == 0 {
		q.cond.Wait() // Wait if the queue is empty
	}
	task := q.tasks[0]
	q.tasks = q.tasks[1:]
	return task, true
}

// IsEmpty returns true if the queue is empty.
func (q *TaskQueue) IsEmpty() bool {
	q.mu.Lock()
	defer q.mu.Unlock()
	return len(q.tasks) == 0
}
