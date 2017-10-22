-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE nodes (
id Int( 255 ) AUTO_INCREMENT NOT NULL,
ip VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
port Int( 255 ) NOT NULL,
name VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
description Text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
created_at DateTime NOT NULL,
update_at DateTime NOT NULL,
status Enum( 'enabled', 'disabled' ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT 'enabled',
PRIMARY KEY ( id ),
CONSTRAINT node UNIQUE( port, ip ),
CONSTRAINT unique_id UNIQUE( id ) )
CHARACTER SET = utf8
COLLATE = utf8_general_ci
ENGINE = InnoDB;


-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE IF EXISTS `nodes` CASCADE;
