-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
alter table message_deliveries
    ADD COLUMN entity_id text DEFAULT NULL,
    ADD CONSTRAINT message_deliveries_2_entities
        FOREIGN KEY (entity_id) REFERENCES entities (id) ON DELETE CASCADE ON UPDATE CASCADE;

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back


