-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
drop table map_devices cascade;
drop table map_elements cascade;
drop table map_images cascade;
drop table map_layers cascade;
drop table map_texts cascade;
drop table maps cascade;
drop table storage cascade;

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back


