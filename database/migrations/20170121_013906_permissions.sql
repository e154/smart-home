-- +migrate Up
CREATE TABLE permissions (
id Int( 32 ) AUTO_INCREMENT NOT NULL,
role_name VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
package_name VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
level_name VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
PRIMARY KEY ( id ) )
CHARACTER SET = utf8
COLLATE = utf8_general_ci
ENGINE = InnoDB;

-- +migrate Down
DROP TABLE IF EXISTS `permissions` CASCADE;
