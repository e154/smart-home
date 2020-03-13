-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE zone_tags (
    id   BIGSERIAL PRIMARY KEY,
    name text NOT NULL
);

CREATE UNIQUE INDEX name_at_zone_tags_unq
    ON zone_tags (name);

ALTER TABLE map_elements
    ADD COLUMN zone_id BIGINT NULL;

ALTER TABLE map_elements
    ADD CONSTRAINT map_elements_2_zone_tags_fk FOREIGN KEY (zone_id)
        REFERENCES zone_tags (id) ON UPDATE CASCADE ON DELETE SET NULL;

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
ALTER TABLE map_elements
    DROP COLUMN IF EXISTS zone_id;

DROP TABLE IF EXISTS zone_tags CASCADE;

