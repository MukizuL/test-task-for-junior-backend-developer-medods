package task

import (
	"context"
	"time"

	"example.com/taskservice/internal/domain/taskdomain"
)

//go:generate mockgen -source=ports.go -destination=mocks/ports.go -package=mocks -mock_names=Repository=MockRepo,Usecase=MockUsecase

type Repository interface {
	Create(ctx context.Context, task *taskdomain.Task) (*taskdomain.Task, error)
	GetByID(ctx context.Context, id int64) (*taskdomain.Task, error)
	Update(ctx context.Context, task *taskdomain.Task) (*taskdomain.Task, error)
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context) ([]taskdomain.Task, error)
	GetDueRecurringTasks(ctx context.Context) ([]taskdomain.RecurringTask, error)
	CreateRecurringTask(ctx context.Context, task *taskdomain.RecurringTask) (*taskdomain.RecurringTask, error)
	GetRecurringTaskByID(ctx context.Context, id int64) (*taskdomain.RecurringTask, error)
	UpdateRecurringTask(ctx context.Context, task *taskdomain.RecurringTask) (*taskdomain.RecurringTask, error)
	DeleteRecurringTask(ctx context.Context, id int64) error
	ListRecurringTasks(ctx context.Context) ([]taskdomain.RecurringTask, error)
	CreateAndUpdateLastRunAt(ctx context.Context, task *taskdomain.Task, next time.Time) error
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
	Title           string
	Description     string
	Status          taskdomain.Status
	DueDate         time.Time
}

type CreateRecurringInput struct {
	Title       string `validate:"required"`
	Description string
	Frequency   taskdomain.Frequency `validate:"required,oneof=daily weekly monthly yearly"`
	Interval    int                  `validate:"required,min=1,max=365"`
	StartDate   time.Time            `validate:"required"`
	EndDate     *time.Time           `validate:"gtfield=StartDate"`
}

type UpdateRecurringInput struct {
	Title       string
	Description string
	Frequency   taskdomain.Frequency
	Interval    int
	StartDate   time.Time
	EndDate     *time.Time
}
