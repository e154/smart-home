-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
create unique index name_at_triggers_unq
    on triggers (name, task_id);
create unique index name_at_actions_unq
    on actions (name, task_id);

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
drop index name_at_triggers_unq cascade;
drop index name_at_actions_unq cascade;

