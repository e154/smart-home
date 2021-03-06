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
const DATA = `
{
  "ws": {
    "read": {
      "actions": [
        "/api/v1/ws",
        "/api/v1/ws/*"
      ],
      "method": "get",
      "description": "stream access"
    }
  },
  "node": {
    "read": {
      "actions": [
        "/api/v1/node",
        "/api/v1/node/[0-9]+"
      ],
      "method": "get",
      "description": ""
    },
    "create": {
      "actions": [
        "/api/v1/node"
      ],
      "method": "post",
      "description": ""
    },
    "update": {
      "actions": [
        "/api/v1/node/[0-9]+"
      ],
      "method": "put",
      "description": ""
    },
    "delete": {
      "actions": [
        "/api/v1/node/[0-9]+"
      ],
      "method": "delete",
      "description": ""
    }
  },
  "device": {
    "read": {
      "actions": [
        "/api/v1/device",
        "/api/v1/device/[0-9]+",
        "/api/v1/device/group",
        "/api/v1/device/[0-9]+/actions",
        "/api/v1/device/search",
        "/api/v1/device/[0-9]+/statuses"
      ],
      "method": "get",
      "description": ""
    },
    "create": {
      "actions": [
        "/api/v1/device"
      ],
      "method": "post",
      "description": ""
    },
    "update": {
      "actions": [
        "/api/v1/device/[0-9]+"
      ],
      "method": "put",
      "description": ""
    },
    "delete": {
      "actions": [
        "/api/v1/device/[0-9]+"
      ],
      "method": "delete",
      "description": ""
    }
  },
  "workflow": {
    "read": {
      "actions": [
        "/api/v1/workflow",
        "/api/v1/workflow/[0-9]+"
      ],
      "method": "get",
      "description": ""
    },
    "create": {
      "actions": [
        "/api/v1/workflow"
      ],
      "method": "post",
      "description": ""
    },
    "update": {
      "actions": [
        "/api/v1/workflow/[0-9]+",
        "/api/v1/workflow/[0-9]+/update_scenario"
      ],
      "method": "put",
      "description": ""
    },
    "delete": {
      "actions": [
        "/api/v1/workflow/[0-9]+"
      ],
      "method": "delete",
      "description": ""
    }
  },
  "flow": {
    "read": {
      "actions": [
        "/api/v1/flow",
        "/api/v1/flow/[0-9]+",
        "/api/v1/flow/[0-9]+/flow",
        "/api/v1/flow/[0-9]+/redactor",
        "/api/v1/flow/[0-9]+/workers",
        "/api/v1/flow/[0-9]+/search"
      ],
      "method": "get",
      "description": ""
    },
    "create": {
      "actions": [
        "/api/v1/flow"
      ],
      "method": "post",
      "description": ""
    },
    "update": {
      "actions": [
        "/api/v1/flow/[0-9]+",
        "/api/v1/flow/[0-9]+/redactor"
      ],
      "method": "put",
      "description": ""
    },
    "delete": {
      "actions": [
        "/api/v1/flow/[0-9]+"
      ],
      "method": "delete",
      "description": ""
    }
  },
  "device_action": {
    "read": {
      "actions": [
        "/api/v1/device_action",
        "/api/v1/device_action/[0-9]+",
        "/api/v1/device_action/search",
        "/api/v1/device_action/get_by_device/[0-9]+"
      ],
      "method": "get",
      "description": ""
    },
    "create": {
      "actions": [
        "/api/v1/device_action"
      ],
      "method": "post",
      "description": ""
    },
    "update": {
      "actions": [
        "/api/v1/device_action/[0-9]+"
      ],
      "method": "put",
      "description": ""
    },
    "delete": {
      "actions": [
        "/api/v1/device_action/[0-9]+"
      ],
      "method": "delete",
      "description": ""
    }
  },
  "worker": {
    "read": {
      "actions": [
        "/api/v1/worker",
        "/api/v1/worker/[0-9]+"
      ],
      "method": "get",
      "description": ""
    },
    "create": {
      "actions": [
        "/api/v1/worker"
      ],
      "method": "post",
      "description": ""
    },
    "update": {
      "actions": [
        "/api/v1/worker/[0-9]+"
      ],
      "method": "put",
      "description": ""
    },
    "delete": {
      "actions": [
        "/api/v1/worker/[0-9]+"
      ],
      "method": "delete",
      "description": ""
    }
  },
  "script": {
    "read": {
      "actions": [
        "/api/v1/script",
        "/api/v1/script/[0-9]+",
        "/api/v1/script/search"
      ],
      "method": "get",
      "description": ""
    },
    "create": {
      "actions": [
        "/api/v1/script"
      ],
      "method": "post",
      "description": ""
    },
    "copy": {
      "actions": [
        "/api/v1/script/[0-9+]/copy"
      ],
      "method": "post",
      "description": ""
    },
    "exec_script": {
      "actions": [
        "/api/v1/script/[0-9]+/exec"
      ],
      "method": "post",
      "description": "execute script"
    },
    "update": {
      "actions": [
        "/api/v1/script/[0-9]+"
      ],
      "method": "put",
      "description": ""
    },
    "delete": {
      "actions": [
        "/api/v1/script/[0-9]+"
      ],
      "method": "delete",
      "description": ""
    }
  },
  "log": {
    "read": {
      "actions": [
        "/api/v1/log",
        "/api/v1/log/[0-9]+"
      ],
      "method": "get",
      "description": ""
    }
  },
  "map": {
    "read_map": {
      "actions": [
        "/api/v1/map",
        "/api/v1/map/[0-9]+",
        "/api/v1/map/[0-9]+/full"
      ],
      "method": "get",
      "description": ""
    },
    "create_map": {
      "actions": [
        "/api/v1/map"
      ],
      "method": "post",
      "description": ""
    },
    "update_map": {
      "actions": [
        "/api/v1/map/[0-9]+"
      ],
      "method": "put",
      "description": ""
    },
    "delete_map": {
      "actions": [
        "/api/v1/map/[0-9]+"
      ],
      "method": "delete",
      "description": ""
    },
    "read_map_layer": {
      "actions": [
        "/api/v1/map_layer",
        "/api/v1/map_layer/[0-9]+"
      ],
      "method": "get",
      "description": ""
    },
    "create_map_layer": {
      "actions": [
        "/api/v1/map_layer"
      ],
      "method": "post",
      "description": ""
    },
    "update_map_layer": {
      "actions": [
        "/api/v1/map_layer/[0-9]+",
        "/api/v1/map_layer/sort"
      ],
      "method": "put",
      "description": ""
    },
    "delete_map_layer": {
      "actions": [
        "/api/v1/map_layer/[0-9]+"
      ],
      "method": "delete",
      "description": ""
    },
    "read_map_element": {
      "actions": [
        "/api/v1/map_element",
        "/api/v1/map_element/[0-9]+",
        "/api/v1/active_elements"
      ],
      "method": "get",
      "description": ""
    },
    "create_map_element": {
      "actions": [
        "/api/v1/map_element"
      ],
      "method": "post",
      "description": ""
    },
    "update_map_element": {
      "actions": [
        "/api/v1/map_element/[0-9]+",
        "/api/v1/map_element/[0-9]+/element_only",
        "/api/v1/map_element/sort"
      ],
      "method": "put",
      "description": ""
    },
    "delete_map_element": {
      "actions": [
        "/api/v1/map_element/[0-9]+"
      ],
      "method": "delete",
      "description": ""
    }
  },
  "device_state": {
    "read": {
      "actions": [
        "/api/v1/device_state",
        "/api/v1/device_state/[0-9]+",
        "/api/v1/device_state/get_by_device/[0-9]+"
      ],
      "method": "get",
      "description": ""
    },
    "create": {
      "actions": [
        "/api/v1/device_state"
      ],
      "method": "post",
      "description": ""
    },
    "update": {
      "actions": [
        "/api/v1/device_state/[0-9]+"
      ],
      "method": "put",
      "description": ""
    },
    "delete": {
      "actions": [
        "/api/v1/device_state/[0-9]+"
      ],
      "method": "delete",
      "description": ""
    }
  },
  "image": {
    "read": {
      "actions": [
        "/api/v1/image",
        "/api/v1/image/[0-9]+"
      ],
      "method": "get",
      "description": ""
    },
    "create": {
      "actions": [
        "/api/v1/image"
      ],
      "method": "post",
      "description": ""
    },
    "upload": {
      "actions": [
        "/api/v1/image/upload"
      ],
      "method": "post",
      "description": ""
    },
    "update": {
      "actions": [
        "/api/v1/image/[0-9]+"
      ],
      "method": "put",
      "description": ""
    },
    "delete": {
      "actions": [
        "/api/v1/image/[0-9]+"
      ],
      "method": "delete",
      "description": ""
    }
  },
  "dashboard": {
    "read": {
      "actions": [
        "/api/v1/dashboard",
        "/api/v1/dashboard/[0-9]+"
      ],
      "method": "get",
      "description": ""
    },
    "create": {
      "actions": [
        "/api/v1/dashboard"
      ],
      "method": "post",
      "description": ""
    },
    "update": {
      "actions": [
        "/api/v1/dashboard/[0-9]+"
      ],
      "method": "put",
      "description": ""
    },
    "delete": {
      "actions": [
        "/api/v1/dashboard/[0-9]+"
      ],
      "method": "delete",
      "description": ""
    }
  },
  "user": {
    "read": {
      "actions": [
        "/api/v1/user",
        "/api/v1/user/[0-9]+",
        "/api/v1/user/search"
      ],
      "method": "get",
      "description": ""
    },
    "create": {
      "actions": [
        "/api/v1/user"
      ],
      "method": "post",
      "description": ""
    },
    "update": {
      "actions": [
        "/api/v1/user/[0-9]+"
      ],
      "method": "put",
      "description": ""
    },
    "update_status": {
      "actions": [
        "/api/v1/user/[0-9]+/update_status"
      ],
      "method": "put",
      "description": ""
    },
    "delete": {
      "actions": [
        "/api/v1/user/[0-9]+"
      ],
      "method": "delete",
      "description": ""
    },
    "read_role": {
      "actions": [
        "/api/v1/role",
        "/api/v1/role/[\\w]+",
        "/api/v1/role/search"
      ],
      "method": "get",
      "description": ""
    },
    "create_role": {
      "actions": [
        "/api/v1/role",
        "/api/v1/role/[\\w]+"
      ],
      "method": "post",
      "description": ""
    },
    "update_role": {
      "actions": [
        "/api/v1/user/[\\w]+"
      ],
      "method": "put",
      "description": ""
    },
    "delete_role": {
      "actions": [
        "/api/v1/role/[\\w]+"
      ],
      "method": "delete",
      "description": ""
    },
    "read_role_access_list": {
      "actions": [
        "/api/v1/role/[\\w]+/access_list",
        "/api/v1/access_list"
      ],
      "method": "get",
      "description": "view role access list"
    },
    "update_role_access_list": {
      "actions": [
        "/api/v1/user/[\\w]+/access_list"
      ],
      "method": "put",
      "description": "update role access list info"
    }
  },
  "scenarios": {
    "read": {
      "actions": [
        "/api/v1/scenario",
        "/api/v1/scenario/[0-9]+"
      ],
      "method": "get",
      "description": ""
    },
    "create": {
      "actions": [
        "/api/v1/scenario"
      ],
      "method": "post",
      "description": ""
    },
    "update": {
      "actions": [
        "/api/v1/scenario/[0-9]+"
      ],
      "method": "put",
      "description": ""
    },
    "delete": {
      "actions": [
        "/api/v1/scenario/[0-9]+"
      ],
      "method": "delete",
      "description": ""
    },
    "read_script": {
      "actions": [
        "/api/v1/scenario_script/[0-9]+",
        "/api/v1/scenario_script"
      ],
      "method": "get",
      "description": ""
    },
    "create_script": {
      "actions": [
        "/api/v1/scenario_script"
      ],
      "method": "post",
      "description": ""
    },
    "update_script": {
      "actions": [
        "/api/v1/scenario_script/[0-9]+"
      ],
      "method": "put",
      "description": ""
    },
    "delete_script": {
      "actions": [
        "/api/v1/scenario_script/[0-9]+"
      ],
      "method": "delete",
      "description": ""
    }
  },
  "gate": {
    "read_gate": {
      "actions": [
        "/api/v1/gate",
        "/api/v1/gate/mobiles"
      ],
      "method": "get",
      "description": ""
    },
    "create_mobile": {
      "actions": [
        "/api/v1/gate/mobile"
      ],
      "method": "post",
      "description": ""
    },
    "update_gate": {
      "actions": [
        "/api/v1/gate"
      ],
      "method": "put",
      "description": ""
    },
    "delete_mobile": {
      "actions": [
        "/api/v1/gate/mobile/[\\w]+"
      ],
      "method": "delete",
      "description": ""
    }
  },
  "map_zone": {
    "read": {
      "actions": [
        "/api/v1/map_zone/[\\w]+"
      ],
      "method": "get",
      "description": ""
    },
    "create": {
      "actions": [
        "/api/v1/map_zone"
      ],
      "method": "post",
      "description": ""
    },
    "delete": {
      "actions": [
        "/api/v1/map_zone/[\\w]+"
      ],
      "method": "delete",
      "description": ""
    }
  },
  "template": {
    "read": {
      "actions": [
        "/api/v1/template/[\\w]+",
        "/api/v1/template_item/[\\w]+",
        "/api/v1/template_items/tree",
        "/api/v1/template_items",
        "/api/v1/templates"
      ],
      "method": "get",
      "description": ""
    },
    "preview": {
      "actions": [
        "/api/v1/templates/preview"
      ],
      "method": "post",
      "description": ""
    },
    "create": {
      "actions": [
        "/api/v1/template",
        "/api/v1/template_item"
      ],
      "method": "post",
      "description": ""
    },
    "update": {
      "actions": [
        "/api/v1/template",
        "/api/v1/template_item/[\\w]+",
        "/api/v1/template_items/status/[\\w]+",
        "/api/v1/template_items/tree"
      ],
      "method": "put",
      "description": ""
    },
    "delete": {
      "actions": [
        "/api/v1/template/[\\w]+",
        "/api/v1/template_item/[\\w]+"
      ],
      "method": "delete",
      "description": ""
    }
  },
  "notifr": {
    "read_config": {
      "actions": [
        "/api/v1/notifr/config"
      ],
      "method": "get",
      "description": ""
    },
    "update_config": {
      "actions": [
        "/api/v1/notifr/config"
      ],
      "method": "put",
      "description": ""
    },
    "create_notify": {
      "actions": [
        "/api/v1/notifr"
      ],
      "method": "post",
      "description": ""
    }
  },
  "mqtt": {
    "close_client": {
      "actions": [
        "/api/v1/mqtt/client/[\\w]+"
      ],
      "method": "delete",
      "description": ""
    },
    "close_topic": {
      "actions": [
        "/api/v1/mqtt/client/[\\w]+/topic"
      ],
      "method": "delete",
      "description": ""
    },
    "read": {
      "actions": [
        "/api/v1/mqtt/client/[\\w]+",
        "/api/v1/mqtt/client/[\\w]+/session",
        "/api/v1/mqtt/client/[\\w]+/subscriptions",
        "/api/v1/mqtt/clients",
        "/api/v1/mqtt/sessions"
      ],
      "method": "get",
      "description": ""
    },
    "publish": {
      "actions": [
        "/api/v1/mqtt/publish"
      ],
      "method": "post",
      "description": ""
    }
  },
  "zigbee2mqtt": {
    "read": {
      "actions": [
        "/api/v1/zigbee2mqtt/[0-9]+",
        "/api/v1/zigbee2mqtts",
        "/api/v1/zigbee2mqtt/[0-9]+/networkmap",
        "/api/v1/zigbee2mqtts/search_device"
      ],
      "method": "get",
      "description": ""
    },
    "create": {
      "actions": [
        "/api/v1/zigbee2mqtt",
        "/api/v1/zigbee2mqtt/[0-9]+/reset",
        "/api/v1/zigbee2mqtt/[0-9]+/device_ban",
        "/api/v1/zigbee2mqtt/[0-9]+/device_whitelist",
        "/api/v1/zigbee2mqtt/[0-9]+/update_networkmap"
      ],
      "method": "post",
      "description": ""
    },
    "update": {
      "actions": [
        "/api/v1/zigbee2mqtt/[0-9]+"
      ],
      "method": "patch",
      "description": ""
    },
    "delete": {
      "actions": [
        "/api/v1/zigbee2mqtt/[0-9]+"
      ],
      "method": "delete",
      "description": ""
    },
    "rename_device": {
      "actions": [
        "/api/v1/zigbee2mqtts/device_rename"
      ],
      "method": "patch",
      "description": ""
    }
  },
  "alexa": {
    "read": {
      "actions": [
        "/api/v1/alexa",
        "/api/v1/alexas/[0-9]+"
      ],
      "method": "get",
      "description": ""
    },
    "create": {
      "actions": [
        "/api/v1/alexa"
      ],
      "method": "post",
      "description": ""
    },
    "update": {
      "actions": [
        "/api/v1/alexa/[0-9]+"
      ],
      "method": "put",
      "description": ""
    },
    "delete": {
      "actions": [
        "/api/v1/alexa/[0-9]+"
      ],
      "method": "delete",
      "description": ""
    }
  },
  "metric": {
    "read": {
      "actions": [
        "/api/v1/metrics",
        "/api/v1/metrics/search",
        "/api/v1/metrics/[0-9]+"
      ],
      "method": "get",
      "description": ""
    },
    "create": {
      "actions": [
        "/api/v1/metric",
        "/api/v1/metric/[0-9]+/data"
      ],
      "method": "post",
      "description": ""
    },
    "update": {
      "actions": [
        "/api/v1/metric/[0-9]+"
      ],
      "method": "put",
      "description": ""
    },
    "delete": {
      "actions": [
        "/api/v1/metric/[0-9]+"
      ],
      "method": "delete",
      "description": ""
    }
  }
}
`
