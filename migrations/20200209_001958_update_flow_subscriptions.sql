-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
DELETE
FROM flow_subscriptions
WHERE topic = '';

ALTER TABLE flow_subscriptions
    ADD CONSTRAINT topic_at_flow_subscriptions_check CHECK (char_length(topic) <= 512 and char_length(topic) > 0);

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
ALTER TABLE flow_subscriptions
    DROP CONSTRAINT IF EXISTS topic_at_flow_subscriptions_check;

