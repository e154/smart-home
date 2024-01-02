-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
create table trigger_entities
(
    trigger_id bigint not null
        constraint trigger_entities_2_triggers_fk
            references triggers
            on update cascade on delete cascade,
    entity_id  text   not null
        constraint trigger_entity_2_entities_fk
            references entities
            on update cascade on delete cascade,
    constraint trigger_entities_pkey
        primary key (trigger_id, entity_id)
);

insert into trigger_entities (trigger_id, entity_id)
select t.id, t.entity_id
from triggers t
where t.entity_id notnull;

alter table triggers
    drop column entity_id cascade;

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
drop table trigger_entities cascade;

