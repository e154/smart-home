package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type Links_20170121_023351 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &Links_20170121_023351{}
	m.Created = "20170121_023351"
	migration.Register("Links_20170121_023351", m)
}

// Run the migrations
func (m *Links_20170121_023351) Up() {
	// use m.SQL("CREATE TABLE ...") to make schema update
	m.SQL("ALTER TABLE `connections` ADD CONSTRAINT `lnk_connections_flows` FOREIGN KEY ( `flow_id` ) REFERENCES `flows`( `id` ) ON DELETE Cascade ON UPDATE Restrict;")
	m.SQL("ALTER TABLE `connections` ADD CONSTRAINT `lnk_connections_flow_elements` FOREIGN KEY ( `element_from` ) REFERENCES `flow_elements`( `uuid` ) ON DELETE Cascade ON UPDATE Cascade;")
	m.SQL("ALTER TABLE `connections` ADD CONSTRAINT `lnk_connections_flow_elements_2` FOREIGN KEY ( `element_to` ) REFERENCES `flow_elements`( `uuid` ) ON DELETE Cascade ON UPDATE Cascade;")
	m.SQL("ALTER TABLE `device_actions` ADD CONSTRAINT `lnk_device_actions_devices` FOREIGN KEY ( `device_id` ) REFERENCES `devices`( `id` ) ON DELETE Cascade ON UPDATE Restrict;")
	m.SQL("ALTER TABLE `device_actions` ADD CONSTRAINT `lnk_scripts_device_actions` FOREIGN KEY ( `script_id` )REFERENCES `scripts`( `id` )ON DELETE Restrict ON UPDATE Cascade;")
	m.SQL("ALTER TABLE `device_states` ADD CONSTRAINT `lnk_devices_device_states` FOREIGN KEY ( `device_id` ) REFERENCES `devices`( `id` ) ON DELETE Cascade ON UPDATE Cascade;")
	m.SQL("ALTER TABLE `devices` ADD CONSTRAINT `lnk_devices_devices` FOREIGN KEY ( `device_id` ) REFERENCES `devices`( `id` ) ON DELETE Restrict ON UPDATE Cascade;")
	m.SQL("ALTER TABLE `devices` ADD CONSTRAINT `lnk_devices_nodes` FOREIGN KEY ( `node_id` ) REFERENCES `nodes`( `id` ) ON DELETE Cascade ON UPDATE Restrict;")
	m.SQL("ALTER TABLE `flow_elements` ADD CONSTRAINT `lnk_flows_flow_elements` FOREIGN KEY ( `flow_link` ) REFERENCES `flows`( `id` ) ON DELETE Restrict ON UPDATE No Action;")
	m.SQL("ALTER TABLE `flow_elements`ADD CONSTRAINT `lnk_flow_elements_flows` FOREIGN KEY ( `flow_id` ) REFERENCES `flows`( `id` ) ON DELETE Cascade ON UPDATE Restrict;")
	m.SQL("ALTER TABLE `flow_elements`ADD CONSTRAINT `lnk_scripts_flow_elements` FOREIGN KEY ( `script_id` )REFERENCES `scripts`( `id` )ON DELETE Restrict ON UPDATE Cascade;")
	m.SQL("ALTER TABLE `flows` ADD CONSTRAINT `lnk_flows_workflows` FOREIGN KEY ( `workflow_id` ) REFERENCES `workflows`( `id` ) ON DELETE Restrict ON UPDATE Restrict;")
	m.SQL("ALTER TABLE `map_device_actions` ADD CONSTRAINT `lnk_device_actions_map_device_actions` FOREIGN KEY ( `device_action_id` ) REFERENCES `device_actions`( `id` ) ON DELETE Cascade ON UPDATE Cascade;")
	m.SQL("ALTER TABLE `map_device_actions` ADD CONSTRAINT `lnk_images_map_device_actions` FOREIGN KEY ( `image_id` ) REFERENCES `images`( `id` ) ON DELETE Restrict ON UPDATE Cascade;")
	m.SQL("ALTER TABLE `map_device_actions` ADD CONSTRAINT `lnk_map_devices_map_device_actions` FOREIGN KEY ( `map_device_id` ) REFERENCES `map_devices`( `id` ) ON DELETE Cascade ON UPDATE Cascade;")
	m.SQL("ALTER TABLE `map_device_states` ADD CONSTRAINT `lnk_device_states_map_device_states` FOREIGN KEY ( `device_state_id` ) REFERENCES `device_states`( `id` ) ON DELETE Restrict ON UPDATE Cascade;")
	m.SQL("ALTER TABLE `map_device_states` ADD CONSTRAINT `lnk_images_map_device_states` FOREIGN KEY ( `image_id` ) REFERENCES `images`( `id` ) ON DELETE Restrict ON UPDATE Cascade;")
	m.SQL("ALTER TABLE `map_device_states` ADD CONSTRAINT `lnk_map_devices_map_device_states` FOREIGN KEY ( `map_device_id` ) REFERENCES `map_devices`( `id` ) ON DELETE Cascade ON UPDATE Cascade;")
	m.SQL("ALTER TABLE `map_devices` ADD CONSTRAINT `lnk_devices_map_devices` FOREIGN KEY ( `device_id` ) REFERENCES `devices`( `id` ) ON DELETE Restrict ON UPDATE Cascade;")
	m.SQL("ALTER TABLE `map_devices` ADD CONSTRAINT `lnk_images_map_devices` FOREIGN KEY ( `image_id` ) REFERENCES `images`( `id` ) ON DELETE Restrict ON UPDATE Cascade;")
	m.SQL("ALTER TABLE `map_elements` ADD CONSTRAINT `lnk_maps_map_elements` FOREIGN KEY ( `map_id` ) REFERENCES `maps`( `id` ) ON DELETE Restrict ON UPDATE Cascade;")
	m.SQL("ALTER TABLE `map_elements` ADD CONSTRAINT `lnk_map_layers_map_elements` FOREIGN KEY ( `layer_id` ) REFERENCES `map_layers`( `id` ) ON DELETE Restrict ON UPDATE Cascade;")
	m.SQL("ALTER TABLE `map_images` ADD CONSTRAINT `lnk_images_map_images` FOREIGN KEY ( `image_id` ) REFERENCES `images`( `id` ) ON DELETE Restrict ON UPDATE Cascade;")
	m.SQL("ALTER TABLE `map_layers` ADD CONSTRAINT `lnk_maps_map_layers` FOREIGN KEY ( `map_id` ) REFERENCES `maps`( `id` ) ON DELETE Restrict ON UPDATE Cascade;")
	m.SQL("ALTER TABLE `permissions` ADD CONSTRAINT `lnk_roles_permissions` FOREIGN KEY ( `role_name` ) REFERENCES `roles`( `name` ) ON DELETE Cascade ON UPDATE Cascade;")
	m.SQL("ALTER TABLE `roles` ADD CONSTRAINT `lnk_roles_roles` FOREIGN KEY ( `parent` ) REFERENCES `roles`( `name` ) ON DELETE Restrict ON UPDATE Cascade;")
	m.SQL("ALTER TABLE `user_metas` ADD CONSTRAINT `lnk_users_users_meta` FOREIGN KEY ( `user_id` ) REFERENCES `users`( `id` ) ON DELETE Cascade ON UPDATE Cascade;")
	m.SQL("ALTER TABLE `users` ADD CONSTRAINT `lnk_images_users` FOREIGN KEY ( `image_id` ) REFERENCES `images`( `id` ) ON DELETE Restrict ON UPDATE Cascade;")
	m.SQL("ALTER TABLE `users` ADD CONSTRAINT `lnk_roles_users` FOREIGN KEY ( `role_name` ) REFERENCES `roles`( `name` ) ON DELETE Restrict ON UPDATE Cascade;")
	m.SQL("ALTER TABLE `users` ADD CONSTRAINT `lnk_users_users` FOREIGN KEY ( `user_id` ) REFERENCES `users`( `id` ) ON DELETE Restrict ON UPDATE Cascade;")
	m.SQL("ALTER TABLE `workers` ADD CONSTRAINT `lnk_workers_device_actions` FOREIGN KEY ( `device_action_id` ) REFERENCES `device_actions`( `id` ) ON DELETE Restrict ON UPDATE Restrict;")
	m.SQL("ALTER TABLE `workers` ADD CONSTRAINT `lnk_workers_flows` FOREIGN KEY ( `flow_id` ) REFERENCES `flows`( `id` ) ON DELETE Cascade ON UPDATE Restrict;")
	m.SQL("ALTER TABLE `workers` ADD CONSTRAINT `lnk_workers_workflows` FOREIGN KEY ( `workflow_id` ) REFERENCES `workflows`( `id` ) ON DELETE Restrict ON UPDATE Restrict;")

}

