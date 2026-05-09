package task

import (
	"time"

	"example.com/taskservice/cmd/worker"
	"example.com/taskservice/internal/types"
	"github.com/go-playground/validator/v10"
)

type Service struct {
	repo       Repository
	validate   *validator.Validate
	schedulers map[types.RecurrenceType]worker.Scheduler
	now        func() time.Time
}

func NewService(repo Repository, validate *validator.Validate, schedulers map[types.RecurrenceType]worker.Scheduler) *Service {
	return &Service{
		repo:       repo,
		validate:   validate,
		schedulers: schedulers,
		now:        func() time.Time { return time.Now().UTC() },
	}
}
