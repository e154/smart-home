-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
create table telegram_chats
(
    entity_id  text   not null
        constraint telegram_chat_2_entities_fk
            references entities (id)
            on update cascade on delete cascade,
    chat_id    bigint not null,
    username   text,
    created_at timestamp with time zone default CURRENT_TIMESTAMP,
    constraint telegram_chats_pkey
        primary key (entity_id, chat_id)
);

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
drop table telegram_chats cascade;

