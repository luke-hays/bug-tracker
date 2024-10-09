-- +goose Up
-- +goose StatementBegin
CREATE TABLE Sessions (
  session_id CHAR(64) PRIMARY KEY,
  account_id BIGINT NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  expires_at  TIMESTAMP NOT NULL,
  FOREIGN KEY (account_id) REFERENCES Accounts(account_id)
);

ALTER TABLE Accounts 
  ADD COLUMN session_id CHAR(64),
  ADD CONSTRAINT session_id FOREIGN KEY (session_id) REFERENCES Sessions(session_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE Accounts DROP CONSTRAINT session_id;
ALTER TABLE Accounts DROP session_id;

DROP TABLE Sessions;
-- +goose StatementEnd
