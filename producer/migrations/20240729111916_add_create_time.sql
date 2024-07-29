-- +goose Up
-- +goose StatementBegin
ALTER TABLE messages ADD COLUMN create_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE messages DROP COLUMN create_time;
-- +goose StatementEnd
