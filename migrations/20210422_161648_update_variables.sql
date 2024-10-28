-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
ALTER TABLE variables
    ADD COLUMN entity_id text default null;
ALTER TABLE variables
    ADD CONSTRAINT variables_2_entities
        FOREIGN KEY (entity_id) REFERENCES entities (id) ON DELETE CASCADE ON UPDATE CASCADE;

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
alter table variables
    drop column entity_id cascade;