// Reverse the migrations
func (m *Links_20170121_023351) Down() {
	// use m.SQL("DROP TABLE ...") to reverse schema update
	m.SQL("ALTER TABLE `connections` DROP FOREIGN KEY `lnk_connections_flows`;")
	m.SQL("ALTER TABLE `connections` DROP FOREIGN KEY `lnk_connections_flow_elements`;")
	m.SQL("ALTER TABLE `connections` DROP FOREIGN KEY `lnk_connections_flow_elements_2`;")
	m.SQL("ALTER TABLE `device_actions` DROP FOREIGN KEY `lnk_device_actions_devices`;")
	m.SQL("ALTER TABLE `device_actions` DROP FOREIGN KEY `lnk_scripts_device_actions`;")
	m.SQL("ALTER TABLE `device_states` DROP FOREIGN KEY `lnk_devices_device_states`;")
	m.SQL("ALTER TABLE `devices` DROP FOREIGN KEY `lnk_devices_devices`;")
	m.SQL("ALTER TABLE `devices` DROP FOREIGN KEY `lnk_devices_nodes`;")
	m.SQL("ALTER TABLE `flow_elements` DROP FOREIGN KEY `lnk_flows_flow_elements`;")
	m.SQL("ALTER TABLE `flow_elements` DROP FOREIGN KEY `lnk_flow_elements_flows`;")
	m.SQL("ALTER TABLE `flow_elements` DROP FOREIGN KEY `lnk_scripts_flow_elements`;")
	m.SQL("ALTER TABLE `flows` DROP FOREIGN KEY `lnk_flows_workflows`;")
	m.SQL("ALTER TABLE `map_device_actions` DROP FOREIGN KEY `lnk_device_actions_map_device_actions`;")
	m.SQL("ALTER TABLE `map_device_actions` DROP FOREIGN KEY `lnk_images_map_device_actions`;")
	m.SQL("ALTER TABLE `map_device_actions` DROP FOREIGN KEY `lnk_map_devices_map_device_actions`;")
	m.SQL("ALTER TABLE `map_device_states` DROP FOREIGN KEY `lnk_device_states_map_device_states`;")
	m.SQL("ALTER TABLE `map_device_states` DROP FOREIGN KEY `lnk_images_map_device_states`;")
	m.SQL("ALTER TABLE `map_device_states` DROP FOREIGN KEY `lnk_map_devices_map_device_states`;")
	m.SQL("ALTER TABLE `map_devices` DROP FOREIGN KEY `lnk_devices_map_devices`;")
	m.SQL("ALTER TABLE `map_devices` DROP FOREIGN KEY `lnk_images_map_devices`;")
	m.SQL("ALTER TABLE `map_elements` DROP FOREIGN KEY `lnk_maps_map_elements`;")
	m.SQL("ALTER TABLE `map_elements` DROP FOREIGN KEY `lnk_map_layers_map_elements`;")
	m.SQL("ALTER TABLE `map_images` DROP FOREIGN KEY `lnk_images_map_images`;")
	m.SQL("ALTER TABLE `map_layers` DROP FOREIGN KEY `lnk_maps_map_layers`;")
	m.SQL("ALTER TABLE `permissions` DROP FOREIGN KEY `lnk_roles_permissions`;")
	m.SQL("ALTER TABLE `roles` DROP FOREIGN KEY `lnk_roles_roles`;")
	m.SQL("ALTER TABLE `user_metas` DROP FOREIGN KEY `lnk_users_users_meta`;")
	m.SQL("ALTER TABLE `users` DROP FOREIGN KEY `lnk_images_users`;")
	m.SQL("ALTER TABLE `users` DROP FOREIGN KEY `lnk_roles_users`;")
	m.SQL("ALTER TABLE `users` DROP FOREIGN KEY `lnk_users_users`;")
	m.SQL("ALTER TABLE `workers` DROP FOREIGN KEY `lnk_workers_device_actions`;")
	m.SQL("ALTER TABLE `workers` DROP FOREIGN KEY `lnk_workers_flows`;")
	m.SQL("ALTER TABLE `workers` DROP FOREIGN KEY `lnk_workers_workflows`;")
}
