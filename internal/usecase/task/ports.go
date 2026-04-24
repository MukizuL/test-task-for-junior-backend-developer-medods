package task

import (
	"context"
	"time"

	taskdomain "example.com/taskservice/internal/domain/task"
)

type Repository interface {
	Create(ctx context.Context, task *taskdomain.Task) (*taskdomain.Task, error)
	GetByID(ctx context.Context, id int64) (*taskdomain.Task, error)
	Update(ctx context.Context, task *taskdomain.Task) (*taskdomain.Task, error)
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context) ([]taskdomain.Task, error)
	GetDueRecurringTasks(ctx context.Context) ([]taskdomain.RecurringTask, error)
	UpdateLastRunAt(ctx context.Context, id int64, next time.Time) error
	CreateRecurringTask(ctx context.Context, task *taskdomain.RecurringTask) (*taskdomain.RecurringTask, error)
	GetRecurringTaskByID(ctx context.Context, id int64) (*taskdomain.RecurringTask, error)
	UpdateRecurringTask(ctx context.Context, task *taskdomain.RecurringTask) (*taskdomain.RecurringTask, error)
	DeleteRecurringTask(ctx context.Context, id int64) error
	ListRecurringTasks(ctx context.Context) ([]taskdomain.RecurringTask, error)
}

type Usecase interface {
	Create(ctx context.Context, input CreateInput) (*taskdomain.Task, error)
	GetByID(ctx context.Context, id int64) (*taskdomain.Task, error)
	Update(ctx context.Context, id int64, input UpdateInput) (*taskdomain.Task, error)
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context) ([]taskdomain.Task, error)
	CreateRecurringTask(ctx context.Context, input CreateRecurringInput) (*taskdomain.RecurringTask, error)
	GetRecurringTaskByID(ctx context.Context, id int64) (*taskdomain.RecurringTask, error)
	UpdateRecurringTask(ctx context.Context, id int64, input UpdateRecurringInput) (*taskdomain.RecurringTask, error)
	DeleteRecurringTask(ctx context.Context, id int64) error
	ListRecurringTasks(ctx context.Context) ([]taskdomain.RecurringTask, error)
}

type CreateInput struct {
	Title       string
	Description string
	Status      taskdomain.Status
	DueDate     time.Time
}

type UpdateInput struct {
	RecurringTaskID int64
	Title           string
	Description     string
	Status          taskdomain.Status
	DueDate         time.Time
}

type CreateRecurringInput struct {
	Title       string
	Description string
	Frequency   taskdomain.Frequency
	Interval    int
	StartDate   time.Time
	EndDate     *time.Time
}

type UpdateRecurringInput struct {
	Title       string
	Description string
	Frequency   taskdomain.Frequency
	Interval    int
	StartDate   time.Time
	EndDate     *time.Time
}
