-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
create type level_type as enum ('Emergency', 'Alert', 'Critical', 'Error', 'Warning', 'Notice', 'Info', 'Debug');

CREATE TABLE logs (
  id         BIGSERIAL PRIMARY KEY,
  body       text        NOT NULL,
  created_at TIMESTAMPTZ NOT NULL,
  level      level_type  NOT NULL DEFAULT 'Info'
);

CREATE INDEX level_2_logs_idx
  ON logs (level);

CREATE INDEX created_at_2_logs_idx
  ON logs (created_at);

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE IF EXISTS logs CASCADE;

