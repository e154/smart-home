-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
create table run_history
(
    id    bigserial not null,
    start timestamp with time zone default CURRENT_TIMESTAMP,
    "end" timestamp with time zone default null
);

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
drop table run_history;
