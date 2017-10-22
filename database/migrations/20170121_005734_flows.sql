-- +migrate Up
CREATE TABLE flows (
id Int( 32 ) AUTO_INCREMENT NOT NULL,
name VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
status Enum( 'enabled', 'disabled' ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT 'enabled',
created_at DateTime NOT NULL,
update_at DateTime NOT NULL,
workflow_id Int( 11 ) NOT NULL,
description Text CHARACTER SET utf8 COLLATE utf8_general_ci NULL,
PRIMARY KEY ( id ),
CONSTRAINT unique_id UNIQUE( id ) )
CHARACTER SET = utf8
COLLATE = utf8_general_ci
ENGINE = InnoDB;

-- +migrate Down
DROP TABLE IF EXISTS `flows` CASCADE;
