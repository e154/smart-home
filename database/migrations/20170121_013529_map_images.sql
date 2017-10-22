-- +migrate Up
CREATE TABLE map_images (
id Int( 32 ) AUTO_INCREMENT NOT NULL,
image_id Int( 32 ) NOT NULL,
style Text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
PRIMARY KEY ( id ),
CONSTRAINT unique_id UNIQUE( id ) )
CHARACTER SET = utf8
COLLATE = utf8_general_ci
ENGINE = InnoDB;

-- +migrate Down
DROP TABLE IF EXISTS `map_images` CASCADE;
