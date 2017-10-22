-- +migrate Up
CREATE TABLE workers (
id Int( 11 ) AUTO_INCREMENT NOT NULL,
flow_id Int( 11 ) NOT NULL,
workflow_id Int( 11 ) NOT NULL,
name VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
device_action_id Int( 11 ) NOT NULL,
time VarChar( 254 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
created_at DateTime NOT NULL,
update_at DateTime NOT NULL,
status Enum( 'enabled', 'disabled' ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT 'enabled',
PRIMARY KEY ( id ) )
CHARACTER SET = utf8
COLLATE = utf8_general_ci
ENGINE = InnoDB;

-- +migrate Down
DROP TABLE IF EXISTS `workers` CASCADE;
