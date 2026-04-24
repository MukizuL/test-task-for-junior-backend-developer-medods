package task

import "time"

type Status string

const (
	StatusNew        Status = "new"
	StatusInProgress Status = "in_progress"
	StatusDone       Status = "done"
)

type Task struct {
	ID              int64     `json:"id"`
	RecurringTaskID *int64    `json:"recurring_task_id"`
	Title           string    `json:"title"`
	Description     string    `json:"description"`
	Status          Status    `json:"status"`
	DueDate         time.Time `json:"due_date"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

func (s Status) Valid() bool {
	switch s {
	case StatusNew, StatusInProgress, StatusDone:
		return true
	default:
		return false
	}
}
