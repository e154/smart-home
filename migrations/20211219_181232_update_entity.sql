-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
alter table entities
    rename column type to plugin;

alter table entities
    add constraint entities_2_plugins_fk
        foreign key (plugin) references plugins (name) on delete set null on update cascade;

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
alter table entities
    drop constraint entities_2_plugins_fk cascade;

alter table entities
    rename column plugin to type;
