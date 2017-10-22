-- +migrate Up
CREATE TABLE users (
id Int( 32 ) AUTO_INCREMENT NOT NULL,
nickname VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
first_name VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
last_name VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
encrypted_password VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
email VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
history Text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
status VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
reset_password_token VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
authentication_token VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
image_id Int( 255 ) NULL,
sign_in_count Int( 32 ) NOT NULL,
current_sign_in_ip VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
last_sign_in_ip VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NULL,
user_id Int( 32 ) NULL,
role_name VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
reset_password_sent_at DateTime NULL,
current_sign_in_at DateTime NULL,
last_sign_in_at DateTime NULL,
created_at DateTime NOT NULL,
update_at DateTime NOT NULL,
deleted DateTime NULL,
PRIMARY KEY ( id ),
CONSTRAINT unique_email UNIQUE( email ),
CONSTRAINT unique_id UNIQUE( id ),
CONSTRAINT unique_nickname UNIQUE( nickname ) )
CHARACTER SET = utf8
COLLATE = utf8_general_ci
ENGINE = InnoDB;

-- +migrate Down
DROP TABLE IF EXISTS `users` CASCADE;
