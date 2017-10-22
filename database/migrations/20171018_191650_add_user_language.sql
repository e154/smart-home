-- +migrate Up
ALTER TABLE `users` ADD COLUMN `lang` VarChar( 255 ) NOT NULL DEFAULT 'en';

-- +migrate Down
ALTER TABLE `users` DROP COLUMN `lang`;

