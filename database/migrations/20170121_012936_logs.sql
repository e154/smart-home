-- +migrate Up
CREATE TABLE logs (
id Int( 32 ) AUTO_INCREMENT NOT NULL,
body Text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
created_at DateTime NOT NULL,
level Enum( 'Emergency', 'Alert', 'Critical', 'Error', 'Warning', 'Notice', 'Info', 'Debug' ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT 'Info',
PRIMARY KEY ( id ),
CONSTRAINT unique_id UNIQUE( id ) )
CHARACTER SET = utf8
COLLATE = utf8_general_ci
ENGINE = InnoDB;

-- +migrate Down
DROP TABLE IF EXISTS `logs` CASCADE;
