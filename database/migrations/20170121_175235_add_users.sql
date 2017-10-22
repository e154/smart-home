-- +migrate Up

INSERT INTO `roles`(`name`,`description`,`parent`,`created_at`,`update_at`) VALUES ( 'demo', '', NULL, '2017-01-15 05:17:09', '2017-01-15 05:17:09' );
INSERT INTO `roles`(`name`,`description`,`parent`,`created_at`,`update_at`) VALUES ( 'user', '', 'demo', '2017-01-15 05:17:36', '2017-01-19 19:39:54' );
INSERT INTO `roles`(`name`,`description`,`parent`,`created_at`,`update_at`) VALUES ( 'admin', '', 'user', '2017-01-15 05:20:58', '2017-01-15 05:21:29' );

INSERT INTO `users`(`id`,`nickname`,`first_name`,`last_name`,`encrypted_password`,`email`,`history`,`status`,`reset_password_token`,`authentication_token`,`image_id`,`sign_in_count`,`current_sign_in_ip`,`last_sign_in_ip`,`user_id`,`role_name`,`reset_password_sent_at`,`current_sign_in_at`,`last_sign_in_at`,`created_at`,`update_at`,`deleted`) VALUES ( '1', 'admin', '', '', 'f6fdffe48c908deb0f4c3bd36c032e72', 'admin@e154.ru', '[]', 'active', '', 'xlzEaHNBbn80OmTfWd1z18XpNUlZikdb4z5fo5YAxlNv3CfWxs', NULL, '111', '127.0.0.1', '127.0.0.1', NULL, 'admin', NULL, '2017-01-21 11:11:26', '2017-01-21 10:54:03', '2017-01-15 05:25:07', '2017-01-21 11:11:26', NULL );
INSERT INTO `users`(`id`,`nickname`,`first_name`,`last_name`,`encrypted_password`,`email`,`history`,`status`,`reset_password_token`,`authentication_token`,`image_id`,`sign_in_count`,`current_sign_in_ip`,`last_sign_in_ip`,`user_id`,`role_name`,`reset_password_sent_at`,`current_sign_in_at`,`last_sign_in_at`,`created_at`,`update_at`,`deleted`) VALUES ( '2', 'demo', '', '', 'c514c91e4ed341f263e458d44b3bb0a7', 'demo@e154.ru', '[]', 'active', '', '5SLTHOzN1hWw6jhgEw0y9JbtwdBIK5mgW3DLt5FYy23zNkVnvW', NULL, '8', '127.0.0.1', '127.0.0.1', NULL, 'demo', NULL, '2017-01-21 11:11:43', '2017-01-20 17:28:23', '2017-01-18 17:13:28', '2017-01-21 11:11:43', NULL );

INSERT INTO `user_metas`(`id`,`user_id`,`key`,`value`) VALUES ( '49', '1', 'phone1', '' );
INSERT INTO `user_metas`(`id`,`user_id`,`key`,`value`) VALUES ( '50', '1', 'phone2', '' );
INSERT INTO `user_metas`(`id`,`user_id`,`key`,`value`) VALUES ( '51', '1', 'phone3', '' );
INSERT INTO `user_metas`(`id`,`user_id`,`key`,`value`) VALUES ( '52', '1', 'autograph', '' );
INSERT INTO `user_metas`(`id`,`user_id`,`key`,`value`) VALUES ( '53', '2', 'phone1', '' );
INSERT INTO `user_metas`(`id`,`user_id`,`key`,`value`) VALUES ( '54', '2', 'phone2', '' );
INSERT INTO `user_metas`(`id`,`user_id`,`key`,`value`) VALUES ( '55', '2', 'phone3', '' );
INSERT INTO `user_metas`(`id`,`user_id`,`key`,`value`) VALUES ( '56', '2', 'autograph', '' );

