-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
create table variable_tags
(
    variable_name text      not null
        constraint variable_name_at_variable_tags_2_variables_fk
            references variables
            on update cascade on delete cascade,
    tag_id        bigserial not null
        constraint variable_tags_2_tags_fk
            references tags
            on update cascade on delete restrict,
    constraint variable_tags_pkey
        primary key (variable_name, tag_id)
);

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
drop table variable_tags cascade;

