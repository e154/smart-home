-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
alter table triggers drop column task_id;
alter table conditions drop column task_id;
alter table actions drop column task_id;

ALTER TABLE triggers ADD PRIMARY KEY (id);
ALTER TABLE conditions ADD PRIMARY KEY (id);
ALTER TABLE actions ADD PRIMARY KEY (id);

alter table triggers
    add column created_at  timestamp with time zone default CURRENT_TIMESTAMP,
    add column updated_at  timestamp with time zone default CURRENT_TIMESTAMP;

alter table conditions
    add column created_at  timestamp with time zone default CURRENT_TIMESTAMP,
    add column updated_at  timestamp with time zone default CURRENT_TIMESTAMP;

alter table actions
    add column created_at  timestamp with time zone default CURRENT_TIMESTAMP,
    add column updated_at  timestamp with time zone default CURRENT_TIMESTAMP;

create table task_triggers
(
    task_id bigint   not null
        constraint task_id_at_tasks_2_tasks_fk
            references tasks
            on update cascade on delete cascade,
    trigger_id bigint not null
        constraint trigger_id_at_task_triggers_2_triggers_fk
            references triggers
            on update cascade on delete cascade,
    constraint task_triggers_pkey
        primary key (task_id, trigger_id)
);

create table task_conditions
(
    task_id bigint   not null
        constraint task_id_at_tasks_2_tasks_fk
            references tasks
            on update cascade on delete cascade,
    condition_id bigint not null
        constraint condition_id_at_task_conditions_2_conditions_fk
            references conditions
            on update cascade on delete cascade,
    constraint task_conditions_pkey
        primary key (task_id, condition_id)
);

create table task_actions
(
    task_id bigint   not null
        constraint task_id_at_tasks_2_tasks_fk
            references tasks
            on update cascade on delete cascade,
    action_id bigint not null
        constraint action_id_at_task_actions_2_actions_fk
            references actions
            on update cascade on delete cascade,
    constraint task_actions_pkey
        primary key (task_id, action_id)
);

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back


