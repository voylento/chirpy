-- +goose Up
-- +goose StatementBegin
CREATE TABLE chirps(
  id          UUID PRIMARY KEY,
  created_at  TIMESTAMP NOT NULL,
  updated_at  TIMESTAMP NOT NULL,
  body        TEXT NOT NULL,
  user_id     UUID NOT NULL,
  CONSTRAINT fk_users
    FOREIGN KEY (user_id)
    REFERENCES  users(id)
    ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE chirps;
-- +goose StatementEnd
