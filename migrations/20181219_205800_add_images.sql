-- +migrate Up
CREATE TABLE images (
  id         BIGSERIAL PRIMARY KEY,
  thumb      text        NOT NULL,
  image      text        NOT NULL,
  mime_type  text        NOT NULL,
  title      text        NOT NULL,
  size       BIGINT      NOT NULL,
  name       text        NOT NULL,
  created_at TIMESTAMPTZ NOT NULL
);

-- +migrate Down
DROP TABLE IF EXISTS images CASCADE;
