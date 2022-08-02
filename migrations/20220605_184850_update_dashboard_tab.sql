-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
alter table dashboard_tabs
    drop column gap;

alter table dashboard_tabs
    add column gap bool default false;

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back


