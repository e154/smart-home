-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
create type nodes_status as enum ('enabled', 'disabled');
create type scripts_lang as enum ('ts', 'coffeescript', 'javascript');
create type level_type as enum ('Emergency', 'Alert', 'Critical', 'Error', 'Warning', 'Notice', 'Info', 'Debug');
create type alexa_skills_status as enum ('enabled', 'disabled');
create type zigbee2mqtt_devices_status as enum ('active', 'banned', 'removed');
create type message_delivery_status as enum ('new', 'in_progress', 'error', 'succeed');
create type message_type as enum ('sms', 'email', 'ui_notify', 'telegram_notify', 'slack');
create type layer_status as enum ('enabled', 'disabled', 'frozen');
create type map_element_status as enum ('enabled', 'disabled', 'frozen');
create type template_type as enum ('item', 'template');
create type template_status as enum ('active', 'inactive');

create table nodes
(
    id                 bigserial                                                not null
        constraint nodes_pkey
            primary key,
    name               text check (char_length(name) <= 255)                    not null,
    description        text                                                     not null,
    created_at         timestamp with time zone default CURRENT_TIMESTAMP,
    updated_at         timestamp with time zone default CURRENT_TIMESTAMP,
    status             nodes_status             default 'enabled'::nodes_status not null,
    login              text,
    encrypted_password text
);

create table scripts
(
    id          bigserial                                                   not null
        constraint scripts_pkey
            primary key,
    lang        scripts_lang             default 'javascript'::scripts_lang not null,
    name        text check (char_length(name) <= 255)                       not null,
    source      text,
    description text,
    compiled    text,
    created_at  timestamp with time zone default CURRENT_TIMESTAMP,
    updated_at  timestamp with time zone default CURRENT_TIMESTAMP
);

create table images
(
    id         bigserial not null
        constraint images_pkey
            primary key,
    thumb      text      not null,
    image      text      not null,
    mime_type  text      not null,
    title      text      not null,
    size       bigint    not null,
    name       text      not null,
    created_at timestamp with time zone default CURRENT_TIMESTAMP
);

create table roles
(
    name        text not null
        constraint roles_pkey
            primary key,
    description text not null,
    parent      text
        constraint parent_at_roles_fk
            references roles
            on update cascade on delete restrict,
    created_at  timestamp with time zone default CURRENT_TIMESTAMP,
    updated_at  timestamp with time zone default CURRENT_TIMESTAMP
);

create table permissions
(
    id           bigserial not null
        constraint permissions_pkey
            primary key,
    role_name    text      not null
        constraint role_at_permissions_fk
            references roles
            on update cascade on delete cascade,
    package_name text      not null,
    level_name   text      not null
);

create unique index permissions_unq
    on permissions (role_name, package_name, level_name);

create table users
(
    id                     bigserial                                                    not null
        constraint users_pkey
            primary key,
    nickname               text check (char_length(nickname) <= 255)                    not null,
    first_name             text check (char_length(first_name) <= 255)                  not null,
    last_name              text check (char_length(last_name) <= 255)                   not null,
    encrypted_password     text                                                         not null,
    email                  text                                                         not null,
    lang                   varchar(2)               default 'en'::character varying     not null,
    history                jsonb                    default '{}'::jsonb,
    status                 text                     default 'active'::character varying not null,
    reset_password_token   text                                                         not null,
    image_id               bigint
        constraint image_at_users_fk
            references images
            on update cascade on delete set null,
    sign_in_count          bigint                                                       not null,
    current_sign_in_ip     text                                                         not null,
    last_sign_in_ip        text,
    user_id                bigint
        constraint user_at_users_fk
            references users
            on update cascade on delete set null,
    role_name              text
        constraint role_name_at_users_fk
            references roles
            on update cascade on delete restrict,
    reset_password_sent_at timestamp with time zone,
    current_sign_in_at     timestamp with time zone,
    last_sign_in_at        timestamp with time zone,
    created_at             timestamp with time zone default CURRENT_TIMESTAMP,
    updated_at             timestamp with time zone default CURRENT_TIMESTAMP,
    deleted_at             timestamp with time zone,
    authentication_token   text
);

create unique index email_2_users_unq
    on users (email);

create unique index nickname_2_users_unq
    on users (nickname);

