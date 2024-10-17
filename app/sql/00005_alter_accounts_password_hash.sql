-- +goose Up
-- +goose StatementBegin
ALTER TABLE Accounts ALTER COLUMN password_hash TYPE VARCHAR(255)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE Accounts ALTER COLUMN password_hash TYPE CHAR(64)
-- +goose StatementEnd
