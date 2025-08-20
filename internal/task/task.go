package task

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Status string

const (
	StatusPending   Status = "pending"
	StatusRunning   Status = "running"
	StatusCompleted Status = "completed"
	StatusFailed    Status = "failed"
)

type Priority int

const (
	PriorityLow Priority = iota
	PriorityNormal
	PriorityHigh
	PriorityCritical
)

type Task struct {
	ID        string                 `json:"id"`
	Type      string                 `json:"type"`
	Payload   map[string]interface{} `json:"payload"`
	Status    Status                 `json:"status"`
	Priority  Priority               `json:"priority"`
	CreatedAt time.Time              `json:"created_at"`
	UpdatedAt time.Time              `json:"updated_at"`
}

func New(taskType string, payload map[string]interface{}) *Task {
	now := time.Now()
	return &Task{
		ID:        uuid.New().String(),
		Type:      taskType,
		Payload:   payload,
		Status:    StatusPending,
		Priority:  PriorityNormal,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func (t *Task) Marshal() ([]byte, error) {
	return json.Marshal(t)
}

func Unmarshal(data []byte) (*Task, error) {
	var t Task
	err := json.Unmarshal(data, &t)
	return &t, err
}
