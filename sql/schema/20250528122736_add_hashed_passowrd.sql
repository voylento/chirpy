-- +goose Up
-- +goose StatementBegin
ALTER TABLE users ADD COLUMN hashed_password TEXT NOT NULL DEFAULT 'UNSET';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users DROP COLUMN hased_password;
-- +goose StatementEnd
