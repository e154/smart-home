-- +migrate Up
CREATE TABLE scripts (
id Int( 32 ) AUTO_INCREMENT NOT NULL,
lang Enum( 'lua', 'coffeescript', 'javascript' ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT 'lua',
name VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
source Text CHARACTER SET utf8 COLLATE utf8_general_ci NULL,
created_at DateTime NOT NULL,
update_at DateTime NOT NULL,
description Text CHARACTER SET utf8 COLLATE utf8_general_ci NULL,
compiled Text CHARACTER SET utf8 COLLATE utf8_general_ci NULL,
PRIMARY KEY ( id ),
CONSTRAINT unique_id UNIQUE( id ),
CONSTRAINT unique_name UNIQUE( name ) )
CHARACTER SET = utf8
COLLATE = utf8_general_ci
ENGINE = InnoDB;

-- +migrate Down
DROP TABLE IF EXISTS `scripts` CASCADE;
