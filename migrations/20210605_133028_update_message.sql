-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
ALTER TABLE messages
    ALTER COLUMN type TYPE text;

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back


