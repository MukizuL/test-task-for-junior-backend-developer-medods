package worker

import (
	"context"
	"time"

	"example.com/taskservice/internal/clock"
	"example.com/taskservice/internal/domain/taskdomain"
	taskUsecase "example.com/taskservice/internal/usecase/task"
)

type Worker struct {
	ctx   context.Context
	repo  taskUsecase.Repository
	clock clock.Clock
}

func New(ctx context.Context, repo taskUsecase.Repository, clock clock.Clock) *Worker {
	return &Worker{
		ctx:   ctx,
		repo:  repo,
		clock: clock,
	}
}

func (w *Worker) Run() error {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-w.ctx.Done():
			return w.ctx.Err()
		case <-ticker.C:
			err := w.RunOnce()
			if err != nil {
				return err
			}
		}
	}
}

func (w *Worker) RunOnce() error {
	now := w.clock.Now()

	recTasks, err := w.repo.GetDueRecurringTasks(w.ctx)
	if err != nil {
		return err
	}

	for _, rt := range recTasks {
		if !isDue(rt, now) {
			continue
		}

		next := nextRun(rt)

		task := taskdomain.Task{
			RecurringTaskID: &rt.ID,
			Title:           rt.Title,
			Description:     rt.Description,
			Status:          taskdomain.StatusNew,
			DueDate:         next,
			CreatedAt:       now,
			UpdatedAt:       now,
		}

		_, err = w.repo.Create(w.ctx, &task)
		if err != nil {
			return err
		}
		err = w.repo.UpdateLastRunAt(w.ctx, rt.ID, next)
		if err != nil {
			return err
		}
	}

	return nil
}

func isDue(rt taskdomain.RecurringTask, now time.Time) bool {
	if rt.LastRunAt == nil {
		return !rt.StartDate.After(now)
	}

	next := calculateNext(rt)
	return !next.After(now)
}

func nextRun(rt taskdomain.RecurringTask) time.Time {
	if rt.LastRunAt == nil {
		return rt.StartDate
	}
	return calculateNext(rt)
}

func calculateNext(rt taskdomain.RecurringTask) time.Time {
	last := *rt.LastRunAt

	switch rt.Frequency {
	case "daily":
		return last.AddDate(0, 0, rt.Interval)
	case "weekly":
		return last.AddDate(0, 0, 7*rt.Interval)
	case "monthly":
		return last.AddDate(0, rt.Interval, 0)
	case "yearly":
		return last.AddDate(rt.Interval, 0, 0)
	default:
		panic("unknown frequency")
	}
}
