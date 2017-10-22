-- +migrate Up
CREATE TABLE map_device_states (
id Int( 32 ) AUTO_INCREMENT NOT NULL,
device_state_id Int( 32 ) NOT NULL,
map_device_id Int( 32 ) NOT NULL,
image_id Int( 32 ) NULL,
created_at DateTime NOT NULL,
update_at DateTime NOT NULL,
style Text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
PRIMARY KEY ( id ) )
CHARACTER SET = utf8
COLLATE = utf8_general_ci
ENGINE = InnoDB;

-- +migrate Down
DROP TABLE IF EXISTS `map_device_states` CASCADE;
