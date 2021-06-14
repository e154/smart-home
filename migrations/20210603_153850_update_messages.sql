-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
alter table messages
    add column payload jsonb default '{}'::jsonb;
alter table messages
    drop column email_from;
alter table messages
    drop column email_subject;
alter table messages
    drop column email_body;
alter table messages
    drop column sms_text;
alter table messages
    drop column ui_text;
alter table messages
    drop column slack_text;
alter table messages
    drop column telegram_text;

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back


