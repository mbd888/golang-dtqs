package queue

import (
    "context"
    "errors"
    
    "golang-dtqs/internal/task"
)

var (
    ErrQueueEmpty   = errors.New("queue is empty")
    ErrTaskNotFound = errors.New("task not found")
)

type Queue interface {
    Enqueue(ctx context.Context, t *task.Task) error
    Dequeue(ctx context.Context) (*task.Task, error)
    Get(ctx context.Context, taskID string) (*task.Task, error)
    Update(ctx context.Context, t *task.Task) error
}