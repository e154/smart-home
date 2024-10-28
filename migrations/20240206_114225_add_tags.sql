-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
create table tags
(
    id   bigserial primary key,
    name text not null
        constraint tag_name_unq
            unique

);

create index tags_name_idx
    on tags (name);

create table entity_tags
(
    entity_id text      not null
        constraint entity_id_at_entity_tags_2_entities_fk
            references entities
            on update cascade on delete cascade,
    tag_id    bigserial not null
        constraint entity_tags_2_tags_fk
            references tags
            on update cascade on delete restrict,
    constraint entity_tags_pkey
        primary key (entity_id, tag_id)
);

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
drop table entity_tags cascade;
drop table tags cascade;

