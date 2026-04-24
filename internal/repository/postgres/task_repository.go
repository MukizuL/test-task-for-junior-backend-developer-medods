package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	taskdomain "example.com/taskservice/internal/domain/task"
)

type Repository struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) Create(ctx context.Context, task *taskdomain.Task) (*taskdomain.Task, error) {
	const query = `
		INSERT INTO tasks (rec_task_id, title, description, status, due_date, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7) ON CONFLICT DO NOTHING
		RETURNING id, rec_task_id, title, description, status, due_date, created_at, updated_at
	`

	row := r.pool.QueryRow(ctx, query, task.RecurringTaskID, task.Title, task.Description, task.Status, task.DueDate, task.CreatedAt, task.UpdatedAt)
	created, err := scanTask(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return created, nil
}

func (r *Repository) GetByID(ctx context.Context, id int64) (*taskdomain.Task, error) {
	const query = `
		SELECT id, title, description, status, due_date, created_at, updated_at
		FROM tasks
		WHERE id = $1
	`

	row := r.pool.QueryRow(ctx, query, id)
	found, err := scanTask(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, taskdomain.ErrNotFound
		}

		return nil, err
	}

	return found, nil
}

func (r *Repository) Update(ctx context.Context, task *taskdomain.Task) (*taskdomain.Task, error) {
	const query = `
		UPDATE tasks
		SET title = $1,
			description = $2,
			status = $3,
			due_date = $4,
			updated_at = $5
		WHERE id = $6
		RETURNING id, title, description, status, due_date, created_at, updated_at
	`

	row := r.pool.QueryRow(ctx, query, task.Title, task.Description, task.Status, task.DueDate, task.UpdatedAt, task.ID)
	updated, err := scanTask(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, taskdomain.ErrNotFound
		}

		return nil, err
	}

	return updated, nil
}

func (r *Repository) Delete(ctx context.Context, id int64) error {
	const query = `DELETE FROM tasks WHERE id = $1`

	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return taskdomain.ErrNotFound
	}

	return nil
}

func (r *Repository) List(ctx context.Context) ([]taskdomain.Task, error) {
	const query = `
		SELECT id, title, description, status, due_date, created_at, updated_at
		FROM tasks
		ORDER BY id DESC
	`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := make([]taskdomain.Task, 0)
	for rows.Next() {
		task, err := scanTask(rows)
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

func (r *Repository) GetDueRecurringTasks(ctx context.Context) ([]taskdomain.RecurringTask, error) {
	now := time.Now().UTC()

	const query = `
		SELECT id, title, description, frequency, interval, start_date, end_date, last_run_at, created_at
		FROM recurring_tasks
		WHERE start_date <= $1 AND (end_date IS NULL OR end_date >= $1)
	`

	rows, err := r.pool.Query(ctx, query, now)
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

func (r *Repository) UpdateLastRunAt(ctx context.Context, id int64, next time.Time) error {
	const query = `
		UPDATE recurring_tasks
		SET last_run_at = $1
		WHERE id = $2 AND (last_run_at IS NULL OR last_run_at < $1)
	`

	_, err := r.pool.Exec(ctx, query, next, id)
	return err
}

func (r *Repository) CreateRecurringTask(ctx context.Context, task *taskdomain.RecurringTask) (*taskdomain.RecurringTask, error) {
	const query = `
		INSERT INTO recurring_tasks (title,	description, frequency,	interval, start_date, end_date,	last_run_at, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
		RETURNING id, title, description, frequency, interval, start_date, end_date, last_run_at, created_at, updated_at
	`

	row := r.pool.QueryRow(
		ctx,
		query,
		task.Title,
		task.Description,
		task.Frequency,
		task.Interval,
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
	const query = `
		SELECT id, title, description, frequency, interval, start_date, end_date, last_run_at, created_at, updated_at
		FROM recurring_tasks
		WHERE id = $1
	`

	row := r.pool.QueryRow(ctx, query, id)
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
	const query = `
		UPDATE recurring_tasks
		SET title = $1,
			description = $2,
			frequency = $3,
			interval = $4,
			start_date = $5,
			end_date = $6,
			updated_at = $7
		WHERE id = $8
		RETURNING id, title, description, frequency, interval, start_date, end_date, last_run_at, created_at, updated_at
	`

	row := r.pool.QueryRow(ctx, query, task.Title, task.Description, task.Frequency, task.Interval, task.StartDate, task.EndDate, task.UpdatedAt, task.ID)
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
	const query = `DELETE FROM recurring_tasks WHERE id = $1`

	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return taskdomain.ErrNotFound
	}

	return nil
}

func (r *Repository) ListRecurringTasks(ctx context.Context) ([]taskdomain.RecurringTask, error) {
	const query = `
		SELECT id, title, description, frequency, interval, start_date, end_date, last_run_at, created_at, updated_at
		FROM recurring_tasks
		ORDER BY id DESC
	`

	rows, err := r.pool.Query(ctx, query)
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
