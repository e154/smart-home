-- +migrate Up
CREATE TABLE flow_elements (
uuid VarChar( 254 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
name VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
graph_settings Text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
flow_id Int( 32 ) NOT NULL,
created_at DateTime NOT NULL,
update_at DateTime NOT NULL,
prototype_type VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT 'default',
status VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NULL,
description Text CHARACTER SET utf8 COLLATE utf8_general_ci NULL,
script_id Int( 11 ) NULL,
flow_link Int( 32 ) NULL,
PRIMARY KEY ( uuid ),
CONSTRAINT id UNIQUE( uuid ),
CONSTRAINT unique_id UNIQUE( uuid ) )
CHARACTER SET = utf8
COLLATE = utf8_general_ci
ENGINE = InnoDB;

-- +migrate Down
DROP TABLE IF EXISTS `flow_elements` CASCADE;
