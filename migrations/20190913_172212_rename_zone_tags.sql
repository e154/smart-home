-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
ALTER TABLE zone_tags
    RENAME TO map_zones;

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
ALTER TABLE map_zones
    RENAME TO zone_tags;

