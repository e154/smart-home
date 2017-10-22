-- +migrate Up
CREATE TABLE message_deliveries (
id Int( 32 ) AUTO_INCREMENT NOT NULL,
message_id Int( 32 ) NOT NULL,
state VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
created_at DateTime NOT NULL,
update_at DateTime NOT NULL,
error_system_code VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
error_system_message Text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
address Text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
PRIMARY KEY ( id ),
CONSTRAINT unique_id UNIQUE( id ) )
CHARACTER SET = utf8
COLLATE = utf8_general_ci
ENGINE = InnoDB;

CREATE INDEX `lnk_messages_message_deliveries` USING BTREE ON `message_deliveries`( `message_id` );
ALTER TABLE `message_deliveries` ADD CONSTRAINT `lnk_messages_message_deliveries` FOREIGN KEY ( `message_id` ) REFERENCES `messages`( `id` ) ON DELETE Cascade ON UPDATE Cascade;

-- +migrate Down
ALTER TABLE `message_deliveries` DROP FOREIGN KEY `lnk_messages_message_deliveries`;
DROP TABLE IF EXISTS `message_deliveries` CASCADE;
