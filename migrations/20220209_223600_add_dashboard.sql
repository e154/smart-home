-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
alter table variables
    rename column autoload to system;

create table dashboards
(
    id          bigserial primary key,
    name        text check (char_length(name) <= 255)        not null,
    description text check (char_length(description) <= 255) null,
    enabled     boolean                  default false,
    area_id     bigint                                       null
        constraint dashboard_2_areas_fk
            references areas
            on update cascade on delete set null,
    created_at  timestamp with time zone default CURRENT_TIMESTAMP,
    updated_at  timestamp with time zone default CURRENT_TIMESTAMP
);

create unique index name_at_dashboards_unq
    on dashboards (name);

create table dashboard_tabs
(
    id           bigserial primary key,
    name         text check (char_length(name) <= 255) not null,
    icon         text check (char_length(name) <= 255)       default '',
    enabled      boolean                                     default false,
    weight       int                                         default 0,
    column_width int                                         default 0,
    gap          int                                         default 0,
    background   text check (char_length(background) <= 255) default '',
    dashboard_id bigint                                not null
        constraint dashboard_2_dashboards_fk
            references dashboards
            on update cascade on delete cascade,
    created_at   timestamp with time zone                    default CURRENT_TIMESTAMP,
    updated_at   timestamp with time zone                    default CURRENT_TIMESTAMP
);

create unique index name_at_dashboard_tabs_unq
    on dashboard_tabs (name, dashboard_id);

create table dashboard_cards
(
    id               bigserial primary key,
    title            text check (char_length(title) <= 255) not null,
    background       text check (char_length(background) <= 255) default '',
    weight           int                                         default 0,
    width            int                                         default 0,
    height           int                                         default 0,
    enabled          boolean                                     default false,
    dashboard_tab_id bigint                                 not null
        constraint dashboard_card_2_dashboard_tab_fk
            references dashboard_tabs
            on update cascade on delete cascade,
    payload          jsonb                                       default '{}'::jsonb,
    created_at       timestamp with time zone                    default CURRENT_TIMESTAMP,
    updated_at       timestamp with time zone                    default CURRENT_TIMESTAMP
);

create table dashboard_card_items
(
    id                bigserial primary key,
    title             text check (char_length(title) <= 255) not null,
    type              text check (char_length(type) <= 255)  not null,
    weight            int                      default 0,
    enabled           boolean                  default false,
    dashboard_card_id bigint                                 not null
        constraint dashboard_card_item_2_dashboard_card_fk
            references dashboard_cards
            on update cascade on delete cascade,
    entity_id         text
        constraint dashboard_card_item_2_entities_fk
            references entities (id)
            on update cascade on delete cascade,
    payload           jsonb                    default '{}'::jsonb,
    hidden            bool                     default false not null,
    frozen            bool                     default false not null,
    created_at        timestamp with time zone default CURRENT_TIMESTAMP,
    updated_at        timestamp with time zone default CURRENT_TIMESTAMP
);

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
drop table if exists dashboard_card_items cascade;
drop table dashboard_cards cascade;
drop table dashboard_tabs cascade;
drop table dashboards cascade;

alter table variables
    rename column system to autoload;
