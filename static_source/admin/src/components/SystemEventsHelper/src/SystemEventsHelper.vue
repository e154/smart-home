<script setup lang="ts">
</script>

<template>
  <div class="tip custom-block w-[100%]">
    <p>System events list</p>

    event_task_completed
    event_created_task_model
    event_removed_task_model
    event_updated_task_model
    event_task_loaded
    event_task_unloaded
    event_trigger_completed
    event_call_trigger
    event_created_trigger_model
    event_removed_trigger_model
    event_updated_trigger_model
    event_trigger_loaded
    event_trigger_unloaded
    event_action_completed
    event_call_task_action
    event_call_action
    event_added_action_model
    event_removed_action_model
    event_updated_action_model
    event_action_loaded
    event_action_unloaded
    event_added_condition_model
    event_removed_condition_model
    event_updated_condition_model

    event_created_backup
    event_removed_backup
    event_uploaded_backup
    event_started_restore

    event_state_changed
    event_last_state_changed
    event_get_last_state
    event_call_entity_action
    event_call_scene
    event_added_actor
    event_created_entity_model
    event_updated_entity_model
    event_updated_metric
    event_entity_unloaded
    event_entity_loaded
    event_entity_set_state

    event_mqtt_new_client

    event_plugin_loaded
    event_plugin_unloaded

    event_created_script_model
    event_updated_script_model
    event_removed_script_model

    event_service_started
    event_service_stopped
    event_service_restarted

    event_update_user_location
    event_direct_message
    event_user_signed_in

    event_updated_variable_model
    event_removed_variable_model

    eventServerVersion
    eventGetServerVersion

    event_created_zigbee2mqtt_model
    event_updated_zigbee2mqtt_model
    event_removed_zigbee2mqtt_model


  </div>
</template>

<style lang="less" scoped>
.light {
  .custom-block.tip {
    padding: 8px 16px;
    background-color: #409eff1a;
    border-radius: 4px;
    border-left: 5px solid var(--el-color-primary);
  }
}

.dark {
  .custom-block.tip {
    padding: 8px 16px;
    background-color: #409eff1a;
    border-radius: 4px;
    border-left: 5px solid var(--el-color-primary);
  }
}

.dYaBhr thead th {
  border: 1px solid ;
  padding: 0 10px;
}

.dYaBhr tbody td {
  border: 1px solid ;
  padding: 0 10px;
}

</style>
