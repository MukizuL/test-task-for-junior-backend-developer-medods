// Package worker implements subroutine which checks database for tasks due for creation and then creates them.
package worker

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"example.com/taskservice/internal/clock"
	"example.com/taskservice/internal/domain/taskdomain"
	"example.com/taskservice/internal/types"
)

var (
	errGettingRecurringTask = errors.New("error getting due recurring tasks")
	errUnknownFrequency     = errors.New("error unknown frequency")
)

type Worker struct {
	repo       Repository
	schedulers map[types.RecurrenceType]Scheduler
	clock      clock.Clock
	logger     *slog.Logger
}

func New(repo Repository, clock clock.Clock, logger *slog.Logger, schedulers map[types.RecurrenceType]Scheduler) *Worker {
	return &Worker{
		repo:       repo,
		schedulers: schedulers,
		clock:      clock,
		logger:     logger,
	}
}

func (w *Worker) Run(ctx context.Context, tick time.Duration) error {
	ticker := time.NewTicker(tick)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			err := w.RunOnce(ctx)
			if err != nil {
				w.logger.Error("error creating next tasks", "error", err)
				continue
			}
		}
	}
}

func (w *Worker) RunOnce(ctx context.Context) error {
	now := w.clock.Now()

	recTasks, err := w.repo.GetDueRecurringTasks(ctx)
	if err != nil {
		return errors.Join(errGettingRecurringTask, err)
	}

	for _, rt := range recTasks {
		scheduler, ok := w.schedulers[rt.Type]
		if !ok {
			w.logger.Error("unknown schedule type", "task_id", rt.ID)
			continue
		}

		next, yes, errLoop := scheduler.IsDue(rt, now)
		if errLoop != nil {
			w.logger.Error("failed computing next recurrence", "task_id", rt.ID, "error", errLoop)
			continue
		}

		if !yes {
			continue
		}

		task := taskdomain.Task{
			RecurringTaskID: &rt.ID,
			Title:           rt.Title,
			Description:     rt.Description,
			Status:          types.StatusNew,
			DueDate:         next,
			CreatedAt:       now,
			UpdatedAt:       now,
		}

		errLoop = w.repo.CreateAndUpdateLastRunAt(ctx, &task, next)
		if errLoop != nil {
			w.logger.Error("failed persisting recurring task", "task_id", rt.ID, "next", next, "error", errLoop)
			continue
		}
	}

	return nil
}
