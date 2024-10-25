-- +goose Up
-- +goose StatementBegin
INSERT INTO BugStatus (status) VALUES ('NEW'), ('IN_PROGRESS'), ('COMPLETED');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM BugStatus;
-- +goose StatementEnd
