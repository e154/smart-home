-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
alter table dashboard_tabs
    add column payload jsonb default '{}'::jsonb;

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back


