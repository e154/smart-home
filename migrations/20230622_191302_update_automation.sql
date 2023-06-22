-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
ALTER TABLE actions
    ADD COLUMN entity_action_id bigint DEFAULT NULL;
ALTER TABLE actions
    ADD CONSTRAINT actions_2_entity_actions
        FOREIGN KEY (entity_action_id) REFERENCES entity_actions (id) ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE actions
    ALTER COLUMN script_id DROP NOT NULL,
    ALTER COLUMN script_id SET DEFAULT NULL;

ALTER TABLE triggers
    ALTER COLUMN script_id DROP NOT NULL,
    ALTER COLUMN script_id SET DEFAULT NULL;

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back


