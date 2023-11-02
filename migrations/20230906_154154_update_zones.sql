-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
alter table areas
--     add column
--         polygon geography(Polygon, 4326) default null,
    add column
        payload jsonb                    default '{}'::jsonb;

-- CREATE INDEX IF NOT EXISTS polygon_geo_gist ON areas USING gist (polygon);

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
alter table areas
--     drop column polygon cascade,
    drop column payload cascade;

