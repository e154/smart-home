-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE `variables` (
	`name` VarChar( 128 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT 'The name of the variable.',
	`value` LongBlob NOT NULL COMMENT 'The value of the variable.',
	`autoload` VarChar( 20 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT 'yes',
	PRIMARY KEY ( `name` ) )
CHARACTER SET = utf8
COLLATE = utf8_general_ci
ENGINE = InnoDB;

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE IF EXISTS `variables` CASCADE;
