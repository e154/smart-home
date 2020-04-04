-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
create table alexa_applications
(
    id             BIGSERIAL PRIMARY KEY,
    application_id text                     not null,
    description    text                     null,
    created_at     timestamp with time zone not null,
    updated_at     timestamp with time zone null
);

CREATE UNIQUE INDEX id_at_alexa_applications_unq
    ON alexa_applications (id);

create table alexa_intents
(
    name                 text                     NOT NULL,
    alexa_application_id bigint                   NOT NULL
        CONSTRAINT alexa_intents_2_alexa_applications_fk REFERENCES alexa_applications (id) ON UPDATE CASCADE ON DELETE CASCADE,
    flow_id              BIGINT                   NOT NULL
        CONSTRAINT flow_alexa_intents_2_flows_fk REFERENCES flows (id) ON UPDATE CASCADE ON DELETE CASCADE,
    description          text                     NULL,
    created_at           timestamp with time zone not null,
    updated_at           timestamp with time zone null
);

CREATE UNIQUE INDEX name_at_alexa_intents_unq
    ON alexa_intents (name, alexa_application_id);

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
drop table alexa_intents cascade;
drop table alexa_applications cascade;

