-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
INSERT INTO `nodes` VALUES ('1', '127.0.0.1', '3000', 'New node', '', '2017-10-30 14:42:32', '2017-10-30 14:42:32', 'enabled');

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
DELETE FROM `nodes` WHERE  id in (1);

