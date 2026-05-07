package taskdomain

import "time"

type Status string

const (
	StatusNew        Status = "new"
	StatusInProgress Status = "in_progress"
	StatusDone       Status = "done"
)

type Task struct {
	ID              int64
	RecurringTaskID *int64
	Title           string
	Description     string
	Status          Status
	DueDate         time.Time
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

func (s Status) Valid() bool {
	switch s {
	case StatusNew, StatusInProgress, StatusDone:
		return true
	default:
		return false
	}
}
