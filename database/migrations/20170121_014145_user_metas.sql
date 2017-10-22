-- +migrate Up
CREATE TABLE `user_metas` ( 
  `id` Int( 32 ) AUTO_INCREMENT NOT NULL,
  `user_id` Int( 32 ) NOT NULL,
  `key` VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
  `value` VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
  PRIMARY KEY ( `id` ),
  CONSTRAINT `unique_id` UNIQUE( `id` ) ) 
  CHARACTER SET = utf8 
  COLLATE = utf8_general_ci 
  ENGINE = InnoDB 
  AUTO_INCREMENT = 57;

-- +migrate Down
DROP TABLE IF EXISTS `user_metas` CASCADE;
