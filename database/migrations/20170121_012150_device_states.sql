-- +migrate Up
CREATE TABLE device_states (
id Int( 11 ) AUTO_INCREMENT NOT NULL,
system_name VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
description Text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
created_at DateTime NOT NULL,
update_at DateTime NOT NULL,
device_id Int( 11 ) NOT NULL,
PRIMARY KEY ( id ),
CONSTRAINT unique_1 UNIQUE( device_id, system_name ) )
CHARACTER SET = utf8
COLLATE = utf8_general_ci
ENGINE = InnoDB;

-- +migrate Down
DROP TABLE IF EXISTS `device_actions` CASCADE;
