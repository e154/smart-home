-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
create type alexa_skills_status as enum ('enabled', 'disabled');

create table alexa_skills
(
    id             bigserial
        constraint alexa_skills_pkey primary key not null,
    skill_id       text                          not null,
    description    text                          null,
    status         alexa_skills_status           NOT NULL DEFAULT 'enabled',
    on_launch      BIGINT                        NULL
        CONSTRAINT on_launch_at_alexa_skills_2_scripts_fk REFERENCES scripts (id) ON UPDATE CASCADE ON DELETE CASCADE,
    on_session_end BIGINT                        NULL
        CONSTRAINT on_session_end_at_alexa_skills_2_scripts_fk REFERENCES scripts (id) ON UPDATE CASCADE ON DELETE CASCADE,
    created_at     timestamp with time zone      not null,
    updated_at     timestamp with time zone      not null
);

CREATE UNIQUE INDEX id_at_alexa_skills_unq
    ON alexa_skills (id);

create table alexa_intents
(
    name           text                     NOT NULL,
    alexa_skill_id bigint                   NOT NULL
        CONSTRAINT alexa_intents_2_alexa_skills_fk REFERENCES alexa_skills (id) ON UPDATE CASCADE ON DELETE CASCADE,
    script_id      BIGINT                   NOT NULL
        CONSTRAINT script_alexa_intents_2_scripts_fk REFERENCES scripts (id) ON UPDATE CASCADE ON DELETE CASCADE,
    description    text                     NULL,
    created_at     timestamp with time zone not null,
    updated_at     timestamp with time zone not null
);

CREATE UNIQUE INDEX name_at_alexa_intents_unq
    ON alexa_intents (name, alexa_skill_id);

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
drop table alexa_intents cascade;
drop table alexa_skills cascade;
drop type alexa_skills_status cascade;

