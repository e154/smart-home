-- +migrate Up
CREATE TABLE map_texts (
id Int( 11 ) AUTO_INCREMENT NOT NULL,
Text Text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
Style Text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
PRIMARY KEY ( id ),
CONSTRAINT unique_id UNIQUE( id ) )
CHARACTER SET = utf8
COLLATE = utf8_general_ci
ENGINE = InnoDB;

-- +migrate Down
DROP TABLE IF EXISTS `map_texts` CASCADE;
