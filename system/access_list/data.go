// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2016-2021, Filippov Alex
//
// This library is free software: you can redistribute it and/or
// modify it under the terms of the GNU Lesser General Public
// License as published by the Free Software Foundation; either
// version 3 of the License, or (at your option) any later version.
//
// This library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU
// Library General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public
// License along with this library.  If not, see
// <https://www.gnu.org/licenses/>.

package access_list

// DATA ...
// todo fix
const DATA = `
{
  "area": {
    "create": {
      "actions": [
        "/api.AreaService/AddArea"
      ],
      "description": ""
    },
    "read": {
      "actions": [
        "/api.AreaService/GetAreaById",
        "/api.AreaService/GetAreaList"
      ],
      "description": ""
    },
    "update": {
      "actions": [
        "/api.AreaService/UpdateArea"
      ],
      "description": ""
    },
    "delete": {
      "actions": [
        "/api.AreaService/DeleteArea"
      ],
      "description": ""
    },
    "search": {
      "actions": [
        "/api.AreaService/SearchArea"
      ],
      "description": ""
    }
  },
  "auth": {
    "read_access_list": {
      "actions": [
        "/api.AuthService/AccessList"
      ],
      "description": ""
    },
    "signin": {
      "actions": [
        "/api.AuthService/Signin"
      ],
      "description": ""
    },
    "signout": {
      "actions": [
        "/api.AuthService/Signout"
      ],
      "description": ""
    }
  },
  "automation": {
    "create": {
      "actions": [
        "/api.AutomationService/AddTask"
      ],
      "description": ""
    },
    "read": {
      "actions": [
        "/api.AutomationService/GetTask",
        "/api.AutomationService/GetTaskList"
      ],
      "description": ""
    },
    "update": {
      "actions": [
        "/api.AutomationService/UpdateTask",
        "/api.AutomationService/DisableTask",
        "/api.AutomationService/EnableTask"
      ],
      "description": ""
    },
    "delete": {
      "actions": [
        "/api.AutomationService/DeleteTask"
      ],
      "description": ""
    }
  },
  "backup": {
    "create": {
      "actions": [
        "/api.BackupService/NewBackup"
      ],
      "description": ""
    },
    "restore": {
      "actions": [
        "/api.BackupService/RestoreBackup"
      ],
      "description": ""
    },
    "read": {
      "actions": [
        "/api.BackupService/GetBackupList"
      ],
      "description": ""
    }
  },
  "dashboard": {
    "create": {
      "actions": [
        "/api.DashboardService/AddDashboard",
        "/api.DashboardService/ImportDashboard"
      ],
      "description": ""
    },
    "read": {
      "actions": [
        "/api.DashboardService/GetDashboardById",
        "/api.DashboardService/GetDashboardList"
      ],
      "description": ""
    },
    "update": {
      "actions": [
        "/api.DashboardService/UpdateDashboard"
      ],
      "description": ""
    },
    "delete": {
      "actions": [
        "/api.DashboardService/DeleteDashboard"
      ],
      "description": ""
    },
    "search": {
      "actions": [
        "/api.DashboardService/SearchDashboard"
      ],
      "description": ""
    }
  },
  "dashboard_card": {
    "create": {
      "actions": [
        "/api.DashboardCardService/AddDashboardCard",
        "/api.DashboardCardService/ImportDashboardCard"
      ],
      "description": ""
    },
    "read": {
      "actions": [
        "/api.DashboardCardService/GetDashboardCardById",
        "/api.DashboardCardService/GetDashboardCardList"
      ],
      "description": ""
    },
    "update": {
      "actions": [
        "/api.DashboardCardService/UpdateDashboardCard"
      ],
      "description": ""
    },
    "delete": {
      "actions": [
        "/api.DashboardCardService/DeleteDashboardCard"
      ],
      "description": ""
    }
  },
  "dashboard_card_item": {
    "create": {
      "actions": [
        "/api.DashboardCardItemService/AddDashboardCardItem"
      ],
      "description": ""
    },
    "read": {
      "actions": [
        "/api.DashboardCardItemService/GetDashboardCardItemById",
        "/api.DashboardCardItemService/GetDashboardCardItemList"
      ],
      "description": ""
    },
    "update": {
      "actions": [
        "/api.DashboardCardItemService/UpdateDashboardCardItem"
      ],
      "description": ""
    },
    "delete": {
      "actions": [
        "/api.DashboardCardItemService/DeleteDashboardCardItem"
      ],
      "description": ""
    }
  },
  "dashboard_tab": {
    "create": {
      "actions": [
        "/api.DashboardTabService/AddDashboardTab"
      ],
      "description": ""
    },
    "read": {
      "actions": [
        "/api.DashboardTabService/GetDashboardTabById",
        "/api.DashboardTabService/GetDashboardTabList"
      ],
      "description": ""
    },
    "update": {
      "actions": [
        " /api.DashboardTabService/UpdateDashboardTab"
      ],
      "description": ""
    },
    "delete": {
      "actions": [
        "/api.DashboardTabService/DeleteDashboardTab"
      ],
      "description": ""
    }
  },
  "developer_tools": {
    "entity_reload": {
      "actions": [
        "/api.DeveloperToolsService/ReloadEntity"
      ],
      "description": ""
    },
    "entity_set_state": {
      "actions": [
        "/api.DeveloperToolsService/EntitySetState"
      ],
      "description": ""
    },
    "task_call_action": {
      "actions": [
        "/api.DeveloperToolsService/TaskCallAction"
      ],
      "description": ""
    },
    "task_call_trigger": {
      "actions": [
        "/api.DeveloperToolsService/TaskCallTrigger"
      ],
      "description": ""
    }
  },
  "entity": {
    "create": {
      "actions": [
        "/api.EntityService/AddEntity",
        "/api.EntityService/ImportEntity"
      ],
      "description": ""
    },
    "read": {
      "actions": [
        "/api.EntityService/GetEntity",
        "/api.EntityService/GetEntityList"
      ],
      "description": ""
    },
    "update": {
      "actions": [
        "/api.EntityService/UpdateEntity"
      ],
      "description": ""
    },
    "delete": {
      "actions": [
        "/api.EntityService/DeleteEntity"
      ],
      "description": ""
    },
    "search": {
      "actions": [
        "/api.EntityService/SearchEntity"
      ],
      "description": ""
    }
  },
  "entity_storage": {
    "read": {
      "actions": [
        "/api.EntityStorageService/GetEntityStorageList"
      ],
      "description": ""
    }
  },
  "image": {
    "create": {
      "actions": [
        "/api.ImageService/AddImage",
        "/api.ImageService/UploadImage"
      ],
      "description": ""
    },
    "read": {
      "actions": [
        "/api.ImageService/GetImageById",
        "/api.ImageService/GetImageList",
        "/api.ImageService/GetImageListByDate",
        "/api.ImageService/GetImageFilterList"
      ],
      "description": ""
    },
    "update": {
      "actions": [
        "/api.ImageService/UpdateImageById"
      ],
      "description": ""
    },
    "delete": {
      "actions": [
        "/api.ImageService/DeleteImageById"
      ],
      "description": ""
    }
  },
  "interact": {
    "create": {
      "actions": [
        "/api.InteractService/EntityCallAction"
      ],
      "description": ""
    }
  },
  "log": {
    "read": {
      "actions": [
        "/api.LogService/GetLogList"
      ],
      "description": ""
    }
  },
  "metric": {
    "read": {
      "actions": [
        "/api.MetricService/GetMetric"
      ],
      "description": ""
    }
  },
  "plugin": {
    "read": {
      "actions": [
        "/api.PluginService/GetPluginList",
        "/api.PluginService/GetPluginOptions"
      ],
      "description": ""
    },
    "update": {
      "actions": [
        "/api.PluginService/DisablePlugin",
        "/api.PluginService/EnablePlugin"
      ],
      "description": ""
    },
    "search": {
      "actions": [
        "/api.PluginService/SearchPlugin,"
      ],
      "description": ""
    }
  },
  "role": {
    "create": {
      "actions": [
        "/api.RoleService/AddRole"
      ],
      "description": ""
    },
    "read": {
      "actions": [
        "/api.RoleService/GetRoleList",
        "/api.RoleService/GetRoleAccessList",
        "/api.RoleService/GetRoleByName"
      ],
      "description": ""
    },
    "update": {
      "actions": [
        "/api.RoleService/UpdateRoleAccessList",
        "/api.RoleService/UpdateRoleByName"
      ],
      "description": ""
    },
    "delete": {
      "actions": [
        "/api.RoleService/DeleteRoleByName"
      ],
      "description": ""
    },
    "search": {
      "actions": [
        "/api.RoleService/SearchRoleByName"
      ],
      "description": ""
    }
  },
  "script": {
    "create": {
      "actions": [
        "/api.ScriptService/AddScript",
        "/api.ScriptService/ExecScriptById",
        "/api.ScriptService/ExecSrcScriptById",
        "/api.ScriptService/CopyScriptById"
      ],
      "description": ""
    },
    "read": {
      "actions": [
        "/api.ScriptService/GetScriptList",
        "/api.ScriptService/GetScriptById"
      ],
      "description": ""
    },
    "update": {
      "actions": [
        "/api.ScriptService/UpdateScriptById"
      ],
      "description": ""
    },
    "delete": {
      "actions": [
        "/api.ScriptService/DeleteScriptById"
      ],
      "description": ""
    },
    "search": {
      "actions": [
        "/api.ScriptService/SearchScript"
      ],
      "description": ""
    }
  },
  "stream": {
    "read": {
      "actions": [
        "/api.StreamService/Subscribe"
      ],
      "description": ""
    }
  },
  "user": {
    "create": {
      "actions": [
        "/api.UserService/AddUserRequest"
      ],
      "description": ""
    },
    "read": {
      "actions": [
        "/api.UserService/GetUserList",
        "/api.UserService/GetUserById"
      ],
      "description": ""
    },
    "update": {
      "actions": [
        "/api.UserService/UpdateUserById"
      ],
      "description": ""
    },
    "delete": {
      "actions": [
        "/api.UserService/DeleteUserById"
      ],
      "description": ""
    }
  },
  "variable": {
    "create": {
      "actions": [
        "/api.VariableService/AddVariable"
      ],
      "description": ""
    },
    "read": {
      "actions": [
        "/api.VariableService/GetVariableByName",
        "/api.VariableService/GetVariableList"
      ],
      "description": ""
    },
    "update": {
      "actions": [
        "/api.VariableService/UpdateVariable"
      ],
      "description": ""
    },
    "delete": {
      "actions": [
        "/api.VariableService/DeleteVariable"
      ],
      "description": ""
    }
  },
  "zigbee2mqtt": {
    "create": {
      "actions": [
        "/api.Zigbee2mqttService/AddZigbee2mqttBridge",
        "/api.Zigbee2mqttService/DeviceWhitelist",
        "/api.Zigbee2mqttService/DeviceRename",
        "/api.Zigbee2mqttService/DeviceBan"
      ],
      "description": ""
    },
    "read": {
      "actions": [
        "/api.Zigbee2mqttService/GetBridgeList",
        "/api.Zigbee2mqttService/GetZigbee2mqttBridge",
        "/api.Zigbee2mqttService/DeviceList",
        "/api.Zigbee2mqttService/Networkmap"
      ],
      "description": ""
    },
    "update": {
      "actions": [
        "/api.Zigbee2mqttService/UpdateBridgeById",
        "/api.Zigbee2mqttService/UpdateNetworkmap",
        "/api.Zigbee2mqttService/ResetBridgeById"
      ],
      "description": ""
    },
    "delete": {
      "actions": [
        "/api.Zigbee2mqttService/DeleteBridgeById"
      ],
      "description": ""
    },
    "search": {
      "actions": [
        "/api.Zigbee2mqttService/SearchDevice"
      ],
      "description": ""
    }
  }
}
`
