-- +goose Up
-- +goose StatementBegin
CREATE TABLE BugStatus (
  status VARCHAR(20) PRIMARY KEY
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE BugStatus;
-- +goose StatementEnd
