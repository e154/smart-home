-- +migrate Up
CREATE TABLE messages (
id Int( 32 ) AUTO_INCREMENT NOT NULL,
type Enum( 'sms', 'email', 'push' ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT 'email',
email_title VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NULL,
email_body Text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
sms_text Text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
ui_text Text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
scheduled_at DateTime NOT NULL,
created_at DateTime NOT NULL,
update_at DateTime NOT NULL,
deleted_at DateTime NULL,
PRIMARY KEY ( id ),
CONSTRAINT unique_id UNIQUE( id ) )
CHARACTER SET = utf8
COLLATE = utf8_general_ci
ENGINE = InnoDB;

-- +migrate Down
DROP TABLE IF EXISTS `messages` CASCADE;
