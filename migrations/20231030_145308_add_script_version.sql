-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
create table script_versions
(
    id         bigserial primary key,
    script_id  bigint                                                      not null
        constraint script_version_2_scripts_fk
            references scripts
            on update cascade on delete cascade,
    lang       scripts_lang             default 'javascript'::scripts_lang not null,
    source     text,
    sum        bytea,
    created_at timestamp with time zone default CURRENT_TIMESTAMP
);

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back


