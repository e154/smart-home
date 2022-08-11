-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
create unique index name_at_areas_unq
    on areas (name);

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back


