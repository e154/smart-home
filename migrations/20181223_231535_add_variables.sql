-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE variables (
  name       VARCHAR(255) NOT NULL PRIMARY KEY,
  value      text         NOT NULL,
  autoload   BOOLEAN DEFAULT FALSE,
  created_at TIMESTAMPTZ  NOT NULL,
  updated_at TIMESTAMPTZ  NOT NULL
);

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE IF EXISTS variables CASCADE;

