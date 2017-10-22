-- +migrate Up
CREATE TABLE connections (
uuid VarChar( 254 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
name VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NULL,
element_from VarChar( 254 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
element_to VarChar( 254 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
flow_id Int( 32 ) NOT NULL,
created_at DateTime NOT NULL,
update_at DateTime NOT NULL,
graph_settings Text CHARACTER SET utf8 COLLATE utf8_general_ci NULL,
point_from Int( 11 ) NOT NULL,
point_to Int( 11 ) NOT NULL,
direction VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
PRIMARY KEY ( uuid ),
CONSTRAINT unique_id UNIQUE( uuid ) )
CHARACTER SET = utf8
COLLATE = utf8_general_ci
ENGINE = InnoDB;

-- +migrate Down
DROP TABLE IF EXISTS `connections` CASCADE;
