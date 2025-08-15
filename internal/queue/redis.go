package queue

import (
    "context"
    "fmt"
    
    "github.com/redis/go-redis/v9"
    "golang-dtqs/internal/task"
)

const (
    queueKey   = "golang-dtqs:queue"
    taskPrefix = "golang-dtqs:task:"
)

type RedisQueue struct {
    client *redis.Client
}

func NewRedisQueue(ctx context.Context, redisURL string) (*RedisQueue, error) {
    opts, err := redis.ParseURL(redisURL)
    if err != nil {
        return nil, fmt.Errorf("invalid redis url: %w", err)
    }
    
    client := redis.NewClient(opts)
    
    // test connection
    if err := client.Ping(ctx).Err(); err != nil {
        return nil, fmt.Errorf("redis connection failed: %w", err)
    }
    
    return &RedisQueue{client: client}, nil
}

func (q *RedisQueue) Enqueue(ctx context.Context, t *task.Task) error {
    data, err := t.Marshal()
    if err != nil {
        return err
    }
    
    // store task data
    err = q.client.Set(ctx, taskPrefix+t.ID, data, 0).Err()
    if err != nil {
        return err
    }
    
    // add to queue
    return q.client.LPush(ctx, queueKey, t.ID).Err()
}

func (q *RedisQueue) Dequeue(ctx context.Context) (*task.Task, error) {
    // pop from queue
    taskID, err := q.client.RPop(ctx, queueKey).Result()
    if err != nil {
        if err == redis.Nil {
            return nil, ErrQueueEmpty
        }
        return nil, err
    }
    
    // get task data
    data, err := q.client.Get(ctx, taskPrefix+taskID).Bytes()
    if err != nil {
        return nil, err
    }
    
    return task.Unmarshal(data)
}

func (q *RedisQueue) Get(ctx context.Context, taskID string) (*task.Task, error) {
    data, err := q.client.Get(ctx, taskPrefix+taskID).Bytes()
    if err != nil {
        if err == redis.Nil {
            return nil, ErrTaskNotFound
        }
        return nil, err
    }
    
    return task.Unmarshal(data)
}

func (q *RedisQueue) Update(ctx context.Context, t *task.Task) error {
    data, err := t.Marshal()
    if err != nil {
        return err
    }
    
    return q.client.Set(ctx, taskPrefix+t.ID, data, 0).Err()
}

func (q *RedisQueue) Close() error {
    return q.client.Close()
}