create index email_2_user_idx
    on users (email);

create table user_metas
(
    id      bigserial not null
        constraint user_metas_pkey
            primary key,
    user_id bigint
        constraint user_at_user_metas_fk
            references users
            on update cascade on delete set null,
    key     text      not null,
    value   text      not null
);

create unique index kay_user_2_user_metas_unq
    on user_metas (key, user_id);

create table variables
(
    name       text check (char_length(name) <= 255) not null
        constraint variables_pkey
            primary key,
    value      text                                  not null,
    autoload   boolean                  default false,
    created_at timestamp with time zone default CURRENT_TIMESTAMP,
    updated_at timestamp with time zone default CURRENT_TIMESTAMP
);

create table maps
(
    id          bigserial                             not null
        constraint maps_pkey
            primary key,
    name        text check (char_length(name) <= 255) not null,
    description text                                  not null,
    options     jsonb                    default '{}'::jsonb,
    created_at  timestamp with time zone default CURRENT_TIMESTAMP,
    updated_at  timestamp with time zone default CURRENT_TIMESTAMP
);

create table map_texts
(
    id    bigserial not null
        constraint map_texts_pkey
            primary key,
    text  text      not null,
    style text default '{}'::text
);

create table map_layers
(
    id          bigserial                          not null
        constraint map_layers_pkey
            primary key,
    name        text                               not null,
    description text                               not null,
    map_id      bigint
        constraint map_layers_2_maps_fk
            references maps
            on update cascade on delete restrict,
    status      layer_status                       not null,
    weight      integer                  default 0 not null,
    created_at  timestamp with time zone default CURRENT_TIMESTAMP,
    updated_at  timestamp with time zone default CURRENT_TIMESTAMP
);

create table map_images
(
    id       bigserial not null
        constraint map_images_pkey
            primary key,
    image_id bigint
        constraint map_images_2_images_fk
            references images
            on update cascade on delete set null,
    style    text default '{}'::text
);

create table map_elements
(
    id             bigserial                                    not null
        constraint map_elements_pkey
            primary key,
    name           text check (char_length(name) <= 255)        not null,
    description    text                                         not null,
    prototype_id   text                                         not null,
    prototype_type text                                         not null,
    graph_settings jsonb                    default '{}'::jsonb not null,
    map_layer_id   bigint
        constraint map_elements_2_map_layers_fk
            references map_layers
            on update cascade on delete restrict,
    map_id         bigint
        constraint map_elements_2_maps_fk
            references maps
            on update cascade on delete restrict,
    status         map_element_status                           not null,
    weight         integer                  default 0           not null,
    created_at     timestamp with time zone default CURRENT_TIMESTAMP,
    updated_at     timestamp with time zone default CURRENT_TIMESTAMP
);

create table map_devices
(
    id         bigserial not null
        constraint map_devices_pkey
            primary key,
    device_id  bigint,
    image_id   bigint
        constraint map_devices_2_images_fk
            references images
            on update cascade on delete set null,
    created_at timestamp with time zone default CURRENT_TIMESTAMP,
    updated_at timestamp with time zone default CURRENT_TIMESTAMP
);

create index device_at_map_devices_idx
    on map_devices (device_id);

create table logs
(
    id         bigserial                                           not null
        constraint logs_pkey
            primary key,
    body       text                                                not null,
    created_at timestamp with time zone default CURRENT_TIMESTAMP,
    level      level_type               default 'Info'::level_type not null
);

create index level_2_logs_idx
    on logs (level);

create index created_at_2_logs_idx
    on logs (created_at);

create table zones
(
    id   bigserial not null
        constraint zone_tags_pkey
            primary key,
    name text      not null
);

create unique index name_at_zone_tags_unq
    on zones (name);

create table templates
(
    name        text not null
        constraint templates_pkey
            primary key,
    description text,
    content     text not null,
    status      template_status          default 'active'::template_status,
    type        template_type            default 'item'::template_type,
    parent      text
        constraint parent_2_templates_fk
            references templates
            on update cascade on delete restrict,
    created_at  timestamp with time zone default CURRENT_TIMESTAMP,
    updated_at  timestamp with time zone default CURRENT_TIMESTAMP
);

create index status_2_templates_idx
    on templates (status);

create index type_2_templates_idx
    on templates (type);

