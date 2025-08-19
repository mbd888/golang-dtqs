package main

import (
    "context"
    "encoding/json"
    "log"
    "net/http"
    "os"
    "strings"

    "golang-dtqs/internal/queue"
    "golang-dtqs/internal/task"
)

var q queue.Queue

func main() {
    ctx := context.Background()

    // Redis URL from env or default
    redisURL := os.Getenv("REDIS_URL")
    if redisURL == "" {
        redisURL = "redis://localhost:6379"
    }

    var err error
    q, err = queue.NewRedisQueue(ctx, redisURL)
    if err != nil {
        log.Fatalf("failed to connect to redis: %v", err)
    }

    // Handlers
    http.HandleFunc("/", rootHandler)
    http.HandleFunc("/health", healthHandler)
    http.HandleFunc("/tasks", loggingMiddleware(tasksHandler))
    http.HandleFunc("/tasks/", loggingMiddleware(taskHandler))

    // Port from env or default
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    log.Printf("Server starting on :%s", port)
    log.Fatal(http.ListenAndServe(":"+port, nil))
}

// Root handler for "/"
func rootHandler(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Golang DTQS API is running"))
}

// Health check endpoint
func healthHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// Task endpoint with method routing
func tasksHandler(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
    case http.MethodPost:
        createTask(w, r)
    default:
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
    }
}

// Get single task by ID
func taskHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }
    
    // Extract ID from path /tasks/{id}
    parts := strings.Split(r.URL.Path, "/")
    if len(parts) != 3 || parts[2] == "" {
        http.Error(w, "Invalid task ID", http.StatusBadRequest)
        return
    }
    
    taskID := parts[2]
    
    // Get task from queue
    t, err := q.Get(r.Context(), taskID)
    if err != nil {
        if err == queue.ErrTaskNotFound {
            http.Error(w, "Task not found", http.StatusNotFound)
        } else {
            log.Printf("Error getting task %s: %v", taskID, err)
            http.Error(w, "Internal server error", http.StatusInternalServerError)
        }
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(t)
}

// Create task handler
func createTask(w http.ResponseWriter, r *http.Request) {
    var req struct {
        Type    string                 `json:"type"`
        Payload map[string]interface{} `json:"payload"`
    }

    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request", http.StatusBadRequest)
        return
    }

    t := task.New(req.Type, req.Payload)

    if err := q.Enqueue(r.Context(), t); err != nil {
        log.Printf("Failed to enqueue task: %v", err)
        http.Error(w, "Failed to enqueue task", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(map[string]string{
        "id":     t.ID,
        "status": string(t.Status),
    })
}

// Simple logging middleware
func loggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        log.Printf("%s %s", r.Method, r.URL.Path)
        next(w, r)
    }
}