package main

import (
    "context"
    "log"
    "os"
    "os/signal"
    
    "golang-dtqs/internal/queue"
    "golang-dtqs/internal/worker"
)

func main() {
    ctx := context.Background()
    
    // connect to redis
    redisURL := os.Getenv("REDIS_URL")
    if redisURL == "" {
        redisURL = "redis://localhost:6379"
    }
    
    q, err := queue.NewRedisQueue(ctx, redisURL)
    if err != nil {
        log.Fatalf("failed to connect to redis: %v", err)
    }
    defer q.Close()
    
    // create worker
    w := worker.NewWorker(1, q)
    
    // start processing
    ctx, cancel := context.WithCancel(ctx)
    go w.Start(ctx)
    
    // wait for interrupt
    c := make(chan os.Signal, 1)
    signal.Notify(c, os.Interrupt)
    <-c
    
    log.Println("shutting down...")
    cancel()
}
