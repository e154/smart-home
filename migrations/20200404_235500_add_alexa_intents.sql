-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
create type alexa_applications_status as enum ('enabled', 'disabled');

create table alexa_applications
(
    id             bigserial
        constraint alexa_applications_pkey primary key not null,
    application_id text                                not null,
    description    text                                null,
    status         alexa_applications_status           NOT NULL DEFAULT 'enabled',
    on_launch      BIGINT                              NULL
        CONSTRAINT on_launch_at_alexa_applications_2_scripts_fk REFERENCES scripts (id) ON UPDATE CASCADE ON DELETE CASCADE,
    on_session_end BIGINT                              NULL
        CONSTRAINT on_session_end_at_alexa_applications_2_scripts_fk REFERENCES scripts (id) ON UPDATE CASCADE ON DELETE CASCADE,
    created_at     timestamp with time zone            not null,
    updated_at     timestamp with time zone            not null
);

CREATE UNIQUE INDEX id_at_alexa_applications_unq
    ON alexa_applications (id);

create table alexa_intents
(
    name                 text                     NOT NULL,
    alexa_application_id bigint                   NOT NULL
        CONSTRAINT alexa_intents_2_alexa_applications_fk REFERENCES alexa_applications (id) ON UPDATE CASCADE ON DELETE CASCADE,
    script_id            BIGINT                   NOT NULL
        CONSTRAINT script_alexa_intents_2_scripts_fk REFERENCES scripts (id) ON UPDATE CASCADE ON DELETE CASCADE,
    description          text                     NULL,
    created_at           timestamp with time zone not null,
    updated_at           timestamp with time zone not null
);

CREATE UNIQUE INDEX name_at_alexa_intents_unq
    ON alexa_intents (name, alexa_application_id);

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
drop table alexa_intents cascade;
drop table alexa_applications cascade;
drop type alexa_applications_status cascade;

