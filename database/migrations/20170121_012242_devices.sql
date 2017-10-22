-- +migrate Up
CREATE TABLE devices (
id Int( 11 ) AUTO_INCREMENT NOT NULL,
name VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
description VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
device_id Int( 11 ) NULL,
created_at DateTime NOT NULL,
update_at DateTime NOT NULL,
node_id Int( 11 ) NULL,
baud Int( 11 ) NULL,
tty VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NULL,
stop_bite Int( 8 ) NULL,
timeout Int( 11 ) NULL,
address Int( 11 ) NULL,
status VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT 'enabled',
sleep Int( 32 ) NOT NULL DEFAULT '0',
PRIMARY KEY ( id ),
CONSTRAINT unique3 UNIQUE( name, device_id ),
CONSTRAINT unique2 UNIQUE( node_id, address ),
CONSTRAINT unique1 UNIQUE( device_id, address ) )
CHARACTER SET = utf8
COLLATE = utf8_general_ci
ENGINE = InnoDB;

-- +migrate Down
DROP TABLE IF EXISTS `devices` CASCADE;
