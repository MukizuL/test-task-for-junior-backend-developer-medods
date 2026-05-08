// Package worker implements subroutine which checks database for tasks due for creation and then creates them.
package worker

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"example.com/taskservice/internal/clock"
	"example.com/taskservice/internal/domain/taskdomain"
	taskUsecase "example.com/taskservice/internal/usecase/task"
)

const (
	FrequencyDaily   taskdomain.Frequency = "daily"
	FrequencyWeekly  taskdomain.Frequency = "weekly"
	FrequencyMonthly taskdomain.Frequency = "monthly"
	FrequencyYearly  taskdomain.Frequency = "yearly"
)

var (
	errGettingRecurringTask = errors.New("error getting due recurring tasks")
	errUnknownFrequency     = errors.New("error unknown frequency")
)

type Worker struct {
	repo   taskUsecase.Repository
	clock  clock.Clock
	logger *slog.Logger
}

func New(repo taskUsecase.Repository, clock clock.Clock, logger *slog.Logger) *Worker {
	return &Worker{
		repo:   repo,
		clock:  clock,
		logger: logger,
	}
}

func (w *Worker) Run(ctx context.Context) error {
	ticker := time.NewTicker(5 * time.Second)
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
		yes, err := isDue(rt, now)
		if err != nil {
			w.logger.Error("error checking due recurring task", "error", err, "task_id", rt.ID)
			continue
		}

		if !yes {
			continue
		}

		next, err := nextRun(rt)
		if err != nil {
			w.logger.Error("error checking next to run", "error", err, "task_id", rt.ID)
			continue
		}

		task := taskdomain.Task{
			RecurringTaskID: &rt.ID,
			Title:           rt.Title,
			Description:     rt.Description,
			Status:          taskdomain.StatusNew,
			DueDate:         next,
			CreatedAt:       now,
			UpdatedAt:       now,
		}

		err = w.repo.CreateAndUpdateLastRunAt(ctx, &task, next)
		if err != nil {
			w.logger.Error("error updating last run", "error", err, "task", task)
			continue
		}
	}

	return nil
}

// isDue reports whether a task should be created
func isDue(rt taskdomain.RecurringTask, now time.Time) (bool, error) {
	if rt.LastRunAt == nil {
		return !rt.StartDate.After(now), nil
	}

	next, err := calculateNext(rt)
	if err != nil {
		return false, err
	}
	return !next.After(now), nil
}

func nextRun(rt taskdomain.RecurringTask) (time.Time, error) {
	if rt.LastRunAt == nil {
		return rt.StartDate, nil
	}
	return calculateNext(rt)
}

// calculateNext returns time the next task should be created at
func calculateNext(rt taskdomain.RecurringTask) (time.Time, error) {
	last := *rt.LastRunAt

	switch rt.Frequency {
	case FrequencyDaily:
		return last.AddDate(0, 0, rt.Interval), nil
	case FrequencyWeekly:
		return last.AddDate(0, 0, 7*rt.Interval), nil
	case FrequencyMonthly:
		return last.AddDate(0, rt.Interval, 0), nil
	case FrequencyYearly:
		return last.AddDate(rt.Interval, 0, 0), nil
	default:
		return time.Time{}, fmt.Errorf("%w: %s", errUnknownFrequency, rt.Frequency)
	}
}
