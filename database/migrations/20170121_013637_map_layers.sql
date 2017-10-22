-- +migrate Up
CREATE TABLE map_layers (
id Int( 11 ) AUTO_INCREMENT NOT NULL,
name VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
description Text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
created_at DateTime NOT NULL,
update_at DateTime NOT NULL,
map_id Int( 11 ) NOT NULL,
status VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
weight Int( 11 ) NOT NULL DEFAULT '0',
PRIMARY KEY ( id ),
CONSTRAINT unique_id UNIQUE( id ) )
CHARACTER SET = utf8
COLLATE = utf8_general_ci
ENGINE = InnoDB;

-- +migrate Down
DROP TABLE IF EXISTS `map_layers` CASCADE;
