-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
create table user_devices
(
    id                bigserial primary key,
    user_id           bigint not null
        constraint user_device_2_users_fk
            references users
            on update cascade on delete cascade,
    push_registration jsonb                    default '{}'::jsonb,
    created_at        timestamp with time zone default CURRENT_TIMESTAMP
);

create unique index push_registration_at_user_devices_unq
    on user_devices (push_registration);

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
drop table user_devices cascade;

