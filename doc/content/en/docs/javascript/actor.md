---
title: "Actor"
linkTitle: "actor"
date: 2021-11-20
description: >

---

In the **Smart Home** project, the JavaScript code includes an accessible object called "Actor," which represents a simplified entity model. The Actor object provides convenient methods and properties for interacting with plugin states.

Device State Management: The Actor object provides methods for reading and setting the device state. You can check the current device state and modify it as needed.

{{< alert color="warning" >}}The object may not be available in all scopes.{{< /alert >}}

The "Actor" object in the **Smart Home** project provides the following methods:

1. `setState(entityStateParams)`: This method allows you to set the state of an actor based on the provided parameters. You pass an `entityStateParams` object as an argument, which contains the actor's state parameters to be set. For example:

```javascript
const entityStateParams = {
  power: true,
  brightness: 80,
};

actor.setState(entityStateParams);
```

2. `getSettings() -> Map[string]any`: This method returns the actor's settings as an associative array. The result is a map, where keys are strings and values can be of any data type. You can use this method to retrieve the actor's settings. Example usage:

```javascript
const settings = actor.getSettings();
console.log(settings);
```

The methods of the "Actor" object provide functionality for setting the actor's state based on provided parameters and retrieving the actor's settings. You can use these methods to work with actors in your **Smart Home** project and configure their states according to your requirements.

----------------

### Actor Object

```coffeescript
# Example of reading device state
Actor
  .setState(entityStateParams)
  .getSettings() -> map[string]any
```

### entityStateParams Object

|   Value   |  Type  | Description                                      |
| --------- | ------ | ------------------------------------------------ |
| new_state | string | Name of the new state (optional)                 |
| attribute_values | map[string]any | Description of the attributes          |
| settings_value | map[string]any | Description of the settings            |
| storage_save | bool   | Flag for database storage                        |

----------------

### Code Example

```coffeescript
settings = Actor.getSettings()
headers = [{'apikey':settings['apikey']}]
url = 'https://webhook.site/2692a589-a5bc-4156-af7d-75875578798f'
res = http.headers(headers).get(url)
```
