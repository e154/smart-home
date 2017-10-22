-- +migrate Up
CREATE TABLE email_templates (
name VarChar( 64 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT 'Название',
description VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NULL COMMENT 'Описание',
content LongText CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT 'Содержимое',
status VarChar( 64 ) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT 'active' COMMENT 'active, unactive',
type VarChar( 64 ) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT 'item' COMMENT 'item, template',
parent VarChar( 64 ) CHARACTER SET utf8 COLLATE utf8_general_ci NULL,
created_at DateTime NOT NULL,
updated_at DateTime NOT NULL,
PRIMARY KEY ( name ),
CONSTRAINT name UNIQUE( name ) )
CHARACTER SET = utf8
COLLATE = utf8_general_ci
ENGINE = InnoDB;

-- +migrate Down
DROP TABLE IF EXISTS `email_templates` CASCADE;
