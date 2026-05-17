package postgres

import (
	"context"
	"errors"
	"time"

	"example.com/taskservice/internal/domain/taskdomain"
	"github.com/jackc/pgx/v5"
)

func (r *Repository) GetDueRecurringTasks(ctx context.Context) ([]taskdomain.RecurringTask, error) {
	ctxPG, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	now := time.Now().UTC()

	const query = `
		SELECT id, title, description, type, config, start_date, end_date, last_run_at, created_at, updated_at
		FROM recurring_tasks
		WHERE (last_run_at IS NULL OR last_run_at < $1) AND (end_date IS NULL OR end_date >= $1)
	`

	rows, err := r.pool.Query(ctxPG, query, now)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []taskdomain.RecurringTask

	for rows.Next() {
		rt, err := scanRecurringTask(rows)
		if err != nil {
			return nil, err
		}

		result = append(result, *rt)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func (r *Repository) CreateRecurringTask(ctx context.Context, task *taskdomain.RecurringTask) (*taskdomain.RecurringTask, error) {
	ctxPG, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	const query = `
		INSERT INTO recurring_tasks (title,	description, type, config, start_date, end_date, last_run_at, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
		RETURNING id, title, description, type,	config, start_date, end_date, last_run_at, created_at, updated_at
	`

	row := r.pool.QueryRow(
		ctxPG,
		query,
		task.Title,
		task.Description,
		task.Type,
		task.Config,
		task.StartDate,
		task.EndDate,
		nil,
		task.CreatedAt,
		task.UpdatedAt,
	)

	created, err := scanRecurringTask(row)
	if err != nil {
		return nil, err
	}

	return created, nil
}

func (r *Repository) GetRecurringTaskByID(ctx context.Context, id int64) (*taskdomain.RecurringTask, error) {
	ctxPG, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	const query = `
		SELECT id, title, description, type, config, start_date, end_date, last_run_at, created_at, updated_at
		FROM recurring_tasks
		WHERE id = $1
	`

	row := r.pool.QueryRow(ctxPG, query, id)
	found, err := scanRecurringTask(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, taskdomain.ErrNotFound
		}

		return nil, err
	}

	return found, nil
}

func (r *Repository) UpdateRecurringTask(ctx context.Context, task *taskdomain.RecurringTask) (*taskdomain.RecurringTask, error) {
	ctxPG, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	const query = `
		UPDATE recurring_tasks
		SET title = $1,
			description = $2,
			type = $3,
			config = $4,
			start_date = $5,
			end_date = $6,
			updated_at = $7
		WHERE id = $8
		RETURNING id, title, description, type, config, start_date, end_date, last_run_at, created_at, updated_at
	`

	row := r.pool.QueryRow(ctxPG, query, task.Title, task.Description, task.Type, task.Config, task.StartDate, task.EndDate, task.UpdatedAt, task.ID)
	updated, err := scanRecurringTask(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, taskdomain.ErrNotFound
		}

		return nil, err
	}

	return updated, nil
}

func (r *Repository) DeleteRecurringTask(ctx context.Context, id int64) error {
	ctxPG, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	const query = `DELETE FROM recurring_tasks WHERE id = $1`

	result, err := r.pool.Exec(ctxPG, query, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return taskdomain.ErrNotFound
	}

	return nil
}

func (r *Repository) ListRecurringTasks(ctx context.Context) ([]taskdomain.RecurringTask, error) {
	ctxPG, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	const query = `
		SELECT id, title, description, type, config, start_date, end_date, last_run_at, created_at, updated_at
		FROM recurring_tasks
		ORDER BY id DESC
	`

	rows, err := r.pool.Query(ctxPG, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := make([]taskdomain.RecurringTask, 0)
	for rows.Next() {
		task, err := scanRecurringTask(rows)
		if err != nil {
			return nil, err
		}

		tasks = append(tasks, *task)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}
