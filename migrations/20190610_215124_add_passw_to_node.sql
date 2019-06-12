-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
alter table nodes
  add column login text null;
alter table nodes
  add column encrypted_password text null;

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
alter table nodes
  drop column if exists login;
alter table nodes
  drop column if exists encrypted_password;

