-- +migrate Up
CREATE TABLE map_device_actions (
id Int( 32 ) AUTO_INCREMENT NOT NULL,
device_action_id Int( 32 ) NOT NULL,
map_device_id Int( 2 ) NOT NULL,
image_id Int( 32 ) NULL,
type VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
created_at DateTime NOT NULL,
update_at DateTime NOT NULL,
PRIMARY KEY ( id ),
CONSTRAINT unique_id UNIQUE( id ) )
CHARACTER SET = utf8
COLLATE = utf8_general_ci
ENGINE = InnoDB;

-- +migrate Down
DROP TABLE IF EXISTS `map_device_actions` CASCADE;
