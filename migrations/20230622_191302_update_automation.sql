-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
ALTER TABLE actions
    ADD COLUMN entity_id text DEFAULT NULL,
    ADD COLUMN entity_action_name text DEFAULT NULL;
ALTER TABLE actions
    ADD CONSTRAINT actions_2_entities
        FOREIGN KEY (entity_id) REFERENCES entities (id) ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE actions
    ALTER COLUMN script_id DROP NOT NULL,
    ALTER COLUMN script_id SET DEFAULT NULL;

ALTER TABLE triggers
    ALTER COLUMN entity_id DROP NOT NULL,
    ALTER COLUMN entity_id SET DEFAULT NULL,
    ALTER COLUMN script_id DROP NOT NULL,
    ALTER COLUMN script_id SET DEFAULT NULL;

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back


