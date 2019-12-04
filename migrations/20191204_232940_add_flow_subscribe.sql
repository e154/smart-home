-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied

create table flow_subscriptions
(
    id         BIGSERIAL PRIMARY KEY,
    flow_id    BIGINT      NOT NULL
        CONSTRAINT flow_subscriptions_2_flows_fk REFERENCES flows (id) ON UPDATE CASCADE ON DELETE CASCADE,
    topic      text        NOT NULL
        CONSTRAINT topic_2_flow_subscriptions_check check (char_length(topic) <= 512),
    created_at TIMESTAMPTZ NOT NULL
);

CREATE UNIQUE INDEX topic_at_flow_subscriptions_unq
    ON flow_subscriptions (topic, flow_id);

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
drop table if exists flow_subscriptions cascade;

