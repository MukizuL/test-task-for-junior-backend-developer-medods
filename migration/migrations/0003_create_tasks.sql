-- +goose Up
-- +goose StatementBegin
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
-- +goose StatementEnd

-- +goose Down
DROP TABLE IF EXISTS tasks;