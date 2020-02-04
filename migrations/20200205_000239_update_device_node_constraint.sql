-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
alter table devices
    drop constraint
        devices_2_nodes_fk;

ALTER TABLE devices
    ADD CONSTRAINT devices_2_nodes_fk FOREIGN KEY (node_id)
        REFERENCES nodes (id) ON UPDATE CASCADE ON DELETE RESTRICT;

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
alter table devices
    drop constraint
        devices_2_nodes_fk;

ALTER TABLE devices
    ADD CONSTRAINT devices_2_nodes_fk FOREIGN KEY (node_id)
        REFERENCES nodes (id) ON UPDATE CASCADE ON DELETE CASCADE;
