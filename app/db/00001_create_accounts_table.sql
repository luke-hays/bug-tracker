-- +goose Up
-- +goose StatementBegin
CREATE TABLE Accounts (
  account_id          SERIAL PRIMARY KEY,
  account_name        VARCHAR(20),
  first_name          VARCHAR(20),
  last_name           VARCHAR(20),
  email               VARCHAR(100),
  password_hash       CHAR(64),
  portrait_image      BYTEA,
  hourly_rate         NUMERIC(9, 2),
  created_at          TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at          TIMESTAMP DEFAULT CURRENT_TIMESTAMP
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE Accounts;
-- +goose StatementEnd
