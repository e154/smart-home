-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
alter table devices
  drop column device_id cascade;
alter table devices
  add column device_id BIGINT default null CONSTRAINT device_2_devices_fk REFERENCES devices (id) on update cascade on delete cascade;

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
alter table devices
  drop column device_id cascade;
alter table devices
  add column device_id smallint NULL;