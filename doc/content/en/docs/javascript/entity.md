---
title: "Устройство"
linkTitle: "entity"
date: 2021-11-20
description: >

---

In the **Smart Home** project, the "Entity" object is the main unit and the primary object that stores the device's state and settings.

An "Entity" represents an abstraction of a specific device in a smart home. It contains information about the device type, its identifier, current state, and other settings.

The key attributes of an "Entity" include:

1. Device Type: Determines the category or class of the device, such as "light fixture," "thermostat," "door lock," etc. This helps classify devices and define the operations and settings available to them.

2. Identifier: A unique identifier assigned to the device, which allows it to be uniquely identified within the smart home system. This can be a numerical value, string, or other identification format.

3. Device State: Contains information about the current state of the device, such as "on/off," "open/closed," "temperature," etc. The state can change as a result of operations or events related to the device.

4. Device Settings: Includes various parameters and settings related to the specific device. For example, these could be brightness settings for a light fixture, preferred temperature for a thermostat, etc.

The "Entity" serves as the main device abstraction in the **Smart Home** project and provides a unified interface for working with different types of devices. This allows convenient management of device states and settings, as well as the creation of flexible automation scenarios and interactions between devices in your smart home.

### Entity Object Properties
```coffeescript
{
  "id": "sensor.device1",
  "type": "sensor",
  "name": "",
  "description": "device description",
  "icon": null,
  "image_url": null,
  "actions": [
    {
      "name": "ENABLE",
      "description": "enable",
      "image_url": null,
      "icon": null
    },
    {
      "name": "DISABLE",
      "description": "disable",
      "image_url": null,
      "icon": null
    }
  ],
  "states": [
    {
      "name": "ENABLED",
      "description": "enabled state",
      "image_url": null,
      "icon": null
    },
    {
      "name": "DISABLED",
      "description": "disabled state",
      "image_url": null,
      "icon": null
    }
  ],
  "state": {
    "name": "ENABLED",
    "description": "enabled state",
    "image_url": null,
    "icon": null
  },
  "attributes": {},
  "area": null,
  "metrics": [],
  "hidden": false
}
```

The abstract "Entity" object in the **Smart Home** project provides the following methods:

1. `setState(state)`: This method allows you to set the device state. You pass the new state as the `state` argument, which will be applied to the device. For example:

```javascript
entity.setState('on');
```

2. `setAttributes(attributes)`: This method is used to set the device attributes. You pass an `attributes` object as an argument, which contains the new attributes to be applied to the device. Example usage:

```javascript
const newAttributes = {
  brightness: 80,
  color: 'red',
};

entity.setAttributes(newAttributes);
```

3. `getAttributes()`: This method allows you to retrieve the current device attributes. It returns an object with the current attribute values. Example usage:

```javascript
const attributes = entity.getAttributes();
console.log(attributes);
```

4. `setMetric(metric, value)`: This method allows you to set the device metric. You pass the

metric name as the `metric` argument and the metric value as the `value` argument. Example usage:

```javascript
entity.setMetric('temperature', 25);
```

5. `callAction(action, params)`: This method allows you to invoke a specific action on the device. You pass the action name as the `action` argument and any required parameters as the `params` argument. Example usage:

```javascript
entity.callAction('turnOn');
```

The methods of the "Entity" object provide convenient capabilities for managing the device's state, attributes, metrics, and executing actions on the device. You can use these methods to create device control logic in your **Smart Home** project.

### Entity Object Methods

```coffeescript
Entity
  .setState(stateName)
  .setAttributes(key, attrs)
  .getAttributes() -> attrs
  .setMetric(name, value)
  .callAction(name, value)
  .getSettings() -> map[string]any
```

| Method | Description |
|--------|-------------|
| `setState` | Sets the current state of the device. `stateName` - name of the future state. |
| `setAttributes` | Sets the attributes of the current device state. `key` - map parameter {key:val} |
| `getAttributes` | Retrieves the current attributes of the device. `attrs` - map parameter {key:val} |
| `setMetric` | Updates the metrics of the device. `name` - metric name, `value` - map parameter {key:val} |
| `callAction` | Executes a command on the device. `actionName` - name of the command, `args` - map parameter {key:val} |
| `getSettings` | Retrieves the settings of the device. `settings` - map parameter {key:val} |

