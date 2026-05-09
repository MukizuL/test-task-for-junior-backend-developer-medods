package handlers

import (
	"encoding/json"
	"time"

	"example.com/taskservice/internal/domain/taskdomain"
	"example.com/taskservice/internal/types"
)

type taskMutationDTO struct {
	Title       string       `json:"title"`
	Description string       `json:"description"`
	Status      types.Status `json:"status"`
	DueDate     time.Time    `json:"due_date"`
}

type taskDTO struct {
	ID          int64        `json:"id"`
	Title       string       `json:"title"`
	Description string       `json:"description"`
	Status      types.Status `json:"status"`
	DueDate     time.Time    `json:"due_date"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
}

// newTaskDTO returns http response-ready struct
func newTaskDTO(task *taskdomain.Task) taskDTO {
	return taskDTO{
		ID:          task.ID,
		Title:       task.Title,
		Description: task.Description,
		Status:      task.Status,
		DueDate:     task.DueDate,
		CreatedAt:   task.CreatedAt,
		UpdatedAt:   task.UpdatedAt,
	}
}

type recurringTaskMutationDTO struct {
	Title       string               `json:"title"`
	Description string               `json:"description"`
	Type        types.RecurrenceType `json:"type"`
	Config      json.RawMessage      `json:"config"`
	StartDate   time.Time            `json:"start_date"`
	EndDate     *time.Time           `json:"end_date"`
}

type recurringTaskDTO struct {
	ID          int64                `json:"id"`
	Title       string               `json:"title"`
	Description string               `json:"description"`
	Type        types.RecurrenceType `json:"type"`
	Config      json.RawMessage      `json:"config"`
	StartDate   time.Time            `json:"start_date"`
	EndDate     *time.Time           `json:"end_date"`
	LastRunAt   *time.Time           `json:"last_run_at"`
	CreatedAt   time.Time            `json:"created_at"`
	UpdatedAt   time.Time            `json:"updated_at"`
}

// newRecurringTaskDTO returns http response-ready struct
func newRecurringTaskDTO(task *taskdomain.RecurringTask) recurringTaskDTO {
	return recurringTaskDTO{
		ID:          task.ID,
		Title:       task.Title,
		Description: task.Description,
		Type:        task.Type,
		Config:      task.Config,
		StartDate:   task.StartDate,
		EndDate:     task.EndDate,
		LastRunAt:   task.LastRunAt,
		CreatedAt:   task.CreatedAt,
		UpdatedAt:   task.UpdatedAt,
	}
}
