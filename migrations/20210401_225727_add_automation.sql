-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
create type condition_type as enum ('or', 'and');

create table tasks
(
    id          bigserial primary key,
    name        text check (char_length(name) <= 255)                 not null,
    description text                                                  null,
    condition   condition_type           default 'or'::condition_type not null,
    enabled     boolean                  default false,
    area_id     bigint                                                null
        constraint trigger_2_areas_fk
            references areas
            on update cascade on delete set null,
    created_at  timestamp with time zone default CURRENT_TIMESTAMP,
    updated_at  timestamp with time zone default CURRENT_TIMESTAMP
);

create table triggers
(
    id          bigserial                                    not null,
    name        text check (char_length(name) <= 255)        not null,
    plugin_name text check (char_length(plugin_name) <= 255) not null,
    task_id     bigint                                       not null
        constraint trigger_2_tasks_fk
            references tasks
            on update cascade on delete cascade,
    entity_id   text                                         null
        constraint trigger_2_entities_fk
            references entities
            on update cascade on delete cascade,
    script_id   bigint                                       not null
        constraint trigger_2_scripts_fk
            references scripts
            on update cascade on delete cascade,
    payload     jsonb default '{}'::jsonb
);

create table conditions
(
    id        bigserial                             not null,
    name      text check (char_length(name) <= 255) not null,
    task_id   bigint                                not null
        constraint condition_2_tasks_fk
            references tasks
            on update cascade on delete cascade,
    script_id bigint                                not null
        constraint condition_2_scripts_fk
            references scripts
            on update cascade on delete cascade
);

create table actions
(
    id        bigserial                             not null,
    name      text check (char_length(name) <= 255) not null,
    task_id   bigint                                not null
        constraint action_2_tasks_fk
            references tasks
            on update cascade on delete cascade,
    script_id bigint                                not null
        constraint action_2_scripts_fk
            references scripts
            on update cascade on delete cascade
);

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
drop table tasks cascade;

