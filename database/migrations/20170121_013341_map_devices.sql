-- +migrate Up
CREATE TABLE map_devices (
id Int( 32 ) AUTO_INCREMENT NOT NULL,
created_at DateTime NOT NULL,
update_at DateTime NOT NULL,
device_id Int( 32 ) NOT NULL,
image_id Int( 32 ) NULL,
PRIMARY KEY ( id ),
CONSTRAINT unique_id UNIQUE( device_id ) )
CHARACTER SET = utf8
COLLATE = utf8_general_ci
ENGINE = InnoDB;

-- +migrate Down
DROP TABLE IF EXISTS `map_devices` CASCADE;
