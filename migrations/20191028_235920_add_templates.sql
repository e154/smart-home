-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
create type template_type as enum ('item', 'template');
create type template_status as enum ('active', 'inactive');

CREATE TABLE templates
(
    name        text primary key         NOT NULL,
    description text                     NULL,
    content     text                     NOT NULL,
    status      template_status               DEFAULT 'active',
    type        template_type            NULL DEFAULT 'item',
    parent      text                     NULL
        CONSTRAINT parent_2_templates_fk REFERENCES templates (name) ON UPDATE CASCADE ON DELETE RESTRICT,
    created_at  TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at  TIMESTAMP WITH TIME ZONE NOT NULL
);

CREATE INDEX status_2_templates_idx
    ON templates (status);

CREATE INDEX type_2_templates_idx
    ON templates (type);

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
drop table templates cascade;
drop type template_type;
drop type template_status;