INSERT INTO `permissions`(`id`,`role_name`,`package_name`,`level_name`) VALUES ( '81', 'admin', 'ws', 'read' );
INSERT INTO `permissions`(`id`,`role_name`,`package_name`,`level_name`) VALUES ( '82', 'admin', 'workflow', 'create' );
INSERT INTO `permissions`(`id`,`role_name`,`package_name`,`level_name`) VALUES ( '83', 'admin', 'workflow', 'delete' );
INSERT INTO `permissions`(`id`,`role_name`,`package_name`,`level_name`) VALUES ( '84', 'admin', 'workflow', 'read' );
INSERT INTO `permissions`(`id`,`role_name`,`package_name`,`level_name`) VALUES ( '86', 'demo', 'ws', 'read' );
INSERT INTO `permissions`(`id`,`role_name`,`package_name`,`level_name`) VALUES ( '88', 'demo', 'map', 'read_map' );
INSERT INTO `permissions`(`id`,`role_name`,`package_name`,`level_name`) VALUES ( '89', 'demo', 'device', 'read' );
INSERT INTO `permissions`(`id`,`role_name`,`package_name`,`level_name`) VALUES ( '90', 'demo', 'node', 'read' );
INSERT INTO `permissions`(`id`,`role_name`,`package_name`,`level_name`) VALUES ( '91', 'demo', 'dashboard', 'read' );
INSERT INTO `permissions`(`id`,`role_name`,`package_name`,`level_name`) VALUES ( '92', 'demo', 'notifr', 'preview_notifr' );
INSERT INTO `permissions`(`id`,`role_name`,`package_name`,`level_name`) VALUES ( '93', 'demo', 'notifr', 'read_notifr_item' );
INSERT INTO `permissions`(`id`,`role_name`,`package_name`,`level_name`) VALUES ( '94', 'demo', 'notifr', 'read_notifr_template' );
INSERT INTO `permissions`(`id`,`role_name`,`package_name`,`level_name`) VALUES ( '95', 'demo', 'script', 'exec_script' );
INSERT INTO `permissions`(`id`,`role_name`,`package_name`,`level_name`) VALUES ( '96', 'demo', 'script', 'read' );
INSERT INTO `permissions`(`id`,`role_name`,`package_name`,`level_name`) VALUES ( '97', 'demo', 'user', 'read' );
INSERT INTO `permissions`(`id`,`role_name`,`package_name`,`level_name`) VALUES ( '98', 'demo', 'user', 'read_role' );
INSERT INTO `permissions`(`id`,`role_name`,`package_name`,`level_name`) VALUES ( '99', 'demo', 'user', 'read_role_access_list' );
INSERT INTO `permissions`(`id`,`role_name`,`package_name`,`level_name`) VALUES ( '100', 'demo', 'worker', 'read' );
INSERT INTO `permissions`(`id`,`role_name`,`package_name`,`level_name`) VALUES ( '101', 'demo', 'device_action', 'read' );
INSERT INTO `permissions`(`id`,`role_name`,`package_name`,`level_name`) VALUES ( '102', 'demo', 'device_state', 'read' );
INSERT INTO `permissions`(`id`,`role_name`,`package_name`,`level_name`) VALUES ( '103', 'demo', 'flow', 'read' );
INSERT INTO `permissions`(`id`,`role_name`,`package_name`,`level_name`) VALUES ( '104', 'demo', 'image', 'read' );
INSERT INTO `permissions`(`id`,`role_name`,`package_name`,`level_name`) VALUES ( '105', 'demo', 'workflow', 'read' );
INSERT INTO `permissions`(`id`,`role_name`,`package_name`,`level_name`) VALUES ( '106', 'demo', 'log', 'read' );
INSERT INTO `permissions`(`id`,`role_name`,`package_name`,`level_name`) VALUES ( '107', 'demo', 'map', 'read_map_element' );
INSERT INTO `permissions`(`id`,`role_name`,`package_name`,`level_name`) VALUES ( '108', 'demo', 'map', 'read_map_layer' );
INSERT INTO `permissions`(`id`,`role_name`,`package_name`,`level_name`) VALUES ( '109', 'user', 'image', 'upload' );
INSERT INTO `permissions`(`id`,`role_name`,`package_name`,`level_name`) VALUES ( '110', 'user', 'image', 'create' );
INSERT INTO `permissions`(`id`,`role_name`,`package_name`,`level_name`) VALUES ( '111', 'user', 'image', 'delete' );
INSERT INTO `permissions`(`id`,`role_name`,`package_name`,`level_name`) VALUES ( '112', 'user', 'image', 'update' );
INSERT INTO `permissions`(`id`,`role_name`,`package_name`,`level_name`) VALUES ( '113', 'user', 'flow', 'create' );
INSERT INTO `permissions`(`id`,`role_name`,`package_name`,`level_name`) VALUES ( '114', 'user', 'flow', 'delete' );
INSERT INTO `permissions`(`id`,`role_name`,`package_name`,`level_name`) VALUES ( '115', 'user', 'flow', 'update' );
INSERT INTO `permissions`(`id`,`role_name`,`package_name`,`level_name`) VALUES ( '116', 'user', 'device_state', 'create' );
INSERT INTO `permissions`(`id`,`role_name`,`package_name`,`level_name`) VALUES ( '117', 'user', 'device_state', 'delete' );
INSERT INTO `permissions`(`id`,`role_name`,`package_name`,`level_name`) VALUES ( '118', 'user', 'device_state', 'update' );
INSERT INTO `permissions`(`id`,`role_name`,`package_name`,`level_name`) VALUES ( '119', 'user', 'device_action', 'delete' );
INSERT INTO `permissions`(`id`,`role_name`,`package_name`,`level_name`) VALUES ( '120', 'user', 'device_action', 'create' );
INSERT INTO `permissions`(`id`,`role_name`,`package_name`,`level_name`) VALUES ( '121', 'user', 'device', 'update' );
INSERT INTO `permissions`(`id`,`role_name`,`package_name`,`level_name`) VALUES ( '122', 'user', 'device', 'create' );
INSERT INTO `permissions`(`id`,`role_name`,`package_name`,`level_name`) VALUES ( '123', 'user', 'device', 'delete' );
INSERT INTO `permissions`(`id`,`role_name`,`package_name`,`level_name`) VALUES ( '124', 'user', 'dashboard', 'create' );
INSERT INTO `permissions`(`id`,`role_name`,`package_name`,`level_name`) VALUES ( '125', 'user', 'dashboard', 'update' );
INSERT INTO `permissions`(`id`,`role_name`,`package_name`,`level_name`) VALUES ( '126', 'user', 'dashboard', 'delete' );
INSERT INTO `permissions`(`id`,`role_name`,`package_name`,`level_name`) VALUES ( '127', 'user', 'script', 'update' );
INSERT INTO `permissions`(`id`,`role_name`,`package_name`,`level_name`) VALUES ( '128', 'user', 'script', 'create' );
INSERT INTO `permissions`(`id`,`role_name`,`package_name`,`level_name`) VALUES ( '129', 'user', 'script', 'delete' );
INSERT INTO `permissions`(`id`,`role_name`,`package_name`,`level_name`) VALUES ( '130', 'user', 'worker', 'delete' );
INSERT INTO `permissions`(`id`,`role_name`,`package_name`,`level_name`) VALUES ( '131', 'user', 'worker', 'create' );
INSERT INTO `permissions`(`id`,`role_name`,`package_name`,`level_name`) VALUES ( '132', 'user', 'worker', 'update' );
INSERT INTO `permissions`(`id`,`role_name`,`package_name`,`level_name`) VALUES ( '133', 'user', 'workflow', 'create' );
INSERT INTO `permissions`(`id`,`role_name`,`package_name`,`level_name`) VALUES ( '134', 'user', 'workflow', 'delete' );
INSERT INTO `permissions`(`id`,`role_name`,`package_name`,`level_name`) VALUES ( '135', 'user', 'workflow', 'update' );

-- +migrate Down
DELETE FROM `user_metas` WHere user_id IN (1,2);
DELETE FROM `users` WHERE  id in (1,2);
DELETE FROM permissions WHERE role_name IN ("admin", "user", "demo");
DELETE FROM roles WHERE  name in ("admin");
DELETE FROM roles WHERE  name in ("user");
DELETE FROM roles WHERE  name in ("demo");

