-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
ALTER TABLE entities
    ADD COLUMN parent text default null;
ALTER TABLE entities
    ADD CONSTRAINT entities_2_entities
        FOREIGN KEY (parent) REFERENCES entities (id) ON DELETE CASCADE ON UPDATE CASCADE;

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
alter table entities drop column parent cascade ;
