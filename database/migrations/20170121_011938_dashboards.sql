-- +migrate Up
CREATE TABLE dashboards (
id Int( 32 ) AUTO_INCREMENT NOT NULL,
name VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
description Text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
created_at DateTime NOT NULL,
update_at DateTime NOT NULL,
widgets Text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
PRIMARY KEY ( id ),
CONSTRAINT unique_id UNIQUE( id ) )
CHARACTER SET = utf8
COLLATE = utf8_general_ci
ENGINE = InnoDB;

-- +migrate Down
DROP TABLE IF EXISTS `dashboards` CASCADE;
