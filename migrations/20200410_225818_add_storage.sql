-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
create table storage
(
    name       text                     not null unique,
    value      jsonb                    not null,
    created_at timestamp with time zone not null,
    updated_at timestamp with time zone not null
);

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
drop table storage;

