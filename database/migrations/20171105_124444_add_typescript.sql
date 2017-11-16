-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
ALTER TABLE `scripts` CHANGE COLUMN `lang` `lang` varchar(255) NOT NULL DEFAULT 'javascript';

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
ALTER TABLE `scripts` CHANGE COLUMN `lang` `lang` Enum( 'lua', 'coffeescript', 'javascript' ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT 'lua';
