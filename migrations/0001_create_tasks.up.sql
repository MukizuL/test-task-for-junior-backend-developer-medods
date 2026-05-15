CREATE TYPE status AS ENUM ('new', 'in_progress', 'done');
CREATE TYPE rec_type AS ENUM ('interval', 'even_odd_days', 'specific_dates');

CREATE TABLE IF NOT EXISTS recurring_tasks (
    id BIGSERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    type rec_type NOT NULL,
    config JSONB NOT NULL,
    start_date TIMESTAMPTZ NOT NULL,
    end_date TIMESTAMPTZ,
    last_run_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS tasks (
    id BIGSERIAL PRIMARY KEY,
    rec_task_id BIGINT,
    title TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    status status NOT NULL,
    due_date TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    FOREIGN KEY (rec_task_id) REFERENCES recurring_tasks(id)
);

CREATE INDEX IF NOT EXISTS idx_tasks_status ON tasks (status);
CREATE UNIQUE INDEX IF NOT EXISTS uniq_task_instance ON tasks (rec_task_id, due_date);