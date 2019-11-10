-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
create type message_type as enum ('sms', 'email', 'ui_notify', 'telegram_notify', 'slack');

create table messages
(
    id            BIGSERIAL PRIMARY KEY,
    type          message_type NOT NULL,
    email_from    text         NULL,
    email_subject text         NULL,
    email_body    text         NULL,
    sms_text      text         NULL,
    ui_text       text         NULL,
    slack_text    text         NULL,
    telegram_text text         NULL,
    created_at    TIMESTAMPTZ  NOT NULL,
    updated_at    TIMESTAMPTZ  NOT NULL
);

create type message_delivery_status as enum ('new', 'in_progress', 'error', 'succeed');

create table message_deliveries
(
    id                   BIGSERIAL PRIMARY KEY,
    message_id           BIGINT                  NOT NULL
        CONSTRAINT message_deliveries_2_messages_fk REFERENCES messages (id) ON UPDATE CASCADE ON DELETE RESTRICT,
    address              text                    NOT NULL,
    status               message_delivery_status NOT NULL,
    error_system_code    text                    NULL
        constraint message_deliveries_error_system_code_check check (char_length(error_system_code) <= 255),
    error_system_message text                    NULL
        constraint message_deliveries_error_system_message_check check (char_length(error_system_message) <= 255),
    created_at           TIMESTAMPTZ             NOT NULL,
    updated_at           TIMESTAMPTZ             NOT NULL
);

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
drop table message_deliveries cascade;
drop table messages cascade;
drop type if exists message_delivery_status;
drop type if exists message_type;