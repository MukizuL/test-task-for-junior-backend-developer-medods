package taskdomain

import (
	"time"

	"example.com/taskservice/internal/types"
)

type Task struct {
	ID              int64
	RecurringTaskID *int64
	Title           string
	Description     string
	Status          types.Status
	DueDate         time.Time
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
