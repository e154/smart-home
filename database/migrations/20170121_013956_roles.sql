-- +migrate Up
CREATE TABLE roles (
name VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
description VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
parent VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NULL,
created_at DateTime NOT NULL,
update_at DateTime NOT NULL,
PRIMARY KEY ( name ),
CONSTRAINT unique_name UNIQUE( name ) )
CHARACTER SET = utf8
COLLATE = utf8_general_ci
ENGINE = InnoDB;

-- +migrate Down
DROP TABLE IF EXISTS `roles` CASCADE;
