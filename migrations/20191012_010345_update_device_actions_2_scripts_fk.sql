-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
alter table flow_elements
    drop constraint
        flow_elements_2_scripts_fk;

ALTER TABLE flow_elements
    ADD CONSTRAINT flow_elements_2_scripts_fk FOREIGN KEY (script_id)
        REFERENCES scripts (id) ON UPDATE CASCADE ON DELETE SET NULL;

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
