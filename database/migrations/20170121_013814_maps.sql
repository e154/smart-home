-- +migrate Up
CREATE TABLE maps (
id Int( 11 ) AUTO_INCREMENT NOT NULL,
created_at DateTime NOT NULL,
update_at DateTime NOT NULL,
name Char( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
description Text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
options Text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
PRIMARY KEY ( id ),
CONSTRAINT unique_id UNIQUE( id ) )
CHARACTER SET = utf8
COLLATE = utf8_general_ci
ENGINE = InnoDB;

-- +migrate Down
DROP TABLE IF EXISTS `maps` CASCADE;
