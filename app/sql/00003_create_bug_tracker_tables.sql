-- +goose Up
-- +goose StatementBegin
CREATE TABLE Bugs (
  bug_id            SERIAL PRIMARY KEY,
  date_reported     DATE NOT NULL,
  summary           VARCHAR(80),
  description       VARCHAR(1000),
  resolution        VARCHAR(1000),
  reported_by       BIGINT NOT NULL,
  assigned_to       BIGINT,
  verified_by       BIGINT,
  status            VARCHAR(20) NOT NULL DEFAULT 'NEW',
  priority          VARCHAR(20),
  hours             NUMERIC(9, 2),
  FOREIGN KEY (reported_by) REFERENCES Accounts(account_id),
  FOREIGN KEY (assigned_to) REFERENCES Accounts(account_id),
  FOREIGN KEY (verified_by) REFERENCES Accounts(account_id),
  FOREIGN KEY (status) REFERENCES BugStatus(status)
);

CREATE TABLE Comments (
  comment_id      SERIAL PRIMARY KEY,
  bug_id          BIGINT NOT NULL,
  author          BIGINT NOT NULL,
  comment_date    TIMESTAMP NOT NULL,
  comment         TEXT NOT NULL,
  FOREIGN KEY (bug_id) REFERENCES Bugs(bug_id),
  FOREIGN KEY (author) REFERENCES Accounts(account_id)
);

CREATE TABLE Screenshots (
  bug_id            BIGINT NOT NULL,
  image_id          BIGINT NOT NULL,
  screenshot_image  BYTEA,
  caption           VARCHAR(100),
  PRIMARY KEY       (bug_id, image_id),
  FOREIGN KEY (bug_id) REFERENCES Bugs(bug_id)
);

CREATE TABLE Tags (
  bug_id        BIGINT NOT NULL,
  tag           VARCHAR(20) NOT NULL,
  PRIMARY KEY   (bug_id, tag),
  FOREIGN KEY (bug_id) REFERENCES Bugs(bug_id)
);

CREATE TABLE Products (
  product_id    SERIAL PRIMARY KEY,
  product_name  VARCHAR(50)
);

CREATE TABLE BugsProducts(
  bug_id          BIGINT NOT NULL,
  product_id      BIGINT NOT NULL,
  PRIMARY KEY     (bug_id, product_id),
  FOREIGN KEY (bug_id) REFERENCES Bugs(bug_id),
  FOREIGN KEY (product_id) REFERENCES Products(product_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Drop foreign key constraints in BugsProducts table
ALTER TABLE BugsProducts DROP CONSTRAINT bug_id;
ALTER TABLE BugsProducts DROP CONSTRAINT product_id;

-- Drop foreign key constraints in Tags table
ALTER TABLE Tags DROP CONSTRAINT bug_id;

-- Drop foreign key constraints in Screenshots table
ALTER TABLE Screenshots DROP CONSTRAINT bug_id;

-- Drop foreign key constraints in Comments table
ALTER TABLE Comments DROP CONSTRAINT bug_id;
ALTER TABLE Comments DROP CONSTRAINT account_id;

-- Drop foreign key constraints in Bugs table
ALTER TABLE Bugs DROP CONSTRAINT reported_by;
ALTER TABLE Bugs DROP CONSTRAINT assigned_to;
ALTER TABLE Bugs DROP CONSTRAINT verified_by;
ALTER TABLE Bugs DROP CONSTRAINT status;

-- Now drop tables in reverse order of dependencies
DROP TABLE BugsProducts;
DROP TABLE Products;
DROP TABLE Tags;
DROP TABLE Screenshots;
DROP TABLE Comments;
DROP TABLE Bugs;
-- +goose StatementEnd
