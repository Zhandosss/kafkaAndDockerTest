-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS messages (
    id UUID NOT NULL DEFAULT uuid_generate_v4(),
    content TEXT NOT NULL,
    is_processed BOOLEAN NOT NULL DEFAULT FALSE,
    primary key (id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS messages;
-- +goose StatementEnd