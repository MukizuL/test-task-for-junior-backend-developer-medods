-- +goose Up
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

-- +goose Down
DROP TABLE IF EXISTS recurring_tasks;