package worker

import (
	"context"
	"time"

	"example.com/taskservice/internal/domain/taskdomain"
)

//go:generate mockgen -source=ports.go -destination=mocks/ports.go -package=mocksWorker -mock_names=Repository=MockRepo

type Repository interface {
	GetDueRecurringTasks(ctx context.Context) ([]taskdomain.RecurringTask, error)
	CreateAndUpdateLastRunAt(ctx context.Context, task *taskdomain.Task, next time.Time) error
}
