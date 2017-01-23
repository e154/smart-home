package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type Index_20170121_021812 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &Index_20170121_021812{}
	m.Created = "20170121_021812"
	migration.Register("Index_20170121_021812", m)
}

// Run the migrations
func (m *Index_20170121_021812) Up() {
	// use m.SQL("CREATE TABLE ...") to make schema update
	m.SQL("CREATE INDEX `lnk_connections_flows` USING BTREE ON `connections`( `flow_id` );")
	m.SQL("CREATE INDEX `lnk_connections_flow_elements` USING BTREE ON `connections`( `element_from` );")
	m.SQL("CREATE INDEX `lnk_connections_flow_elements_2` USING BTREE ON `connections`( `element_to` );")
	m.SQL("CREATE INDEX `lnk_device_actions_devices` USING BTREE ON `device_actions`( `device_id` );")
	m.SQL("CREATE INDEX `lnk_scripts_device_actions` USING BTREE ON `device_actions`( `script_id` );")
	m.SQL("CREATE INDEX `index_address` USING BTREE ON `devices`( `address` );")
	m.SQL("CREATE INDEX `lnk_email_template_email_template` USING BTREE ON `email_templates`( `parent` );")
	m.SQL("CREATE INDEX `lnk_flows_flow_elements` USING BTREE ON `flow_elements`( `flow_link` );")
	m.SQL("CREATE INDEX `lnk_flow_elements_flows` USING BTREE ON `flow_elements`( `flow_id` );")
	m.SQL("CREATE INDEX `lnk_scripts_flow_elements` USING BTREE ON `flow_elements`( `script_id` );")
	m.SQL("CREATE INDEX `lnk_flows_workflows` USING BTREE ON `flows`( `workflow_id` );")
	m.SQL("CREATE INDEX `lnk_device_actions_map_device_actions` USING BTREE ON `map_device_actions`( `device_action_id` );")
	m.SQL("CREATE INDEX `lnk_images_map_device_actions` USING BTREE ON `map_device_actions`( `image_id` );")
	m.SQL("CREATE INDEX `lnk_map_devices_map_device_actions` USING BTREE ON `map_device_actions`( `map_device_id` );")
	m.SQL("CREATE INDEX `lnk_device_states_map_device_states` USING BTREE ON `map_device_states`( `device_state_id` );")
	m.SQL("CREATE INDEX `lnk_images_map_device_states` USING BTREE ON `map_device_states`( `image_id` );")
	m.SQL("CREATE INDEX `lnk_map_devices_map_device_states` USING BTREE ON `map_device_states`( `map_device_id` );")
	m.SQL("CREATE INDEX `lnk_images_map_devices` USING BTREE ON `map_devices`( `image_id` );")
	m.SQL("CREATE INDEX `lnk_maps_map_elements` USING BTREE ON `map_elements`( `map_id` );")
	m.SQL("CREATE INDEX `lnk_map_layers_map_elements` USING BTREE ON `map_elements`( `layer_id` );")
	m.SQL("CREATE INDEX `lnk_images_map_images` USING BTREE ON `map_images`( `image_id` );")
	m.SQL("CREATE INDEX `lnk_maps_map_layers` USING BTREE ON `map_layers`( `map_id` );")
	m.SQL("CREATE INDEX `lnk_roles_permissions` USING BTREE ON `permissions`( `role_name` );")
	m.SQL("CREATE INDEX `lnk_roles_roles` USING BTREE ON `roles`( `parent` );")
	m.SQL("CREATE INDEX `lnk_users_users_meta` USING BTREE ON `user_metas`( `user_id` );")
	m.SQL("CREATE INDEX `lnk_images_users` USING BTREE ON `users`( `image_id` );")
	m.SQL("CREATE INDEX `lnk_roles_users` USING BTREE ON `users`( `role_name` );")
	m.SQL("CREATE INDEX `lnk_users_users` USING BTREE ON `users`( `user_id` );")
	m.SQL("CREATE INDEX `lnk_workers_device_actions` USING BTREE ON `workers`( `device_action_id` );")
	m.SQL("CREATE INDEX `lnk_workers_flows` USING BTREE ON `workers`( `flow_id` );")
	m.SQL("CREATE INDEX `lnk_workers_workflows` USING BTREE ON `workers`( `workflow_id` );")
}

// Reverse the migrations
func (m *Index_20170121_021812) Down() {
	// use m.SQL("DROP TABLE ...") to reverse schema update
	m.SQL("DROP INDEX `lnk_connections_flows` ON `connections`;")
	m.SQL("DROP INDEX `lnk_connections_flow_elements` ON `connections`;")
	m.SQL("DROP INDEX `lnk_connections_flow_elements_2` ON `connections`;")
	m.SQL("DROP INDEX `lnk_scripts_device_actions` ON `device_actions`;")
	m.SQL("DROP INDEX `lnk_device_actions_devices` ON `device_actions`;")
	m.SQL("DROP INDEX `index_address` ON `devices`;")
	m.SQL("DROP INDEX `lnk_email_template_email_template` ON `email_templates`;")
	m.SQL("DROP INDEX `lnk_flows_flow_elements` ON `flow_elements`;")
	m.SQL("DROP INDEX `lnk_flow_elements_flows` ON `flow_elements`;")
	m.SQL("DROP INDEX `lnk_scripts_flow_elements` ON `flow_elements`;")
	m.SQL("DROP INDEX `lnk_flows_workflows` ON `flows`;")
	m.SQL("DROP INDEX `lnk_device_actions_map_device_actions` ON `map_device_actions`;")
	m.SQL("DROP INDEX `lnk_images_map_device_actions` ON `map_device_actions`;")
	m.SQL("DROP INDEX `lnk_map_devices_map_device_actions` ON `map_device_actions`;")
	m.SQL("DROP INDEX `lnk_device_states_map_device_states` ON `map_device_states`;")
	m.SQL("DROP INDEX `lnk_images_map_device_states` ON `map_device_states`;")
	m.SQL("DROP INDEX `lnk_map_devices_map_device_states` ON `map_device_states`;")
	m.SQL("DROP INDEX `lnk_images_map_devices` ON `map_devices`;")
	m.SQL("DROP INDEX `lnk_maps_map_elements` ON `map_elements`;")
	m.SQL("DROP INDEX `lnk_map_layers_map_elements` ON `map_elements`;")
	m.SQL("DROP INDEX `lnk_images_map_images` ON `map_images`;")
	m.SQL("DROP INDEX `lnk_maps_map_layers` ON `map_layers`;")
	m.SQL("DROP INDEX `lnk_roles_permissions` ON `permissions`;")
	m.SQL("DROP INDEX `lnk_roles_roles` ON `roles`;")
	m.SQL("DROP INDEX `lnk_users_users_meta` ON `user_metas`;")
	m.SQL("DROP INDEX `lnk_images_users` ON `users`;")
	m.SQL("DROP INDEX `lnk_roles_users` ON `users`;")
	m.SQL("DROP INDEX `lnk_users_users` ON `users`;")
	m.SQL("DROP INDEX `lnk_workers_device_actions` ON `workers`;")
	m.SQL("DROP INDEX `lnk_workers_flows` ON `workers`;")
	m.SQL("DROP INDEX `lnk_workers_workflows` ON `workers`;")
}
