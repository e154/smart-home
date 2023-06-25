-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
alter table dashboard_cards
    add column hidden    bool NOT NULL DEFAULT FALSE,
    ADD COLUMN entity_id text DEFAULT NULL,
    ADD CONSTRAINT dashboard_cards_2_entities
        FOREIGN KEY (entity_id) REFERENCES entities (id) ON DELETE CASCADE ON UPDATE CASCADE;

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back