create table messages
(
    id            bigserial    not null
        constraint messages_pkey
            primary key,
    type          message_type not null,
    email_from    text,
    email_subject text,
    email_body    text,
    sms_text      text,
    ui_text       text,
    slack_text    text,
    telegram_text text,
    created_at    timestamp with time zone default CURRENT_TIMESTAMP,
    updated_at    timestamp with time zone default CURRENT_TIMESTAMP
);

create table message_deliveries
(
    id                   bigserial               not null
        constraint message_deliveries_pkey
            primary key,
    message_id           bigint                  not null
        constraint message_deliveries_2_messages_fk
            references messages
            on update cascade on delete restrict,
    address              text                    not null,
    status               message_delivery_status not null,
    error_system_code    text
        constraint message_deliveries_error_system_code_check
            check (char_length(error_system_code) <= 255),
    error_system_message text
        constraint message_deliveries_error_system_message_check
            check (char_length(error_system_message) <= 255),
    created_at           timestamp with time zone default CURRENT_TIMESTAMP,
    updated_at           timestamp with time zone default CURRENT_TIMESTAMP
);

create table zigbee2mqtt
(
    id                 bigserial not null
        constraint zigbee2mqtt_pkey
            primary key,
    name               text,
    login              text,
    encrypted_password text,
    permit_join        boolean                  default true,
    base_topic         text                     default 'zigbee2mqtt'::text,
    created_at         timestamp with time zone default CURRENT_TIMESTAMP,
    updated_at         timestamp with time zone default CURRENT_TIMESTAMP
);

create unique index base_topic_at_zigbee2mqtt_unq
    on zigbee2mqtt (base_topic);

create index login_at_zigbee2mqtt_idx
    on zigbee2mqtt (login);

create table zigbee2mqtt_devices
(
    id             text not null
        constraint zigbee2mqtt_devices_pkey
            primary key,
    zigbee2mqtt_id bigint
        constraint zigbee2mqtt_devices_2_zigbee2mqtt_fk
            references zigbee2mqtt
            on update cascade on delete restrict,
    name           text,
    type           text,
    model          text,
    description    text,
    manufacturer   text,
    functions      text[],
    status         zigbee2mqtt_devices_status default 'active'::zigbee2mqtt_devices_status,
    created_at     timestamp with time zone   default CURRENT_TIMESTAMP,
    updated_at     timestamp with time zone   default CURRENT_TIMESTAMP
);

create table alexa_skills
(
    id             bigserial                                                       not null
        constraint alexa_skills_pkey
            primary key,
    skill_id       text                                                            not null,
    description    text,
    status         alexa_skills_status      default 'enabled'::alexa_skills_status not null,
    on_launch      bigint
        constraint on_launch_at_alexa_skills_2_scripts_fk
            references scripts
            on update cascade on delete cascade,
    on_session_end bigint
        constraint on_session_end_at_alexa_skills_2_scripts_fk
            references scripts
            on update cascade on delete cascade,
    created_at     timestamp with time zone default CURRENT_TIMESTAMP,
    updated_at     timestamp with time zone default CURRENT_TIMESTAMP
);

create unique index id_at_alexa_skills_unq
    on alexa_skills (id);

create table alexa_intents
(
    name           text   not null,
    alexa_skill_id bigint not null
        constraint alexa_intents_2_alexa_skills_fk
            references alexa_skills
            on update cascade on delete cascade,
    script_id      bigint not null
        constraint script_alexa_intents_2_scripts_fk
            references scripts
            on update cascade on delete cascade,
    description    text,
    created_at     timestamp with time zone default CURRENT_TIMESTAMP,
    updated_at     timestamp with time zone default CURRENT_TIMESTAMP
);

create unique index name_at_alexa_intents_unq
    on alexa_intents (name, alexa_skill_id);

create table storage
(
    name       text  not null
        constraint storage_name_key
            unique,
    value      jsonb not null,
    created_at timestamp with time zone default CURRENT_TIMESTAMP,
    updated_at timestamp with time zone default CURRENT_TIMESTAMP
);

