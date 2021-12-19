-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
alter table entities
    rename column type to plugin_name;

alter table entities
    add constraint entities_2_plugins_fk
        foreign key (plugin_name) references plugins (name) on delete set null on update cascade;

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
alter table entities
    drop constraint entities_2_plugins_fk cascade;

alter table entities
    rename column plugin_name to type;
