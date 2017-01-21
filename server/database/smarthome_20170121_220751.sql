-- Valentina Studio --
-- MySQL dump --
-- ---------------------------------------------------------


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
-- ---------------------------------------------------------


-- CREATE TABLE "connections" ------------------------------
CREATE TABLE `connections` ( 
	`uuid` VarChar( 254 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
	`name` VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NULL,
	`element_from` VarChar( 254 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
	`element_to` VarChar( 254 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
	`flow_id` Int( 32 ) NOT NULL,
	`created_at` DateTime NOT NULL,
	`update_at` DateTime NOT NULL,
	`graph_settings` Text CHARACTER SET utf8 COLLATE utf8_general_ci NULL,
	`point_from` Int( 11 ) NOT NULL,
	`point_to` Int( 11 ) NOT NULL,
	`direction` VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
	PRIMARY KEY ( `uuid` ),
	CONSTRAINT `unique_id` UNIQUE( `uuid` ) )
CHARACTER SET = utf8
COLLATE = utf8_general_ci
ENGINE = InnoDB;
-- ---------------------------------------------------------


-- CREATE TABLE "dashboards" -------------------------------
CREATE TABLE `dashboards` ( 
	`id` Int( 32 ) AUTO_INCREMENT NOT NULL,
	`name` VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
	`description` Text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
	`created_at` DateTime NOT NULL,
	`update_at` DateTime NOT NULL,
	`widgets` Text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
	PRIMARY KEY ( `id` ),
	CONSTRAINT `unique_id` UNIQUE( `id` ) )
CHARACTER SET = utf8
COLLATE = utf8_general_ci
ENGINE = InnoDB
AUTO_INCREMENT = 1;
-- ---------------------------------------------------------


-- CREATE TABLE "device_actions" ---------------------------
CREATE TABLE `device_actions` ( 
	`id` Int( 11 ) AUTO_INCREMENT NOT NULL,
	`device_id` Int( 11 ) NULL,
	`name` VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
	`description` VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NULL,
	`created_at` DateTime NOT NULL,
	`update_at` DateTime NOT NULL,
	`script_id` Int( 32 ) NULL,
	PRIMARY KEY ( `id` ),
	CONSTRAINT `unique_id` UNIQUE( `id` ) )
CHARACTER SET = utf8
COLLATE = utf8_general_ci
ENGINE = InnoDB
AUTO_INCREMENT = 6;
-- ---------------------------------------------------------


-- CREATE TABLE "device_states" ----------------------------
CREATE TABLE `device_states` ( 
	`id` Int( 11 ) AUTO_INCREMENT NOT NULL,
	`system_name` VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
	`description` Text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
	`created_at` DateTime NOT NULL,
	`update_at` DateTime NOT NULL,
	`device_id` Int( 11 ) NOT NULL,
	PRIMARY KEY ( `id` ),
	CONSTRAINT `unique_1` UNIQUE( `device_id`, `system_name` ) )
CHARACTER SET = utf8
COLLATE = utf8_general_ci
ENGINE = InnoDB
AUTO_INCREMENT = 4;
-- ---------------------------------------------------------


-- CREATE TABLE "devices" ----------------------------------
CREATE TABLE `devices` ( 
	`id` Int( 11 ) AUTO_INCREMENT NOT NULL,
	`name` VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
	`description` VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
	`device_id` Int( 11 ) NULL,
	`created_at` DateTime NOT NULL,
	`update_at` DateTime NOT NULL,
	`node_id` Int( 11 ) NULL,
	`baud` Int( 11 ) NULL,
	`tty` VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NULL,
	`stop_bite` Int( 8 ) NULL,
	`timeout` Int( 11 ) NULL,
	`address` Int( 11 ) NULL,
	`status` VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT 'enabled',
	`sleep` Int( 32 ) NOT NULL DEFAULT '0',
	PRIMARY KEY ( `id` ),
	CONSTRAINT `unique1` UNIQUE( `device_id`, `address` ),
	CONSTRAINT `unique2` UNIQUE( `node_id`, `address` ),
	CONSTRAINT `unique3` UNIQUE( `name`, `device_id` ) )
CHARACTER SET = utf8
COLLATE = utf8_general_ci
ENGINE = InnoDB
AUTO_INCREMENT = 8;
-- ---------------------------------------------------------


-- CREATE TABLE "email_templates" --------------------------
CREATE TABLE `email_templates` ( 
	`name` VarChar( 64 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT 'Название',
	`description` VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NULL COMMENT 'Описание',
	`content` LongText CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT 'Содержимое',
	`status` VarChar( 64 ) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT 'active' COMMENT 'active, unactive',
	`type` VarChar( 64 ) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT 'item' COMMENT 'item, template',
	`parent` VarChar( 64 ) CHARACTER SET utf8 COLLATE utf8_general_ci NULL,
	`created_at` DateTime NOT NULL,
	`updated_at` DateTime NOT NULL,
	PRIMARY KEY ( `name` ),
	CONSTRAINT `name` UNIQUE( `name` ) )
CHARACTER SET = utf8
COLLATE = utf8_general_ci
ENGINE = InnoDB;
-- ---------------------------------------------------------


-- CREATE TABLE "flow_elements" ----------------------------
CREATE TABLE `flow_elements` ( 
	`uuid` VarChar( 254 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
	`name` VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
	`graph_settings` Text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
	`flow_id` Int( 32 ) NOT NULL,
	`created_at` DateTime NOT NULL,
	`update_at` DateTime NOT NULL,
	`prototype_type` VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT 'default',
	`status` VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NULL,
	`description` Text CHARACTER SET utf8 COLLATE utf8_general_ci NULL,
	`script_id` Int( 11 ) NULL,
	`flow_link` Int( 32 ) NULL,
	PRIMARY KEY ( `uuid` ),
	CONSTRAINT `id` UNIQUE( `uuid` ),
	CONSTRAINT `unique_id` UNIQUE( `uuid` ) )
CHARACTER SET = utf8
COLLATE = utf8_general_ci
ENGINE = InnoDB;
-- ---------------------------------------------------------


-- CREATE TABLE "flows" ------------------------------------
CREATE TABLE `flows` ( 
	`id` Int( 32 ) AUTO_INCREMENT NOT NULL,
	`name` VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
	`status` Enum( 'enabled', 'disabled' ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT 'enabled',
	`created_at` DateTime NOT NULL,
	`update_at` DateTime NOT NULL,
	`workflow_id` Int( 11 ) NOT NULL,
	`description` Text CHARACTER SET utf8 COLLATE utf8_general_ci NULL,
	PRIMARY KEY ( `id` ),
	CONSTRAINT `unique_id` UNIQUE( `id` ) )
CHARACTER SET = utf8
COLLATE = utf8_general_ci
ENGINE = InnoDB
AUTO_INCREMENT = 2;
-- ---------------------------------------------------------


-- CREATE TABLE "images" -----------------------------------
CREATE TABLE `images` ( 
	`id` Int( 32 ) AUTO_INCREMENT NOT NULL,
	`thumb` VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
	`image` VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
	`mime_type` VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
	`title` VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
	`size` Int( 11 ) NOT NULL,
	`name` VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
	`created_at` DateTime NOT NULL,
	PRIMARY KEY ( `id` ),
	CONSTRAINT `id` UNIQUE( `id` ) )
CHARACTER SET = utf8
COLLATE = utf8_general_ci
ENGINE = InnoDB
AUTO_INCREMENT = 10;
-- ---------------------------------------------------------


-- CREATE TABLE "logs" -------------------------------------
CREATE TABLE `logs` ( 
	`id` Int( 32 ) AUTO_INCREMENT NOT NULL,
	`body` Text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
	`created_at` DateTime NOT NULL,
	`level` Enum( 'Emergency', 'Alert', 'Critical', 'Error', 'Warning', 'Notice', 'Info', 'Debug' ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT 'Info',
	PRIMARY KEY ( `id` ),
	CONSTRAINT `unique_id` UNIQUE( `id` ) )
CHARACTER SET = utf8
COLLATE = utf8_general_ci
ENGINE = InnoDB
AUTO_INCREMENT = 266;
-- ---------------------------------------------------------


-- CREATE TABLE "map_device_actions" -----------------------
CREATE TABLE `map_device_actions` ( 
	`id` Int( 32 ) AUTO_INCREMENT NOT NULL,
	`device_action_id` Int( 32 ) NOT NULL,
	`map_device_id` Int( 2 ) NOT NULL,
	`image_id` Int( 32 ) NULL,
	`type` VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
	`created_at` DateTime NOT NULL,
	`update_at` DateTime NOT NULL,
	PRIMARY KEY ( `id` ),
	CONSTRAINT `unique_id` UNIQUE( `id` ) )
CHARACTER SET = utf8
COLLATE = utf8_general_ci
ENGINE = InnoDB
AUTO_INCREMENT = 76;
-- ---------------------------------------------------------


-- CREATE TABLE "map_device_states" ------------------------
CREATE TABLE `map_device_states` ( 
	`id` Int( 32 ) AUTO_INCREMENT NOT NULL,
	`device_state_id` Int( 32 ) NOT NULL,
	`map_device_id` Int( 32 ) NOT NULL,
	`image_id` Int( 32 ) NULL,
	`created_at` DateTime NOT NULL,
	`update_at` DateTime NOT NULL,
	`style` Text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
	PRIMARY KEY ( `id` ) )
CHARACTER SET = utf8
COLLATE = utf8_general_ci
ENGINE = InnoDB
AUTO_INCREMENT = 76;
-- ---------------------------------------------------------


-- CREATE TABLE "map_devices" ------------------------------
CREATE TABLE `map_devices` ( 
	`id` Int( 32 ) AUTO_INCREMENT NOT NULL,
	`created_at` DateTime NOT NULL,
	`update_at` DateTime NOT NULL,
	`device_id` Int( 32 ) NOT NULL,
	`image_id` Int( 32 ) NULL,
	PRIMARY KEY ( `id` ),
	CONSTRAINT `unique_id` UNIQUE( `device_id` ) )
CHARACTER SET = utf8
COLLATE = utf8_general_ci
ENGINE = InnoDB
AUTO_INCREMENT = 14;
-- ---------------------------------------------------------


-- CREATE TABLE "map_elements" -----------------------------
CREATE TABLE `map_elements` ( 
	`id` Int( 32 ) AUTO_INCREMENT NOT NULL,
	`name` VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
	`description` Text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
	`created_at` DateTime NOT NULL,
	`update_at` DateTime NOT NULL,
	`prototype_type` VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
	`layer_id` Int( 32 ) NOT NULL,
	`map_id` Int( 32 ) NOT NULL,
	`graph_settings` Text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
	`status` VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
	`weight` Int( 11 ) NOT NULL DEFAULT '0',
	`prototype_id` Int( 32 ) NOT NULL,
	PRIMARY KEY ( `id` ),
	CONSTRAINT `unique_id` UNIQUE( `id` ) )
CHARACTER SET = utf8
COLLATE = utf8_general_ci
ENGINE = InnoDB
AUTO_INCREMENT = 7;
-- ---------------------------------------------------------


-- CREATE TABLE "map_images" -------------------------------
CREATE TABLE `map_images` ( 
	`id` Int( 32 ) AUTO_INCREMENT NOT NULL,
	`image_id` Int( 32 ) NOT NULL,
	`style` Text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
	PRIMARY KEY ( `id` ),
	CONSTRAINT `unique_id` UNIQUE( `id` ) )
CHARACTER SET = utf8
COLLATE = utf8_general_ci
ENGINE = InnoDB
AUTO_INCREMENT = 4;
-- ---------------------------------------------------------


-- CREATE TABLE "map_layers" -------------------------------
CREATE TABLE `map_layers` ( 
	`id` Int( 11 ) AUTO_INCREMENT NOT NULL,
	`name` VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
	`description` Text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
	`created_at` DateTime NOT NULL,
	`update_at` DateTime NOT NULL,
	`map_id` Int( 11 ) NOT NULL,
	`status` VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
	`weight` Int( 11 ) NOT NULL DEFAULT '0',
	PRIMARY KEY ( `id` ),
	CONSTRAINT `unique_id` UNIQUE( `id` ) )
CHARACTER SET = utf8
COLLATE = utf8_general_ci
ENGINE = InnoDB
AUTO_INCREMENT = 3;
-- ---------------------------------------------------------


-- CREATE TABLE "map_texts" --------------------------------
CREATE TABLE `map_texts` ( 
	`id` Int( 11 ) AUTO_INCREMENT NOT NULL,
	`Text` Text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
	`Style` Text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
	PRIMARY KEY ( `id` ),
	CONSTRAINT `unique_id` UNIQUE( `id` ) )
CHARACTER SET = utf8
COLLATE = utf8_general_ci
ENGINE = InnoDB
AUTO_INCREMENT = 14;
-- ---------------------------------------------------------


-- CREATE TABLE "maps" -------------------------------------
CREATE TABLE `maps` ( 
	`id` Int( 11 ) AUTO_INCREMENT NOT NULL,
	`created_at` DateTime NOT NULL,
	`update_at` DateTime NOT NULL,
	`name` Char( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
	`description` Text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
	`options` Text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
	PRIMARY KEY ( `id` ),
	CONSTRAINT `unique_id` UNIQUE( `id` ) )
CHARACTER SET = utf8
COLLATE = utf8_general_ci
ENGINE = InnoDB
AUTO_INCREMENT = 2;
-- ---------------------------------------------------------


-- CREATE TABLE "migrations" -------------------------------
CREATE TABLE `migrations` ( 
	`id_migration` Int( 10 ) UNSIGNED AUTO_INCREMENT NOT NULL COMMENT 'surrogate key',
	`name` VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NULL COMMENT 'migration name, unique',
	`created_at` Timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'date migrated or rolled back',
	`statements` LongText CHARACTER SET utf8 COLLATE utf8_general_ci NULL COMMENT 'SQL statements for this migration',
	`rollback_statements` LongText CHARACTER SET utf8 COLLATE utf8_general_ci NULL COMMENT 'SQL statment for rolling back migration',
	`status` Enum( 'update', 'rollback' ) CHARACTER SET utf8 COLLATE utf8_general_ci NULL COMMENT 'update indicates it is a normal migration while rollback means this migration is rolled back',
	PRIMARY KEY ( `id_migration` ) )
CHARACTER SET = utf8
COLLATE = utf8_general_ci
ENGINE = InnoDB
AUTO_INCREMENT = 32;
-- ---------------------------------------------------------


-- CREATE TABLE "nodes" ------------------------------------
CREATE TABLE `nodes` ( 
	`id` Int( 255 ) AUTO_INCREMENT NOT NULL,
	`ip` VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
	`port` Int( 255 ) NOT NULL,
	`name` VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
	`description` Text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
	`created_at` DateTime NOT NULL,
	`update_at` DateTime NOT NULL,
	`status` Enum( 'enabled', 'disabled' ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT 'enabled',
	PRIMARY KEY ( `id` ),
	CONSTRAINT `node` UNIQUE( `port`, `ip` ),
	CONSTRAINT `unique_id` UNIQUE( `id` ) )
CHARACTER SET = utf8
COLLATE = utf8_general_ci
ENGINE = InnoDB
AUTO_INCREMENT = 3;
-- ---------------------------------------------------------


-- CREATE TABLE "permissions" ------------------------------
CREATE TABLE `permissions` ( 
	`id` Int( 32 ) AUTO_INCREMENT NOT NULL,
	`role_name` VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
	`package_name` VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
	`level_name` VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
	PRIMARY KEY ( `id` ) )
CHARACTER SET = utf8
COLLATE = utf8_general_ci
ENGINE = InnoDB
AUTO_INCREMENT = 136;
-- ---------------------------------------------------------


-- CREATE TABLE "roles" ------------------------------------
CREATE TABLE `roles` ( 
	`name` VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
	`description` VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
	`parent` VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NULL,
	`created_at` DateTime NOT NULL,
	`update_at` DateTime NOT NULL,
	PRIMARY KEY ( `name` ),
	CONSTRAINT `unique_name` UNIQUE( `name` ) )
CHARACTER SET = utf8
COLLATE = utf8_general_ci
ENGINE = InnoDB;
-- ---------------------------------------------------------


-- CREATE TABLE "scripts" ----------------------------------
CREATE TABLE `scripts` ( 
	`id` Int( 32 ) AUTO_INCREMENT NOT NULL,
	`lang` Enum( 'lua', 'coffeescript', 'javascript' ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT 'lua',
	`name` VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
	`source` Text CHARACTER SET utf8 COLLATE utf8_general_ci NULL,
	`created_at` DateTime NOT NULL,
	`update_at` DateTime NOT NULL,
	`description` Text CHARACTER SET utf8 COLLATE utf8_general_ci NULL,
	`compiled` Text CHARACTER SET utf8 COLLATE utf8_general_ci NULL,
	PRIMARY KEY ( `id` ),
	CONSTRAINT `unique_id` UNIQUE( `id` ),
	CONSTRAINT `unique_name` UNIQUE( `name` ) )
CHARACTER SET = utf8
COLLATE = utf8_general_ci
ENGINE = InnoDB
AUTO_INCREMENT = 7;
-- ---------------------------------------------------------


-- CREATE TABLE "user_metas" -------------------------------
CREATE TABLE `user_metas` ( 
	`id` Int( 32 ) AUTO_INCREMENT NOT NULL,
	`user_id` Int( 32 ) NOT NULL,
	`key` VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
	`value` VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
	PRIMARY KEY ( `id` ),
	CONSTRAINT `unique_id` UNIQUE( `id` ) )
CHARACTER SET = utf8
COLLATE = utf8_general_ci
ENGINE = InnoDB
AUTO_INCREMENT = 57;
-- ---------------------------------------------------------


-- CREATE TABLE "users" ------------------------------------
CREATE TABLE `users` ( 
	`id` Int( 32 ) AUTO_INCREMENT NOT NULL,
	`nickname` VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
	`first_name` VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
	`last_name` VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
	`encrypted_password` VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
	`email` VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
	`history` Text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
	`status` VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
	`reset_password_token` VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
	`authentication_token` VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
	`image_id` Int( 255 ) NULL,
	`sign_in_count` Int( 32 ) NOT NULL,
	`current_sign_in_ip` VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
	`last_sign_in_ip` VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NULL,
	`user_id` Int( 32 ) NULL,
	`role_name` VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
	`reset_password_sent_at` DateTime NULL,
	`current_sign_in_at` DateTime NULL,
	`last_sign_in_at` DateTime NULL,
	`created_at` DateTime NOT NULL,
	`update_at` DateTime NOT NULL,
	`deleted` DateTime NULL,
	PRIMARY KEY ( `id` ),
	CONSTRAINT `unique_email` UNIQUE( `email` ),
	CONSTRAINT `unique_id` UNIQUE( `id` ),
	CONSTRAINT `unique_nickname` UNIQUE( `nickname` ) )
CHARACTER SET = utf8
COLLATE = utf8_general_ci
ENGINE = InnoDB
AUTO_INCREMENT = 3;
-- ---------------------------------------------------------


-- CREATE TABLE "workers" ----------------------------------
CREATE TABLE `workers` ( 
	`id` Int( 11 ) AUTO_INCREMENT NOT NULL,
	`flow_id` Int( 11 ) NOT NULL,
	`workflow_id` Int( 11 ) NOT NULL,
	`name` VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
	`device_action_id` Int( 11 ) NOT NULL,
	`time` VarChar( 254 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
	`created_at` DateTime NOT NULL,
	`update_at` DateTime NOT NULL,
	`status` Enum( 'enabled', 'disabled' ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT 'enabled',
	PRIMARY KEY ( `id` ) )
CHARACTER SET = utf8
COLLATE = utf8_general_ci
ENGINE = InnoDB
AUTO_INCREMENT = 2;
-- ---------------------------------------------------------


-- CREATE TABLE "workflows" --------------------------------
CREATE TABLE `workflows` ( 
	`id` Int( 255 ) AUTO_INCREMENT NOT NULL,
	`name` VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
	`status` Enum( 'enabled', 'disabled' ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT 'enabled',
	`created_at` DateTime NOT NULL,
	`update_at` DateTime NOT NULL,
	`description` Text CHARACTER SET utf8 COLLATE utf8_general_ci NULL,
	PRIMARY KEY ( `id` ),
	CONSTRAINT `unique_name` UNIQUE( `name` ) )
CHARACTER SET = utf8
COLLATE = utf8_general_ci
ENGINE = InnoDB
AUTO_INCREMENT = 2;
-- ---------------------------------------------------------


-- Dump data of "connections" ------------------------------
INSERT INTO `connections`(`uuid`,`name`,`element_from`,`element_to`,`flow_id`,`created_at`,`update_at`,`graph_settings`,`point_from`,`point_to`,`direction`) VALUES ( '30d0dde2-959b-4e3b-89fe-5edb7b87ee4b', '', '251c8846-db07-44a4-a1d3-8e3c5a885fac', 'c192ad47-f3ab-498b-882c-dbea7f2ee5c0', '1', '2017-01-21 14:08:18', '2017-01-21 14:08:18', '', '4', '3', '' );
INSERT INTO `connections`(`uuid`,`name`,`element_from`,`element_to`,`flow_id`,`created_at`,`update_at`,`graph_settings`,`point_from`,`point_to`,`direction`) VALUES ( 'f0ccc4ae-e228-484a-a199-ea7accc67b58', '', 'b602d3b4-1a32-4035-bc78-cc2b5e8778e7', '251c8846-db07-44a4-a1d3-8e3c5a885fac', '1', '2017-01-21 14:08:18', '2017-01-21 14:08:18', '', '1', '10', '' );
-- ---------------------------------------------------------


-- Dump data of "dashboards" -------------------------------
-- ---------------------------------------------------------


-- Dump data of "device_actions" ---------------------------
INSERT INTO `device_actions`(`id`,`device_id`,`name`,`description`,`created_at`,`update_at`,`script_id`) VALUES ( '1', '1', 'Проверка состояния', 'Какое-то действие', '2017-01-21 13:32:28', '2017-01-21 13:32:28', '1' );
INSERT INTO `device_actions`(`id`,`device_id`,`name`,`description`,`created_at`,`update_at`,`script_id`) VALUES ( '3', '1', 'Включить', 'Какое-то действие', '2017-01-21 13:39:41', '2017-01-21 13:45:00', '4' );
INSERT INTO `device_actions`(`id`,`device_id`,`name`,`description`,`created_at`,`update_at`,`script_id`) VALUES ( '4', '1', 'Выключить', 'Какое-то действие', '2017-01-21 13:39:56', '2017-01-21 13:39:56', '5' );
-- ---------------------------------------------------------


-- Dump data of "device_states" ----------------------------
INSERT INTO `device_states`(`id`,`system_name`,`description`,`created_at`,`update_at`,`device_id`) VALUES ( '1', 'ENABLED', 'Включено', '2017-01-21 12:27:46', '2017-01-21 12:27:46', '1' );
INSERT INTO `device_states`(`id`,`system_name`,`description`,`created_at`,`update_at`,`device_id`) VALUES ( '2', 'DISABLED', 'Выключено', '2017-01-21 12:28:07', '2017-01-21 12:28:07', '1' );
INSERT INTO `device_states`(`id`,`system_name`,`description`,`created_at`,`update_at`,`device_id`) VALUES ( '3', 'ERROR', 'Ошибка', '2017-01-21 12:28:19', '2017-01-21 12:28:19', '1' );
-- ---------------------------------------------------------


-- Dump data of "devices" ----------------------------------
INSERT INTO `devices`(`id`,`name`,`description`,`device_id`,`created_at`,`update_at`,`node_id`,`baud`,`tty`,`stop_bite`,`timeout`,`address`,`status`,`sleep`) VALUES ( '1', 'Группа розеток', 'Группа розеток', NULL, '2017-01-21 12:16:24', '2017-01-21 13:41:07', '1', '19200', '', '2', '1000', NULL, 'enabled', '50' );
INSERT INTO `devices`(`id`,`name`,`description`,`device_id`,`created_at`,`update_at`,`node_id`,`baud`,`tty`,`stop_bite`,`timeout`,`address`,`status`,`sleep`) VALUES ( '2', 'р1', '', '1', '2017-01-21 13:31:54', '2017-01-21 13:54:52', NULL, '0', '', '0', '0', '1', 'enabled', '0' );
INSERT INTO `devices`(`id`,`name`,`description`,`device_id`,`created_at`,`update_at`,`node_id`,`baud`,`tty`,`stop_bite`,`timeout`,`address`,`status`,`sleep`) VALUES ( '6', 'р2', '', '1', '2017-01-21 13:58:25', '2017-01-21 13:58:25', NULL, '0', '', '0', '0', '2', 'enabled', '0' );
-- ---------------------------------------------------------


-- Dump data of "email_templates" --------------------------
INSERT INTO `email_templates`(`name`,`description`,`content`,`status`,`type`,`parent`,`created_at`,`updated_at`) VALUES ( 'body', '', '[body:content]', 'active', 'item', 'message', '2014-06-21 21:56:07', '2015-04-14 15:24:39' );
INSERT INTO `email_templates`(`name`,`description`,`content`,`status`,`type`,`parent`,`created_at`,`updated_at`) VALUES ( 'callout', '', '<table class="row callout">
    <tr>
        <td class="wrapper last">

            <table class="twelve columns">
                <tr>
                    <td class="panel">
                        [callout:content]
                    </td>
                    <td class="expander"></td>
                </tr>
            </table>

        </td>
    </tr>
</table>', 'active', 'item', 'main', '2014-06-15 17:45:31', '2015-04-12 01:40:03' );
INSERT INTO `email_templates`(`name`,`description`,`content`,`status`,`type`,`parent`,`created_at`,`updated_at`) VALUES ( 'contacts', '', '<table class="six columns">
    <tr>
        <td class="last right-text-pad">
        <h3>Контакты:</h3>
            [contacts:content]
        </td>
        <td class="expander"></td>
    </tr>
</table>', 'active', 'item', 'footer', '2014-06-15 17:45:20', '2015-04-12 01:40:03' );
INSERT INTO `email_templates`(`name`,`description`,`content`,`status`,`type`,`parent`,`created_at`,`updated_at`) VALUES ( 'facebook', '', '<table class="tiny-button facebook">
    <tr>
        <td>
            <a href="#">Facebook</a>
        </td>
    </tr>
</table>

<br>', 'active', 'item', 'social', '2014-06-15 17:46:59', '2015-04-13 01:50:01' );
INSERT INTO `email_templates`(`name`,`description`,`content`,`status`,`type`,`parent`,`created_at`,`updated_at`) VALUES ( 'footer', '', '<table class="row footer">
    <tr>

        <td class="wrapper ">

            [social:block]

        </td>
        <td class="wrapper last">

            [contacts:block]

        </td>

    </tr>
</table>', 'active', 'item', 'main', '2014-06-15 17:45:07', '2015-04-12 01:40:03' );
INSERT INTO `email_templates`(`name`,`description`,`content`,`status`,`type`,`parent`,`created_at`,`updated_at`) VALUES ( 'google', '', '<table class="tiny-button google-plus">
    <tr>
        <td>
            <a href="#">Google +</a>
        </td>
    </tr>
</table>

<br>', 'active', 'item', 'social', '2014-06-15 17:47:17', '2015-04-12 01:40:01' );
INSERT INTO `email_templates`(`name`,`description`,`content`,`status`,`type`,`parent`,`created_at`,`updated_at`) VALUES ( 'header', '', '<table class="row header">
    <tr>
        <td class="center" align="center">
            <center>

                <table class="container">
                    <tr>
                        <td class="wrapper last">

                            <table class="twelve columns">
                                <tr>
                                    <td class="six sub-columns">
                                        <img alt="Облачная типография, Календари-домики, Листовки и Флаеры, Карманные календари, Визитки, Буклеты, Пластиковые карты" src="data:image/png;base64,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">
                                    </td>
                                    <td class="six sub-columns last" style="text-align:right; vertical-align:middle;">
                                    </td>
                                    <td class="expander"></td>
                                </tr>
                            </table>

                        </td>
                    </tr>
                </table>

            </center>
        </td>
    </tr>
</table>', 'active', 'item', 'main', '2014-06-15 17:44:55', '2015-04-16 23:43:22' );
INSERT INTO `email_templates`(`name`,`description`,`content`,`status`,`type`,`parent`,`created_at`,`updated_at`) VALUES ( 'main', 'Основной слой', '<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Strict//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-strict.dtd">
<html xmlns="http://www.w3.org/1999/xhtml">
<head>
  <meta http-equiv="Content-Type" content="text/html; charset=UTF-8">
  <meta name="viewport" content="width=device-width"/>
  <style type="text/css">

#outlook a {
  padding:0;
}

body{
  width:100% !important;
  min-width: 100%;
  -webkit-text-size-adjust:100%;
  -ms-text-size-adjust:100%;
  margin:0;
  padding:0;
}

.ExternalClass {
  width:100%;
}

.ExternalClass,
.ExternalClass p,
.ExternalClass span,
.ExternalClass font,
.ExternalClass td,
.ExternalClass div {
  line-height: 100%;
}

#backgroundTable {
  margin:0;
  padding:0;
  width:100% !important;
  line-height: 100% !important;
}

img {
  outline:none;
  text-decoration:none;
  -ms-interpolation-mode: bicubic;
  width: auto;
  max-width: 100%;
  float: left;
  clear: both;
  display: block;
}

center {
  width: 100%;
  min-width: 580px;
}

a img {
  border: none;
}

p {
  margin: 0 0 0 10px;
}

table {
  border-spacing: 0;
  border-collapse: collapse;
}

td {
  word-break: break-word;
  -webkit-hyphens: auto;
  -moz-hyphens: auto;
  hyphens: auto;
  border-collapse: collapse !important;
}

table, tr, td {
  padding: 0;
  vertical-align: top;
  text-align: left;
}

hr {
  color: #d9d9d9;
  background-color: #d9d9d9;
  height: 1px;
  border: none;
}

/* Responsive Grid */

table.body {
  height: 100%;
  width: 100%;
}

table.container {
  width: 580px;
  margin: 0 auto;
  text-align: inherit;
}

table.row {
  padding: 0;
  width: 100%;
  position: relative;
}

table.container table.row {
  display: block;
}

td.wrapper {
  padding: 10px 20px 0 0;
  position: relative;
}

table.columns,
table.column {
  margin: 0 auto;
}

table.columns td,
table.column td {
  padding: 0 0 10px;
}

table.columns td.sub-columns,
table.column td.sub-columns,
table.columns td.sub-column,
table.column td.sub-column {
  padding-right: 10px;
}

td.sub-column, td.sub-columns {
  min-width: 0;
}

table.row td.last,
table.container td.last {
  padding-right: 0;
}

table.one { width: 30px; }
table.two { width: 80px; }
table.three { width: 130px; }
table.four { width: 180px; }
table.five { width: 230px; }
table.six { width: 280px; }
table.seven { width: 330px; }
table.eight { width: 380px; }
table.nine { width: 430px; }
table.ten { width: 480px; }
table.eleven { width: 530px; }
table.twelve { width: 580px; }

table.one center { min-width: 30px; }
table.two center { min-width: 80px; }
table.three center { min-width: 130px; }
table.four center { min-width: 180px; }
table.five center { min-width: 230px; }
table.six center { min-width: 280px; }
table.seven center { min-width: 330px; }
table.eight center { min-width: 380px; }
table.nine center { min-width: 430px; }
table.ten center { min-width: 480px; }
table.eleven center { min-width: 530px; }
table.twelve center { min-width: 580px; }

table.one .panel center { min-width: 10px; }
table.two .panel center { min-width: 60px; }
table.three .panel center { min-width: 110px; }
table.four .panel center { min-width: 160px; }
table.five .panel center { min-width: 210px; }
table.six .panel center { min-width: 260px; }
table.seven .panel center { min-width: 310px; }
table.eight .panel center { min-width: 360px; }
table.nine .panel center { min-width: 410px; }
table.ten .panel center { min-width: 460px; }
table.eleven .panel center { min-width: 510px; }
table.twelve .panel center { min-width: 560px; }

.body .columns td.one,
.body .column td.one { width: 8.333333%; }
.body .columns td.two,
.body .column td.two { width: 16.666666%; }
.body .columns td.three,
.body .column td.three { width: 25%; }
.body .columns td.four,
.body .column td.four { width: 33.333333%; }
.body .columns td.five,
.body .column td.five { width: 41.666666%; }
.body .columns td.six,
.body .column td.six { width: 50%; }
.body .columns td.seven,
.body .column td.seven { width: 58.333333%; }
.body .columns td.eight,
.body .column td.eight { width: 66.666666%; }
.body .columns td.nine,
.body .column td.nine { width: 75%; }
.body .columns td.ten,
.body .column td.ten { width: 83.333333%; }
.body .columns td.eleven,
.body .column td.eleven { width: 91.666666%; }
.body .columns td.twelve,
.body .column td.twelve { width: 100%; }

td.offset-by-one { padding-left: 50px; }
td.offset-by-two { padding-left: 100px; }
td.offset-by-three { padding-left: 150px; }
td.offset-by-four { padding-left: 200px; }
td.offset-by-five { padding-left: 250px; }
td.offset-by-six { padding-left: 300px; }
td.offset-by-seven { padding-left: 350px; }
td.offset-by-eight { padding-left: 400px; }
td.offset-by-nine { padding-left: 450px; }
td.offset-by-ten { padding-left: 500px; }
td.offset-by-eleven { padding-left: 550px; }

td.expander {
  visibility: hidden;
  width: 0;
  padding: 0 !important;
}

table.columns .text-pad,
table.column .text-pad {
  padding-left: 10px;
  padding-right: 10px;
}

table.columns .left-text-pad,
table.columns .text-pad-left,
table.column .left-text-pad,
table.column .text-pad-left {
  padding-left: 10px;
}

table.columns .right-text-pad,
table.columns .text-pad-right,
table.column .right-text-pad,
table.column .text-pad-right {
  padding-right: 10px;
}

/* Block Grid */

.block-grid {
  width: 100%;
  max-width: 580px;
}

.block-grid td {
  display: inline-block;
  padding:10px;
}

.two-up td {
  width:270px;
}

.three-up td {
  width:173px;
}

.four-up td {
  width:125px;
}

.five-up td {
  width:96px;
}

.six-up td {
  width:76px;
}

.seven-up td {
  width:62px;
}

.eight-up td {
  width:52px;
}

/* Alignment & Visibility Classes */

table.center, td.center {
  text-align: center;
}

h1.center,
h2.center,
h3.center,
h4.center,
h5.center,
h6.center {
  text-align: center;
}

span.center {
  display: block;
  width: 100%;
  text-align: center;
}

img.center {
  margin: 0 auto;
  float: none;
}

.show-for-small,
.hide-for-desktop {
  display: none;
}

/* Typography */

body, table.body, h1, h2, h3, h4, h5, h6, p, td {
  color: #222222;
  font-family: "Helvetica", "Arial", sans-serif;
  font-weight: normal;
  padding:0;
  margin: 0;
  text-align: left;
  line-height: 1.3;
}

h1, h2, h3, h4, h5, h6 {
  word-break: normal;
}

/*h1 {font-size: 40px;}*/
h1 {font-size: 30px;}
/*h2 {font-size: 36px;}*/
h2 {font-size: 26px;}
h3 {font-size: 32px;}
h4 {font-size: 28px;}
h5 {font-size: 27px;}
h6 {font-size: 20px;}
body, table.body, p, td {font-size: 14px;line-height:19px;}

p.lead, p.lede, p.leed {
  font-size: 18px;
  line-height:21px;
}

p {
  margin-bottom: 10px;
}

small {
  font-size: 10px;
}

a {
  color: #2ba6cb;
  text-decoration: none;
}

a:hover {
  color: #2795b6 !important;
}

a:active {
  color: #2795b6 !important;
}

a:visited {
  color: #2ba6cb !important;
}

h1 a,
h2 a,
h3 a,
h4 a,
h5 a,
h6 a {
  color: #2ba6cb;
}

h1 a:active,
h2 a:active,
h3 a:active,
h4 a:active,
h5 a:active,
h6 a:active {
  color: #2ba6cb !important;
}

h1 a:visited,
h2 a:visited,
h3 a:visited,
h4 a:visited,
h5 a:visited,
h6 a:visited {
  color: #2ba6cb !important;
}

/* Panels */

.panel {
  background: #f2f2f2;
  border: 1px solid #d9d9d9;
  padding: 10px !important;
}

.sub-grid table {
  width: 100%;
}

.sub-grid td.sub-columns {
  padding-bottom: 0;
}

/* Buttons */

table.button,
table.tiny-button,
table.small-button,
table.medium-button,
table.large-button {
  width: 100%;
  overflow: hidden;
}

table.button td,
table.tiny-button td,
table.small-button td,
table.medium-button td,
table.large-button td {
  display: block;
  width: auto !important;
  text-align: center;
  background: #2ba6cb;
  border: 1px solid #2284a1;
  color: #ffffff;
  padding: 8px 0;
}

table.tiny-button td {
  padding: 5px 0 4px;
}

table.small-button td {
  padding: 8px 0 7px;
}

table.medium-button td {
  padding: 12px 0 10px;
}

table.large-button td {
  padding: 21px 0 18px;
}

table.button td a,
table.tiny-button td a,
table.small-button td a,
table.medium-button td a,
table.large-button td a {
  font-weight: bold;
  text-decoration: none;
  font-family: Helvetica, Arial, sans-serif;
  color: #ffffff;
  font-size: 16px;
}

table.tiny-button td a {
  font-size: 12px;
  font-weight: normal;
}

table.small-button td a {
  font-size: 16px;
}

table.medium-button td a {
  font-size: 20px;
}

table.large-button td a {
  font-size: 24px;
}

table.button:hover td,
table.button:visited td,
table.button:active td {
  background: #2795b6 !important;
}

table.button:hover td a,
table.button:visited td a,
table.button:active td a {
  color: #fff !important;
}

table.button:hover td,
table.tiny-button:hover td,
table.small-button:hover td,
table.medium-button:hover td,
table.large-button:hover td {
  background: #2795b6 !important;
}

table.button:hover td a,
table.button:active td a,
table.button td a:visited,
table.tiny-button:hover td a,
table.tiny-button:active td a,
table.tiny-button td a:visited,
table.small-button:hover td a,
table.small-button:active td a,
table.small-button td a:visited,
table.medium-button:hover td a,
table.medium-button:active td a,
table.medium-button td a:visited,
table.large-button:hover td a,
table.large-button:active td a,
table.large-button td a:visited {
  color: #ffffff !important;
}

table.secondary td {
  background: #e9e9e9;
  border-color: #d0d0d0;
  color: #555;
}

table.secondary td a {
  color: #555;
}

table.secondary:hover td {
  background: #d0d0d0 !important;
  color: #555;
}

table.secondary:hover td a,
table.secondary td a:visited,
table.secondary:active td a {
  color: #555 !important;
}

table.success td {
  background: #5da423;
  border-color: #457a1a;
}

table.success:hover td {
  background: #457a1a !important;
}

table.alert td {
  background: #c60f13;
  border-color: #970b0e;
}

table.alert:hover td {
  background: #970b0e !important;
}

table.radius td {
  -webkit-border-radius: 3px;
  -moz-border-radius: 3px;
  border-radius: 3px;
}

table.round td {
  -webkit-border-radius: 500px;
  -moz-border-radius: 500px;
  border-radius: 500px;
}

/* Outlook First */

body.outlook p {
  display: inline !important;
}

/*  Media Queries */

@media only screen and (max-width: 600px) {

  table[class="body"] img {
    width: auto !important;
    height: auto !important;
  }

  table[class="body"] center {
    min-width: 0 !important;
  }

  table[class="body"] .container {
    width: 95% !important;
  }

  table[class="body"] .row {
    width: 100% !important;
    display: block !important;
  }

  table[class="body"] .wrapper {
    display: block !important;
    padding-right: 0 !important;
  }

  table[class="body"] .columns,
  table[class="body"] .column {
    table-layout: fixed !important;
    float: none !important;
    width: 100% !important;
    padding-right: 0 !important;
    padding-left: 0 !important;
    display: block !important;
  }

  table[class="body"] .wrapper.first .columns,
  table[class="body"] .wrapper.first .column {
    display: table !important;
  }

  table[class="body"] table.columns td,
  table[class="body"] table.column td {
    width: 100% !important;
  }

  table[class="body"] .columns td.one,
  table[class="body"] .column td.one { width: 8.333333% !important; }
  table[class="body"] .columns td.two,
  table[class="body"] .column td.two { width: 16.666666% !important; }
  table[class="body"] .columns td.three,
  table[class="body"] .column td.three { width: 25% !important; }
  table[class="body"] .columns td.four,
  table[class="body"] .column td.four { width: 33.333333% !important; }
  table[class="body"] .columns td.five,
  table[class="body"] .column td.five { width: 41.666666% !important; }
  table[class="body"] .columns td.six,
  table[class="body"] .column td.six { width: 50% !important; }
  table[class="body"] .columns td.seven,
  table[class="body"] .column td.seven { width: 58.333333% !important; }
  table[class="body"] .columns td.eight,
  table[class="body"] .column td.eight { width: 66.666666% !important; }
  table[class="body"] .columns td.nine,
  table[class="body"] .column td.nine { width: 75% !important; }
  table[class="body"] .columns td.ten,
  table[class="body"] .column td.ten { width: 83.333333% !important; }
  table[class="body"] .columns td.eleven,
  table[class="body"] .column td.eleven { width: 91.666666% !important; }
  table[class="body"] .columns td.twelve,
  table[class="body"] .column td.twelve { width: 100% !important; }

  table[class="body"] td.offset-by-one,
  table[class="body"] td.offset-by-two,
  table[class="body"] td.offset-by-three,
  table[class="body"] td.offset-by-four,
  table[class="body"] td.offset-by-five,
  table[class="body"] td.offset-by-six,
  table[class="body"] td.offset-by-seven,
  table[class="body"] td.offset-by-eight,
  table[class="body"] td.offset-by-nine,
  table[class="body"] td.offset-by-ten,
  table[class="body"] td.offset-by-eleven {
    padding-left: 0 !important;
  }

  table[class="body"] table.columns td.expander {
    width: 1px !important;
  }

  table[class="body"] .right-text-pad,
  table[class="body"] .text-pad-right {
    padding-left: 10px !important;
  }

  table[class="body"] .left-text-pad,
  table[class="body"] .text-pad-left {
    padding-right: 10px !important;
  }

  table[class="body"] .hide-for-small,
  table[class="body"] .show-for-desktop {
    display: none !important;
  }

  table[class="body"] .show-for-small,
  table[class="body"] .hide-for-desktop {
    display: inherit !important;
  }
}

  </style>
  <style type="text/css">

    table.facebook td {
      background: #3b5998;
      border-color: #2d4473;
    }

    table.facebook:hover td {
      background: #2d4473 !important;
    }

    table.twitter td {
      background: #00acee;
      border-color: #0087bb;
    }

    table.twitter:hover td {
      background: #0087bb !important;
    }

    table.google-plus td {
      background-color: #DB4A39;
      border-color: #CC0000;
    }

    table.google-plus:hover td {
      background: #CC0000 !important;
    }

    .template-label {
      color: #ffffff;
      font-weight: bold;
      font-size: 11px;
    }

    .callout .wrapper {
      padding-bottom: 20px;
    }

    .callout .panel {
      background: #ECF8FF;
      border-color: #b9e5ff;
    }

    .header {
      background: #394e63;
      min-height:100px;
    }

    .footer .wrapper {
      background: #ebebeb;
    }

    .footer h5 {
      padding-bottom: 10px;
    }

    table.columns .text-pad {
      padding-left: 10px;
      padding-right: 10px;
    }

    table.columns .left-text-pad {
      padding-left: 10px;
    }

    table.columns .right-text-pad {
      padding-right: 10px;
    }

    @media only screen and (max-width: 600px) {

      table[class="body"] .right-text-pad {
        padding-left: 10px !important;
      }

      table[class="body"] .left-text-pad {
        padding-right: 10px !important;
      }
    }

  </style>
</head>
<body>
  <table class="body">
    <tr>
      <td class="center" align="center" valign="top">
        <center>

            [header:block]

          <table class="container">
            <tr>
              <td>

                [message:block]

                [callout:block]

                [footer:block]

                [privacy:block]

              <!-- container end below -->
              </td>
            </tr>
          </table>

        </center>
      </td>
    </tr>
  </table>
</body>
</html>', 'active', 'item', NULL, '2014-06-15 17:44:12', '2015-04-16 23:47:23' );
INSERT INTO `email_templates`(`name`,`description`,`content`,`status`,`type`,`parent`,`created_at`,`updated_at`) VALUES ( 'message', '', '<table class="row">
  <tr>
      <td class="wrapper last">

          <table class="twelve columns">
              <tr>
                  <td>
                      [title:block]
                      [body:block]
                  </td>
                  <td class="expander"></td>
              </tr>
          </table>

      </td>
  </tr>
</table>', 'active', 'item', 'main', '2014-06-20 09:15:51', '2015-04-12 01:40:00' );
INSERT INTO `email_templates`(`name`,`description`,`content`,`status`,`type`,`parent`,`created_at`,`updated_at`) VALUES ( 'password_reset', 'Восстановление пароля', '{"items":["body","header"],"title":"Смена логина или пароля для [user:name:last] [user:name:first] на сайте [site:name]","fields":[{"name":"body","value":"[user:name:last] [user:name:first],<br><br>Запрос на сброс пароля для вашего аккаунта был сделан на сайте [site:name]. <br><br>Вы можете сейчас войти на сайт, кликнув на ссылку или скопировав и вставив её в браузер: <br><br>[user:one-time-login-url] <br><br>Это одноразовая ссылка для входа и она перебросит Вас на страницу, где Вы сможете изменить пароль. Ссылка истекает через 1 сутки и ничего не случится, если она не будет использована. <br><br>-- Команда сайта [site:name]"}]}', 'active', 'template', NULL, '2014-06-20 18:35:00', '2015-04-16 17:27:39' );
INSERT INTO `email_templates`(`name`,`description`,`content`,`status`,`type`,`parent`,`created_at`,`updated_at`) VALUES ( 'privacy', '', '<table class="row">
    <tr>
        <td class="wrapper last">

            <table class="twelve columns">
                <tr>
                    <td align="center">
                        <center>
                            [privacy:content]
                        </center>
                    </td>
                    <td class="expander"></td>
                </tr>
            </table>
        </td>
    </tr>
</table>', 'active', 'item', 'main', '2014-06-15 17:45:42', '2015-04-12 01:38:28' );
INSERT INTO `email_templates`(`name`,`description`,`content`,`status`,`type`,`parent`,`created_at`,`updated_at`) VALUES ( 'register_admin_created', 'Добро пожаловать (новый пользователь создан администратором)', '{"items":["header","body"],"title":"Администратор создал для вас учётную запись","fields":[{"name":"body","value":"[user:name:last] [user:name:first],<br><br>Администратор системы [site:name] создал для вас аккаунт. Можете войти в систему, кликнув на ссылку или скопировав и вставив её в адресную строку браузера:<br><br>[user:one-time-login-url] <br><br>Эта одноразовая ссылка для входа в систему направит вас на страницу задания своего пароля.<br><br>После установки пароля вы сможете входить в систему через страницу<br>[site:login-url]<br><br>-- Команда [site:name]"}]}', 'active', 'template', '', '2014-06-20 18:27:05', '2017-01-15 17:03:11' );
INSERT INTO `email_templates`(`name`,`description`,`content`,`status`,`type`,`parent`,`created_at`,`updated_at`) VALUES ( 'social', '', '<table class="six columns">
    <tr>
        <td class="left-text-pad">

            <h3>Мы в сети:</h3>

            [facebook:block]
            [twitter:block]
            [vk:block]
            [google:block]

</td>
        <td class="expander"></td>
    </tr>
</table>', 'active', 'item', 'footer', '2014-06-15 17:46:47', '2015-06-23 14:10:14' );
INSERT INTO `email_templates`(`name`,`description`,`content`,`status`,`type`,`parent`,`created_at`,`updated_at`) VALUES ( 'status_activated', 'Учётная запись активирована', '{"items":["header","body"],"title":"Детали учётной записи для [user:name:last] [user:name:first] на [site:name] (одобрено)","fields":[{"name":"body","value":"[user:name:last] [user:name:first],<br><br>Ваш аккаунт на сайте [site:name] был активирован.<br><br>Вы можете войти на сайт, кликнув на ссылку или скопировав и вставив её в Ваш браузер: <br><br>[user:one-time-login-url] <br><br>Это одноразовая ссылка для входа и она перебросит Вас на страницу, где Вы сможете установить свой пароль.<br><br>После установки Вашего пароля, вы сможете заходить на странице [site:login-url].<br><br>-- Команда сайта [site:name]"}]}', 'active', 'template', '', '2014-06-20 18:31:06', '2017-01-15 06:07:52' );
INSERT INTO `email_templates`(`name`,`description`,`content`,`status`,`type`,`parent`,`created_at`,`updated_at`) VALUES ( 'title', '', '<h1>[title:content]</h1>
', 'active', 'item', 'message', '2014-06-21 21:54:48', '2015-04-14 15:20:18' );
INSERT INTO `email_templates`(`name`,`description`,`content`,`status`,`type`,`parent`,`created_at`,`updated_at`) VALUES ( 'twitter', '', '<table class="tiny-button twitter">
    <tr>
        <td>
            <a href="#">Twitter</a>
        </td>
    </tr>
</table>

<br>', 'active', 'item', 'social', '2014-06-15 17:47:08', '2015-04-13 01:47:52' );
INSERT INTO `email_templates`(`name`,`description`,`content`,`status`,`type`,`parent`,`created_at`,`updated_at`) VALUES ( 'vk', '', '<table class="tiny-button vk">
    <tr>
        <td>
            <a href="#">vk</a>
        </td>
    </tr>
</table>

<br>', 'active', 'item', 'social', '2014-06-15 17:50:15', '2015-04-13 01:51:30' );
-- ---------------------------------------------------------


-- Dump data of "flow_elements" ----------------------------
INSERT INTO `flow_elements`(`uuid`,`name`,`graph_settings`,`flow_id`,`created_at`,`update_at`,`prototype_type`,`status`,`description`,`script_id`,`flow_link`) VALUES ( '251c8846-db07-44a4-a1d3-8e3c5a885fac', 'invert power', '{"position":{"top":135,"left":280}}', '1', '2017-01-21 14:08:18', '2017-01-21 14:08:18', 'Task', 'enabled', '', '6', NULL );
INSERT INTO `flow_elements`(`uuid`,`name`,`graph_settings`,`flow_id`,`created_at`,`update_at`,`prototype_type`,`status`,`description`,`script_id`,`flow_link`) VALUES ( 'b602d3b4-1a32-4035-bc78-cc2b5e8778e7', 'start event', '{"position":{"top":155,"left":100}}', '1', '2017-01-21 14:08:18', '2017-01-21 14:08:18', 'MessageHandler', 'enabled', '', NULL, NULL );
INSERT INTO `flow_elements`(`uuid`,`name`,`graph_settings`,`flow_id`,`created_at`,`update_at`,`prototype_type`,`status`,`description`,`script_id`,`flow_link`) VALUES ( 'c192ad47-f3ab-498b-882c-dbea7f2ee5c0', 'end event', '{"position":{"top":155,"left":510}}', '1', '2017-01-21 14:08:18', '2017-01-21 14:08:18', 'MessageEmitter', 'enabled', '', NULL, NULL );
-- ---------------------------------------------------------


-- Dump data of "flows" ------------------------------------
INSERT INTO `flows`(`id`,`name`,`status`,`created_at`,`update_at`,`workflow_id`,`description`) VALUES ( '1', 'Розетки', 'enabled', '0000-00-00 00:00:00', '2017-01-21 14:08:18', '1', '' );
-- ---------------------------------------------------------


-- Dump data of "images" -----------------------------------
INSERT INTO `images`(`id`,`thumb`,`image`,`mime_type`,`title`,`size`,`name`,`created_at`) VALUES ( '1', '', '531ae6400425797a88afd36a888c69e5.svg', 'image/svg+xml', '', '8484', 'map-chematic-original.svg', '2017-01-21 14:17:59' );
INSERT INTO `images`(`id`,`thumb`,`image`,`mime_type`,`title`,`size`,`name`,`created_at`) VALUES ( '2', '', 'e883a3b833eee0c15926633d7dc91ea2.svg', 'image/svg+xml', '', '7326', 'socket_v1_def.svg', '2017-01-21 14:22:06' );
INSERT INTO `images`(`id`,`thumb`,`image`,`mime_type`,`title`,`size`,`name`,`created_at`) VALUES ( '3', '', '98d8902925c4f3ccb394182d30966fad.svg', 'image/svg+xml', '', '7326', 'socket_v1_g.svg', '2017-01-21 14:22:06' );
INSERT INTO `images`(`id`,`thumb`,`image`,`mime_type`,`title`,`size`,`name`,`created_at`) VALUES ( '4', '', 'a75a6dd113828adccdc68c676eb2f2e5.svg', 'image/svg+xml', '', '7326', 'socket_v1_r.svg', '2017-01-21 14:22:06' );
INSERT INTO `images`(`id`,`thumb`,`image`,`mime_type`,`title`,`size`,`name`,`created_at`) VALUES ( '5', '', '01a1780763aa1eedec1a9d92b235d260.svg', 'image/svg+xml', '', '7326', 'socket_v1_b.svg', '2017-01-21 14:22:06' );
INSERT INTO `images`(`id`,`thumb`,`image`,`mime_type`,`title`,`size`,`name`,`created_at`) VALUES ( '6', '', '20bd64f497a2d9650cd239f8de6545e3.svg', 'image/svg+xml', '', '8484', 'map-chematic-original.svg', '2017-01-21 14:22:07' );
INSERT INTO `images`(`id`,`thumb`,`image`,`mime_type`,`title`,`size`,`name`,`created_at`) VALUES ( '7', '', '1249feecd0c2ba6f06b383b6bb5ea389.svg', 'image/svg+xml', '', '3212', 'button_v1_refresh.svg', '2017-01-21 14:22:07' );
INSERT INTO `images`(`id`,`thumb`,`image`,`mime_type`,`title`,`size`,`name`,`created_at`) VALUES ( '8', '', 'cf784e9045151fbf211f1fe602a1b3fd.svg', 'image/svg+xml', '', '2398', 'button_v1_on.svg', '2017-01-21 14:22:07' );
INSERT INTO `images`(`id`,`thumb`,`image`,`mime_type`,`title`,`size`,`name`,`created_at`) VALUES ( '9', '', '8039b98f01765d7ff69f1b24ef56f3a9.svg', 'image/svg+xml', '', '2518', 'button_v1_off.svg', '2017-01-21 14:22:07' );
-- ---------------------------------------------------------


-- Dump data of "logs" -------------------------------------
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '1', '[I] AppPath: /home/delta54/workspace/golang/smart-home/server/src ', '2017-01-21 12:03:34', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '2', '[I] Development mode enabled ', '2017-01-21 12:03:34', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '3', '[I] Filters initialize... ', '2017-01-21 12:03:34', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '4', '[I] Crontab initialize... ', '2017-01-21 12:03:34', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '5', '[I] Notifr initialize... ', '2017-01-21 12:03:34', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '6', '[I] Telemetry initialize... ', '2017-01-21 12:03:34', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '7', '[I] Core initialize... ', '2017-01-21 12:03:34', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '8', '[I] Starting.... ', '2017-01-21 12:03:34', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '9', '[I] new sockjs session established, from ip: 127.0.0.1', '2017-01-21 12:03:35', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '10', '[W] rbac: <QuerySeter> no row found', '2017-01-21 12:03:39', 'Warning' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '11', '[I] sockjs session from ip: 127.0.0.1 closed', '2017-01-21 12:03:40', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '12', '[I] new sockjs session established, from ip: 127.0.0.1', '2017-01-21 12:03:41', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '13', '[W] rbac: <QuerySeter> no row found', '2017-01-21 12:03:45', 'Warning' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '14', '[I] sockjs session from ip: 127.0.0.1 closed', '2017-01-21 12:05:02', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '15', '[W] rbac: <QuerySeter> no row found', '2017-01-21 12:05:03', 'Warning' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '16', '[I] new sockjs session established, from ip: 127.0.0.1', '2017-01-21 12:05:03', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '17', '[I] sockjs session from ip: 127.0.0.1 closed', '2017-01-21 12:05:04', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '18', '[I] new sockjs session established, from ip: 127.0.0.1', '2017-01-21 12:05:04', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '19', '[I] sockjs session from ip: 127.0.0.1 closed', '2017-01-21 12:05:06', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '20', '[I] new sockjs session established, from ip: 127.0.0.1', '2017-01-21 12:05:07', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '21', '[W] rbac: <QuerySeter> no row found', '2017-01-21 12:05:09', 'Warning' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '22', '[I] sockjs session from ip: 127.0.0.1 closed', '2017-01-21 12:05:10', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '23', '[I] new sockjs session established, from ip: 127.0.0.1', '2017-01-21 12:05:11', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '24', '[I] sockjs session from ip: 127.0.0.1 closed', '2017-01-21 12:05:13', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '25', '[I] new sockjs session established, from ip: 127.0.0.1', '2017-01-21 12:05:14', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '26', '[W] rbac: <QuerySeter> no row found', '2017-01-21 12:05:15', 'Warning' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '27', '[I] sockjs session from ip: 127.0.0.1 closed', '2017-01-21 12:05:56', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '28', '[W] rbac: <QuerySeter> no row found', '2017-01-21 12:05:57', 'Warning' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '29', '[I] new sockjs session established, from ip: 127.0.0.1', '2017-01-21 12:05:57', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '30', '[I] sockjs session from ip: 127.0.0.1 closed', '2017-01-21 12:05:59', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '31', '[I] new sockjs session established, from ip: 127.0.0.1', '2017-01-21 12:06:00', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '32', '[W] rbac: <QuerySeter> no row found', '2017-01-21 12:06:03', 'Warning' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '33', '[I] sockjs session from ip: 127.0.0.1 closed', '2017-01-21 12:06:47', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '34', '[I] new sockjs session established, from ip: 127.0.0.1', '2017-01-21 12:06:49', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '35', '[W] rbac: <QuerySeter> no row found', '2017-01-21 12:06:51', 'Warning' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '36', '[I] sockjs session from ip: 127.0.0.1 closed', '2017-01-21 12:07:01', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '37', '[W] rbac: <QuerySeter> no row found', '2017-01-21 12:07:03', 'Warning' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '38', '[I] new sockjs session established, from ip: 127.0.0.1', '2017-01-21 12:07:14', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '39', '[W] rbac: <QuerySeter> no row found', '2017-01-21 12:07:15', 'Warning' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '40', '[I] Add node: "darkastar"', '2017-01-21 12:07:45', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '41', '[I] Node dial tcp 127.0.0.1:3000 ... ok', '2017-01-21 12:07:45', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '42', '[W] rbac: <QuerySeter> no row found', '2017-01-21 12:07:51', 'Warning' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '43', '[I] Add node: "ganimed"', '2017-01-21 12:08:01', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '44', '[W] rbac: <QuerySeter> no row found', '2017-01-21 12:08:03', 'Warning' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '45', '[E] Node error dial tcp 10.0.36.2:3000: i/o timeout', '2017-01-21 12:08:21', 'Error' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '46', '[W] rbac: <QuerySeter> no row found', '2017-01-21 12:08:27', 'Warning' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '47', '[I] Reload node: "ganimed"', '2017-01-21 12:08:35', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '48', '[W] rbac: <QuerySeter> no row found', '2017-01-21 12:08:39', 'Warning' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '49', '[I] Reload node: "darkastar"', '2017-01-21 12:09:40', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '50', '[I] Node dial tcp 127.0.0.1:3001 ... ok', '2017-01-21 12:09:41', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '51', '[W] rbac: <QuerySeter> no row found', '2017-01-21 12:09:45', 'Warning' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '52', '[I] Reload node: "ganimed"', '2017-01-21 12:09:46', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '53', '[W] rbac: <QuerySeter> no row found', '2017-01-21 12:09:51', 'Warning' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '54', '[I] Node dial tcp 10.0.36.1:3001 ... ok', '2017-01-21 12:10:30', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '55', '[W] rbac: <QuerySeter> no row found', '2017-01-21 12:10:33', 'Warning' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '56', '[I] sockjs session from ip: 127.0.0.1 closed', '2017-01-21 12:10:50', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '57', '[W] rbac: <QuerySeter> no row found', '2017-01-21 12:10:51', 'Warning' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '58', '[I] new sockjs session established, from ip: 127.0.0.1', '2017-01-21 12:10:51', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '59', '[W] rbac: <QuerySeter> no row found', '2017-01-21 12:10:57', 'Warning' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '60', '[I] sockjs session from ip: 127.0.0.1 closed', '2017-01-21 12:13:41', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '61', '[I] new sockjs session established, from ip: 127.0.0.1', '2017-01-21 12:13:43', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '62', '[I] sockjs session from ip: 127.0.0.1 closed', '2017-01-21 12:13:44', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '63', '[I] new sockjs session established, from ip: 127.0.0.1', '2017-01-21 12:13:44', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '64', '[W] rbac: <QuerySeter> no row found', '2017-01-21 12:13:45', 'Warning' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '65', '[I] sockjs session from ip: 127.0.0.1 closed', '2017-01-21 12:13:46', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '66', '[I] new sockjs session established, from ip: 127.0.0.1', '2017-01-21 12:13:46', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '67', '[W] rbac: <QuerySeter> no row found', '2017-01-21 12:13:51', 'Warning' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '68', '[I] sockjs session from ip: 127.0.0.1 closed', '2017-01-21 12:14:00', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '69', '[I] new sockjs session established, from ip: 127.0.0.1', '2017-01-21 12:14:01', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '70', '[W] rbac: <QuerySeter> no row found', '2017-01-21 12:14:03', 'Warning' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '71', '[I] sockjs session from ip: 127.0.0.1 closed', '2017-01-21 12:14:18', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '72', '[I] new sockjs session established, from ip: 127.0.0.1', '2017-01-21 12:14:19', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '73', '[I] sockjs session from ip: 127.0.0.1 closed', '2017-01-21 12:14:20', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '74', '[I] new sockjs session established, from ip: 127.0.0.1', '2017-01-21 12:14:20', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '75', '[W] rbac: <QuerySeter> no row found', '2017-01-21 12:14:21', 'Warning' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '76', '[I] sockjs session from ip: 127.0.0.1 closed', '2017-01-21 12:14:40', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '77', '[I] new sockjs session established, from ip: 127.0.0.1', '2017-01-21 12:14:41', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '78', '[I] sockjs session from ip: 127.0.0.1 closed', '2017-01-21 12:14:41', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '79', '[I] new sockjs session established, from ip: 127.0.0.1', '2017-01-21 12:14:42', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '80', '[W] rbac: <QuerySeter> no row found', '2017-01-21 12:14:45', 'Warning' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '81', '[I] sockjs session from ip: 127.0.0.1 closed', '2017-01-21 12:15:21', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '82', '[I] new sockjs session established, from ip: 127.0.0.1', '2017-01-21 12:15:22', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '83', '[W] rbac: <QuerySeter> no row found', '2017-01-21 12:15:27', 'Warning' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '84', '[I] sockjs session from ip: 127.0.0.1 closed', '2017-01-21 12:20:02', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '85', '[I] new sockjs session established, from ip: 127.0.0.1', '2017-01-21 12:20:03', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '86', '[I] sockjs session from ip: 127.0.0.1 closed', '2017-01-21 12:20:34', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '87', '[I] new sockjs session established, from ip: 127.0.0.1', '2017-01-21 12:20:35', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '88', '[I] sockjs session from ip: 127.0.0.1 closed', '2017-01-21 12:34:38', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '89', '[I] new sockjs session established, from ip: 127.0.0.1', '2017-01-21 12:34:39', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '90', '[I] sockjs session from ip: 127.0.0.1 closed', '2017-01-21 12:35:04', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '91', '[I] new sockjs session established, from ip: 127.0.0.1', '2017-01-21 12:35:05', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '92', '[I] sockjs session from ip: 127.0.0.1 closed', '2017-01-21 12:35:11', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '93', '[I] new sockjs session established, from ip: 127.0.0.1', '2017-01-21 12:35:12', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '94', '[I] sockjs session from ip: 127.0.0.1 closed', '2017-01-21 12:35:31', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '95', '[I] new sockjs session established, from ip: 127.0.0.1', '2017-01-21 12:35:32', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '96', '[I] sockjs session from ip: 127.0.0.1 closed', '2017-01-21 12:40:15', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '97', '[I] new sockjs session established, from ip: 127.0.0.1', '2017-01-21 12:40:16', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '98', '[I] sockjs session from ip: 127.0.0.1 closed', '2017-01-21 12:41:54', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '99', '[I] new sockjs session established, from ip: 127.0.0.1', '2017-01-21 12:41:55', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '100', '[I] sockjs session from ip: 127.0.0.1 closed', '2017-01-21 12:42:25', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '101', '[I] new sockjs session established, from ip: 127.0.0.1', '2017-01-21 12:42:26', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '102', '[I] sockjs session from ip: 127.0.0.1 closed', '2017-01-21 12:42:55', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '103', '[I] new sockjs session established, from ip: 127.0.0.1', '2017-01-21 12:42:56', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '104', '[I] sockjs session from ip: 127.0.0.1 closed', '2017-01-21 12:43:45', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '105', '[I] new sockjs session established, from ip: 127.0.0.1', '2017-01-21 12:43:45', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '106', '[I] sockjs session from ip: 127.0.0.1 closed', '2017-01-21 12:44:59', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '107', '[I] new sockjs session established, from ip: 127.0.0.1', '2017-01-21 12:45:00', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '108', '[I] sockjs session from ip: 127.0.0.1 closed', '2017-01-21 12:45:07', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '109', '[I] new sockjs session established, from ip: 127.0.0.1', '2017-01-21 12:45:07', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '110', '[I] sockjs session from ip: 127.0.0.1 closed', '2017-01-21 12:45:30', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '111', '[I] new sockjs session established, from ip: 127.0.0.1', '2017-01-21 12:45:31', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '112', '[I] sockjs session from ip: 127.0.0.1 closed', '2017-01-21 12:48:23', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '113', '[I] new sockjs session established, from ip: 127.0.0.1', '2017-01-21 12:48:25', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '114', '[I] sockjs session from ip: 127.0.0.1 closed', '2017-01-21 12:50:32', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '115', '[I] new sockjs session established, from ip: 127.0.0.1', '2017-01-21 12:50:33', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '116', '[I] sockjs session from ip: 127.0.0.1 closed', '2017-01-21 12:53:02', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '117', '[I] new sockjs session established, from ip: 127.0.0.1', '2017-01-21 12:53:03', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '118', '[I] sockjs session from ip: 127.0.0.1 closed', '2017-01-21 12:53:15', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '119', '[I] new sockjs session established, from ip: 127.0.0.1', '2017-01-21 12:53:17', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '120', '[I] sockjs session from ip: 127.0.0.1 closed', '2017-01-21 12:53:18', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '121', '[I] new sockjs session established, from ip: 127.0.0.1', '2017-01-21 12:53:19', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '122', '[I] sockjs session from ip: 127.0.0.1 closed', '2017-01-21 12:53:31', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '123', '[I] new sockjs session established, from ip: 127.0.0.1', '2017-01-21 12:53:32', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '124', '[I] sockjs session from ip: 127.0.0.1 closed', '2017-01-21 12:53:47', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '125', '[I] new sockjs session established, from ip: 127.0.0.1', '2017-01-21 12:53:48', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '126', '[I] sockjs session from ip: 127.0.0.1 closed', '2017-01-21 12:54:42', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '127', '[I] new sockjs session established, from ip: 127.0.0.1', '2017-01-21 12:54:43', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '128', '[I] sockjs session from ip: 127.0.0.1 closed', '2017-01-21 12:55:01', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '129', '[I] new sockjs session established, from ip: 127.0.0.1', '2017-01-21 12:55:02', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '130', '[I] sockjs session from ip: 127.0.0.1 closed', '2017-01-21 12:55:04', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '131', '[I] new sockjs session established, from ip: 127.0.0.1', '2017-01-21 12:55:05', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '132', '[I] sockjs session from ip: 127.0.0.1 closed', '2017-01-21 13:25:49', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '133', '[I] new sockjs session established, from ip: 127.0.0.1', '2017-01-21 13:25:50', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '134', '[I] sockjs session from ip: 127.0.0.1 closed', '2017-01-21 13:26:39', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '135', '[I] new sockjs session established, from ip: 127.0.0.1', '2017-01-21 13:26:40', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '136', '[I] Reload node: "darkstar"', '2017-01-21 13:27:14', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '137', '[I] Node dial tcp 127.0.0.1:3001 ... ok', '2017-01-21 13:27:15', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '138', '[I] new sockjs session established, from ip: 127.0.0.1', '2017-01-21 13:30:58', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '139', '[I] sockjs session from ip: 127.0.0.1 closed', '2017-01-21 13:38:55', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '140', '[I] new sockjs session established, from ip: 127.0.0.1', '2017-01-21 13:38:56', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '141', '[I] AppPath: /home/delta54/workspace/golang/smart-home/server/src ', '2017-01-21 13:41:48', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '142', '[I] Development mode enabled ', '2017-01-21 13:41:48', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '143', '[I] Filters initialize... ', '2017-01-21 13:41:48', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '144', '[I] Crontab initialize... ', '2017-01-21 13:41:48', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '145', '[I] Notifr initialize... ', '2017-01-21 13:41:48', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '146', '[I] Telemetry initialize... ', '2017-01-21 13:41:49', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '147', '[I] Core initialize... ', '2017-01-21 13:41:49', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '148', '[I] Add node: "darkstar"', '2017-01-21 13:41:49', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '149', '[I] Add node: "ganimed"', '2017-01-21 13:41:49', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '150', '[I] Node dial tcp 127.0.0.1:3001 ... ok', '2017-01-21 13:41:49', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '151', '[I] Starting.... ', '2017-01-21 13:41:49', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '152', '[I] Node dial tcp 10.0.36.1:3001 ... ok', '2017-01-21 13:41:49', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '153', '[I] new sockjs session established, from ip: 127.0.0.1', '2017-01-21 13:41:49', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '154', '[I] sockjs session from ip: 127.0.0.1 closed', '2017-01-21 13:42:25', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '155', '[I] new sockjs session established, from ip: 127.0.0.1', '2017-01-21 13:42:26', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '156', '[I] sockjs session from ip: 127.0.0.1 closed', '2017-01-21 13:43:28', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '157', '[I] new sockjs session established, from ip: 127.0.0.1', '2017-01-21 13:43:29', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '158', '[I] sockjs session from ip: 127.0.0.1 closed', '2017-01-21 13:44:38', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '159', '[I] new sockjs session established, from ip: 127.0.0.1', '2017-01-21 13:44:42', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '160', '[I] AppPath: /home/delta54/workspace/golang/smart-home/server/src ', '2017-01-21 13:50:45', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '161', '[I] Development mode enabled ', '2017-01-21 13:50:45', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '162', '[I] Filters initialize... ', '2017-01-21 13:50:45', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '163', '[I] Crontab initialize... ', '2017-01-21 13:50:45', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '164', '[I] Notifr initialize... ', '2017-01-21 13:50:45', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '165', '[I] Telemetry initialize... ', '2017-01-21 13:50:45', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '166', '[I] Core initialize... ', '2017-01-21 13:50:45', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '167', '[I] Add node: "darkstar"', '2017-01-21 13:50:45', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '168', '[I] Add node: "ganimed"', '2017-01-21 13:50:45', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '169', '[I] Node dial tcp 127.0.0.1:3001 ... ok', '2017-01-21 13:50:45', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '170', '[I] Starting.... ', '2017-01-21 13:50:45', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '171', '[I] Node dial tcp 10.0.36.1:3001 ... ok', '2017-01-21 13:50:45', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '172', '[I] new sockjs session established, from ip: 127.0.0.1', '2017-01-21 13:50:46', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '173', '[I] AppPath: /home/delta54/workspace/golang/smart-home/server/src ', '2017-01-21 13:52:12', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '174', '[I] Development mode enabled ', '2017-01-21 13:52:12', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '175', '[I] Filters initialize... ', '2017-01-21 13:52:12', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '176', '[I] Crontab initialize... ', '2017-01-21 13:52:12', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '177', '[I] Notifr initialize... ', '2017-01-21 13:52:12', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '178', '[I] Telemetry initialize... ', '2017-01-21 13:52:12', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '179', '[I] Core initialize... ', '2017-01-21 13:52:13', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '180', '[I] Add node: "darkstar"', '2017-01-21 13:52:13', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '181', '[I] Add node: "ganimed"', '2017-01-21 13:52:13', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '182', '[I] Node dial tcp 127.0.0.1:3001 ... ok', '2017-01-21 13:52:13', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '183', '[I] Starting.... ', '2017-01-21 13:52:13', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '184', '[I] Node dial tcp 10.0.36.1:3001 ... ok', '2017-01-21 13:52:13', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '185', '[I] new sockjs session established, from ip: 127.0.0.1', '2017-01-21 13:52:13', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '186', '[I] AppPath: /home/delta54/workspace/golang/smart-home/server/src ', '2017-01-21 13:53:28', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '187', '[I] Development mode enabled ', '2017-01-21 13:53:28', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '188', '[I] Filters initialize... ', '2017-01-21 13:53:28', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '189', '[I] Crontab initialize... ', '2017-01-21 13:53:28', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '190', '[I] Notifr initialize... ', '2017-01-21 13:53:28', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '191', '[I] Telemetry initialize... ', '2017-01-21 13:53:28', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '192', '[I] Core initialize... ', '2017-01-21 13:53:28', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '193', '[I] Add node: "darkstar"', '2017-01-21 13:53:28', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '194', '[I] Add node: "ganimed"', '2017-01-21 13:53:28', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '195', '[I] Node dial tcp 127.0.0.1:3001 ... ok', '2017-01-21 13:53:28', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '196', '[I] Starting.... ', '2017-01-21 13:53:28', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '197', '[I] Node dial tcp 10.0.36.1:3001 ... ok', '2017-01-21 13:53:28', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '198', '[I] new sockjs session established, from ip: 127.0.0.1', '2017-01-21 13:53:29', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '199', '[I] AppPath: /home/delta54/workspace/golang/smart-home/server/src ', '2017-01-21 13:54:18', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '200', '[I] Development mode enabled ', '2017-01-21 13:54:18', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '201', '[I] Filters initialize... ', '2017-01-21 13:54:18', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '202', '[I] Crontab initialize... ', '2017-01-21 13:54:19', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '203', '[I] Notifr initialize... ', '2017-01-21 13:54:19', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '204', '[I] Telemetry initialize... ', '2017-01-21 13:54:19', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '205', '[I] Core initialize... ', '2017-01-21 13:54:19', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '206', '[I] Add node: "darkstar"', '2017-01-21 13:54:19', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '207', '[I] Add node: "ganimed"', '2017-01-21 13:54:19', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '208', '[I] Node dial tcp 127.0.0.1:3001 ... ok', '2017-01-21 13:54:19', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '209', '[I] Starting.... ', '2017-01-21 13:54:19', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '210', '[I] Node dial tcp 10.0.36.1:3001 ... ok', '2017-01-21 13:54:19', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '211', '[I] new sockjs session established, from ip: 127.0.0.1', '2017-01-21 13:54:19', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '212', '[I] AppPath: /home/delta54/workspace/golang/smart-home/server/src ', '2017-01-21 13:57:12', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '213', '[I] Development mode enabled ', '2017-01-21 13:57:12', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '214', '[I] Filters initialize... ', '2017-01-21 13:57:12', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '215', '[I] Crontab initialize... ', '2017-01-21 13:57:12', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '216', '[I] Notifr initialize... ', '2017-01-21 13:57:12', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '217', '[I] Telemetry initialize... ', '2017-01-21 13:57:12', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '218', '[I] Core initialize... ', '2017-01-21 13:57:13', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '219', '[I] Add node: "darkstar"', '2017-01-21 13:57:13', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '220', '[I] Add node: "ganimed"', '2017-01-21 13:57:13', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '221', '[I] Node dial tcp 127.0.0.1:3001 ... ok', '2017-01-21 13:57:13', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '222', '[I] Starting.... ', '2017-01-21 13:57:13', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '223', '[I] Node dial tcp 10.0.36.1:3001 ... ok', '2017-01-21 13:57:13', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '224', '[I] new sockjs session established, from ip: 127.0.0.1', '2017-01-21 13:57:13', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '225', '[I] sockjs session from ip: 127.0.0.1 closed', '2017-01-21 13:57:51', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '226', '[I] new sockjs session established, from ip: 127.0.0.1', '2017-01-21 13:57:53', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '227', '[I] Add workflow: Новый процесс ', '2017-01-21 14:00:00', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '228', '[I] Add flow: Розетки ', '2017-01-21 14:00:36', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '229', '[I] Remove flow: Розетки ', '2017-01-21 14:02:12', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '230', '[I] Add flow: Розетки ', '2017-01-21 14:02:12', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '231', '[I] Add worker: "Действие"', '2017-01-21 14:02:12', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '232', '[I] Remove flow: Розетки ', '2017-01-21 14:06:55', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '233', '[I] Remove worker: "Действие"', '2017-01-21 14:06:55', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '234', '[I] Add flow: Розетки ', '2017-01-21 14:06:55', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '235', '[I] Add worker: "Действие"', '2017-01-21 14:06:55', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '236', '[I] sockjs session from ip: 127.0.0.1 closed', '2017-01-21 14:07:10', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '237', '[I] new sockjs session established, from ip: 127.0.0.1', '2017-01-21 14:07:12', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '238', '[I] Remove flow: Розетки ', '2017-01-21 14:08:18', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '239', '[I] Remove worker: "Действие"', '2017-01-21 14:08:18', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '240', '[I] Add flow: Розетки ', '2017-01-21 14:08:18', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '241', '[I] Add worker: "Действие"', '2017-01-21 14:08:18', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '242', '[I] sockjs session from ip: 127.0.0.1 closed', '2017-01-21 14:09:31', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '243', '[I] new sockjs session established, from ip: 127.0.0.1', '2017-01-21 14:09:33', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '244', '[I] sockjs session from ip: 127.0.0.1 closed', '2017-01-21 14:09:39', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '245', '[I] new sockjs session established, from ip: 127.0.0.1', '2017-01-21 14:09:41', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '246', '[I] sockjs session from ip: 127.0.0.1 closed', '2017-01-21 14:10:53', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '247', '[I] new sockjs session established, from ip: 127.0.0.1', '2017-01-21 14:10:54', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '248', '[I] sockjs session from ip: 127.0.0.1 closed', '2017-01-21 14:11:22', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '249', '[I] new sockjs session established, from ip: 127.0.0.1', '2017-01-21 14:11:24', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '250', '[I] sockjs session from ip: 127.0.0.1 closed', '2017-01-21 14:11:50', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '251', '[I] new sockjs session established, from ip: 127.0.0.1', '2017-01-21 14:11:52', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '252', '[I] sockjs session from ip: 127.0.0.1 closed', '2017-01-21 14:12:30', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '253', '[I] new sockjs session established, from ip: 127.0.0.1', '2017-01-21 14:12:31', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '254', '[I] sockjs session from ip: 127.0.0.1 closed', '2017-01-21 14:23:54', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '255', '[I] new sockjs session established, from ip: 127.0.0.1', '2017-01-21 14:23:55', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '256', '[I] sockjs session from ip: 127.0.0.1 closed', '2017-01-21 14:25:48', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '257', '[I] new sockjs session established, from ip: 127.0.0.1', '2017-01-21 14:25:50', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '258', '[I] sockjs session from ip: 127.0.0.1 closed', '2017-01-21 14:27:54', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '259', '[I] new sockjs session established, from ip: 127.0.0.1', '2017-01-21 14:27:56', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '260', '[I] sockjs session from ip: 127.0.0.1 closed', '2017-01-21 14:28:37', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '261', '[I] new sockjs session established, from ip: 127.0.0.1', '2017-01-21 14:28:39', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '262', '[I] sockjs session from ip: 127.0.0.1 closed', '2017-01-21 14:28:40', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '263', '[I] new sockjs session established, from ip: 127.0.0.1', '2017-01-21 14:28:41', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '264', '[I] sockjs session from ip: 127.0.0.1 closed', '2017-01-21 14:28:42', 'Info' );
INSERT INTO `logs`(`id`,`body`,`created_at`,`level`) VALUES ( '265', '[I] new sockjs session established, from ip: 127.0.0.1', '2017-01-21 14:28:44', 'Info' );
-- ---------------------------------------------------------


-- Dump data of "map_device_actions" -----------------------
INSERT INTO `map_device_actions`(`id`,`device_action_id`,`map_device_id`,`image_id`,`type`,`created_at`,`update_at`) VALUES ( '67', '1', '13', '7', '', '2017-01-21 14:26:16', '2017-01-21 14:26:16' );
INSERT INTO `map_device_actions`(`id`,`device_action_id`,`map_device_id`,`image_id`,`type`,`created_at`,`update_at`) VALUES ( '68', '3', '13', '8', '', '2017-01-21 14:26:16', '2017-01-21 14:26:16' );
INSERT INTO `map_device_actions`(`id`,`device_action_id`,`map_device_id`,`image_id`,`type`,`created_at`,`update_at`) VALUES ( '69', '4', '13', '9', '', '2017-01-21 14:26:16', '2017-01-21 14:26:16' );
INSERT INTO `map_device_actions`(`id`,`device_action_id`,`map_device_id`,`image_id`,`type`,`created_at`,`update_at`) VALUES ( '73', '1', '10', '7', '', '2017-01-21 14:26:31', '2017-01-21 14:26:31' );
INSERT INTO `map_device_actions`(`id`,`device_action_id`,`map_device_id`,`image_id`,`type`,`created_at`,`update_at`) VALUES ( '74', '3', '10', '8', '', '2017-01-21 14:26:31', '2017-01-21 14:26:31' );
INSERT INTO `map_device_actions`(`id`,`device_action_id`,`map_device_id`,`image_id`,`type`,`created_at`,`update_at`) VALUES ( '75', '4', '10', '9', '', '2017-01-21 14:26:31', '2017-01-21 14:26:31' );
-- ---------------------------------------------------------


-- Dump data of "map_device_states" ------------------------
INSERT INTO `map_device_states`(`id`,`device_state_id`,`map_device_id`,`image_id`,`created_at`,`update_at`,`style`) VALUES ( '67', '1', '13', '3', '2017-01-21 14:26:16', '2017-01-21 14:26:16', '' );
INSERT INTO `map_device_states`(`id`,`device_state_id`,`map_device_id`,`image_id`,`created_at`,`update_at`,`style`) VALUES ( '68', '2', '13', '2', '2017-01-21 14:26:16', '2017-01-21 14:26:16', '' );
INSERT INTO `map_device_states`(`id`,`device_state_id`,`map_device_id`,`image_id`,`created_at`,`update_at`,`style`) VALUES ( '69', '3', '13', '4', '2017-01-21 14:26:16', '2017-01-21 14:26:16', '' );
INSERT INTO `map_device_states`(`id`,`device_state_id`,`map_device_id`,`image_id`,`created_at`,`update_at`,`style`) VALUES ( '73', '1', '10', '3', '2017-01-21 14:26:31', '2017-01-21 14:26:31', '' );
INSERT INTO `map_device_states`(`id`,`device_state_id`,`map_device_id`,`image_id`,`created_at`,`update_at`,`style`) VALUES ( '74', '2', '10', '2', '2017-01-21 14:26:31', '2017-01-21 14:26:31', '' );
INSERT INTO `map_device_states`(`id`,`device_state_id`,`map_device_id`,`image_id`,`created_at`,`update_at`,`style`) VALUES ( '75', '3', '10', '4', '2017-01-21 14:26:31', '2017-01-21 14:26:31', '' );
-- ---------------------------------------------------------


-- Dump data of "map_devices" ------------------------------
INSERT INTO `map_devices`(`id`,`created_at`,`update_at`,`device_id`,`image_id`) VALUES ( '10', '2017-01-21 14:26:31', '2017-01-21 14:26:31', '2', '5' );
INSERT INTO `map_devices`(`id`,`created_at`,`update_at`,`device_id`,`image_id`) VALUES ( '13', '2017-01-21 14:26:16', '2017-01-21 14:26:16', '6', '5' );
-- ---------------------------------------------------------


-- Dump data of "map_elements" -----------------------------
INSERT INTO `map_elements`(`id`,`name`,`description`,`created_at`,`update_at`,`prototype_type`,`layer_id`,`map_id`,`graph_settings`,`status`,`weight`,`prototype_id`) VALUES ( '3', 'Новый элемент №1', '', '0000-00-00 00:00:00', '2017-01-21 14:18:29', 'image', '1', '1', '{"width":1336.453125,"height":1543.109375,"position":{"top":0,"left":0}}', 'frozen', '1', '3' );
INSERT INTO `map_elements`(`id`,`name`,`description`,`created_at`,`update_at`,`prototype_type`,`layer_id`,`map_id`,`graph_settings`,`status`,`weight`,`prototype_id`) VALUES ( '4', 'Спальня', '', '0000-00-00 00:00:00', '2017-01-21 14:20:54', 'text', '1', '1', '{"width":null,"height":null,"position":{"top":1017,"left":170}}', 'frozen', '0', '13' );
INSERT INTO `map_elements`(`id`,`name`,`description`,`created_at`,`update_at`,`prototype_type`,`layer_id`,`map_id`,`graph_settings`,`status`,`weight`,`prototype_id`) VALUES ( '5', 'р1', '', '0000-00-00 00:00:00', '2017-01-21 14:26:31', 'device', '2', '1', '{"width":35.984375,"height":36,"position":{"top":1149,"left":380}}', 'enabled', '1', '10' );
INSERT INTO `map_elements`(`id`,`name`,`description`,`created_at`,`update_at`,`prototype_type`,`layer_id`,`map_id`,`graph_settings`,`status`,`weight`,`prototype_id`) VALUES ( '6', 'р2', '', '0000-00-00 00:00:00', '2017-01-21 14:26:16', 'device', '2', '1', '{"width":31.953125,"height":32,"position":{"top":854,"left":64}}', 'enabled', '0', '13' );
-- ---------------------------------------------------------


-- Dump data of "map_images" -------------------------------
INSERT INTO `map_images`(`id`,`image_id`,`style`) VALUES ( '3', '1', '' );
-- ---------------------------------------------------------


-- Dump data of "map_layers" -------------------------------
INSERT INTO `map_layers`(`id`,`name`,`description`,`created_at`,`update_at`,`map_id`,`status`,`weight`) VALUES ( '1', 'Фон', '', '0000-00-00 00:00:00', '2017-01-21 14:14:57', '1', 'enabled', '1' );
INSERT INTO `map_layers`(`id`,`name`,`description`,`created_at`,`update_at`,`map_id`,`status`,`weight`) VALUES ( '2', 'Розетки', '', '0000-00-00 00:00:00', '2017-01-21 14:21:12', '1', 'enabled', '0' );
-- ---------------------------------------------------------


-- Dump data of "map_texts" --------------------------------
INSERT INTO `map_texts`(`id`,`Text`,`Style`) VALUES ( '13', 'Спальня', '{"font-size":"34px","color": "gray"}' );
-- ---------------------------------------------------------


-- Dump data of "maps" -------------------------------------
INSERT INTO `maps`(`id`,`created_at`,`update_at`,`name`,`description`,`options`) VALUES ( '1', '0000-00-00 00:00:00', '2017-01-21 14:18:54', 'Дом', '', '{"zoom":1}' );
-- ---------------------------------------------------------


-- Dump data of "migrations" -------------------------------
INSERT INTO `migrations`(`id_migration`,`name`,`created_at`,`statements`,`rollback_statements`,`status`) VALUES ( '1', 'Nodes_20170121_004649', '2017-01-21 19:02:59', 'CREATE TABLE nodes (
	id Int( 255 ) AUTO_INCREMENT NOT NULL,
	ip VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
	port Int( 255 ) NOT NULL,
	name VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
	description Text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
	created_at DateTime NOT NULL,
	update_at DateTime NOT NULL,
	status Enum( \'enabled\', \'disabled\' ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT \'enabled\',
	PRIMARY KEY ( id ),
	CONSTRAINT node UNIQUE( port, ip ),
	CONSTRAINT unique_id UNIQUE( id ) )
	CHARACTER SET = utf8
	COLLATE = utf8_general_ci
	ENGINE = InnoDB;', NULL, 'update' );
INSERT INTO `migrations`(`id_migration`,`name`,`created_at`,`statements`,`rollback_statements`,`status`) VALUES ( '2', 'Workflows_20170121_005244', '2017-01-21 19:02:59', 'CREATE TABLE workflows (
	id Int( 255 ) AUTO_INCREMENT NOT NULL,
	name VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
	status Enum( \'enabled\', \'disabled\' ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT \'enabled\',
	created_at DateTime NOT NULL,
	update_at DateTime NOT NULL,
	description Text CHARACTER SET utf8 COLLATE utf8_general_ci NULL,
	PRIMARY KEY ( id ),
	CONSTRAINT unique_name UNIQUE( name ) )
	CHARACTER SET = utf8
	COLLATE = utf8_general_ci
	ENGINE = InnoDB;', NULL, 'update' );
INSERT INTO `migrations`(`id_migration`,`name`,`created_at`,`statements`,`rollback_statements`,`status`) VALUES ( '3', 'Flows_20170121_005734', '2017-01-21 19:02:59', 'CREATE TABLE flows (
	id Int( 32 ) AUTO_INCREMENT NOT NULL,
	name VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
	status Enum( \'enabled\', \'disabled\' ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT \'enabled\',
	created_at DateTime NOT NULL,
	update_at DateTime NOT NULL,
	workflow_id Int( 11 ) NOT NULL,
	description Text CHARACTER SET utf8 COLLATE utf8_general_ci NULL,
	PRIMARY KEY ( id ),
	CONSTRAINT unique_id UNIQUE( id ) )
	CHARACTER SET = utf8
	COLLATE = utf8_general_ci
	ENGINE = InnoDB;', NULL, 'update' );
INSERT INTO `migrations`(`id_migration`,`name`,`created_at`,`statements`,`rollback_statements`,`status`) VALUES ( '4', 'FlowElements_20170121_010456', '2017-01-21 19:03:00', 'CREATE TABLE flow_elements (
	uuid VarChar( 254 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
	name VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
	graph_settings Text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
	flow_id Int( 32 ) NOT NULL,
	created_at DateTime NOT NULL,
	update_at DateTime NOT NULL,
	prototype_type VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT \'default\',
	status VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NULL,
	description Text CHARACTER SET utf8 COLLATE utf8_general_ci NULL,
	script_id Int( 11 ) NULL,
	flow_link Int( 32 ) NULL,
	PRIMARY KEY ( uuid ),
	CONSTRAINT id UNIQUE( uuid ),
	CONSTRAINT unique_id UNIQUE( uuid ) )
	CHARACTER SET = utf8
	COLLATE = utf8_general_ci
	ENGINE = InnoDB;', NULL, 'update' );
INSERT INTO `migrations`(`id_migration`,`name`,`created_at`,`statements`,`rollback_statements`,`status`) VALUES ( '5', 'Connections_20170121_011406', '2017-01-21 19:03:00', 'CREATE TABLE connections (
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
	ENGINE = InnoDB;', NULL, 'update' );
INSERT INTO `migrations`(`id_migration`,`name`,`created_at`,`statements`,`rollback_statements`,`status`) VALUES ( '6', 'Dashboards_20170121_011938', '2017-01-21 19:03:00', 'CREATE TABLE dashboards (
	id Int( 32 ) AUTO_INCREMENT NOT NULL,
	name VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
	description Text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
	created_at DateTime NOT NULL,
	update_at DateTime NOT NULL,
	widgets Text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
	PRIMARY KEY ( id ),
	CONSTRAINT unique_id UNIQUE( id ) )
	CHARACTER SET = utf8
	COLLATE = utf8_general_ci
	ENGINE = InnoDB;', NULL, 'update' );
INSERT INTO `migrations`(`id_migration`,`name`,`created_at`,`statements`,`rollback_statements`,`status`) VALUES ( '7', 'DeviceActions_20170121_012048', '2017-01-21 19:03:00', 'CREATE TABLE device_actions (
	id Int( 11 ) AUTO_INCREMENT NOT NULL,
	device_id Int( 11 ) NULL,
	name VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
	description VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NULL,
	created_at DateTime NOT NULL,
	update_at DateTime NOT NULL,
	script_id Int( 32 ) NULL,
	PRIMARY KEY ( id ),
	CONSTRAINT unique_id UNIQUE( id ) )
	CHARACTER SET = utf8
	COLLATE = utf8_general_ci
	ENGINE = InnoDB;', NULL, 'update' );
INSERT INTO `migrations`(`id_migration`,`name`,`created_at`,`statements`,`rollback_statements`,`status`) VALUES ( '8', 'DeviceStates_20170121_012150', '2017-01-21 19:03:00', 'CREATE TABLE device_states (
	id Int( 11 ) AUTO_INCREMENT NOT NULL,
	system_name VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
	description Text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
	created_at DateTime NOT NULL,
	update_at DateTime NOT NULL,
	device_id Int( 11 ) NOT NULL,
	PRIMARY KEY ( id ),
	CONSTRAINT unique_1 UNIQUE( device_id, system_name ) )
	CHARACTER SET = utf8
	COLLATE = utf8_general_ci
	ENGINE = InnoDB;', NULL, 'update' );
INSERT INTO `migrations`(`id_migration`,`name`,`created_at`,`statements`,`rollback_statements`,`status`) VALUES ( '9', 'Devices_20170121_012242', '2017-01-21 19:03:00', 'CREATE TABLE devices (
	id Int( 11 ) AUTO_INCREMENT NOT NULL,
	name VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
	description VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
	device_id Int( 11 ) NULL,
	created_at DateTime NOT NULL,
	update_at DateTime NOT NULL,
	node_id Int( 11 ) NULL,
	baud Int( 11 ) NULL,
	tty VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NULL,
	stop_bite Int( 8 ) NULL,
	timeout Int( 11 ) NULL,
	address Int( 11 ) NULL,
	status VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT \'enabled\',
	sleep Int( 32 ) NOT NULL DEFAULT \'0\',
	PRIMARY KEY ( id ),
	CONSTRAINT unique3 UNIQUE( name, device_id ),
	CONSTRAINT unique2 UNIQUE( node_id, address ),
	CONSTRAINT unique1 UNIQUE( device_id, address ) )
	CHARACTER SET = utf8
	COLLATE = utf8_general_ci
	ENGINE = InnoDB;', NULL, 'update' );
INSERT INTO `migrations`(`id_migration`,`name`,`created_at`,`statements`,`rollback_statements`,`status`) VALUES ( '10', 'EmailTemplates_20170121_012357', '2017-01-21 19:03:00', 'CREATE TABLE email_templates (
	name VarChar( 64 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT \'Название\',
	description VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NULL COMMENT \'Описание\',
	content LongText CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT \'Содержимое\',
	status VarChar( 64 ) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT \'active\' COMMENT \'active, unactive\',
	type VarChar( 64 ) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT \'item\' COMMENT \'item, template\',
	parent VarChar( 64 ) CHARACTER SET utf8 COLLATE utf8_general_ci NULL,
	created_at DateTime NOT NULL,
	updated_at DateTime NOT NULL,
	PRIMARY KEY ( name ),
	CONSTRAINT name UNIQUE( name ) )
	CHARACTER SET = utf8
	COLLATE = utf8_general_ci
	ENGINE = InnoDB;', NULL, 'update' );
INSERT INTO `migrations`(`id_migration`,`name`,`created_at`,`statements`,`rollback_statements`,`status`) VALUES ( '11', 'Images_20170121_012702', '2017-01-21 19:03:00', 'CREATE TABLE images (
	id Int( 32 ) AUTO_INCREMENT NOT NULL,
	thumb VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
	image VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
	mime_type VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
	title VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
	size Int( 11 ) NOT NULL,
	name VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
	created_at DateTime NOT NULL,
	PRIMARY KEY ( id ),
	CONSTRAINT id UNIQUE( id ) )
	CHARACTER SET = utf8
	COLLATE = utf8_general_ci
	ENGINE = InnoDB;', NULL, 'update' );
INSERT INTO `migrations`(`id_migration`,`name`,`created_at`,`statements`,`rollback_statements`,`status`) VALUES ( '12', 'Logs_20170121_012936', '2017-01-21 19:03:00', 'CREATE TABLE logs (
	id Int( 32 ) AUTO_INCREMENT NOT NULL,
	body Text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
	created_at DateTime NOT NULL,
	level Enum( \'Emergency\', \'Alert\', \'Critical\', \'Error\', \'Warning\', \'Notice\', \'Info\', \'Debug\' ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT \'Info\',
	PRIMARY KEY ( id ),
	CONSTRAINT unique_id UNIQUE( id ) )
	CHARACTER SET = utf8
	COLLATE = utf8_general_ci
	ENGINE = InnoDB;', NULL, 'update' );
INSERT INTO `migrations`(`id_migration`,`name`,`created_at`,`statements`,`rollback_statements`,`status`) VALUES ( '13', 'MapDeviceActions_20170121_013100', '2017-01-21 19:03:00', 'CREATE TABLE map_device_actions (
	id Int( 32 ) AUTO_INCREMENT NOT NULL,
	device_action_id Int( 32 ) NOT NULL,
	map_device_id Int( 2 ) NOT NULL,
	image_id Int( 32 ) NULL,
	type VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
	created_at DateTime NOT NULL,
	update_at DateTime NOT NULL,
	PRIMARY KEY ( id ),
	CONSTRAINT unique_id UNIQUE( id ) )
	CHARACTER SET = utf8
	COLLATE = utf8_general_ci
	ENGINE = InnoDB;', NULL, 'update' );
INSERT INTO `migrations`(`id_migration`,`name`,`created_at`,`statements`,`rollback_statements`,`status`) VALUES ( '14', 'MapDeviceStates_20170121_013247', '2017-01-21 19:03:00', 'CREATE TABLE map_device_states (
	id Int( 32 ) AUTO_INCREMENT NOT NULL,
	device_state_id Int( 32 ) NOT NULL,
	map_device_id Int( 32 ) NOT NULL,
	image_id Int( 32 ) NULL,
	created_at DateTime NOT NULL,
	update_at DateTime NOT NULL,
	style Text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
	PRIMARY KEY ( id ) )
	CHARACTER SET = utf8
	COLLATE = utf8_general_ci
	ENGINE = InnoDB;', NULL, 'update' );
INSERT INTO `migrations`(`id_migration`,`name`,`created_at`,`statements`,`rollback_statements`,`status`) VALUES ( '15', 'MapDevices_20170121_013341', '2017-01-21 19:03:01', 'CREATE TABLE map_devices (
	id Int( 32 ) AUTO_INCREMENT NOT NULL,
	created_at DateTime NOT NULL,
	update_at DateTime NOT NULL,
	device_id Int( 32 ) NOT NULL,
	image_id Int( 32 ) NULL,
	PRIMARY KEY ( id ),
	CONSTRAINT unique_id UNIQUE( device_id ) )
	CHARACTER SET = utf8
	COLLATE = utf8_general_ci
	ENGINE = InnoDB;', NULL, 'update' );
INSERT INTO `migrations`(`id_migration`,`name`,`created_at`,`statements`,`rollback_statements`,`status`) VALUES ( '16', 'MapElements_20170121_013436', '2017-01-21 19:03:01', 'CREATE TABLE map_elements (
	id Int( 32 ) AUTO_INCREMENT NOT NULL,
	name VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
	description Text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
	created_at DateTime NOT NULL,
	update_at DateTime NOT NULL,
	prototype_type VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
	layer_id Int( 32 ) NOT NULL,
	map_id Int( 32 ) NOT NULL,
	graph_settings Text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
	status VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
	weight Int( 11 ) NOT NULL DEFAULT \'0\',
	prototype_id Int( 32 ) NOT NULL,
	PRIMARY KEY ( id ),
	CONSTRAINT unique_id UNIQUE( id ) )
	CHARACTER SET = utf8
	COLLATE = utf8_general_ci
	ENGINE = InnoDB;', NULL, 'update' );
INSERT INTO `migrations`(`id_migration`,`name`,`created_at`,`statements`,`rollback_statements`,`status`) VALUES ( '17', 'MapImages_20170121_013529', '2017-01-21 19:03:01', 'CREATE TABLE map_images (
	id Int( 32 ) AUTO_INCREMENT NOT NULL,
	image_id Int( 32 ) NOT NULL,
	style Text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
	PRIMARY KEY ( id ),
	CONSTRAINT unique_id UNIQUE( id ) )
	CHARACTER SET = utf8
	COLLATE = utf8_general_ci
	ENGINE = InnoDB;', NULL, 'update' );
INSERT INTO `migrations`(`id_migration`,`name`,`created_at`,`statements`,`rollback_statements`,`status`) VALUES ( '18', 'MapLayers_20170121_013637', '2017-01-21 19:03:01', 'CREATE TABLE map_layers (
	id Int( 11 ) AUTO_INCREMENT NOT NULL,
	name VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
	description Text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
	created_at DateTime NOT NULL,
	update_at DateTime NOT NULL,
	map_id Int( 11 ) NOT NULL,
	status VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
	weight Int( 11 ) NOT NULL DEFAULT \'0\',
	PRIMARY KEY ( id ),
	CONSTRAINT unique_id UNIQUE( id ) )
	CHARACTER SET = utf8
	COLLATE = utf8_general_ci
	ENGINE = InnoDB;', NULL, 'update' );
INSERT INTO `migrations`(`id_migration`,`name`,`created_at`,`statements`,`rollback_statements`,`status`) VALUES ( '19', 'MapTexts_20170121_013724', '2017-01-21 19:03:01', 'CREATE TABLE map_texts (
	id Int( 11 ) AUTO_INCREMENT NOT NULL,
	Text Text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
	Style Text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
	PRIMARY KEY ( id ),
	CONSTRAINT unique_id UNIQUE( id ) )
	CHARACTER SET = utf8
	COLLATE = utf8_general_ci
	ENGINE = InnoDB;', NULL, 'update' );
INSERT INTO `migrations`(`id_migration`,`name`,`created_at`,`statements`,`rollback_statements`,`status`) VALUES ( '20', 'Maps_20170121_013814', '2017-01-21 19:03:01', 'CREATE TABLE maps (
	id Int( 11 ) AUTO_INCREMENT NOT NULL,
	created_at DateTime NOT NULL,
	update_at DateTime NOT NULL,
	name Char( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
	description Text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
	options Text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
	PRIMARY KEY ( id ),
	CONSTRAINT unique_id UNIQUE( id ) )
	CHARACTER SET = utf8
	COLLATE = utf8_general_ci
	ENGINE = InnoDB;', NULL, 'update' );
INSERT INTO `migrations`(`id_migration`,`name`,`created_at`,`statements`,`rollback_statements`,`status`) VALUES ( '21', 'Permissions_20170121_013906', '2017-01-21 19:03:01', 'CREATE TABLE permissions (
	id Int( 32 ) AUTO_INCREMENT NOT NULL,
	role_name VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
	package_name VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
	level_name VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
	PRIMARY KEY ( id ) )
	CHARACTER SET = utf8
	COLLATE = utf8_general_ci
	ENGINE = InnoDB;', NULL, 'update' );
INSERT INTO `migrations`(`id_migration`,`name`,`created_at`,`statements`,`rollback_statements`,`status`) VALUES ( '22', 'Roles_20170121_013956', '2017-01-21 19:03:01', 'CREATE TABLE roles (
	name VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
	description VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
	parent VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NULL,
	created_at DateTime NOT NULL,
	update_at DateTime NOT NULL,
	PRIMARY KEY ( name ),
	CONSTRAINT unique_name UNIQUE( name ) )
	CHARACTER SET = utf8
	COLLATE = utf8_general_ci
	ENGINE = InnoDB;', NULL, 'update' );
INSERT INTO `migrations`(`id_migration`,`name`,`created_at`,`statements`,`rollback_statements`,`status`) VALUES ( '23', 'Scripts_20170121_014057', '2017-01-21 19:03:01', 'CREATE TABLE scripts (
	id Int( 32 ) AUTO_INCREMENT NOT NULL,
	lang Enum( \'lua\', \'coffeescript\', \'javascript\' ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT \'lua\',
	name VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
	source Text CHARACTER SET utf8 COLLATE utf8_general_ci NULL,
	created_at DateTime NOT NULL,
	update_at DateTime NOT NULL,
	description Text CHARACTER SET utf8 COLLATE utf8_general_ci NULL,
	compiled Text CHARACTER SET utf8 COLLATE utf8_general_ci NULL,
	PRIMARY KEY ( id ),
	CONSTRAINT unique_id UNIQUE( id ),
	CONSTRAINT unique_name UNIQUE( name ) )
	CHARACTER SET = utf8
	COLLATE = utf8_general_ci
	ENGINE = InnoDB;', NULL, 'update' );
INSERT INTO `migrations`(`id_migration`,`name`,`created_at`,`statements`,`rollback_statements`,`status`) VALUES ( '24', 'UserMetas_20170121_014145', '2017-01-21 19:03:01', 'CREATE TABLE `user_metas` ( `id` Int( 32 ) AUTO_INCREMENT NOT NULL,`user_id` Int( 32 ) NOT NULL,`key` VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,`value` VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,PRIMARY KEY ( `id` ),CONSTRAINT `unique_id` UNIQUE( `id` ) ) CHARACTER SET = utf8 COLLATE = utf8_general_ci ENGINE = InnoDB AUTO_INCREMENT = 57;', NULL, 'update' );
INSERT INTO `migrations`(`id_migration`,`name`,`created_at`,`statements`,`rollback_statements`,`status`) VALUES ( '25', 'Users_20170121_014234', '2017-01-21 19:03:01', 'CREATE TABLE users (
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
	ENGINE = InnoDB;', NULL, 'update' );
INSERT INTO `migrations`(`id_migration`,`name`,`created_at`,`statements`,`rollback_statements`,`status`) VALUES ( '26', 'Workers_20170121_014331', '2017-01-21 19:03:01', 'CREATE TABLE workers (
	id Int( 11 ) AUTO_INCREMENT NOT NULL,
	flow_id Int( 11 ) NOT NULL,
	workflow_id Int( 11 ) NOT NULL,
	name VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
	device_action_id Int( 11 ) NOT NULL,
	time VarChar( 254 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
	created_at DateTime NOT NULL,
	update_at DateTime NOT NULL,
	status Enum( \'enabled\', \'disabled\' ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT \'enabled\',
	PRIMARY KEY ( id ) )
	CHARACTER SET = utf8
	COLLATE = utf8_general_ci
	ENGINE = InnoDB;', NULL, 'update' );
INSERT INTO `migrations`(`id_migration`,`name`,`created_at`,`statements`,`rollback_statements`,`status`) VALUES ( '27', 'Index_20170121_021812', '2017-01-21 19:03:07', 'CREATE INDEX `lnk_connections_flows` USING BTREE ON `connections`( `flow_id` );; CREATE INDEX `lnk_connections_flow_elements` USING BTREE ON `connections`( `element_from` );; CREATE INDEX `lnk_connections_flow_elements_2` USING BTREE ON `connections`( `element_to` );; CREATE INDEX `lnk_device_actions_devices` USING BTREE ON `device_actions`( `device_id` );; CREATE INDEX `lnk_scripts_device_actions` USING BTREE ON `device_actions`( `script_id` );; CREATE INDEX `index_address` USING BTREE ON `devices`( `address` );; CREATE INDEX `lnk_email_template_email_template` USING BTREE ON `email_templates`( `parent` );; CREATE INDEX `lnk_flows_flow_elements` USING BTREE ON `flow_elements`( `flow_link` );; CREATE INDEX `lnk_flow_elements_flows` USING BTREE ON `flow_elements`( `flow_id` );; CREATE INDEX `lnk_scripts_flow_elements` USING BTREE ON `flow_elements`( `script_id` );; CREATE INDEX `lnk_flows_workflows` USING BTREE ON `flows`( `workflow_id` );; CREATE INDEX `lnk_device_actions_map_device_actions` USING BTREE ON `map_device_actions`( `device_action_id` );; CREATE INDEX `lnk_images_map_device_actions` USING BTREE ON `map_device_actions`( `image_id` );; CREATE INDEX `lnk_map_devices_map_device_actions` USING BTREE ON `map_device_actions`( `map_device_id` );; CREATE INDEX `lnk_device_states_map_device_states` USING BTREE ON `map_device_states`( `device_state_id` );; CREATE INDEX `lnk_images_map_device_states` USING BTREE ON `map_device_states`( `image_id` );; CREATE INDEX `lnk_map_devices_map_device_states` USING BTREE ON `map_device_states`( `map_device_id` );; CREATE INDEX `lnk_images_map_devices` USING BTREE ON `map_devices`( `image_id` );; CREATE INDEX `lnk_maps_map_elements` USING BTREE ON `map_elements`( `map_id` );; CREATE INDEX `lnk_map_layers_map_elements` USING BTREE ON `map_elements`( `layer_id` );; CREATE INDEX `lnk_images_map_images` USING BTREE ON `map_images`( `image_id` );; CREATE INDEX `lnk_maps_map_layers` USING BTREE ON `map_layers`( `map_id` );; CREATE INDEX `lnk_roles_permissions` USING BTREE ON `permissions`( `role_name` );; CREATE INDEX `lnk_roles_roles` USING BTREE ON `roles`( `parent` );; CREATE INDEX `lnk_users_users_meta` USING BTREE ON `user_metas`( `user_id` );; CREATE INDEX `lnk_images_users` USING BTREE ON `users`( `image_id` );; CREATE INDEX `lnk_roles_users` USING BTREE ON `users`( `role_name` );; CREATE INDEX `lnk_users_users` USING BTREE ON `users`( `user_id` );; CREATE INDEX `lnk_workers_device_actions` USING BTREE ON `workers`( `device_action_id` );; CREATE INDEX `lnk_workers_flows` USING BTREE ON `workers`( `flow_id` );; CREATE INDEX `lnk_workers_workflows` USING BTREE ON `workers`( `workflow_id` );', NULL, 'update' );
INSERT INTO `migrations`(`id_migration`,`name`,`created_at`,`statements`,`rollback_statements`,`status`) VALUES ( '28', 'Links_20170121_023351', '2017-01-21 19:03:17', 'ALTER TABLE `connections` ADD CONSTRAINT `lnk_connections_flows` FOREIGN KEY ( `flow_id` ) REFERENCES `flows`( `id` ) ON DELETE Cascade ON UPDATE Restrict;; ALTER TABLE `connections` ADD CONSTRAINT `lnk_connections_flow_elements` FOREIGN KEY ( `element_from` ) REFERENCES `flow_elements`( `uuid` ) ON DELETE Cascade ON UPDATE Cascade;; ALTER TABLE `connections` ADD CONSTRAINT `lnk_connections_flow_elements_2` FOREIGN KEY ( `element_to` ) REFERENCES `flow_elements`( `uuid` ) ON DELETE Cascade ON UPDATE Cascade;; ALTER TABLE `device_actions` ADD CONSTRAINT `lnk_device_actions_devices` FOREIGN KEY ( `device_id` ) REFERENCES `devices`( `id` ) ON DELETE Cascade ON UPDATE Restrict;; ALTER TABLE `device_actions` ADD CONSTRAINT `lnk_scripts_device_actions` FOREIGN KEY ( `script_id` )REFERENCES `scripts`( `id` )ON DELETE Restrict ON UPDATE Cascade;; ALTER TABLE `device_states` ADD CONSTRAINT `lnk_devices_device_states` FOREIGN KEY ( `device_id` ) REFERENCES `devices`( `id` ) ON DELETE Cascade ON UPDATE Cascade;; ALTER TABLE `devices` ADD CONSTRAINT `lnk_devices_devices` FOREIGN KEY ( `device_id` ) REFERENCES `devices`( `id` ) ON DELETE Restrict ON UPDATE Cascade;; ALTER TABLE `devices` ADD CONSTRAINT `lnk_devices_nodes` FOREIGN KEY ( `node_id` ) REFERENCES `nodes`( `id` ) ON DELETE Cascade ON UPDATE Restrict;; ALTER TABLE `flow_elements` ADD CONSTRAINT `lnk_flows_flow_elements` FOREIGN KEY ( `flow_link` ) REFERENCES `flows`( `id` ) ON DELETE Restrict ON UPDATE No Action;; ALTER TABLE `flow_elements`ADD CONSTRAINT `lnk_flow_elements_flows` FOREIGN KEY ( `flow_id` ) REFERENCES `flows`( `id` ) ON DELETE Cascade ON UPDATE Restrict;; ALTER TABLE `flow_elements`ADD CONSTRAINT `lnk_scripts_flow_elements` FOREIGN KEY ( `script_id` )REFERENCES `scripts`( `id` )ON DELETE Restrict ON UPDATE Cascade;; ALTER TABLE `flows` ADD CONSTRAINT `lnk_flows_workflows` FOREIGN KEY ( `workflow_id` ) REFERENCES `workflows`( `id` ) ON DELETE Restrict ON UPDATE Restrict;; ALTER TABLE `map_device_actions` ADD CONSTRAINT `lnk_device_actions_map_device_actions` FOREIGN KEY ( `device_action_id` ) REFERENCES `device_actions`( `id` ) ON DELETE Cascade ON UPDATE Cascade;; ALTER TABLE `map_device_actions` ADD CONSTRAINT `lnk_images_map_device_actions` FOREIGN KEY ( `image_id` ) REFERENCES `images`( `id` ) ON DELETE Restrict ON UPDATE Cascade;; ALTER TABLE `map_device_actions` ADD CONSTRAINT `lnk_map_devices_map_device_actions` FOREIGN KEY ( `map_device_id` ) REFERENCES `map_devices`( `id` ) ON DELETE Cascade ON UPDATE Cascade;; ALTER TABLE `map_device_states` ADD CONSTRAINT `lnk_device_states_map_device_states` FOREIGN KEY ( `device_state_id` ) REFERENCES `device_states`( `id` ) ON DELETE Restrict ON UPDATE Cascade;; ALTER TABLE `map_device_states` ADD CONSTRAINT `lnk_images_map_device_states` FOREIGN KEY ( `image_id` ) REFERENCES `images`( `id` ) ON DELETE Restrict ON UPDATE Cascade;; ALTER TABLE `map_device_states` ADD CONSTRAINT `lnk_map_devices_map_device_states` FOREIGN KEY ( `map_device_id` ) REFERENCES `map_devices`( `id` ) ON DELETE Cascade ON UPDATE Cascade;; ALTER TABLE `map_devices` ADD CONSTRAINT `lnk_devices_map_devices` FOREIGN KEY ( `device_id` ) REFERENCES `devices`( `id` ) ON DELETE Restrict ON UPDATE Cascade;; ALTER TABLE `map_devices` ADD CONSTRAINT `lnk_images_map_devices` FOREIGN KEY ( `image_id` ) REFERENCES `images`( `id` ) ON DELETE Restrict ON UPDATE Cascade;; ALTER TABLE `map_elements` ADD CONSTRAINT `lnk_maps_map_elements` FOREIGN KEY ( `map_id` ) REFERENCES `maps`( `id` ) ON DELETE Restrict ON UPDATE Cascade;; ALTER TABLE `map_elements` ADD CONSTRAINT `lnk_map_layers_map_elements` FOREIGN KEY ( `layer_id` ) REFERENCES `map_layers`( `id` ) ON DELETE Restrict ON UPDATE Cascade;; ALTER TABLE `map_images` ADD CONSTRAINT `lnk_images_map_images` FOREIGN KEY ( `image_id` ) REFERENCES `images`( `id` ) ON DELETE Restrict ON UPDATE Cascade;; ALTER TABLE `map_layers` ADD CONSTRAINT `lnk_maps_map_layers` FOREIGN KEY ( `map_id` ) REFERENCES `maps`( `id` ) ON DELETE Restrict ON UPDATE Cascade;; ALTER TABLE `permissions` ADD CONSTRAINT `lnk_roles_permissions` FOREIGN KEY ( `role_name` ) REFERENCES `roles`( `name` ) ON DELETE Cascade ON UPDATE Cascade;; ALTER TABLE `roles` ADD CONSTRAINT `lnk_roles_roles` FOREIGN KEY ( `parent` ) REFERENCES `roles`( `name` ) ON DELETE Restrict ON UPDATE Cascade;; ALTER TABLE `user_metas` ADD CONSTRAINT `lnk_users_users_meta` FOREIGN KEY ( `user_id` ) REFERENCES `users`( `id` ) ON DELETE Cascade ON UPDATE Cascade;; ALTER TABLE `users` ADD CONSTRAINT `lnk_images_users` FOREIGN KEY ( `image_id` ) REFERENCES `images`( `id` ) ON DELETE Restrict ON UPDATE Cascade;; ALTER TABLE `users` ADD CONSTRAINT `lnk_roles_users` FOREIGN KEY ( `role_name` ) REFERENCES `roles`( `name` ) ON DELETE Restrict ON UPDATE Cascade;; ALTER TABLE `users` ADD CONSTRAINT `lnk_users_users` FOREIGN KEY ( `user_id` ) REFERENCES `users`( `id` ) ON DELETE Restrict ON UPDATE Cascade;; ALTER TABLE `workers` ADD CONSTRAINT `lnk_workers_device_actions` FOREIGN KEY ( `device_action_id` ) REFERENCES `device_actions`( `id` ) ON DELETE Restrict ON UPDATE Restrict;; ALTER TABLE `workers` ADD CONSTRAINT `lnk_workers_flows` FOREIGN KEY ( `flow_id` ) REFERENCES `flows`( `id` ) ON DELETE Cascade ON UPDATE Restrict;; ALTER TABLE `workers` ADD CONSTRAINT `lnk_workers_workflows` FOREIGN KEY ( `workflow_id` ) REFERENCES `workflows`( `id` ) ON DELETE Restrict ON UPDATE Restrict;', NULL, 'update' );
INSERT INTO `migrations`(`id_migration`,`name`,`created_at`,`statements`,`rollback_statements`,`status`) VALUES ( '29', 'AddUsers_20170121_175235', '2017-01-21 19:03:19', 'INSERT INTO `roles`(`name`,`description`,`parent`,`created_at`,`update_at`) VALUES ( \'demo\', \'\', NULL, \'2017-01-15 05:17:09\', \'2017-01-15 05:17:09\' );; INSERT INTO `roles`(`name`,`description`,`parent`,`created_at`,`update_at`) VALUES ( \'user\', \'\', \'demo\', \'2017-01-15 05:17:36\', \'2017-01-19 19:39:54\' );; INSERT INTO `roles`(`name`,`description`,`parent`,`created_at`,`update_at`) VALUES ( \'admin\', \'\', \'user\', \'2017-01-15 05:20:58\', \'2017-01-15 05:21:29\' );; INSERT INTO `users`(`id`,`nickname`,`first_name`,`last_name`,`encrypted_password`,`email`,`history`,`status`,`reset_password_token`,`authentication_token`,`image_id`,`sign_in_count`,`current_sign_in_ip`,`last_sign_in_ip`,`user_id`,`role_name`,`reset_password_sent_at`,`current_sign_in_at`,`last_sign_in_at`,`created_at`,`update_at`,`deleted`) VALUES ( \'1\', \'admin\', \'\', \'\', \'f6fdffe48c908deb0f4c3bd36c032e72\', \'admin@e154.ru\', \'[]\', \'active\', \'\', \'xlzEaHNBbn80OmTfWd1z18XpNUlZikdb4z5fo5YAxlNv3CfWxs\', NULL, \'111\', \'127.0.0.1\', \'127.0.0.1\', NULL, \'admin\', NULL, \'2017-01-21 11:11:26\', \'2017-01-21 10:54:03\', \'2017-01-15 05:25:07\', \'2017-01-21 11:11:26\', NULL );; INSERT INTO `users`(`id`,`nickname`,`first_name`,`last_name`,`encrypted_password`,`email`,`history`,`status`,`reset_password_token`,`authentication_token`,`image_id`,`sign_in_count`,`current_sign_in_ip`,`last_sign_in_ip`,`user_id`,`role_name`,`reset_password_sent_at`,`current_sign_in_at`,`last_sign_in_at`,`created_at`,`update_at`,`deleted`) VALUES ( \'2\', \'demo\', \'\', \'\', \'c514c91e4ed341f263e458d44b3bb0a7\', \'demo@e154.ru\', \'[]\', \'active\', \'\', \'5SLTHOzN1hWw6jhgEw0y9JbtwdBIK5mgW3DLt5FYy23zNkVnvW\', NULL, \'8\', \'127.0.0.1\', \'127.0.0.1\', NULL, \'demo\', NULL, \'2017-01-21 11:11:43\', \'2017-01-20 17:28:23\', \'2017-01-18 17:13:28\', \'2017-01-21 11:11:43\', NULL );; INSERT INTO `user_metas`(`id`,`user_id`,`key`,`value`) VALUES ( \'49\', \'1\', \'phone1\', \'\' );; INSERT INTO `user_metas`(`id`,`user_id`,`key`,`value`) VALUES ( \'50\', \'1\', \'phone2\', \'\' );; INSERT INTO `user_metas`(`id`,`user_id`,`key`,`value`) VALUES ( \'51\', \'1\', \'phone3\', \'\' );; INSERT INTO `user_metas`(`id`,`user_id`,`key`,`value`) VALUES ( \'52\', \'1\', \'autograph\', \'\' );; INSERT INTO `user_metas`(`id`,`user_id`,`key`,`value`) VALUES ( \'53\', \'2\', \'phone1\', \'\' );; INSERT INTO `user_metas`(`id`,`user_id`,`key`,`value`) VALUES ( \'54\', \'2\', \'phone2\', \'\' );; INSERT INTO `user_metas`(`id`,`user_id`,`key`,`value`) VALUES ( \'55\', \'2\', \'phone3\', \'\' );; INSERT INTO `user_metas`(`id`,`user_id`,`key`,`value`) VALUES ( \'56\', \'2\', \'autograph\', \'\' );; INSERT INTO `permissions`(`id`,`role_name`,`package_name`,`level_name`) VALUES ( \'81\', \'admin\', \'ws\', \'read\' );; INSERT INTO `permissions`(`id`,`role_name`,`package_name`,`level_name`) VALUES ( \'82\', \'admin\', \'workflow\', \'create\' );; INSERT INTO `permissions`(`id`,`role_name`,`package_name`,`level_name`) VALUES ( \'83\', \'admin\', \'workflow\', \'delete\' );; INSERT INTO `permissions`(`id`,`role_name`,`package_name`,`level_name`) VALUES ( \'84\', \'admin\', \'workflow\', \'read\' );; INSERT INTO `permissions`(`id`,`role_name`,`package_name`,`level_name`) VALUES ( \'86\', \'demo\', \'ws\', \'read\' );; INSERT INTO `permissions`(`id`,`role_name`,`package_name`,`level_name`) VALUES ( \'88\', \'demo\', \'map\', \'read_map\' );; INSERT INTO `permissions`(`id`,`role_name`,`package_name`,`level_name`) VALUES ( \'89\', \'demo\', \'device\', \'read\' );; INSERT INTO `permissions`(`id`,`role_name`,`package_name`,`level_name`) VALUES ( \'90\', \'demo\', \'node\', \'read\' );; INSERT INTO `permissions`(`id`,`role_name`,`package_name`,`level_name`) VALUES ( \'91\', \'demo\', \'dashboard\', \'read\' );; INSERT INTO `permissions`(`id`,`role_name`,`package_name`,`level_name`) VALUES ( \'92\', \'demo\', \'notifr\', \'preview_notifr\' );; INSERT INTO `permissions`(`id`,`role_name`,`package_name`,`level_name`) VALUES ( \'93\', \'demo\', \'notifr\', \'read_notifr_item\' );; INSERT INTO `permissions`(`id`,`role_name`,`package_name`,`level_name`) VALUES ( \'94\', \'demo\', \'notifr\', \'read_notifr_template\' );; INSERT INTO `permissions`(`id`,`role_name`,`package_name`,`level_name`) VALUES ( \'95\', \'demo\', \'script\', \'exec_script\' );; INSERT INTO `permissions`(`id`,`role_name`,`package_name`,`level_name`) VALUES ( \'96\', \'demo\', \'script\', \'read\' );; INSERT INTO `permissions`(`id`,`role_name`,`package_name`,`level_name`) VALUES ( \'97\', \'demo\', \'user\', \'read\' );; INSERT INTO `permissions`(`id`,`role_name`,`package_name`,`level_name`) VALUES ( \'98\', \'demo\', \'user\', \'read_role\' );; INSERT INTO `permissions`(`id`,`role_name`,`package_name`,`level_name`) VALUES ( \'99\', \'demo\', \'user\', \'read_role_access_list\' );; INSERT INTO `permissions`(`id`,`role_name`,`package_name`,`level_name`) VALUES ( \'100\', \'demo\', \'worker\', \'read\' );; INSERT INTO `permissions`(`id`,`role_name`,`package_name`,`level_name`) VALUES ( \'101\', \'demo\', \'device_action\', \'read\' );; INSERT INTO `permissions`(`id`,`role_name`,`package_name`,`level_name`) VALUES ( \'102\', \'demo\', \'device_state\', \'read\' );; INSERT INTO `permissions`(`id`,`role_name`,`package_name`,`level_name`) VALUES ( \'103\', \'demo\', \'flow\', \'read\' );; INSERT INTO `permissions`(`id`,`role_name`,`package_name`,`level_name`) VALUES ( \'104\', \'demo\', \'image\', \'read\' );; INSERT INTO `permissions`(`id`,`role_name`,`package_name`,`level_name`) VALUES ( \'105\', \'demo\', \'workflow\', \'read\' );; INSERT INTO `permissions`(`id`,`role_name`,`package_name`,`level_name`) VALUES ( \'106\', \'demo\', \'log\', \'read\' );; INSERT INTO `permissions`(`id`,`role_name`,`package_name`,`level_name`) VALUES ( \'107\', \'demo\', \'map\', \'read_map_element\' );; INSERT INTO `permissions`(`id`,`role_name`,`package_name`,`level_name`) VALUES ( \'108\', \'demo\', \'map\', \'read_map_layer\' );; INSERT INTO `permissions`(`id`,`role_name`,`package_name`,`level_name`) VALUES ( \'109\', \'user\', \'image\', \'upload\' );; INSERT INTO `permissions`(`id`,`role_name`,`package_name`,`level_name`) VALUES ( \'110\', \'user\', \'image\', \'create\' );; INSERT INTO `permissions`(`id`,`role_name`,`package_name`,`level_name`) VALUES ( \'111\', \'user\', \'image\', \'delete\' );; INSERT INTO `permissions`(`id`,`role_name`,`package_name`,`level_name`) VALUES ( \'112\', \'user\', \'image\', \'update\' );; INSERT INTO `permissions`(`id`,`role_name`,`package_name`,`level_name`) VALUES ( \'113\', \'user\', \'flow\', \'create\' );; INSERT INTO `permissions`(`id`,`role_name`,`package_name`,`level_name`) VALUES ( \'114\', \'user\', \'flow\', \'delete\' );; INSERT INTO `permissions`(`id`,`role_name`,`package_name`,`level_name`) VALUES ( \'115\', \'user\', \'flow\', \'update\' );; INSERT INTO `permissions`(`id`,`role_name`,`package_name`,`level_name`) VALUES ( \'116\', \'user\', \'device_state\', \'create\' );; INSERT INTO `permissions`(`id`,`role_name`,`package_name`,`level_name`) VALUES ( \'117\', \'user\', \'device_state\', \'delete\' );; INSERT INTO `permissions`(`id`,`role_name`,`package_name`,`level_name`) VALUES ( \'118\', \'user\', \'device_state\', \'update\' );; INSERT INTO `permissions`(`id`,`role_name`,`package_name`,`level_name`) VALUES ( \'119\', \'user\', \'device_action\', \'delete\' );; INSERT INTO `permissions`(`id`,`role_name`,`package_name`,`level_name`) VALUES ( \'120\', \'user\', \'device_action\', \'create\' );; INSERT INTO `permissions`(`id`,`role_name`,`package_name`,`level_name`) VALUES ( \'121\', \'user\', \'device\', \'update\' );; INSERT INTO `permissions`(`id`,`role_name`,`package_name`,`level_name`) VALUES ( \'122\', \'user\', \'device\', \'create\' );; INSERT INTO `permissions`(`id`,`role_name`,`package_name`,`level_name`) VALUES ( \'123\', \'user\', \'device\', \'delete\' );; INSERT INTO `permissions`(`id`,`role_name`,`package_name`,`level_name`) VALUES ( \'124\', \'user\', \'dashboard\', \'create\' );; INSERT INTO `permissions`(`id`,`role_name`,`package_name`,`level_name`) VALUES ( \'125\', \'user\', \'dashboard\', \'update\' );; INSERT INTO `permissions`(`id`,`role_name`,`package_name`,`level_name`) VALUES ( \'126\', \'user\', \'dashboard\', \'delete\' );; INSERT INTO `permissions`(`id`,`role_name`,`package_name`,`level_name`) VALUES ( \'127\', \'user\', \'script\', \'update\' );; INSERT INTO `permissions`(`id`,`role_name`,`package_name`,`level_name`) VALUES ( \'128\', \'user\', \'script\', \'create\' );; INSERT INTO `permissions`(`id`,`role_name`,`package_name`,`level_name`) VALUES ( \'129\', \'user\', \'script\', \'delete\' );; INSERT INTO `permissions`(`id`,`role_name`,`package_name`,`level_name`) VALUES ( \'130\', \'user\', \'worker\', \'delete\' );; INSERT INTO `permissions`(`id`,`role_name`,`package_name`,`level_name`) VALUES ( \'131\', \'user\', \'worker\', \'create\' );; INSERT INTO `permissions`(`id`,`role_name`,`package_name`,`level_name`) VALUES ( \'132\', \'user\', \'worker\', \'update\' );; INSERT INTO `permissions`(`id`,`role_name`,`package_name`,`level_name`) VALUES ( \'133\', \'user\', \'workflow\', \'create\' );; INSERT INTO `permissions`(`id`,`role_name`,`package_name`,`level_name`) VALUES ( \'134\', \'user\', \'workflow\', \'delete\' );; INSERT INTO `permissions`(`id`,`role_name`,`package_name`,`level_name`) VALUES ( \'135\', \'user\', \'workflow\', \'update\' );', NULL, 'update' );
INSERT INTO `migrations`(`id_migration`,`name`,`created_at`,`statements`,`rollback_statements`,`status`) VALUES ( '30', 'EmailTemplates_20170121_200137', '2017-01-21 20:13:09', 'INSERT INTO email_templates(name,description,content,status,type,parent,created_at,updated_at) VALUES ( \'body\', \'\', \'[body:content]\', \'active\', \'item\', \'message\', \'2014-06-21 21:56:07\', \'2015-04-14 15:24:39\' );; INSERT INTO email_templates(name,description,content,status,type,parent,created_at,updated_at) VALUES ( \'callout\', \'\', \'<table class="row callout">
    <tr>
        <td class="wrapper last">

            <table class="twelve columns">
                <tr>
                    <td class="panel">
                        [callout:content]
                    </td>
                    <td class="expander"></td>
                </tr>
            </table>

        </td>
    </tr>
</table>\', \'active\', \'item\', \'main\', \'2014-06-15 17:45:31\', \'2015-04-12 01:40:03\' );; INSERT INTO email_templates(name,description,content,status,type,parent,created_at,updated_at) VALUES ( \'contacts\', \'\', \'<table class="six columns">
    <tr>
        <td class="last right-text-pad">
        <h3>Контакты:</h3>
            [contacts:content]
        </td>
        <td class="expander"></td>
    </tr>
</table>\', \'active\', \'item\', \'footer\', \'2014-06-15 17:45:20\', \'2015-04-12 01:40:03\' );; INSERT INTO email_templates(name,description,content,status,type,parent,created_at,updated_at) VALUES ( \'facebook\', \'\', \'<table class="tiny-button facebook">
    <tr>
        <td>
            <a href="#">Facebook</a>
        </td>
    </tr>
</table>

<br>\', \'active\', \'item\', \'social\', \'2014-06-15 17:46:59\', \'2015-04-13 01:50:01\' );; INSERT INTO email_templates(name,description,content,status,type,parent,created_at,updated_at) VALUES ( \'footer\', \'\', \'<table class="row footer">
    <tr>

        <td class="wrapper ">

            [social:block]

        </td>
        <td class="wrapper last">

            [contacts:block]

        </td>

    </tr>
</table>\', \'active\', \'item\', \'main\', \'2014-06-15 17:45:07\', \'2015-04-12 01:40:03\' );; INSERT INTO email_templates(name,description,content,status,type,parent,created_at,updated_at) VALUES ( \'google\', \'\', \'<table class="tiny-button google-plus">
    <tr>
        <td>
            <a href="#">Google +</a>
        </td>
    </tr>
</table>

<br>\', \'active\', \'item\', \'social\', \'2014-06-15 17:47:17\', \'2015-04-12 01:40:01\' );; INSERT INTO email_templates(name,description,content,status,type,parent,created_at,updated_at) VALUES ( \'header\', \'\', \'<table class="row header">
    <tr>
        <td class="center" align="center">
            <center>

                <table class="container">
                    <tr>
                        <td class="wrapper last">

                            <table class="twelve columns">
                                <tr>
                                    <td class="six sub-columns">
                                        <img alt="Облачная типография, Календари-домики, Листовки и Флаеры, Карманные календари, Визитки, Буклеты, Пластиковые карты" src="data:image/png;base64,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">
                                    </td>
                                    <td class="six sub-columns last" style="text-align:right; vertical-align:middle;">
                                    </td>
                                    <td class="expander"></td>
                                </tr>
                            </table>

                        </td>
                    </tr>
                </table>

            </center>
        </td>
    </tr>
</table>\', \'active\', \'item\', \'main\', \'2014-06-15 17:44:55\', \'2015-04-16 23:43:22\' );; INSERT INTO email_templates(name,description,content,status,type,parent,created_at,updated_at) VALUES ( \'main\', \'Основной слой\', \'<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Strict//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-strict.dtd">
<html xmlns="http://www.w3.org/1999/xhtml">
<head>
  <meta http-equiv="Content-Type" content="text/html; charset=UTF-8">
  <meta name="viewport" content="width=device-width"/>
  <style type="text/css">

#outlook a {
  padding:0;
}

body{
  width:100% !important;
  min-width: 100%;
  -webkit-text-size-adjust:100%;
  -ms-text-size-adjust:100%;
  margin:0;
  padding:0;
}

.ExternalClass {
  width:100%;
}

.ExternalClass,
.ExternalClass p,
.ExternalClass span,
.ExternalClass font,
.ExternalClass td,
.ExternalClass div {
  line-height: 100%;
}

#backgroundTable {
  margin:0;
  padding:0;
  width:100% !important;
  line-height: 100% !important;
}

img {
  outline:none;
  text-decoration:none;
  -ms-interpolation-mode: bicubic;
  width: auto;
  max-width: 100%;
  float: left;
  clear: both;
  display: block;
}

center {
  width: 100%;
  min-width: 580px;
}

a img {
  border: none;
}

p {
  margin: 0 0 0 10px;
}

table {
  border-spacing: 0;
  border-collapse: collapse;
}

td {
  word-break: break-word;
  -webkit-hyphens: auto;
  -moz-hyphens: auto;
  hyphens: auto;
  border-collapse: collapse !important;
}

table, tr, td {
  padding: 0;
  vertical-align: top;
  text-align: left;
}

hr {
  color: #d9d9d9;
  background-color: #d9d9d9;
  height: 1px;
  border: none;
}

/* Responsive Grid */

table.body {
  height: 100%;
  width: 100%;
}

table.container {
  width: 580px;
  margin: 0 auto;
  text-align: inherit;
}

table.row {
  padding: 0;
  width: 100%;
  position: relative;
}

table.container table.row {
  display: block;
}

td.wrapper {
  padding: 10px 20px 0 0;
  position: relative;
}

table.columns,
table.column {
  margin: 0 auto;
}

table.columns td,
table.column td {
  padding: 0 0 10px;
}

table.columns td.sub-columns,
table.column td.sub-columns,
table.columns td.sub-column,
table.column td.sub-column {
  padding-right: 10px;
}

td.sub-column, td.sub-columns {
  min-width: 0;
}

table.row td.last,
table.container td.last {
  padding-right: 0;
}

table.one { width: 30px; }
table.two { width: 80px; }
table.three { width: 130px; }
table.four { width: 180px; }
table.five { width: 230px; }
table.six { width: 280px; }
table.seven { width: 330px; }
table.eight { width: 380px; }
table.nine { width: 430px; }
table.ten { width: 480px; }
table.eleven { width: 530px; }
table.twelve { width: 580px; }

table.one center { min-width: 30px; }
table.two center { min-width: 80px; }
table.three center { min-width: 130px; }
table.four center { min-width: 180px; }
table.five center { min-width: 230px; }
table.six center { min-width: 280px; }
table.seven center { min-width: 330px; }
table.eight center { min-width: 380px; }
table.nine center { min-width: 430px; }
table.ten center { min-width: 480px; }
table.eleven center { min-width: 530px; }
table.twelve center { min-width: 580px; }

table.one .panel center { min-width: 10px; }
table.two .panel center { min-width: 60px; }
table.three .panel center { min-width: 110px; }
table.four .panel center { min-width: 160px; }
table.five .panel center { min-width: 210px; }
table.six .panel center { min-width: 260px; }
table.seven .panel center { min-width: 310px; }
table.eight .panel center { min-width: 360px; }
table.nine .panel center { min-width: 410px; }
table.ten .panel center { min-width: 460px; }
table.eleven .panel center { min-width: 510px; }
table.twelve .panel center { min-width: 560px; }

.body .columns td.one,
.body .column td.one { width: 8.333333%; }
.body .columns td.two,
.body .column td.two { width: 16.666666%; }
.body .columns td.three,
.body .column td.three { width: 25%; }
.body .columns td.four,
.body .column td.four { width: 33.333333%; }
.body .columns td.five,
.body .column td.five { width: 41.666666%; }
.body .columns td.six,
.body .column td.six { width: 50%; }
.body .columns td.seven,
.body .column td.seven { width: 58.333333%; }
.body .columns td.eight,
.body .column td.eight { width: 66.666666%; }
.body .columns td.nine,
.body .column td.nine { width: 75%; }
.body .columns td.ten,
.body .column td.ten { width: 83.333333%; }
.body .columns td.eleven,
.body .column td.eleven { width: 91.666666%; }
.body .columns td.twelve,
.body .column td.twelve { width: 100%; }

td.offset-by-one { padding-left: 50px; }
td.offset-by-two { padding-left: 100px; }
td.offset-by-three { padding-left: 150px; }
td.offset-by-four { padding-left: 200px; }
td.offset-by-five { padding-left: 250px; }
td.offset-by-six { padding-left: 300px; }
td.offset-by-seven { padding-left: 350px; }
td.offset-by-eight { padding-left: 400px; }
td.offset-by-nine { padding-left: 450px; }
td.offset-by-ten { padding-left: 500px; }
td.offset-by-eleven { padding-left: 550px; }

td.expander {
  visibility: hidden;
  width: 0;
  padding: 0 !important;
}

table.columns .text-pad,
table.column .text-pad {
  padding-left: 10px;
  padding-right: 10px;
}

table.columns .left-text-pad,
table.columns .text-pad-left,
table.column .left-text-pad,
table.column .text-pad-left {
  padding-left: 10px;
}

table.columns .right-text-pad,
table.columns .text-pad-right,
table.column .right-text-pad,
table.column .text-pad-right {
  padding-right: 10px;
}

/* Block Grid */

.block-grid {
  width: 100%;
  max-width: 580px;
}

.block-grid td {
  display: inline-block;
  padding:10px;
}

.two-up td {
  width:270px;
}

.three-up td {
  width:173px;
}

.four-up td {
  width:125px;
}

.five-up td {
  width:96px;
}

.six-up td {
  width:76px;
}

.seven-up td {
  width:62px;
}

.eight-up td {
  width:52px;
}

/* Alignment & Visibility Classes */

table.center, td.center {
  text-align: center;
}

h1.center,
h2.center,
h3.center,
h4.center,
h5.center,
h6.center {
  text-align: center;
}

span.center {
  display: block;
  width: 100%;
  text-align: center;
}

img.center {
  margin: 0 auto;
  float: none;
}

.show-for-small,
.hide-for-desktop {
  display: none;
}

/* Typography */

body, table.body, h1, h2, h3, h4, h5, h6, p, td {
  color: #222222;
  font-family: "Helvetica", "Arial", sans-serif;
  font-weight: normal;
  padding:0;
  margin: 0;
  text-align: left;
  line-height: 1.3;
}

h1, h2, h3, h4, h5, h6 {
  word-break: normal;
}

/*h1 {font-size: 40px;}*/
h1 {font-size: 30px;}
/*h2 {font-size: 36px;}*/
h2 {font-size: 26px;}
h3 {font-size: 32px;}
h4 {font-size: 28px;}
h5 {font-size: 27px;}
h6 {font-size: 20px;}
body, table.body, p, td {font-size: 14px;line-height:19px;}

p.lead, p.lede, p.leed {
  font-size: 18px;
  line-height:21px;
}

p {
  margin-bottom: 10px;
}

small {
  font-size: 10px;
}

a {
  color: #2ba6cb;
  text-decoration: none;
}

a:hover {
  color: #2795b6 !important;
}

a:active {
  color: #2795b6 !important;
}

a:visited {
  color: #2ba6cb !important;
}

h1 a,
h2 a,
h3 a,
h4 a,
h5 a,
h6 a {
  color: #2ba6cb;
}

h1 a:active,
h2 a:active,
h3 a:active,
h4 a:active,
h5 a:active,
h6 a:active {
  color: #2ba6cb !important;
}

h1 a:visited,
h2 a:visited,
h3 a:visited,
h4 a:visited,
h5 a:visited,
h6 a:visited {
  color: #2ba6cb !important;
}

/* Panels */

.panel {
  background: #f2f2f2;
  border: 1px solid #d9d9d9;
  padding: 10px !important;
}

.sub-grid table {
  width: 100%;
}

.sub-grid td.sub-columns {
  padding-bottom: 0;
}

/* Buttons */

table.button,
table.tiny-button,
table.small-button,
table.medium-button,
table.large-button {
  width: 100%;
  overflow: hidden;
}

table.button td,
table.tiny-button td,
table.small-button td,
table.medium-button td,
table.large-button td {
  display: block;
  width: auto !important;
  text-align: center;
  background: #2ba6cb;
  border: 1px solid #2284a1;
  color: #ffffff;
  padding: 8px 0;
}

table.tiny-button td {
  padding: 5px 0 4px;
}

table.small-button td {
  padding: 8px 0 7px;
}

table.medium-button td {
  padding: 12px 0 10px;
}

table.large-button td {
  padding: 21px 0 18px;
}

table.button td a,
table.tiny-button td a,
table.small-button td a,
table.medium-button td a,
table.large-button td a {
  font-weight: bold;
  text-decoration: none;
  font-family: Helvetica, Arial, sans-serif;
  color: #ffffff;
  font-size: 16px;
}

table.tiny-button td a {
  font-size: 12px;
  font-weight: normal;
}

table.small-button td a {
  font-size: 16px;
}

table.medium-button td a {
  font-size: 20px;
}

table.large-button td a {
  font-size: 24px;
}

table.button:hover td,
table.button:visited td,
table.button:active td {
  background: #2795b6 !important;
}

table.button:hover td a,
table.button:visited td a,
table.button:active td a {
  color: #fff !important;
}

table.button:hover td,
table.tiny-button:hover td,
table.small-button:hover td,
table.medium-button:hover td,
table.large-button:hover td {
  background: #2795b6 !important;
}

table.button:hover td a,
table.button:active td a,
table.button td a:visited,
table.tiny-button:hover td a,
table.tiny-button:active td a,
table.tiny-button td a:visited,
table.small-button:hover td a,
table.small-button:active td a,
table.small-button td a:visited,
table.medium-button:hover td a,
table.medium-button:active td a,
table.medium-button td a:visited,
table.large-button:hover td a,
table.large-button:active td a,
table.large-button td a:visited {
  color: #ffffff !important;
}

table.secondary td {
  background: #e9e9e9;
  border-color: #d0d0d0;
  color: #555;
}

table.secondary td a {
  color: #555;
}

table.secondary:hover td {
  background: #d0d0d0 !important;
  color: #555;
}

table.secondary:hover td a,
table.secondary td a:visited,
table.secondary:active td a {
  color: #555 !important;
}

table.success td {
  background: #5da423;
  border-color: #457a1a;
}

table.success:hover td {
  background: #457a1a !important;
}

table.alert td {
  background: #c60f13;
  border-color: #970b0e;
}

table.alert:hover td {
  background: #970b0e !important;
}

table.radius td {
  -webkit-border-radius: 3px;
  -moz-border-radius: 3px;
  border-radius: 3px;
}

table.round td {
  -webkit-border-radius: 500px;
  -moz-border-radius: 500px;
  border-radius: 500px;
}

/* Outlook First */

body.outlook p {
  display: inline !important;
}

/*  Media Queries */

@media only screen and (max-width: 600px) {

  table[class="body"] img {
    width: auto !important;
    height: auto !important;
  }

  table[class="body"] center {
    min-width: 0 !important;
  }

  table[class="body"] .container {
    width: 95% !important;
  }

  table[class="body"] .row {
    width: 100% !important;
    display: block !important;
  }

  table[class="body"] .wrapper {
    display: block !important;
    padding-right: 0 !important;
  }

  table[class="body"] .columns,
  table[class="body"] .column {
    table-layout: fixed !important;
    float: none !important;
    width: 100% !important;
    padding-right: 0 !important;
    padding-left: 0 !important;
    display: block !important;
  }

  table[class="body"] .wrapper.first .columns,
  table[class="body"] .wrapper.first .column {
    display: table !important;
  }

  table[class="body"] table.columns td,
  table[class="body"] table.column td {
    width: 100% !important;
  }

  table[class="body"] .columns td.one,
  table[class="body"] .column td.one { width: 8.333333% !important; }
  table[class="body"] .columns td.two,
  table[class="body"] .column td.two { width: 16.666666% !important; }
  table[class="body"] .columns td.three,
  table[class="body"] .column td.three { width: 25% !important; }
  table[class="body"] .columns td.four,
  table[class="body"] .column td.four { width: 33.333333% !important; }
  table[class="body"] .columns td.five,
  table[class="body"] .column td.five { width: 41.666666% !important; }
  table[class="body"] .columns td.six,
  table[class="body"] .column td.six { width: 50% !important; }
  table[class="body"] .columns td.seven,
  table[class="body"] .column td.seven { width: 58.333333% !important; }
  table[class="body"] .columns td.eight,
  table[class="body"] .column td.eight { width: 66.666666% !important; }
  table[class="body"] .columns td.nine,
  table[class="body"] .column td.nine { width: 75% !important; }
  table[class="body"] .columns td.ten,
  table[class="body"] .column td.ten { width: 83.333333% !important; }
  table[class="body"] .columns td.eleven,
  table[class="body"] .column td.eleven { width: 91.666666% !important; }
  table[class="body"] .columns td.twelve,
  table[class="body"] .column td.twelve { width: 100% !important; }

  table[class="body"] td.offset-by-one,
  table[class="body"] td.offset-by-two,
  table[class="body"] td.offset-by-three,
  table[class="body"] td.offset-by-four,
  table[class="body"] td.offset-by-five,
  table[class="body"] td.offset-by-six,
  table[class="body"] td.offset-by-seven,
  table[class="body"] td.offset-by-eight,
  table[class="body"] td.offset-by-nine,
  table[class="body"] td.offset-by-ten,
  table[class="body"] td.offset-by-eleven {
    padding-left: 0 !important;
  }

  table[class="body"] table.columns td.expander {
    width: 1px !important;
  }

  table[class="body"] .right-text-pad,
  table[class="body"] .text-pad-right {
    padding-left: 10px !important;
  }

  table[class="body"] .left-text-pad,
  table[class="body"] .text-pad-left {
    padding-right: 10px !important;
  }

  table[class="body"] .hide-for-small,
  table[class="body"] .show-for-desktop {
    display: none !important;
  }

  table[class="body"] .show-for-small,
  table[class="body"] .hide-for-desktop {
    display: inherit !important;
  }
}

  </style>
  <style type="text/css">

    table.facebook td {
      background: #3b5998;
      border-color: #2d4473;
    }

    table.facebook:hover td {
      background: #2d4473 !important;
    }

    table.twitter td {
      background: #00acee;
      border-color: #0087bb;
    }

    table.twitter:hover td {
      background: #0087bb !important;
    }

    table.google-plus td {
      background-color: #DB4A39;
      border-color: #CC0000;
    }

    table.google-plus:hover td {
      background: #CC0000 !important;
    }

    .template-label {
      color: #ffffff;
      font-weight: bold;
      font-size: 11px;
    }

    .callout .wrapper {
      padding-bottom: 20px;
    }

    .callout .panel {
      background: #ECF8FF;
      border-color: #b9e5ff;
    }

    .header {
      background: #394e63;
      min-height:100px;
    }

    .footer .wrapper {
      background: #ebebeb;
    }

    .footer h5 {
      padding-bottom: 10px;
    }

    table.columns .text-pad {
      padding-left: 10px;
      padding-right: 10px;
    }

    table.columns .left-text-pad {
      padding-left: 10px;
    }

    table.columns .right-text-pad {
      padding-right: 10px;
    }

    @media only screen and (max-width: 600px) {

      table[class="body"] .right-text-pad {
        padding-left: 10px !important;
      }

      table[class="body"] .left-text-pad {
        padding-right: 10px !important;
      }
    }

  </style>
</head>
<body>
  <table class="body">
    <tr>
      <td class="center" align="center" valign="top">
        <center>

            [header:block]

          <table class="container">
            <tr>
              <td>

                [message:block]

                [callout:block]

                [footer:block]

                [privacy:block]

              <!-- container end below -->
              </td>
            </tr>
          </table>

        </center>
      </td>
    </tr>
  </table>
</body>
</html>\', \'active\', \'item\', NULL, \'2014-06-15 17:44:12\', \'2015-04-16 23:47:23\' );; INSERT INTO email_templates(name,description,content,status,type,parent,created_at,updated_at) VALUES ( \'message\', \'\', \'<table class="row">
  <tr>
      <td class="wrapper last">

          <table class="twelve columns">
              <tr>
                  <td>
                      [title:block]
                      [body:block]
                  </td>
                  <td class="expander"></td>
              </tr>
          </table>

      </td>
  </tr>
</table>\', \'active\', \'item\', \'main\', \'2014-06-20 09:15:51\', \'2015-04-12 01:40:00\' );; INSERT INTO email_templates(name,description,content,status,type,parent,created_at,updated_at) VALUES ( \'password_reset\', \'Восстановление пароля\', \'{"items":["body","header"],"title":"Смена логина или пароля для [user:name:last] [user:name:first] на сайте [site:name]","fields":[{"name":"body","value":"[user:name:last] [user:name:first],<br><br>Запрос на сброс пароля для вашего аккаунта был сделан на сайте [site:name]. <br><br>Вы можете сейчас войти на сайт, кликнув на ссылку или скопировав и вставив её в браузер: <br><br>[user:one-time-login-url] <br><br>Это одноразовая ссылка для входа и она перебросит Вас на страницу, где Вы сможете изменить пароль. Ссылка истекает через 1 сутки и ничего не случится, если она не будет использована. <br><br>-- Команда сайта [site:name]"}]}\', \'active\', \'template\', NULL, \'2014-06-20 18:35:00\', \'2015-04-16 17:27:39\' );; INSERT INTO email_templates(name,description,content,status,type,parent,created_at,updated_at) VALUES ( \'privacy\', \'\', \'<table class="row">
    <tr>
        <td class="wrapper last">

            <table class="twelve columns">
                <tr>
                    <td align="center">
                        <center>
                            [privacy:content]
                        </center>
                    </td>
                    <td class="expander"></td>
                </tr>
            </table>
        </td>
    </tr>
</table>\', \'active\', \'item\', \'main\', \'2014-06-15 17:45:42\', \'2015-04-12 01:38:28\' );; INSERT INTO email_templates(name,description,content,status,type,parent,created_at,updated_at) VALUES ( \'register_admin_created\', \'Добро пожаловать (новый пользователь создан администратором)\', \'{"items":["header","body"],"title":"Администратор создал для вас учётную запись","fields":[{"name":"body","value":"[user:name:last] [user:name:first],<br><br>Администратор системы [site:name] создал для вас аккаунт. Можете войти в систему, кликнув на ссылку или скопировав и вставив её в адресную строку браузера:<br><br>[user:one-time-login-url] <br><br>Эта одноразовая ссылка для входа в систему направит вас на страницу задания своего пароля.<br><br>После установки пароля вы сможете входить в систему через страницу<br>[site:login-url]<br><br>-- Команда [site:name]"}]}\', \'active\', \'template\', \'\', \'2014-06-20 18:27:05\', \'2017-01-15 17:03:11\' );; INSERT INTO email_templates(name,description,content,status,type,parent,created_at,updated_at) VALUES ( \'social\', \'\', \'<table class="six columns">
    <tr>
        <td class="left-text-pad">

            <h3>Мы в сети:</h3>

            [facebook:block]
            [twitter:block]
            [vk:block]
            [google:block]

</td>
        <td class="expander"></td>
    </tr>
</table>\', \'active\', \'item\', \'footer\', \'2014-06-15 17:46:47\', \'2015-06-23 14:10:14\' );; INSERT INTO email_templates(name,description,content,status,type,parent,created_at,updated_at) VALUES ( \'status_activated\', \'Учётная запись активирована\', \'{"items":["header","body"],"title":"Детали учётной записи для [user:name:last] [user:name:first] на [site:name] (одобрено)","fields":[{"name":"body","value":"[user:name:last] [user:name:first],<br><br>Ваш аккаунт на сайте [site:name] был активирован.<br><br>Вы можете войти на сайт, кликнув на ссылку или скопировав и вставив её в Ваш браузер: <br><br>[user:one-time-login-url] <br><br>Это одноразовая ссылка для входа и она перебросит Вас на страницу, где Вы сможете установить свой пароль.<br><br>После установки Вашего пароля, вы сможете заходить на странице [site:login-url].<br><br>-- Команда сайта [site:name]"}]}\', \'active\', \'template\', \'\', \'2014-06-20 18:31:06\', \'2017-01-15 06:07:52\' );; INSERT INTO email_templates(name,description,content,status,type,parent,created_at,updated_at) VALUES ( \'title\', \'\', \'<h1>[title:content]</h1>
\', \'active\', \'item\', \'message\', \'2014-06-21 21:54:48\', \'2015-04-14 15:20:18\' );; INSERT INTO email_templates(name,description,content,status,type,parent,created_at,updated_at) VALUES ( \'twitter\', \'\', \'<table class="tiny-button twitter">
    <tr>
        <td>
            <a href="#">Twitter</a>
        </td>
    </tr>
</table>

<br>\', \'active\', \'item\', \'social\', \'2014-06-15 17:47:08\', \'2015-04-13 01:47:52\' );; INSERT INTO email_templates(name,description,content,status,type,parent,created_at,updated_at) VALUES ( \'vk\', \'\', \'<table class="tiny-button vk">
    <tr>
        <td>
            <a href="#">vk</a>
        </td>
    </tr>
</table>

<br>\', \'active\', \'item\', \'social\', \'2014-06-15 17:50:15\', \'2015-04-13 01:51:30\' );', 'DELETE FROM `email_templates`', 'rollback' );
INSERT INTO `migrations`(`id_migration`,`name`,`created_at`,`statements`,`rollback_statements`,`status`) VALUES ( '31', 'EmailTemplates_20170121_200137', '2017-01-21 20:13:16', 'INSERT INTO email_templates(name,description,content,status,type,parent,created_at,updated_at) VALUES ( \'body\', \'\', \'[body:content]\', \'active\', \'item\', \'message\', \'2014-06-21 21:56:07\', \'2015-04-14 15:24:39\' );; INSERT INTO email_templates(name,description,content,status,type,parent,created_at,updated_at) VALUES ( \'callout\', \'\', \'<table class="row callout">
    <tr>
        <td class="wrapper last">

            <table class="twelve columns">
                <tr>
                    <td class="panel">
                        [callout:content]
                    </td>
                    <td class="expander"></td>
                </tr>
            </table>

        </td>
    </tr>
</table>\', \'active\', \'item\', \'main\', \'2014-06-15 17:45:31\', \'2015-04-12 01:40:03\' );; INSERT INTO email_templates(name,description,content,status,type,parent,created_at,updated_at) VALUES ( \'contacts\', \'\', \'<table class="six columns">
    <tr>
        <td class="last right-text-pad">
        <h3>Контакты:</h3>
            [contacts:content]
        </td>
        <td class="expander"></td>
    </tr>
</table>\', \'active\', \'item\', \'footer\', \'2014-06-15 17:45:20\', \'2015-04-12 01:40:03\' );; INSERT INTO email_templates(name,description,content,status,type,parent,created_at,updated_at) VALUES ( \'facebook\', \'\', \'<table class="tiny-button facebook">
    <tr>
        <td>
            <a href="#">Facebook</a>
        </td>
    </tr>
</table>

<br>\', \'active\', \'item\', \'social\', \'2014-06-15 17:46:59\', \'2015-04-13 01:50:01\' );; INSERT INTO email_templates(name,description,content,status,type,parent,created_at,updated_at) VALUES ( \'footer\', \'\', \'<table class="row footer">
    <tr>

        <td class="wrapper ">

            [social:block]

        </td>
        <td class="wrapper last">

            [contacts:block]

        </td>

    </tr>
</table>\', \'active\', \'item\', \'main\', \'2014-06-15 17:45:07\', \'2015-04-12 01:40:03\' );; INSERT INTO email_templates(name,description,content,status,type,parent,created_at,updated_at) VALUES ( \'google\', \'\', \'<table class="tiny-button google-plus">
    <tr>
        <td>
            <a href="#">Google +</a>
        </td>
    </tr>
</table>

<br>\', \'active\', \'item\', \'social\', \'2014-06-15 17:47:17\', \'2015-04-12 01:40:01\' );; INSERT INTO email_templates(name,description,content,status,type,parent,created_at,updated_at) VALUES ( \'header\', \'\', \'<table class="row header">
    <tr>
        <td class="center" align="center">
            <center>

                <table class="container">
                    <tr>
                        <td class="wrapper last">

                            <table class="twelve columns">
                                <tr>
                                    <td class="six sub-columns">
                                        <img alt="Облачная типография, Календари-домики, Листовки и Флаеры, Карманные календари, Визитки, Буклеты, Пластиковые карты" src="data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAQIAAAAkCAYAAABizTTPAAAABHNCSVQICAgIfAhkiAAAAAlwSFlzAAAN1wAADdcBQiibeAAAABl0RVh0U29mdHdhcmUAd3d3Lmlua3NjYXBlLm9yZ5vuPBoAAA5dSURBVHic7Z15lBXFFYe/GcYdQYdFURhxQcSVRAGjcY1bXI4mLokxLjFqzKIecQlRY4xrjMZo1MQ1iqgRjx6XxN24gAvGhSjigiCIiAoioIiKMDd//LpDv+rqftXvvWFmPO87p89M9avt1au+devWreoGM6MM6wH7AFsDvaOrO/AxMAv4EPgv8E9gUrnM6nxt6Al8OxE24J52qkudKmnIEQQ/Ak4DNimQ31vAVcDlwFfVVa1OjRiMBPTCGue7M/DvRLgV6FLjMuosIxo99zYGxgK3UEwIAAwA/gS8AnynuqrVqZLuwF+AF4CV27kudTo4riAYADxOqcpXCRsBDwPfqzKfOpVxCPAGcBz1UbpOAE2J/3sDD0V/XZ5DNoD3omsW0ANoia7dgW2cNI3ArcCuwFM1rXWdPNYHbm7vStTpXCQFwQhgXefz54FfIy0hj7OBPYAL0Jw0ZkVgFNI0FldV0zp16rQZ8dSgGTja+WwSsBvlhUDMg8D2wJvO/f7ADyqsX506dZYBsSDYG+jqfHYwMK9gfp8C+wOfOfeHF69anTp1lhWxINjQuT8beKnCPCcCdzr3tgBWqDC/OnXqtDGxjWA95361jkGPAIclwl2AgWhZsVK6IsNkb+TENANpICE0oynKisjXYXYV9XBZDlgDGU+bkf/EPGAO8H4Ny4nLGgisiqZgH9c4/7agJ7A2apv4d/ukxmW0AL3Qkul8ZMyeCSypcTkxK0ZlroV+g3eBuYFpuwHrAKsBU1A9a0kvYE3U7p8CH0RltOYligWB26EGVlmZR4Bno3znRH+/zIl/C3qQQD9iLES6IPvCscB2nnTPAJcCdyDPtiRdgZ8CxyDfiCRzo/qdAEwu+23S9AIOB3aJ6pW1Tv8+MAa4Cbg/MO8HE/9PAX4Z/X8UMsr2SXz+JnAxMupeGN1byZPnaEodvM6lbVdyVgGOR7/dFp7PX0VtcgPwUYVlbIXaZlckaFzmoiXsW4F7C+T7R2DzRHhvlhq69wR+DnyX9LLsROAyYCSwyPmsCfgx6ovfcj5bgDxzTwL+U6CeSQaiZ2QvZJh3+QitCN5AqRPYUswMMzvR0uwffbYsrpmJct+J7q1tZmM89fIxysyaEvltZmZvBKT7zMxOKFDPVczsIjNbEFivJE+Z2VoBZSQZH90bXibvcwvW5YCAepS7dnbyXBLd3830G4bwvpntVLDctczsvmJf18aZ2VaB+T/ppF3ezLqZ2e2BZT0exY/z62v67cux2MwutNJ+XO5a1cyuMLOvAutmZvYvM2tx84r/2cjMWp0EH5vZngUqVUtB0GLqJEW4OsprmOkBL8IxAXVsNnWoaphmZj3KlJNkvJkNNT1kWSw0PXxFaCtBcICZLSpYl8Vmtn1gmduZ2eyC+ccstLDv7QqCnmb2SsGyHovy6m9mHxRMe3FAHWOBOL5g3jEzzWxwMr/kXoOHkZrl8izyD3gU+LxC1aUcM1mq8s5CKnVSpXweuT2/B2yA1PFNnTyWIF+GmxJ5fQGMQ4bPmUiFOgBY3Uk7H6mX7mpHkifR8miSBcBtSNWdEZXXHzn17BPV1eV6pOZnkZzivIpU+m/kxB8NnA6cEoW7oRWfJCOjusVcQ+XG4Bh3rwFIhU76pjyNnNFmoPbdGDmfuR6tk5E6nte/BqGp4GrO/bnAE8iVegLyhRmCfqsWJ64B3wfuzinH/Z0fRsvoMW9G96ZF+Q8lre6DpkWns3SasQSp/i8BbyM7wf6kpzWtSL1/O6eOqwEvkrbtLQDuAsajqUq/qH7bI2/fJPNRv5oKkJQwg8xsbo4UWWhmD5jZ8WY2IFBqVaIRJHnX/KNFF5PkdEmOnE9l1LObmb3mSXtETv0O8sQfbWbdc9I0mdmxpnZLstDMVs5Jl8UsMzvJzLY0syEmLeZFM9vLSb++J23PnPJqpREkecY0Pcsa1ad50uRNRVcys6meNGPNrE9Omhs8aeaZ2ihUI4iZb9IoGjxpjrO01pYMv27SVH195BFPWefl1A9T33N5yMzWyYjfaGYXeNI8G9UBN8HOZva5J4GP6aa5+VFmtkGZilciCN4zs35l0k3IqNv9pi+flW5D0w+b5Pqc+M87cV8wsxUDv9sIT/12yInvY7p55nXR5XbM9hYEd5vZCmXSbu1Jl9f+J3ni/9XC5tM/y6hjEUGw0BxV2nPd5WsMM5to+X2lh6WF3Jic+Ht6yrgmoB0wCbLFTtqjzdKCAJO9oJK58LtmdomFG2WSl08QHB6Q7heedF+Z2cCAtDc46Z7JiNfXU8a+Bb7bSpa2v+yXE99HXnz3ak9BMNfMegemd41vr2fEW8mkDSUZZ9IKQ+t6lZO+1bL7iE8QnBtQxq6edGZmewSkPc1JMzsn7uNO3Glm1rVAW1zvpH/NzBp825DfALYFjkDzu1D6Aiei+fwrwA4F0rrMQXsUyjHRc+9m0m7OPl52wu7cM6Y7OmPhETRvmwHcF5B/zOek/Qm6FUj/Pp3nwI8bkY0nhCeccA9fJNSPejn3TqeYj8AZlC7pNQBHFkh/aUAcX198mtLl4CxC++LGwI7OveHINhDKWZS23SBgW58gIIo4Ep1KNBitXz9HGaeEBJuhH3ok2V8qj3GBZfkcdtxGzcL1nVglI95EtHa8GzIC9qP4BirXaaSIIHiatI9ER+XRAnHfccKuATdmJyc8DXisQDmgdXRXmLr5ZjGJMF+HD0n32Ur7YhOwvCeea6w2svwCsnkXDfZJtssSBEleBn6HhEIv4IfAdcjZpRyHAQ+Q/ZBlMT4wnm8vRJ61NUmoV2I1NKMzGdZ07hc5KKSIVtbeTCgQd4YTbsLfT1xHsjFUJhifcMLfzCjPJbQvLiE9Mof2xVBPS3er/yRk/S+K+522afJGy+ZjtFw1Ogq3oGWkvdDSnbtxCSRA7kbLRqEaRai7po+pgfFq6X66HFoqHIiWaQaipZnNkRpaDe4D05Ep0il9R9n5BqZ+TrjSZU83XRckoMsNaEU33iUJ7Yuhz8UgJ9wI/CG8Ov+nvxPeoKggcJmO5oU3ok1Fh6H5mLt+uwsaGd3NSFlUM1oXmS9VyjDkbroJmrcNoHT9vJbMaaN824K20LJc20GoDcLFp95nTUeSVLMvotZ90W2LAei8kKrzDZkahPIlcC0aDX0P/Cmee1nkOfa0J3sih5NxwJnIIWQQ+UJgCWnf8yIsiylMLWglfGQLZWXSeycqfTB96VYNSNeR+mJzW+XrCoIGqj/j7gvgIOQFlmQY/mPQfHQ049jyyIPwPtIGGx/TgNuRhrQG2lRSKR2tLZYlvulbpf3TZ3wLmcp0pPZ3B5zFaACu9lrUhKyp/dFD2gstAV5eZYVbgXOQoTDJelSu2rUn15J9ytJ0JPTGoQf+ZdLzyvoBopURd9TkWRZFVlyS+NJVuvOxvfiIUgPnCHRqeNU0ofl8ctvl0FpkjHyhXfqhB6YzsRelZyuARqrfAP8gzJjnqnS1nJJ93ZmHtKqYvhXm40vX2QTBHLRHIabStkjRCLzu3NsJWcGrxTdf7GwND0vPA4j5CgmHiwi36LtGnrqGEI7rqDOkwnwGO+GF1P6lL22N6yi3Va0ybiR9atDa6C1H1eLbMRfie9CRaEArHkmeRIc8hLI+abW0rVYYvo6MdcI7Utmxd7s74c54xP7TTnhrwlY+XI5Bxu5D0TtM+jTiP1HlbHQMUzW4qwTxMVWdiV6ktaMiQgC0HdmlLTUCn9djLTS89sL1nGtGW8mLsAFpT8LQpeyOhHuieBPwq4J5rIFcpn+PtuyPBS5rRG66tzqRW1CHX4PKGEHpHm6iwmu9vNTWzCf9YBVZwukD/NZzvy01Ap+6m+XH3xkYS9pj8XzClv5iLqfUsWsJ2rff2XiNtIfkyaTPJcjjVNJLsrfGRquLSJ8puClyYTyZcFVsTWRhv8C5Px/4W3BVOw5fkrah7I1/KcqlGXlg+gRHkU5cFJ8gqJlRqZ24xAm3IG/VnmXSrQBcjbxek4yitgfYLksucsLdUFv4DsFxORStCiZ5Dbi3MRE4hPSI3S0q+EOkNRyMjDUtqDP3RUaYI9E6+2TSp+8siSpQiU90R+AFJ7wZcCXZ+wUa0ElPL+E/cBXazjEE/Eawq9FvMBj5QfgO++zIjCQ9JdsZ+cxvnZFmfbSse4xzfwo6WLWzcj961pJshtoia0dlV6SZ/p1SzcjQoamtSRX1TtRAV3gy6o6EgHsEVggnoPcmdlbOBPaj1ChzFFo5uAl1rFlIG1oHzV+TJ8kuQrvlkqNSW47QhoRX0vGpJaprzKmkR5aOjKEl3OcpdV/vi47Sm4q+86tIAAxD7+pw93ksQH24s3hrZnEs2jSVfB9JV3QM3gVoEIr3VgxA9hGf9nQu0TZpd656JdoxdQXF5h0+ZiEh4EqvzsYM1PCjnft9KO/nPRU4EK2FJ49N3zZKX+v3HsTcTL4HpLt5pTMwCz3g95D2dVk3ug7MSf8qEtIhZ1V0dOajnYh3kj73ozcadNzpkMt5aJAD/I4tD6DNNGejlyMUZSE6HHMjOr8QiLkd+Anhm0g+QVbZwcixagqlO9G6UJl2Fcp15HuHugdZdhY+QB3/NMJ3Bc5HNoahfD2EQMwcNAUdTjF7x2QkJM5I3kyeYuyjAdgSbbbZEY1iPdEctxGpbLPRyPYKssQ+RHFHjZMoNaDdgSR4OVYhvUx5KWGdZCClD+M8yp9E04wEwu6oXVZHbbQItcM49P3vIL2V+kAkYGPeQi92cTnLCV+LTm+uhCHoxRqD0EgxFwmkCcCfK8wzZl30kpeYVjR4hNKLtLPW+YRv0FodnUi8B/qevZE1/DPk9j0VjZi3Uaw/HkHpNt3Hke9ICCPQW5BiRhHmO9OT9DLgOYRvle8K7IsEwzZotS/2XfkCPfzjUHs8jGf17n+sJrYxUbQ5hgAAAABJRU5ErkJggg==">
                                    </td>
                                    <td class="six sub-columns last" style="text-align:right; vertical-align:middle;">
                                    </td>
                                    <td class="expander"></td>
                                </tr>
                            </table>

                        </td>
                    </tr>
                </table>

            </center>
        </td>
    </tr>
</table>\', \'active\', \'item\', \'main\', \'2014-06-15 17:44:55\', \'2015-04-16 23:43:22\' );; INSERT INTO email_templates(name,description,content,status,type,parent,created_at,updated_at) VALUES ( \'main\', \'Основной слой\', \'<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Strict//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-strict.dtd">
<html xmlns="http://www.w3.org/1999/xhtml">
<head>
  <meta http-equiv="Content-Type" content="text/html; charset=UTF-8">
  <meta name="viewport" content="width=device-width"/>
  <style type="text/css">

#outlook a {
  padding:0;
}

body{
  width:100% !important;
  min-width: 100%;
  -webkit-text-size-adjust:100%;
  -ms-text-size-adjust:100%;
  margin:0;
  padding:0;
}

.ExternalClass {
  width:100%;
}

.ExternalClass,
.ExternalClass p,
.ExternalClass span,
.ExternalClass font,
.ExternalClass td,
.ExternalClass div {
  line-height: 100%;
}

#backgroundTable {
  margin:0;
  padding:0;
  width:100% !important;
  line-height: 100% !important;
}

img {
  outline:none;
  text-decoration:none;
  -ms-interpolation-mode: bicubic;
  width: auto;
  max-width: 100%;
  float: left;
  clear: both;
  display: block;
}

center {
  width: 100%;
  min-width: 580px;
}

a img {
  border: none;
}

p {
  margin: 0 0 0 10px;
}

table {
  border-spacing: 0;
  border-collapse: collapse;
}

td {
  word-break: break-word;
  -webkit-hyphens: auto;
  -moz-hyphens: auto;
  hyphens: auto;
  border-collapse: collapse !important;
}

table, tr, td {
  padding: 0;
  vertical-align: top;
  text-align: left;
}

hr {
  color: #d9d9d9;
  background-color: #d9d9d9;
  height: 1px;
  border: none;
}

/* Responsive Grid */

table.body {
  height: 100%;
  width: 100%;
}

table.container {
  width: 580px;
  margin: 0 auto;
  text-align: inherit;
}

table.row {
  padding: 0;
  width: 100%;
  position: relative;
}

table.container table.row {
  display: block;
}

td.wrapper {
  padding: 10px 20px 0 0;
  position: relative;
}

table.columns,
table.column {
  margin: 0 auto;
}

table.columns td,
table.column td {
  padding: 0 0 10px;
}

table.columns td.sub-columns,
table.column td.sub-columns,
table.columns td.sub-column,
table.column td.sub-column {
  padding-right: 10px;
}

td.sub-column, td.sub-columns {
  min-width: 0;
}

table.row td.last,
table.container td.last {
  padding-right: 0;
}

table.one { width: 30px; }
table.two { width: 80px; }
table.three { width: 130px; }
table.four { width: 180px; }
table.five { width: 230px; }
table.six { width: 280px; }
table.seven { width: 330px; }
table.eight { width: 380px; }
table.nine { width: 430px; }
table.ten { width: 480px; }
table.eleven { width: 530px; }
table.twelve { width: 580px; }

table.one center { min-width: 30px; }
table.two center { min-width: 80px; }
table.three center { min-width: 130px; }
table.four center { min-width: 180px; }
table.five center { min-width: 230px; }
table.six center { min-width: 280px; }
table.seven center { min-width: 330px; }
table.eight center { min-width: 380px; }
table.nine center { min-width: 430px; }
table.ten center { min-width: 480px; }
table.eleven center { min-width: 530px; }
table.twelve center { min-width: 580px; }

table.one .panel center { min-width: 10px; }
table.two .panel center { min-width: 60px; }
table.three .panel center { min-width: 110px; }
table.four .panel center { min-width: 160px; }
table.five .panel center { min-width: 210px; }
table.six .panel center { min-width: 260px; }
table.seven .panel center { min-width: 310px; }
table.eight .panel center { min-width: 360px; }
table.nine .panel center { min-width: 410px; }
table.ten .panel center { min-width: 460px; }
table.eleven .panel center { min-width: 510px; }
table.twelve .panel center { min-width: 560px; }

.body .columns td.one,
.body .column td.one { width: 8.333333%; }
.body .columns td.two,
.body .column td.two { width: 16.666666%; }
.body .columns td.three,
.body .column td.three { width: 25%; }
.body .columns td.four,
.body .column td.four { width: 33.333333%; }
.body .columns td.five,
.body .column td.five { width: 41.666666%; }
.body .columns td.six,
.body .column td.six { width: 50%; }
.body .columns td.seven,
.body .column td.seven { width: 58.333333%; }
.body .columns td.eight,
.body .column td.eight { width: 66.666666%; }
.body .columns td.nine,
.body .column td.nine { width: 75%; }
.body .columns td.ten,
.body .column td.ten { width: 83.333333%; }
.body .columns td.eleven,
.body .column td.eleven { width: 91.666666%; }
.body .columns td.twelve,
.body .column td.twelve { width: 100%; }

td.offset-by-one { padding-left: 50px; }
td.offset-by-two { padding-left: 100px; }
td.offset-by-three { padding-left: 150px; }
td.offset-by-four { padding-left: 200px; }
td.offset-by-five { padding-left: 250px; }
td.offset-by-six { padding-left: 300px; }
td.offset-by-seven { padding-left: 350px; }
td.offset-by-eight { padding-left: 400px; }
td.offset-by-nine { padding-left: 450px; }
td.offset-by-ten { padding-left: 500px; }
td.offset-by-eleven { padding-left: 550px; }

td.expander {
  visibility: hidden;
  width: 0;
  padding: 0 !important;
}

table.columns .text-pad,
table.column .text-pad {
  padding-left: 10px;
  padding-right: 10px;
}

table.columns .left-text-pad,
table.columns .text-pad-left,
table.column .left-text-pad,
table.column .text-pad-left {
  padding-left: 10px;
}

table.columns .right-text-pad,
table.columns .text-pad-right,
table.column .right-text-pad,
table.column .text-pad-right {
  padding-right: 10px;
}

/* Block Grid */

.block-grid {
  width: 100%;
  max-width: 580px;
}

.block-grid td {
  display: inline-block;
  padding:10px;
}

.two-up td {
  width:270px;
}

.three-up td {
  width:173px;
}

.four-up td {
  width:125px;
}

.five-up td {
  width:96px;
}

.six-up td {
  width:76px;
}

.seven-up td {
  width:62px;
}

.eight-up td {
  width:52px;
}

/* Alignment & Visibility Classes */

table.center, td.center {
  text-align: center;
}

h1.center,
h2.center,
h3.center,
h4.center,
h5.center,
h6.center {
  text-align: center;
}

span.center {
  display: block;
  width: 100%;
  text-align: center;
}

img.center {
  margin: 0 auto;
  float: none;
}

.show-for-small,
.hide-for-desktop {
  display: none;
}

/* Typography */

body, table.body, h1, h2, h3, h4, h5, h6, p, td {
  color: #222222;
  font-family: "Helvetica", "Arial", sans-serif;
  font-weight: normal;
  padding:0;
  margin: 0;
  text-align: left;
  line-height: 1.3;
}

h1, h2, h3, h4, h5, h6 {
  word-break: normal;
}

/*h1 {font-size: 40px;}*/
h1 {font-size: 30px;}
/*h2 {font-size: 36px;}*/
h2 {font-size: 26px;}
h3 {font-size: 32px;}
h4 {font-size: 28px;}
h5 {font-size: 27px;}
h6 {font-size: 20px;}
body, table.body, p, td {font-size: 14px;line-height:19px;}

p.lead, p.lede, p.leed {
  font-size: 18px;
  line-height:21px;
}

p {
  margin-bottom: 10px;
}

small {
  font-size: 10px;
}

a {
  color: #2ba6cb;
  text-decoration: none;
}

a:hover {
  color: #2795b6 !important;
}

a:active {
  color: #2795b6 !important;
}

a:visited {
  color: #2ba6cb !important;
}

h1 a,
h2 a,
h3 a,
h4 a,
h5 a,
h6 a {
  color: #2ba6cb;
}

h1 a:active,
h2 a:active,
h3 a:active,
h4 a:active,
h5 a:active,
h6 a:active {
  color: #2ba6cb !important;
}

h1 a:visited,
h2 a:visited,
h3 a:visited,
h4 a:visited,
h5 a:visited,
h6 a:visited {
  color: #2ba6cb !important;
}

/* Panels */

.panel {
  background: #f2f2f2;
  border: 1px solid #d9d9d9;
  padding: 10px !important;
}

.sub-grid table {
  width: 100%;
}

.sub-grid td.sub-columns {
  padding-bottom: 0;
}

/* Buttons */

table.button,
table.tiny-button,
table.small-button,
table.medium-button,
table.large-button {
  width: 100%;
  overflow: hidden;
}

table.button td,
table.tiny-button td,
table.small-button td,
table.medium-button td,
table.large-button td {
  display: block;
  width: auto !important;
  text-align: center;
  background: #2ba6cb;
  border: 1px solid #2284a1;
  color: #ffffff;
  padding: 8px 0;
}

table.tiny-button td {
  padding: 5px 0 4px;
}

table.small-button td {
  padding: 8px 0 7px;
}

table.medium-button td {
  padding: 12px 0 10px;
}

table.large-button td {
  padding: 21px 0 18px;
}

table.button td a,
table.tiny-button td a,
table.small-button td a,
table.medium-button td a,
table.large-button td a {
  font-weight: bold;
  text-decoration: none;
  font-family: Helvetica, Arial, sans-serif;
  color: #ffffff;
  font-size: 16px;
}

table.tiny-button td a {
  font-size: 12px;
  font-weight: normal;
}

table.small-button td a {
  font-size: 16px;
}

table.medium-button td a {
  font-size: 20px;
}

table.large-button td a {
  font-size: 24px;
}

table.button:hover td,
table.button:visited td,
table.button:active td {
  background: #2795b6 !important;
}

table.button:hover td a,
table.button:visited td a,
table.button:active td a {
  color: #fff !important;
}

table.button:hover td,
table.tiny-button:hover td,
table.small-button:hover td,
table.medium-button:hover td,
table.large-button:hover td {
  background: #2795b6 !important;
}

table.button:hover td a,
table.button:active td a,
table.button td a:visited,
table.tiny-button:hover td a,
table.tiny-button:active td a,
table.tiny-button td a:visited,
table.small-button:hover td a,
table.small-button:active td a,
table.small-button td a:visited,
table.medium-button:hover td a,
table.medium-button:active td a,
table.medium-button td a:visited,
table.large-button:hover td a,
table.large-button:active td a,
table.large-button td a:visited {
  color: #ffffff !important;
}

table.secondary td {
  background: #e9e9e9;
  border-color: #d0d0d0;
  color: #555;
}

table.secondary td a {
  color: #555;
}

table.secondary:hover td {
  background: #d0d0d0 !important;
  color: #555;
}

table.secondary:hover td a,
table.secondary td a:visited,
table.secondary:active td a {
  color: #555 !important;
}

table.success td {
  background: #5da423;
  border-color: #457a1a;
}

table.success:hover td {
  background: #457a1a !important;
}

table.alert td {
  background: #c60f13;
  border-color: #970b0e;
}

table.alert:hover td {
  background: #970b0e !important;
}

table.radius td {
  -webkit-border-radius: 3px;
  -moz-border-radius: 3px;
  border-radius: 3px;
}

table.round td {
  -webkit-border-radius: 500px;
  -moz-border-radius: 500px;
  border-radius: 500px;
}

/* Outlook First */

body.outlook p {
  display: inline !important;
}

/*  Media Queries */

@media only screen and (max-width: 600px) {

  table[class="body"] img {
    width: auto !important;
    height: auto !important;
  }

  table[class="body"] center {
    min-width: 0 !important;
  }

  table[class="body"] .container {
    width: 95% !important;
  }

  table[class="body"] .row {
    width: 100% !important;
    display: block !important;
  }

  table[class="body"] .wrapper {
    display: block !important;
    padding-right: 0 !important;
  }

  table[class="body"] .columns,
  table[class="body"] .column {
    table-layout: fixed !important;
    float: none !important;
    width: 100% !important;
    padding-right: 0 !important;
    padding-left: 0 !important;
    display: block !important;
  }

  table[class="body"] .wrapper.first .columns,
  table[class="body"] .wrapper.first .column {
    display: table !important;
  }

  table[class="body"] table.columns td,
  table[class="body"] table.column td {
    width: 100% !important;
  }

  table[class="body"] .columns td.one,
  table[class="body"] .column td.one { width: 8.333333% !important; }
  table[class="body"] .columns td.two,
  table[class="body"] .column td.two { width: 16.666666% !important; }
  table[class="body"] .columns td.three,
  table[class="body"] .column td.three { width: 25% !important; }
  table[class="body"] .columns td.four,
  table[class="body"] .column td.four { width: 33.333333% !important; }
  table[class="body"] .columns td.five,
  table[class="body"] .column td.five { width: 41.666666% !important; }
  table[class="body"] .columns td.six,
  table[class="body"] .column td.six { width: 50% !important; }
  table[class="body"] .columns td.seven,
  table[class="body"] .column td.seven { width: 58.333333% !important; }
  table[class="body"] .columns td.eight,
  table[class="body"] .column td.eight { width: 66.666666% !important; }
  table[class="body"] .columns td.nine,
  table[class="body"] .column td.nine { width: 75% !important; }
  table[class="body"] .columns td.ten,
  table[class="body"] .column td.ten { width: 83.333333% !important; }
  table[class="body"] .columns td.eleven,
  table[class="body"] .column td.eleven { width: 91.666666% !important; }
  table[class="body"] .columns td.twelve,
  table[class="body"] .column td.twelve { width: 100% !important; }

  table[class="body"] td.offset-by-one,
  table[class="body"] td.offset-by-two,
  table[class="body"] td.offset-by-three,
  table[class="body"] td.offset-by-four,
  table[class="body"] td.offset-by-five,
  table[class="body"] td.offset-by-six,
  table[class="body"] td.offset-by-seven,
  table[class="body"] td.offset-by-eight,
  table[class="body"] td.offset-by-nine,
  table[class="body"] td.offset-by-ten,
  table[class="body"] td.offset-by-eleven {
    padding-left: 0 !important;
  }

  table[class="body"] table.columns td.expander {
    width: 1px !important;
  }

  table[class="body"] .right-text-pad,
  table[class="body"] .text-pad-right {
    padding-left: 10px !important;
  }

  table[class="body"] .left-text-pad,
  table[class="body"] .text-pad-left {
    padding-right: 10px !important;
  }

  table[class="body"] .hide-for-small,
  table[class="body"] .show-for-desktop {
    display: none !important;
  }

  table[class="body"] .show-for-small,
  table[class="body"] .hide-for-desktop {
    display: inherit !important;
  }
}

  </style>
  <style type="text/css">

    table.facebook td {
      background: #3b5998;
      border-color: #2d4473;
    }

    table.facebook:hover td {
      background: #2d4473 !important;
    }

    table.twitter td {
      background: #00acee;
      border-color: #0087bb;
    }

    table.twitter:hover td {
      background: #0087bb !important;
    }

    table.google-plus td {
      background-color: #DB4A39;
      border-color: #CC0000;
    }

    table.google-plus:hover td {
      background: #CC0000 !important;
    }

    .template-label {
      color: #ffffff;
      font-weight: bold;
      font-size: 11px;
    }

    .callout .wrapper {
      padding-bottom: 20px;
    }

    .callout .panel {
      background: #ECF8FF;
      border-color: #b9e5ff;
    }

    .header {
      background: #394e63;
      min-height:100px;
    }

    .footer .wrapper {
      background: #ebebeb;
    }

    .footer h5 {
      padding-bottom: 10px;
    }

    table.columns .text-pad {
      padding-left: 10px;
      padding-right: 10px;
    }

    table.columns .left-text-pad {
      padding-left: 10px;
    }

    table.columns .right-text-pad {
      padding-right: 10px;
    }

    @media only screen and (max-width: 600px) {

      table[class="body"] .right-text-pad {
        padding-left: 10px !important;
      }

      table[class="body"] .left-text-pad {
        padding-right: 10px !important;
      }
    }

  </style>
</head>
<body>
  <table class="body">
    <tr>
      <td class="center" align="center" valign="top">
        <center>

            [header:block]

          <table class="container">
            <tr>
              <td>

                [message:block]

                [callout:block]

                [footer:block]

                [privacy:block]

              <!-- container end below -->
              </td>
            </tr>
          </table>

        </center>
      </td>
    </tr>
  </table>
</body>
</html>\', \'active\', \'item\', NULL, \'2014-06-15 17:44:12\', \'2015-04-16 23:47:23\' );; INSERT INTO email_templates(name,description,content,status,type,parent,created_at,updated_at) VALUES ( \'message\', \'\', \'<table class="row">
  <tr>
      <td class="wrapper last">

          <table class="twelve columns">
              <tr>
                  <td>
                      [title:block]
                      [body:block]
                  </td>
                  <td class="expander"></td>
              </tr>
          </table>

      </td>
  </tr>
</table>\', \'active\', \'item\', \'main\', \'2014-06-20 09:15:51\', \'2015-04-12 01:40:00\' );; INSERT INTO email_templates(name,description,content,status,type,parent,created_at,updated_at) VALUES ( \'password_reset\', \'Восстановление пароля\', \'{"items":["body","header"],"title":"Смена логина или пароля для [user:name:last] [user:name:first] на сайте [site:name]","fields":[{"name":"body","value":"[user:name:last] [user:name:first],<br><br>Запрос на сброс пароля для вашего аккаунта был сделан на сайте [site:name]. <br><br>Вы можете сейчас войти на сайт, кликнув на ссылку или скопировав и вставив её в браузер: <br><br>[user:one-time-login-url] <br><br>Это одноразовая ссылка для входа и она перебросит Вас на страницу, где Вы сможете изменить пароль. Ссылка истекает через 1 сутки и ничего не случится, если она не будет использована. <br><br>-- Команда сайта [site:name]"}]}\', \'active\', \'template\', NULL, \'2014-06-20 18:35:00\', \'2015-04-16 17:27:39\' );; INSERT INTO email_templates(name,description,content,status,type,parent,created_at,updated_at) VALUES ( \'privacy\', \'\', \'<table class="row">
    <tr>
        <td class="wrapper last">

            <table class="twelve columns">
                <tr>
                    <td align="center">
                        <center>
                            [privacy:content]
                        </center>
                    </td>
                    <td class="expander"></td>
                </tr>
            </table>
        </td>
    </tr>
</table>\', \'active\', \'item\', \'main\', \'2014-06-15 17:45:42\', \'2015-04-12 01:38:28\' );; INSERT INTO email_templates(name,description,content,status,type,parent,created_at,updated_at) VALUES ( \'register_admin_created\', \'Добро пожаловать (новый пользователь создан администратором)\', \'{"items":["header","body"],"title":"Администратор создал для вас учётную запись","fields":[{"name":"body","value":"[user:name:last] [user:name:first],<br><br>Администратор системы [site:name] создал для вас аккаунт. Можете войти в систему, кликнув на ссылку или скопировав и вставив её в адресную строку браузера:<br><br>[user:one-time-login-url] <br><br>Эта одноразовая ссылка для входа в систему направит вас на страницу задания своего пароля.<br><br>После установки пароля вы сможете входить в систему через страницу<br>[site:login-url]<br><br>-- Команда [site:name]"}]}\', \'active\', \'template\', \'\', \'2014-06-20 18:27:05\', \'2017-01-15 17:03:11\' );; INSERT INTO email_templates(name,description,content,status,type,parent,created_at,updated_at) VALUES ( \'social\', \'\', \'<table class="six columns">
    <tr>
        <td class="left-text-pad">

            <h3>Мы в сети:</h3>

            [facebook:block]
            [twitter:block]
            [vk:block]
            [google:block]

</td>
        <td class="expander"></td>
    </tr>
</table>\', \'active\', \'item\', \'footer\', \'2014-06-15 17:46:47\', \'2015-06-23 14:10:14\' );; INSERT INTO email_templates(name,description,content,status,type,parent,created_at,updated_at) VALUES ( \'status_activated\', \'Учётная запись активирована\', \'{"items":["header","body"],"title":"Детали учётной записи для [user:name:last] [user:name:first] на [site:name] (одобрено)","fields":[{"name":"body","value":"[user:name:last] [user:name:first],<br><br>Ваш аккаунт на сайте [site:name] был активирован.<br><br>Вы можете войти на сайт, кликнув на ссылку или скопировав и вставив её в Ваш браузер: <br><br>[user:one-time-login-url] <br><br>Это одноразовая ссылка для входа и она перебросит Вас на страницу, где Вы сможете установить свой пароль.<br><br>После установки Вашего пароля, вы сможете заходить на странице [site:login-url].<br><br>-- Команда сайта [site:name]"}]}\', \'active\', \'template\', \'\', \'2014-06-20 18:31:06\', \'2017-01-15 06:07:52\' );; INSERT INTO email_templates(name,description,content,status,type,parent,created_at,updated_at) VALUES ( \'title\', \'\', \'<h1>[title:content]</h1>
\', \'active\', \'item\', \'message\', \'2014-06-21 21:54:48\', \'2015-04-14 15:20:18\' );; INSERT INTO email_templates(name,description,content,status,type,parent,created_at,updated_at) VALUES ( \'twitter\', \'\', \'<table class="tiny-button twitter">
    <tr>
        <td>
            <a href="#">Twitter</a>
        </td>
    </tr>
</table>

<br>\', \'active\', \'item\', \'social\', \'2014-06-15 17:47:08\', \'2015-04-13 01:47:52\' );; INSERT INTO email_templates(name,description,content,status,type,parent,created_at,updated_at) VALUES ( \'vk\', \'\', \'<table class="tiny-button vk">
    <tr>
        <td>
            <a href="#">vk</a>
        </td>
    </tr>
</table>

<br>\', \'active\', \'item\', \'social\', \'2014-06-15 17:50:15\', \'2015-04-13 01:51:30\' );', NULL, 'update' );
-- ---------------------------------------------------------


-- Dump data of "nodes" ------------------------------------
INSERT INTO `nodes`(`id`,`ip`,`port`,`name`,`description`,`created_at`,`update_at`,`status`) VALUES ( '1', '127.0.0.1', '3001', 'darkstar', '', '2017-01-21 12:07:45', '2017-01-21 13:27:14', 'enabled' );
INSERT INTO `nodes`(`id`,`ip`,`port`,`name`,`description`,`created_at`,`update_at`,`status`) VALUES ( '2', '10.0.36.1', '3001', 'ganimed', '', '2017-01-21 12:08:01', '2017-01-21 12:09:46', 'enabled' );
-- ---------------------------------------------------------


-- Dump data of "permissions" ------------------------------
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
-- ---------------------------------------------------------


-- Dump data of "roles" ------------------------------------
INSERT INTO `roles`(`name`,`description`,`parent`,`created_at`,`update_at`) VALUES ( 'admin', '', 'user', '2017-01-15 05:20:58', '2017-01-15 05:21:29' );
INSERT INTO `roles`(`name`,`description`,`parent`,`created_at`,`update_at`) VALUES ( 'demo', '', NULL, '2017-01-15 05:17:09', '2017-01-15 05:17:09' );
INSERT INTO `roles`(`name`,`description`,`parent`,`created_at`,`update_at`) VALUES ( 'user', '', 'demo', '2017-01-15 05:17:36', '2017-01-19 19:39:54' );
-- ---------------------------------------------------------


-- Dump data of "scripts" ----------------------------------
INSERT INTO `scripts`(`id`,`lang`,`name`,`source`,`created_at`,`update_at`,`description`,`compiled`) VALUES ( '1', 'coffeescript', 'Розетка/Проверка состояния', '# Worker
# -----------------------------
# result:
# [31 247 1 0 0 68 65] 

main =()->
    
    FUNCTION = 3
    DEVICE_ADDR = message.device.address
    
    # print "send msg from dev:", DEVICE_ADDR
    
    # send message to modbus channel
    request.baud = message.device.baud
    request.result = true
    request.device = message.device.tty
    request.timeout = 4000000
    request.stopBits = 2
    request.sleep = message.device.sleep
    request.command = [DEVICE_ADDR,FUNCTION,0,0,0,5]
    
    from_node = node_send \'modbus\', message.node, request
    if from_node.error
        message.setError from_node.error
        message.device_state \'ERROR\'
        # print "#{message.device.name} - error: #{from_node.error}"
        return false
    
    # decode hex string to byte array
    # result = SmartJs.hex2arr(from_node.result)
    # print "dev:",DEVICE_ADDR ,"result:",result
 
    if from_node.result != ""
        result = SmartJs.hex2arr(from_node.result)
        # проверка питания
        if result[2] == 1
            message.device_state \'ENABLED\'
        else if result[2] == 0
            message.device_state \'DISABLED\'
            
    return from_node.result', '2017-01-21 13:30:46', '2017-01-21 13:30:46', '', 'var main;

main = function() {
  var DEVICE_ADDR, FUNCTION, from_node, result;
  FUNCTION = 3;
  DEVICE_ADDR = message.device.address;
  request.baud = message.device.baud;
  request.result = true;
  request.device = message.device.tty;
  request.timeout = 4000000;
  request.stopBits = 2;
  request.sleep = message.device.sleep;
  request.command = [DEVICE_ADDR, FUNCTION, 0, 0, 0, 5];
  from_node = node_send(\'modbus\', message.node, request);
  if (from_node.error) {
    message.setError(from_node.error);
    message.device_state(\'ERROR\');
    return false;
  }
  if (from_node.result !== "") {
    result = SmartJs.hex2arr(from_node.result);
    if (result[2] === 1) {
      message.device_state(\'ENABLED\');
    } else if (result[2] === 0) {
      message.device_state(\'DISABLED\');
    }
  }
  return from_node.result;
};
' );
INSERT INTO `scripts`(`id`,`lang`,`name`,`source`,`created_at`,`update_at`,`description`,`compiled`) VALUES ( '2', 'coffeescript', 'Розетка/Проверка питания', 'main =->
    # print "Проверка питания"
        
    if message.error != ""
        message.device_state \'ERROR\'
        return false
        
    if message.data[\'result\'] != ""
        result = SmartJs.hex2arr(message.data[\'result\'])
        
        # проверка питания
        if result[2] == 1
            message.device_state \'ENABLED\'
        else if result[2] == 0
            message.device_state \'DISABLED\'
        return false', '2017-01-21 13:35:25', '2017-01-21 13:35:25', '', 'var main;

main = function() {
  var result;
  if (message.error !== "") {
    message.device_state(\'ERROR\');
    return false;
  }
  if (message.data[\'result\'] !== "") {
    result = SmartJs.hex2arr(message.data[\'result\']);
    if (result[2] === 1) {
      message.device_state(\'ENABLED\');
    } else if (result[2] === 0) {
      message.device_state(\'DISABLED\');
    }
    return false;
  }
};
' );
INSERT INTO `scripts`(`id`,`lang`,`name`,`source`,`created_at`,`update_at`,`description`,`compiled`) VALUES ( '3', 'coffeescript', 'Розетка/Проверка температуры', 'tempFormat = (t, f)->
    p = [0,10,100,1000,10000]
    t/p[f]
    
main =->
    
    # print "Проверка температуры"

    if message.error != ""
        message.dir = false
        message.device_state \'ERROR\'
        return false
    
    max_temp = 0
    cur_temp = 0
    if message.data[\'result\'] != ""
        r = SmartJs.hex2arr(message.data[\'result\'])
        max_temp = tempFormat((r[5] << 4) | r[6], 1)
        cur_temp = tempFormat((r[0] << 4) | r[1], 1)
        # print "current temp=",cur_temp
    
    return cur_temp<=max_temp', '2017-01-21 13:36:46', '2017-01-21 13:36:46', '', 'var main, tempFormat;

tempFormat = function(t, f) {
  var p;
  p = [0, 10, 100, 1000, 10000];
  return t / p[f];
};

main = function() {
  var cur_temp, max_temp, r;
  if (message.error !== "") {
    message.dir = false;
    message.device_state(\'ERROR\');
    return false;
  }
  max_temp = 0;
  cur_temp = 0;
  if (message.data[\'result\'] !== "") {
    r = SmartJs.hex2arr(message.data[\'result\']);
    max_temp = tempFormat((r[5] << 4) | r[6], 1);
    cur_temp = tempFormat((r[0] << 4) | r[1], 1);
  }
  return cur_temp <= max_temp;
};
' );
INSERT INTO `scripts`(`id`,`lang`,`name`,`source`,`created_at`,`update_at`,`description`,`compiled`) VALUES ( '4', 'coffeescript', 'Розетка/Включить', 'main =()->
    FUNCTION = 4
    DEVICE_ADDR = device.address || message.device.address
    if !DEVICE_ADDR?
        return
    
    # print "send msg from dev:", DEVICE_ADDR
    
    # send message to modbus channel
    request.baud = message.device.baud
    request.result = true
    request.device = message.device.tty
    request.timeout = 4000000
    request.stopBits = 2
    request.sleep = message.device.sleep
    request.command = [DEVICE_ADDR,FUNCTION,0,1,0,0]

    data = node_send \'modbus\', message.node, request
    if data.error != "" 
        # print "#{message.device.name} - error: #{data.error}"
        message.device_state \'ERROR\'
        return false
  
    if data.result != ""
        result = SmartJs.hex2arr(data.result)
        # проверка питания
        if result[2] == 1
            message.device_state \'ENABLED\'
        else if result[2] == 0
            message.device_state \'DISABLED\'
            
    message.device_state \'ENABLED\'
    return true', '2017-01-21 13:37:42', '2017-01-21 13:37:42', '', 'var main;

main = function() {
  var DEVICE_ADDR, FUNCTION, data, result;
  FUNCTION = 4;
  DEVICE_ADDR = device.address || message.device.address;
  if (DEVICE_ADDR == null) {
    return;
  }
  request.baud = message.device.baud;
  request.result = true;
  request.device = message.device.tty;
  request.timeout = 4000000;
  request.stopBits = 2;
  request.sleep = message.device.sleep;
  request.command = [DEVICE_ADDR, FUNCTION, 0, 1, 0, 0];
  data = node_send(\'modbus\', message.node, request);
  if (data.error !== "") {
    message.device_state(\'ERROR\');
    return false;
  }
  if (data.result !== "") {
    result = SmartJs.hex2arr(data.result);
    if (result[2] === 1) {
      message.device_state(\'ENABLED\');
    } else if (result[2] === 0) {
      message.device_state(\'DISABLED\');
    }
  }
  message.device_state(\'ENABLED\');
  return true;
};
' );
INSERT INTO `scripts`(`id`,`lang`,`name`,`source`,`created_at`,`update_at`,`description`,`compiled`) VALUES ( '5', 'coffeescript', 'Розетка/Выключить', 'main =()->
    
    FUNCTION = 4
    DEVICE_ADDR = device.address || message.device.address
    if !DEVICE_ADDR?
        return
    
    print "send msg from dev:", DEVICE_ADDR
    
    # send message to modbus channel
    request.baud = message.device.baud
    request.result = true
    request.device = message.device.tty
    request.timeout = 4000000
    request.stopBits = 2
    request.sleep = message.device.sleep
    request.command = [DEVICE_ADDR,FUNCTION,0,0,0,0]
    
    data = node_send \'modbus\', message.node, request
    if data.error != "" 
        print "#{message.device.name} - error: #{data.error}"
        message.device_state \'ERROR\'
        return false
    
    if data.result != ""
        result = SmartJs.hex2arr(data.result)
        # проверка питания
        if result[2] == 1
            message.device_state \'ENABLED\'
        else if result[2] == 0
            message.device_state \'DISABLED\'
    
    
    message.device_state \'DISABLED\'
    return true', '2017-01-21 13:38:34', '2017-01-21 13:42:46', '', 'var main;

main = function() {
  var DEVICE_ADDR, FUNCTION, data, result;
  FUNCTION = 4;
  DEVICE_ADDR = device.address || message.device.address;
  if (DEVICE_ADDR == null) {
    return;
  }
  print("send msg from dev:", DEVICE_ADDR);
  request.baud = message.device.baud;
  request.result = true;
  request.device = message.device.tty;
  request.timeout = 4000000;
  request.stopBits = 2;
  request.sleep = message.device.sleep;
  request.command = [DEVICE_ADDR, FUNCTION, 0, 0, 0, 0];
  data = node_send(\'modbus\', message.node, request);
  if (data.error !== "") {
    print(message.device.name + " - error: " + data.error);
    message.device_state(\'ERROR\');
    return false;
  }
  if (data.result !== "") {
    result = SmartJs.hex2arr(data.result);
    if (result[2] === 1) {
      message.device_state(\'ENABLED\');
    } else if (result[2] === 0) {
      message.device_state(\'DISABLED\');
    }
  }
  message.device_state(\'DISABLED\');
  return true;
};
' );
INSERT INTO `scripts`(`id`,`lang`,`name`,`source`,`created_at`,`update_at`,`description`,`compiled`) VALUES ( '6', 'coffeescript', 'Розетка/Инверт', 'main =()->
    
    FUNCTION = 4
    DEVICE_ADDR = device.address || message.device.address
    if !DEVICE_ADDR?
        return
     
    # print "send msg from dev:", DEVICE_ADDR
    
    q = 0
    result = SmartJs.hex2arr(message.data[\'result\'])
    if result[2] != 1
        q = 1
    
    # send message to modbus channel
    request.baud = message.device.baud
    request.result = true
    request.device = message.device.tty
    request.timeout = 1000000
    request.stopBits = 2
    request.sleep = message.device.sleep
    request.command = [DEVICE_ADDR,FUNCTION,0,q,0,0]
    
    # посыл команды выключения
    # print "переключаем"
    from_node = node_send \'modbus\', message.node, request
    if from_node.error != "" 
        # print "#{device.name} - error: #{from_node.error}"
        message.device_state \'ERROR\'
        return false
        
    # print "from_node", from_node.result
    
    # проверка вернувшихся данных
    result = SmartJs.hex2arr(from_node.result)
    if result[2] == q
        # print "dev #{DEVICE_ADDR} ok..."
        if q
            message.device_state \'ENABLED\'
        else
            message.device_state \'DISABLED\'
        return true
    
    message.device_state \'ERROR\'
    return false', '2017-01-21 14:05:54', '2017-01-21 14:05:54', '', 'var main;

main = function() {
  var DEVICE_ADDR, FUNCTION, from_node, q, result;
  FUNCTION = 4;
  DEVICE_ADDR = device.address || message.device.address;
  if (DEVICE_ADDR == null) {
    return;
  }
  q = 0;
  result = SmartJs.hex2arr(message.data[\'result\']);
  if (result[2] !== 1) {
    q = 1;
  }
  request.baud = message.device.baud;
  request.result = true;
  request.device = message.device.tty;
  request.timeout = 1000000;
  request.stopBits = 2;
  request.sleep = message.device.sleep;
  request.command = [DEVICE_ADDR, FUNCTION, 0, q, 0, 0];
  from_node = node_send(\'modbus\', message.node, request);
  if (from_node.error !== "") {
    message.device_state(\'ERROR\');
    return false;
  }
  result = SmartJs.hex2arr(from_node.result);
  if (result[2] === q) {
    if (q) {
      message.device_state(\'ENABLED\');
    } else {
      message.device_state(\'DISABLED\');
    }
    return true;
  }
  message.device_state(\'ERROR\');
  return false;
};
' );
-- ---------------------------------------------------------


-- Dump data of "user_metas" -------------------------------
INSERT INTO `user_metas`(`id`,`user_id`,`key`,`value`) VALUES ( '49', '1', 'phone1', '' );
INSERT INTO `user_metas`(`id`,`user_id`,`key`,`value`) VALUES ( '50', '1', 'phone2', '' );
INSERT INTO `user_metas`(`id`,`user_id`,`key`,`value`) VALUES ( '51', '1', 'phone3', '' );
INSERT INTO `user_metas`(`id`,`user_id`,`key`,`value`) VALUES ( '52', '1', 'autograph', '' );
INSERT INTO `user_metas`(`id`,`user_id`,`key`,`value`) VALUES ( '53', '2', 'phone1', '' );
INSERT INTO `user_metas`(`id`,`user_id`,`key`,`value`) VALUES ( '54', '2', 'phone2', '' );
INSERT INTO `user_metas`(`id`,`user_id`,`key`,`value`) VALUES ( '55', '2', 'phone3', '' );
INSERT INTO `user_metas`(`id`,`user_id`,`key`,`value`) VALUES ( '56', '2', 'autograph', '' );
-- ---------------------------------------------------------


-- Dump data of "users" ------------------------------------
INSERT INTO `users`(`id`,`nickname`,`first_name`,`last_name`,`encrypted_password`,`email`,`history`,`status`,`reset_password_token`,`authentication_token`,`image_id`,`sign_in_count`,`current_sign_in_ip`,`last_sign_in_ip`,`user_id`,`role_name`,`reset_password_sent_at`,`current_sign_in_at`,`last_sign_in_at`,`created_at`,`update_at`,`deleted`) VALUES ( '1', 'admin', '', '', 'f6fdffe48c908deb0f4c3bd36c032e72', 'admin@e154.ru', '[{"ip":"127.0.0.1","time":"2017-01-21T19:07:13.499096151+07:00"}]', 'active', '', '4j0lE0J46U94tAabPhWCuMKF6y7Plr8TMCzqThFDnkwlvcdvZd', NULL, '112', '127.0.0.1', '127.0.0.1', NULL, 'admin', NULL, '2017-01-21 12:07:13', '2017-01-21 11:11:26', '2017-01-15 05:25:07', '2017-01-21 12:07:13', NULL );
INSERT INTO `users`(`id`,`nickname`,`first_name`,`last_name`,`encrypted_password`,`email`,`history`,`status`,`reset_password_token`,`authentication_token`,`image_id`,`sign_in_count`,`current_sign_in_ip`,`last_sign_in_ip`,`user_id`,`role_name`,`reset_password_sent_at`,`current_sign_in_at`,`last_sign_in_at`,`created_at`,`update_at`,`deleted`) VALUES ( '2', 'demo', '', '', 'c514c91e4ed341f263e458d44b3bb0a7', 'demo@e154.ru', '[]', 'active', '', '5SLTHOzN1hWw6jhgEw0y9JbtwdBIK5mgW3DLt5FYy23zNkVnvW', NULL, '8', '127.0.0.1', '127.0.0.1', NULL, 'demo', NULL, '2017-01-21 11:11:43', '2017-01-20 17:28:23', '2017-01-18 17:13:28', '2017-01-21 11:11:43', NULL );
-- ---------------------------------------------------------


-- Dump data of "workers" ----------------------------------
INSERT INTO `workers`(`id`,`flow_id`,`workflow_id`,`name`,`device_action_id`,`time`,`created_at`,`update_at`,`status`) VALUES ( '1', '1', '1', 'Действие', '1', '0,5,10,15,20,25,30,35,40,45,50,55 * * * * *', '2017-01-21 14:02:12', '2017-01-21 14:08:18', 'enabled' );
-- ---------------------------------------------------------


-- Dump data of "workflows" --------------------------------
INSERT INTO `workflows`(`id`,`name`,`status`,`created_at`,`update_at`,`description`) VALUES ( '1', 'Спальня', 'enabled', '2017-01-21 14:00:00', '2017-01-21 14:00:16', '' );
-- ---------------------------------------------------------


-- CREATE INDEX "lnk_connections_flows" --------------------
CREATE INDEX `lnk_connections_flows` USING BTREE ON `connections`( `flow_id` );
-- ---------------------------------------------------------


-- CREATE INDEX "lnk_connections_flow_elements" ------------
CREATE INDEX `lnk_connections_flow_elements` USING BTREE ON `connections`( `element_from` );
-- ---------------------------------------------------------


-- CREATE INDEX "lnk_connections_flow_elements_2" ----------
CREATE INDEX `lnk_connections_flow_elements_2` USING BTREE ON `connections`( `element_to` );
-- ---------------------------------------------------------


-- CREATE INDEX "lnk_device_actions_devices" ---------------
CREATE INDEX `lnk_device_actions_devices` USING BTREE ON `device_actions`( `device_id` );
-- ---------------------------------------------------------


-- CREATE INDEX "lnk_scripts_device_actions" ---------------
CREATE INDEX `lnk_scripts_device_actions` USING BTREE ON `device_actions`( `script_id` );
-- ---------------------------------------------------------


-- CREATE INDEX "index_address" ----------------------------
CREATE INDEX `index_address` USING BTREE ON `devices`( `address` );
-- ---------------------------------------------------------


-- CREATE INDEX "lnk_email_template_email_template" --------
CREATE INDEX `lnk_email_template_email_template` USING BTREE ON `email_templates`( `parent` );
-- ---------------------------------------------------------


-- CREATE INDEX "lnk_flows_flow_elements" ------------------
CREATE INDEX `lnk_flows_flow_elements` USING BTREE ON `flow_elements`( `flow_link` );
-- ---------------------------------------------------------


-- CREATE INDEX "lnk_flow_elements_flows" ------------------
CREATE INDEX `lnk_flow_elements_flows` USING BTREE ON `flow_elements`( `flow_id` );
-- ---------------------------------------------------------


-- CREATE INDEX "lnk_scripts_flow_elements" ----------------
CREATE INDEX `lnk_scripts_flow_elements` USING BTREE ON `flow_elements`( `script_id` );
-- ---------------------------------------------------------


-- CREATE INDEX "lnk_flows_workflows" ----------------------
CREATE INDEX `lnk_flows_workflows` USING BTREE ON `flows`( `workflow_id` );
-- ---------------------------------------------------------


-- CREATE INDEX "lnk_device_actions_map_device_actions" ----
CREATE INDEX `lnk_device_actions_map_device_actions` USING BTREE ON `map_device_actions`( `device_action_id` );
-- ---------------------------------------------------------


-- CREATE INDEX "lnk_images_map_device_actions" ------------
CREATE INDEX `lnk_images_map_device_actions` USING BTREE ON `map_device_actions`( `image_id` );
-- ---------------------------------------------------------


-- CREATE INDEX "lnk_map_devices_map_device_actions" -------
CREATE INDEX `lnk_map_devices_map_device_actions` USING BTREE ON `map_device_actions`( `map_device_id` );
-- ---------------------------------------------------------


-- CREATE INDEX "lnk_device_states_map_device_states" ------
CREATE INDEX `lnk_device_states_map_device_states` USING BTREE ON `map_device_states`( `device_state_id` );
-- ---------------------------------------------------------


-- CREATE INDEX "lnk_images_map_device_states" -------------
CREATE INDEX `lnk_images_map_device_states` USING BTREE ON `map_device_states`( `image_id` );
-- ---------------------------------------------------------


-- CREATE INDEX "lnk_map_devices_map_device_states" --------
CREATE INDEX `lnk_map_devices_map_device_states` USING BTREE ON `map_device_states`( `map_device_id` );
-- ---------------------------------------------------------


-- CREATE INDEX "lnk_images_map_devices" -------------------
CREATE INDEX `lnk_images_map_devices` USING BTREE ON `map_devices`( `image_id` );
-- ---------------------------------------------------------


-- CREATE INDEX "lnk_maps_map_elements" --------------------
CREATE INDEX `lnk_maps_map_elements` USING BTREE ON `map_elements`( `map_id` );
-- ---------------------------------------------------------


-- CREATE INDEX "lnk_map_layers_map_elements" --------------
CREATE INDEX `lnk_map_layers_map_elements` USING BTREE ON `map_elements`( `layer_id` );
-- ---------------------------------------------------------


-- CREATE INDEX "lnk_images_map_images" --------------------
CREATE INDEX `lnk_images_map_images` USING BTREE ON `map_images`( `image_id` );
-- ---------------------------------------------------------


-- CREATE INDEX "lnk_maps_map_layers" ----------------------
CREATE INDEX `lnk_maps_map_layers` USING BTREE ON `map_layers`( `map_id` );
-- ---------------------------------------------------------


-- CREATE INDEX "lnk_roles_permissions" --------------------
CREATE INDEX `lnk_roles_permissions` USING BTREE ON `permissions`( `role_name` );
-- ---------------------------------------------------------


-- CREATE INDEX "lnk_roles_roles" --------------------------
CREATE INDEX `lnk_roles_roles` USING BTREE ON `roles`( `parent` );
-- ---------------------------------------------------------


-- CREATE INDEX "lnk_users_users_meta" ---------------------
CREATE INDEX `lnk_users_users_meta` USING BTREE ON `user_metas`( `user_id` );
-- ---------------------------------------------------------


-- CREATE INDEX "lnk_images_users" -------------------------
CREATE INDEX `lnk_images_users` USING BTREE ON `users`( `image_id` );
-- ---------------------------------------------------------


-- CREATE INDEX "lnk_roles_users" --------------------------
CREATE INDEX `lnk_roles_users` USING BTREE ON `users`( `role_name` );
-- ---------------------------------------------------------


-- CREATE INDEX "lnk_users_users" --------------------------
CREATE INDEX `lnk_users_users` USING BTREE ON `users`( `user_id` );
-- ---------------------------------------------------------


-- CREATE INDEX "lnk_workers_device_actions" ---------------
CREATE INDEX `lnk_workers_device_actions` USING BTREE ON `workers`( `device_action_id` );
-- ---------------------------------------------------------


-- CREATE INDEX "lnk_workers_flows" ------------------------
CREATE INDEX `lnk_workers_flows` USING BTREE ON `workers`( `flow_id` );
-- ---------------------------------------------------------


-- CREATE INDEX "lnk_workers_workflows" --------------------
CREATE INDEX `lnk_workers_workflows` USING BTREE ON `workers`( `workflow_id` );
-- ---------------------------------------------------------


-- CREATE LINK "lnk_connections_flow_elements_2" -----------
ALTER TABLE `connections`
	ADD CONSTRAINT `lnk_connections_flow_elements_2` FOREIGN KEY ( `element_to` )
	REFERENCES `flow_elements`( `uuid` )
	ON DELETE Cascade
	ON UPDATE Cascade;
-- ---------------------------------------------------------


-- CREATE LINK "lnk_connections_flows" ---------------------
ALTER TABLE `connections`
	ADD CONSTRAINT `lnk_connections_flows` FOREIGN KEY ( `flow_id` )
	REFERENCES `flows`( `id` )
	ON DELETE Cascade
	ON UPDATE Restrict;
-- ---------------------------------------------------------


-- CREATE LINK "lnk_connections_flow_elements" -------------
ALTER TABLE `connections`
	ADD CONSTRAINT `lnk_connections_flow_elements` FOREIGN KEY ( `element_from` )
	REFERENCES `flow_elements`( `uuid` )
	ON DELETE Cascade
	ON UPDATE Cascade;
-- ---------------------------------------------------------


-- CREATE LINK "lnk_scripts_device_actions" ----------------
ALTER TABLE `device_actions`
	ADD CONSTRAINT `lnk_scripts_device_actions` FOREIGN KEY ( `script_id` )
	REFERENCES `scripts`( `id` )
	ON DELETE Restrict
	ON UPDATE Cascade;
-- ---------------------------------------------------------


-- CREATE LINK "lnk_device_actions_devices" ----------------
ALTER TABLE `device_actions`
	ADD CONSTRAINT `lnk_device_actions_devices` FOREIGN KEY ( `device_id` )
	REFERENCES `devices`( `id` )
	ON DELETE Cascade
	ON UPDATE Restrict;
-- ---------------------------------------------------------


-- CREATE LINK "lnk_devices_device_states" -----------------
ALTER TABLE `device_states`
	ADD CONSTRAINT `lnk_devices_device_states` FOREIGN KEY ( `device_id` )
	REFERENCES `devices`( `id` )
	ON DELETE Cascade
	ON UPDATE Cascade;
-- ---------------------------------------------------------


-- CREATE LINK "lnk_devices_nodes" -------------------------
ALTER TABLE `devices`
	ADD CONSTRAINT `lnk_devices_nodes` FOREIGN KEY ( `node_id` )
	REFERENCES `nodes`( `id` )
	ON DELETE Cascade
	ON UPDATE Restrict;
-- ---------------------------------------------------------


-- CREATE LINK "lnk_devices_devices" -----------------------
ALTER TABLE `devices`
	ADD CONSTRAINT `lnk_devices_devices` FOREIGN KEY ( `device_id` )
	REFERENCES `devices`( `id` )
	ON DELETE Restrict
	ON UPDATE Cascade;
-- ---------------------------------------------------------


-- CREATE LINK "lnk_scripts_flow_elements" -----------------
ALTER TABLE `flow_elements`
	ADD CONSTRAINT `lnk_scripts_flow_elements` FOREIGN KEY ( `script_id` )
	REFERENCES `scripts`( `id` )
	ON DELETE Restrict
	ON UPDATE Cascade;
-- ---------------------------------------------------------


-- CREATE LINK "lnk_flows_flow_elements" -------------------
ALTER TABLE `flow_elements`
	ADD CONSTRAINT `lnk_flows_flow_elements` FOREIGN KEY ( `flow_link` )
	REFERENCES `flows`( `id` )
	ON DELETE Restrict
	ON UPDATE No Action;
-- ---------------------------------------------------------


-- CREATE LINK "lnk_flow_elements_flows" -------------------
ALTER TABLE `flow_elements`
	ADD CONSTRAINT `lnk_flow_elements_flows` FOREIGN KEY ( `flow_id` )
	REFERENCES `flows`( `id` )
	ON DELETE Cascade
	ON UPDATE Restrict;
-- ---------------------------------------------------------


-- CREATE LINK "lnk_flows_workflows" -----------------------
ALTER TABLE `flows`
	ADD CONSTRAINT `lnk_flows_workflows` FOREIGN KEY ( `workflow_id` )
	REFERENCES `workflows`( `id` )
	ON DELETE Restrict
	ON UPDATE Restrict;
-- ---------------------------------------------------------


-- CREATE LINK "lnk_map_devices_map_device_actions" --------
ALTER TABLE `map_device_actions`
	ADD CONSTRAINT `lnk_map_devices_map_device_actions` FOREIGN KEY ( `map_device_id` )
	REFERENCES `map_devices`( `id` )
	ON DELETE Cascade
	ON UPDATE Cascade;
-- ---------------------------------------------------------


-- CREATE LINK "lnk_device_actions_map_device_actions" -----
ALTER TABLE `map_device_actions`
	ADD CONSTRAINT `lnk_device_actions_map_device_actions` FOREIGN KEY ( `device_action_id` )
	REFERENCES `device_actions`( `id` )
	ON DELETE Cascade
	ON UPDATE Cascade;
-- ---------------------------------------------------------


-- CREATE LINK "lnk_images_map_device_actions" -------------
ALTER TABLE `map_device_actions`
	ADD CONSTRAINT `lnk_images_map_device_actions` FOREIGN KEY ( `image_id` )
	REFERENCES `images`( `id` )
	ON DELETE Restrict
	ON UPDATE Cascade;
-- ---------------------------------------------------------


-- CREATE LINK "lnk_map_devices_map_device_states" ---------
ALTER TABLE `map_device_states`
	ADD CONSTRAINT `lnk_map_devices_map_device_states` FOREIGN KEY ( `map_device_id` )
	REFERENCES `map_devices`( `id` )
	ON DELETE Cascade
	ON UPDATE Cascade;
-- ---------------------------------------------------------


-- CREATE LINK "lnk_device_states_map_device_states" -------
ALTER TABLE `map_device_states`
	ADD CONSTRAINT `lnk_device_states_map_device_states` FOREIGN KEY ( `device_state_id` )
	REFERENCES `device_states`( `id` )
	ON DELETE Restrict
	ON UPDATE Cascade;
-- ---------------------------------------------------------


-- CREATE LINK "lnk_images_map_device_states" --------------
ALTER TABLE `map_device_states`
	ADD CONSTRAINT `lnk_images_map_device_states` FOREIGN KEY ( `image_id` )
	REFERENCES `images`( `id` )
	ON DELETE Restrict
	ON UPDATE Cascade;
-- ---------------------------------------------------------


-- CREATE LINK "lnk_images_map_devices" --------------------
ALTER TABLE `map_devices`
	ADD CONSTRAINT `lnk_images_map_devices` FOREIGN KEY ( `image_id` )
	REFERENCES `images`( `id` )
	ON DELETE Restrict
	ON UPDATE Cascade;
-- ---------------------------------------------------------


-- CREATE LINK "lnk_devices_map_devices" -------------------
ALTER TABLE `map_devices`
	ADD CONSTRAINT `lnk_devices_map_devices` FOREIGN KEY ( `device_id` )
	REFERENCES `devices`( `id` )
	ON DELETE Restrict
	ON UPDATE Cascade;
-- ---------------------------------------------------------


-- CREATE LINK "lnk_map_layers_map_elements" ---------------
ALTER TABLE `map_elements`
	ADD CONSTRAINT `lnk_map_layers_map_elements` FOREIGN KEY ( `layer_id` )
	REFERENCES `map_layers`( `id` )
	ON DELETE Restrict
	ON UPDATE Cascade;
-- ---------------------------------------------------------


-- CREATE LINK "lnk_maps_map_elements" ---------------------
ALTER TABLE `map_elements`
	ADD CONSTRAINT `lnk_maps_map_elements` FOREIGN KEY ( `map_id` )
	REFERENCES `maps`( `id` )
	ON DELETE Restrict
	ON UPDATE Cascade;
-- ---------------------------------------------------------


-- CREATE LINK "lnk_images_map_images" ---------------------
ALTER TABLE `map_images`
	ADD CONSTRAINT `lnk_images_map_images` FOREIGN KEY ( `image_id` )
	REFERENCES `images`( `id` )
	ON DELETE Restrict
	ON UPDATE Cascade;
-- ---------------------------------------------------------


-- CREATE LINK "lnk_maps_map_layers" -----------------------
ALTER TABLE `map_layers`
	ADD CONSTRAINT `lnk_maps_map_layers` FOREIGN KEY ( `map_id` )
	REFERENCES `maps`( `id` )
	ON DELETE Restrict
	ON UPDATE Cascade;
-- ---------------------------------------------------------


-- CREATE LINK "lnk_roles_permissions" ---------------------
ALTER TABLE `permissions`
	ADD CONSTRAINT `lnk_roles_permissions` FOREIGN KEY ( `role_name` )
	REFERENCES `roles`( `name` )
	ON DELETE Cascade
	ON UPDATE Cascade;
-- ---------------------------------------------------------


-- CREATE LINK "lnk_roles_roles" ---------------------------
ALTER TABLE `roles`
	ADD CONSTRAINT `lnk_roles_roles` FOREIGN KEY ( `parent` )
	REFERENCES `roles`( `name` )
	ON DELETE Restrict
	ON UPDATE Cascade;
-- ---------------------------------------------------------


-- CREATE LINK "lnk_users_users_meta" ----------------------
ALTER TABLE `user_metas`
	ADD CONSTRAINT `lnk_users_users_meta` FOREIGN KEY ( `user_id` )
	REFERENCES `users`( `id` )
	ON DELETE Cascade
	ON UPDATE Cascade;
-- ---------------------------------------------------------


-- CREATE LINK "lnk_users_users" ---------------------------
ALTER TABLE `users`
	ADD CONSTRAINT `lnk_users_users` FOREIGN KEY ( `user_id` )
	REFERENCES `users`( `id` )
	ON DELETE Restrict
	ON UPDATE Cascade;
-- ---------------------------------------------------------


-- CREATE LINK "lnk_images_users" --------------------------
ALTER TABLE `users`
	ADD CONSTRAINT `lnk_images_users` FOREIGN KEY ( `image_id` )
	REFERENCES `images`( `id` )
	ON DELETE Restrict
	ON UPDATE Cascade;
-- ---------------------------------------------------------


-- CREATE LINK "lnk_roles_users" ---------------------------
ALTER TABLE `users`
	ADD CONSTRAINT `lnk_roles_users` FOREIGN KEY ( `role_name` )
	REFERENCES `roles`( `name` )
	ON DELETE Restrict
	ON UPDATE Cascade;
-- ---------------------------------------------------------


-- CREATE LINK "lnk_workers_workflows" ---------------------
ALTER TABLE `workers`
	ADD CONSTRAINT `lnk_workers_workflows` FOREIGN KEY ( `workflow_id` )
	REFERENCES `workflows`( `id` )
	ON DELETE Restrict
	ON UPDATE Restrict;
-- ---------------------------------------------------------


-- CREATE LINK "lnk_workers_device_actions" ----------------
ALTER TABLE `workers`
	ADD CONSTRAINT `lnk_workers_device_actions` FOREIGN KEY ( `device_action_id` )
	REFERENCES `device_actions`( `id` )
	ON DELETE Restrict
	ON UPDATE Restrict;
-- ---------------------------------------------------------


-- CREATE LINK "lnk_workers_flows" -------------------------
ALTER TABLE `workers`
	ADD CONSTRAINT `lnk_workers_flows` FOREIGN KEY ( `flow_id` )
	REFERENCES `flows`( `id` )
	ON DELETE Cascade
	ON UPDATE Restrict;
-- ---------------------------------------------------------


/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
-- ---------------------------------------------------------