create table metrics
(
    id          bigserial                                     not null
        constraint metrics_pkey
            primary key,
    name        text                                          not null,
    description text,
    options     jsonb                    default '{}'::jsonb  not null,
    type        text                     default 'line'::text not null,
    created_at  timestamp with time zone default CURRENT_TIMESTAMP,
    updated_at  timestamp with time zone default CURRENT_TIMESTAMP
);

create index name_at_metrics_unq
    on metrics (name);

create table metric_bucket
(
    value     jsonb                    not null,
    metric_id bigint                   not null
        constraint metric_id_at_metric_bucket_2_metric_fk
            references metrics
            on update cascade on delete cascade,
    time      timestamp with time zone not null,
    constraint metric_bucket_pkey
        primary key (time, metric_id)
);

create index time_at_metric_bucket_idx
    on metric_bucket (time);

create table areas
(
    id          bigserial not null
        constraint areas_pkey
            primary key,
    name        text      not null,
    description text                     default '',
    created_at  timestamp with time zone default CURRENT_TIMESTAMP,
    updated_at  timestamp with time zone default CURRENT_TIMESTAMP
);

create table entities
(
    id          text   not null
        constraint entities_pkey
            primary key,
    type        text   not null,
    description text,
    area_id     bigint null
        constraint area_id_at_entities_2_areas_fk
            references areas
            on update cascade on delete set null,
    icon        text,
    image_id    bigint
        constraint image_id_at_entities_2_images_fk
            references images
            on update cascade on delete set null,
    payload     jsonb                    default '{}'::jsonb,
    auto_load   boolean                  default false,
    created_at  timestamp with time zone default CURRENT_TIMESTAMP,
    updated_at  timestamp with time zone default CURRENT_TIMESTAMP
);

create index name_at_entities_idx
    on entities (type);

create table entity_metrics
(
    entity_id text   not null
        constraint entity_id_at_entity_metric_2_entities_fk
            references entities
            on update cascade on delete cascade,
    metric_id bigint not null
        constraint metric_id_at_entity_metric_2_metric_fk
            references metrics
            on update cascade on delete cascade,
    constraint entity_metrics_pkey
        primary key (entity_id, metric_id)
);

create table entity_states
(
    id          bigserial not null
        constraint entity_states_pkey
            primary key,
    name        text      not null,
    description text,
    icon        text,
    entity_id   text
        constraint entity_states_2_entity_fk
            references entities
            on update cascade on delete cascade,
    image_id    bigint
        constraint entity_states_2_images_fk
            references images
            on update cascade on delete set null,
    style       text                     default '{}'::text,
    created_at  timestamp with time zone default CURRENT_TIMESTAMP,
    updated_at  timestamp with time zone default CURRENT_TIMESTAMP
);

create unique index name_at_entity_states_unq
    on entity_states (name, entity_id);

create table entity_actions
(
    id          bigserial not null
        constraint entity_actions_pkey
            primary key,
    name        text      not null,
    description text,
    icon        text,
    entity_id   text
        constraint entity_actions_2_entities_fk
            references entities
            on update cascade on delete cascade,
    image_id    bigint
        constraint entity_actions_2_images_fk
            references images
            on update cascade on delete set null,
    script_id   bigint
        constraint entity_actions_2_scripts_fk
            references scripts
            on update cascade on delete set null,
    type        text      not null,
    created_at  timestamp with time zone default CURRENT_TIMESTAMP,
    updated_at  timestamp with time zone default CURRENT_TIMESTAMP
);

create unique index name_at_entity_actions_unq
    on entity_actions (name, entity_id);

create table entity_storage
(
    id         bigserial                                    not null
        constraint entity_storage_pkey
            primary key,
    entity_id  text
        constraint entity_storage_2_entities_fk
            references entities
            on update cascade on delete cascade,
    state      text                                         not null,
    attributes jsonb                    default '{}'::jsonb not null,
    created_at timestamp with time zone default CURRENT_TIMESTAMP
);

create index created_at_map_entity_storage_idx
    on entity_storage (created_at);

create table entity_scripts
(
    id        bigserial         not null
        constraint entity_scripts_pkey
            primary key,
    entity_id text              not null
        constraint entity_id_at_entity_scripts_2_entities_fk
            references entities
            on update cascade on delete cascade,
    script_id bigint            not null
        constraint entity_scripts_2_scripts_fk
            references scripts
            on update cascade on delete restrict,
    weight    integer default 0 not null
);

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back


