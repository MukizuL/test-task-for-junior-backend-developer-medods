package task

import (
	"context"
	"encoding/json"
	"time"

	"example.com/taskservice/internal/domain/taskdomain"
	"example.com/taskservice/internal/types"
)

//go:generate mockgen -source=ports.go -destination=mocks/ports.go -package=mocks -mock_names=Repository=MockRepo,Usecase=MockUsecase

type Repository interface {
	Create(ctx context.Context, task *taskdomain.Task) (*taskdomain.Task, error)
	GetByID(ctx context.Context, id int64) (*taskdomain.Task, error)
	Update(ctx context.Context, task *taskdomain.Task) (*taskdomain.Task, error)
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context) ([]taskdomain.Task, error)
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
	Title       string `validate:"required"`
	Description string
	Status      taskdomain.Status `validate:"required,oneof=new in_progress done"`
	DueDate     time.Time         `validate:"required,gt"`
}

type UpdateInput struct {
	RecurringTaskID int64
	Title           string `validate:"required"`
	Description     string
	Status          taskdomain.Status `validate:"required,oneof=new in_progress done"`
	DueDate         time.Time         `validate:"required,gt"`
}

type CreateRecurringInput struct {
	Title       string `validate:"required"`
	Description string
	Type        types.RecurrenceType `validate:"required,oneof=interval odd_days even_days yearly_on"`
	Config      json.RawMessage      `validate:"required"`
	StartDate   time.Time            `validate:"required"`
	EndDate     *time.Time           `validate:"gtfield=StartDate"`
}

type UpdateRecurringInput struct {
	Title       string `validate:"required"`
	Description string
	Type        types.RecurrenceType `validate:"required,oneof=interval odd_days even_days yearly_on"`
	Config      json.RawMessage      `validate:"required"`
	StartDate   time.Time            `validate:"required"`
	EndDate     *time.Time           `validate:"gtfield=StartDate"`
}
