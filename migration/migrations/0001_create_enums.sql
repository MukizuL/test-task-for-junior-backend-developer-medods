-- +goose Up
-- +goose StatementBegin
CREATE TYPE status AS ENUM ('new', 'in_progress', 'done');
CREATE TYPE rec_type AS ENUM ('interval', 'even_odd_days', 'specific_dates');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TYPE IF EXISTS status;
DROP TYPE IF EXISTS rec_type;
-- +goose StatementEnd
