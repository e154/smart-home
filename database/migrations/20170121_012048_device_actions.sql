-- +migrate Up
CREATE TABLE device_actions (
id Int( 11 ) AUTO_INCREMENT NOT NULL,
device_id Int( 11 ) NULL,
name VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
description VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NULL,
created_at DateTime NOT NULL,
update_at DateTime NOT NULL,
script_id Int( 32 ) NULL,
PRIMARY KEY ( id ),
CONSTRAINT unique_id UNIQUE( id ) )
CHARACTER SET = utf8
COLLATE = utf8_general_ci
ENGINE = InnoDB;

-- +migrate Down
DROP TABLE IF EXISTS `device_actions` CASCADE;
