package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type AddUsers_20170121_175235 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &AddUsers_20170121_175235{}
	m.Created = "20170121_175235"
	migration.Register("AddUsers_20170121_175235", m)
}

// Run the migrations
func (m *AddUsers_20170121_175235) Up() {
	// use m.SQL("CREATE TABLE ...") to make schema update

	m.SQL("INSERT INTO `roles`(`name`,`description`,`parent`,`created_at`,`update_at`) VALUES ( 'demo', '', NULL, '2017-01-15 05:17:09', '2017-01-15 05:17:09' );")
	m.SQL("INSERT INTO `roles`(`name`,`description`,`parent`,`created_at`,`update_at`) VALUES ( 'user', '', 'demo', '2017-01-15 05:17:36', '2017-01-19 19:39:54' );")
	m.SQL("INSERT INTO `roles`(`name`,`description`,`parent`,`created_at`,`update_at`) VALUES ( 'admin', '', 'user', '2017-01-15 05:20:58', '2017-01-15 05:21:29' );")

	m.SQL("INSERT INTO `users`(`id`,`nickname`,`first_name`,`last_name`,`encrypted_password`,`email`,`history`,`status`,`reset_password_token`,`authentication_token`,`image_id`,`sign_in_count`,`current_sign_in_ip`,`last_sign_in_ip`,`user_id`,`role_name`,`reset_password_sent_at`,`current_sign_in_at`,`last_sign_in_at`,`created_at`,`update_at`,`deleted`) VALUES ( '1', 'admin', '', '', 'f6fdffe48c908deb0f4c3bd36c032e72', 'admin@e154.ru', '[]', 'active', '', 'xlzEaHNBbn80OmTfWd1z18XpNUlZikdb4z5fo5YAxlNv3CfWxs', NULL, '111', '127.0.0.1', '127.0.0.1', NULL, 'admin', NULL, '2017-01-21 11:11:26', '2017-01-21 10:54:03', '2017-01-15 05:25:07', '2017-01-21 11:11:26', NULL );")
	m.SQL("INSERT INTO `users`(`id`,`nickname`,`first_name`,`last_name`,`encrypted_password`,`email`,`history`,`status`,`reset_password_token`,`authentication_token`,`image_id`,`sign_in_count`,`current_sign_in_ip`,`last_sign_in_ip`,`user_id`,`role_name`,`reset_password_sent_at`,`current_sign_in_at`,`last_sign_in_at`,`created_at`,`update_at`,`deleted`) VALUES ( '2', 'demo', '', '', 'c514c91e4ed341f263e458d44b3bb0a7', 'demo@e154.ru', '[]', 'active', '', '5SLTHOzN1hWw6jhgEw0y9JbtwdBIK5mgW3DLt5FYy23zNkVnvW', NULL, '8', '127.0.0.1', '127.0.0.1', NULL, 'demo', NULL, '2017-01-21 11:11:43', '2017-01-20 17:28:23', '2017-01-18 17:13:28', '2017-01-21 11:11:43', NULL );")

	m.SQL("INSERT INTO `user_metas`(`id`,`user_id`,`key`,`value`) VALUES ( '49', '1', 'phone1', '' );")
	m.SQL("INSERT INTO `user_metas`(`id`,`user_id`,`key`,`value`) VALUES ( '50', '1', 'phone2', '' );")
	m.SQL("INSERT INTO `user_metas`(`id`,`user_id`,`key`,`value`) VALUES ( '51', '1', 'phone3', '' );")
	m.SQL("INSERT INTO `user_metas`(`id`,`user_id`,`key`,`value`) VALUES ( '52', '1', 'autograph', '' );")
	m.SQL("INSERT INTO `user_metas`(`id`,`user_id`,`key`,`value`) VALUES ( '53', '2', 'phone1', '' );")
	m.SQL("INSERT INTO `user_metas`(`id`,`user_id`,`key`,`value`) VALUES ( '54', '2', 'phone2', '' );")
	m.SQL("INSERT INTO `user_metas`(`id`,`user_id`,`key`,`value`) VALUES ( '55', '2', 'phone3', '' );")
	m.SQL("INSERT INTO `user_metas`(`id`,`user_id`,`key`,`value`) VALUES ( '56', '2', 'autograph', '' );")

	m.SQL("INSERT INTO `permissions`(`id`,`role_name`,`package_name`,`level_name`) VALUES ( '81', 'admin', 'ws', 'read' );")
	m.SQL("INSERT INTO `permissions`(`id`,`role_name`,`package_name`,`level_name`) VALUES ( '82', 'admin', 'workflow', 'create' );")
	m.SQL("INSERT INTO `permissions`(`id`,`role_name`,`package_name`,`level_name`) VALUES ( '83', 'admin', 'workflow', 'delete' );")
	m.SQL("INSERT INTO `permissions`(`id`,`role_name`,`package_name`,`level_name`) VALUES ( '84', 'admin', 'workflow', 'read' );")
	m.SQL("INSERT INTO `permissions`(`id`,`role_name`,`package_name`,`level_name`) VALUES ( '86', 'demo', 'ws', 'read' );")
	m.SQL("INSERT INTO `permissions`(`id`,`role_name`,`package_name`,`level_name`) VALUES ( '88', 'demo', 'map', 'read_map' );")
	m.SQL("INSERT INTO `permissions`(`id`,`role_name`,`package_name`,`level_name`) VALUES ( '89', 'demo', 'device', 'read' );")
	m.SQL("INSERT INTO `permissions`(`id`,`role_name`,`package_name`,`level_name`) VALUES ( '90', 'demo', 'node', 'read' );")
	m.SQL("INSERT INTO `permissions`(`id`,`role_name`,`package_name`,`level_name`) VALUES ( '91', 'demo', 'dashboard', 'read' );")
	m.SQL("INSERT INTO `permissions`(`id`,`role_name`,`package_name`,`level_name`) VALUES ( '92', 'demo', 'notifr', 'preview_notifr' );")
	m.SQL("INSERT INTO `permissions`(`id`,`role_name`,`package_name`,`level_name`) VALUES ( '93', 'demo', 'notifr', 'read_notifr_item' );")
	m.SQL("INSERT INTO `permissions`(`id`,`role_name`,`package_name`,`level_name`) VALUES ( '94', 'demo', 'notifr', 'read_notifr_template' );")
	m.SQL("INSERT INTO `permissions`(`id`,`role_name`,`package_name`,`level_name`) VALUES ( '95', 'demo', 'script', 'exec_script' );")
	m.SQL("INSERT INTO `permissions`(`id`,`role_name`,`package_name`,`level_name`) VALUES ( '96', 'demo', 'script', 'read' );")
	m.SQL("INSERT INTO `permissions`(`id`,`role_name`,`package_name`,`level_name`) VALUES ( '97', 'demo', 'user', 'read' );")
	m.SQL("INSERT INTO `permissions`(`id`,`role_name`,`package_name`,`level_name`) VALUES ( '98', 'demo', 'user', 'read_role' );")
	m.SQL("INSERT INTO `permissions`(`id`,`role_name`,`package_name`,`level_name`) VALUES ( '99', 'demo', 'user', 'read_role_access_list' );")
	m.SQL("INSERT INTO `permissions`(`id`,`role_name`,`package_name`,`level_name`) VALUES ( '100', 'demo', 'worker', 'read' );")
	m.SQL("INSERT INTO `permissions`(`id`,`role_name`,`package_name`,`level_name`) VALUES ( '101', 'demo', 'device_action', 'read' );")
	m.SQL("INSERT INTO `permissions`(`id`,`role_name`,`package_name`,`level_name`) VALUES ( '102', 'demo', 'device_state', 'read' );")
	m.SQL("INSERT INTO `permissions`(`id`,`role_name`,`package_name`,`level_name`) VALUES ( '103', 'demo', 'flow', 'read' );")
	m.SQL("INSERT INTO `permissions`(`id`,`role_name`,`package_name`,`level_name`) VALUES ( '104', 'demo', 'image', 'read' );")
	m.SQL("INSERT INTO `permissions`(`id`,`role_name`,`package_name`,`level_name`) VALUES ( '105', 'demo', 'workflow', 'read' );")
	m.SQL("INSERT INTO `permissions`(`id`,`role_name`,`package_name`,`level_name`) VALUES ( '106', 'demo', 'log', 'read' );")
	m.SQL("INSERT INTO `permissions`(`id`,`role_name`,`package_name`,`level_name`) VALUES ( '107', 'demo', 'map', 'read_map_element' );")
	m.SQL("INSERT INTO `permissions`(`id`,`role_name`,`package_name`,`level_name`) VALUES ( '108', 'demo', 'map', 'read_map_layer' );")
	m.SQL("INSERT INTO `permissions`(`id`,`role_name`,`package_name`,`level_name`) VALUES ( '109', 'user', 'image', 'upload' );")
	m.SQL("INSERT INTO `permissions`(`id`,`role_name`,`package_name`,`level_name`) VALUES ( '110', 'user', 'image', 'create' );")
	m.SQL("INSERT INTO `permissions`(`id`,`role_name`,`package_name`,`level_name`) VALUES ( '111', 'user', 'image', 'delete' );")
	m.SQL("INSERT INTO `permissions`(`id`,`role_name`,`package_name`,`level_name`) VALUES ( '112', 'user', 'image', 'update' );")
	m.SQL("INSERT INTO `permissions`(`id`,`role_name`,`package_name`,`level_name`) VALUES ( '113', 'user', 'flow', 'create' );")
	m.SQL("INSERT INTO `permissions`(`id`,`role_name`,`package_name`,`level_name`) VALUES ( '114', 'user', 'flow', 'delete' );")
	m.SQL("INSERT INTO `permissions`(`id`,`role_name`,`package_name`,`level_name`) VALUES ( '115', 'user', 'flow', 'update' );")
	m.SQL("INSERT INTO `permissions`(`id`,`role_name`,`package_name`,`level_name`) VALUES ( '116', 'user', 'device_state', 'create' );")
	m.SQL("INSERT INTO `permissions`(`id`,`role_name`,`package_name`,`level_name`) VALUES ( '117', 'user', 'device_state', 'delete' );")
	m.SQL("INSERT INTO `permissions`(`id`,`role_name`,`package_name`,`level_name`) VALUES ( '118', 'user', 'device_state', 'update' );")
	m.SQL("INSERT INTO `permissions`(`id`,`role_name`,`package_name`,`level_name`) VALUES ( '119', 'user', 'device_action', 'delete' );")
	m.SQL("INSERT INTO `permissions`(`id`,`role_name`,`package_name`,`level_name`) VALUES ( '120', 'user', 'device_action', 'create' );")
	m.SQL("INSERT INTO `permissions`(`id`,`role_name`,`package_name`,`level_name`) VALUES ( '121', 'user', 'device', 'update' );")
	m.SQL("INSERT INTO `permissions`(`id`,`role_name`,`package_name`,`level_name`) VALUES ( '122', 'user', 'device', 'create' );")
	m.SQL("INSERT INTO `permissions`(`id`,`role_name`,`package_name`,`level_name`) VALUES ( '123', 'user', 'device', 'delete' );")
	m.SQL("INSERT INTO `permissions`(`id`,`role_name`,`package_name`,`level_name`) VALUES ( '124', 'user', 'dashboard', 'create' );")
	m.SQL("INSERT INTO `permissions`(`id`,`role_name`,`package_name`,`level_name`) VALUES ( '125', 'user', 'dashboard', 'update' );")
	m.SQL("INSERT INTO `permissions`(`id`,`role_name`,`package_name`,`level_name`) VALUES ( '126', 'user', 'dashboard', 'delete' );")
	m.SQL("INSERT INTO `permissions`(`id`,`role_name`,`package_name`,`level_name`) VALUES ( '127', 'user', 'script', 'update' );")
	m.SQL("INSERT INTO `permissions`(`id`,`role_name`,`package_name`,`level_name`) VALUES ( '128', 'user', 'script', 'create' );")
	m.SQL("INSERT INTO `permissions`(`id`,`role_name`,`package_name`,`level_name`) VALUES ( '129', 'user', 'script', 'delete' );")
	m.SQL("INSERT INTO `permissions`(`id`,`role_name`,`package_name`,`level_name`) VALUES ( '130', 'user', 'worker', 'delete' );")
	m.SQL("INSERT INTO `permissions`(`id`,`role_name`,`package_name`,`level_name`) VALUES ( '131', 'user', 'worker', 'create' );")
	m.SQL("INSERT INTO `permissions`(`id`,`role_name`,`package_name`,`level_name`) VALUES ( '132', 'user', 'worker', 'update' );")
	m.SQL("INSERT INTO `permissions`(`id`,`role_name`,`package_name`,`level_name`) VALUES ( '133', 'user', 'workflow', 'create' );")
	m.SQL("INSERT INTO `permissions`(`id`,`role_name`,`package_name`,`level_name`) VALUES ( '134', 'user', 'workflow', 'delete' );")
	m.SQL("INSERT INTO `permissions`(`id`,`role_name`,`package_name`,`level_name`) VALUES ( '135', 'user', 'workflow', 'update' );")
}

// Reverse the migrations
func (m *AddUsers_20170121_175235) Down() {
	// use m.SQL("DROP TABLE ...") to reverse schema update
	m.SQL("DELETE FROM `user_metas` WHere user_id IN (1,2)")
	m.SQL("DELETE FROM `users` WHERE  id in (1,2)")
	m.SQL(`DELETE FROM permissions WHERE role_name IN ("admin", "user", "demo")`)
	m.SQL(`DELETE FROM roles WHERE  name in ("admin")`)
	m.SQL(`DELETE FROM roles WHERE  name in ("user")`)
	m.SQL(`DELETE FROM roles WHERE  name in ("demo")`)
}
