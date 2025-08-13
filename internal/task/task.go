package task

import (
    "encoding/json"
    "time"
)

type Status string

const (
    StatusPending   Status = "pending"
    StatusRunning   Status = "running" 
    StatusCompleted Status = "completed"
    StatusFailed    Status = "failed"
)

type Task struct {
    ID        string                 `json:"id"`
    Type      string                 `json:"type"`
    Payload   map[string]interface{} `json:"payload"`
    Status    Status                 `json:"status"`
    CreatedAt time.Time              `json:"created_at"`
    UpdatedAt time.Time              `json:"updated_at"`
}

func New(taskType string, payload map[string]interface{}) *Task {
    now := time.Now()
    return &Task{
        ID:        generateID(), // TODO: implement proper ID generation
        Type:      taskType,
        Payload:   payload,
        Status:    StatusPending,
        CreatedAt: now,
        UpdatedAt: now,
    }
}

func generateID() string {
    // temporary - will add uuid later
    return time.Now().Format("20060102150405")
}

func (t *Task) Marshal() ([]byte, error) {
    return json.Marshal(t)
}

func Unmarshal(data []byte) (*Task, error) {
    var t Task
    err := json.Unmarshal(data, &t)
    return &t, err
}