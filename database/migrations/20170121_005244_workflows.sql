-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied

CREATE TABLE workflows (
id Int( 255 ) AUTO_INCREMENT NOT NULL,
name VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
status Enum( 'enabled', 'disabled' ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT 'enabled',
created_at DateTime NOT NULL,
update_at DateTime NOT NULL,
description Text CHARACTER SET utf8 COLLATE utf8_general_ci NULL,
PRIMARY KEY ( id ),
CONSTRAINT unique_name UNIQUE( name ) )
CHARACTER SET = utf8
COLLATE = utf8_general_ci
ENGINE = InnoDB;


-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE IF EXISTS `workflows` CASCADE;
