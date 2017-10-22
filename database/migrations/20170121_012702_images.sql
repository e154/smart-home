-- +migrate Up
CREATE TABLE images (
id Int( 32 ) AUTO_INCREMENT NOT NULL,
thumb VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
image VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
mime_type VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
title VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
size Int( 11 ) NOT NULL,
name VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
created_at DateTime NOT NULL,
PRIMARY KEY ( id ),
CONSTRAINT id UNIQUE( id ) )
CHARACTER SET = utf8
COLLATE = utf8_general_ci
ENGINE = InnoDB;

-- +migrate Down
DROP TABLE IF EXISTS `images` CASCADE;
