---
title: "entityManager"
linkTitle: "entityManager"
date: 2021-11-20 
description: >

---

In the **Smart Home** project, there is a device manager called "Entity Manager" that provides convenient access to devices, their settings, and current state.

{{< alert color="success" >}}This function is available in any system script.{{< /alert >}}

The "Entity Manager" device manager offers a set of methods and functions for managing devices and retrieving information about them. Some key features of the "Entity Manager" include:

1. `getEntity(id) -> Entity`: This method allows you to retrieve an "Entity" object by its identifier.
   You pass the device identifier as the `id` argument, and the method returns the corresponding "Entity" object.

   ```javascript
   const entity = entityManager.getEntity(id);
   ```

2. `setState(id, stateName)`: This method is used to set the state of a device with the given identifier.
   You pass the device identifier as the `id` argument and the state name as the `stateName` argument. For example:

   ```javascript
   entityManager.setState(id, 'on');
   ```

3. `setAttributes(id, attrs)`: This method allows you to set attributes for the device with the given identifier.
   You pass the device identifier as the `id` argument and an object with the new attributes as the `attrs` argument.

   ```javascript
   const newAttributes = {
     brightness: 80,
     color: 'red',
   };

   entityManager.setAttributes(id, newAttributes);
   ```

4. `setMetric(id, name, value)`: This method is used to set a metric for the device with the given identifier.
   You pass the device identifier as the `id` argument, the metric name as the `name` argument, and the metric value as the `value` argument.

   ```javascript
   entityManager.setMetric(id, 'temperature', 25);
   ```

5. `callAction(id, actionName, args)`: This method allows you to invoke a specific action on the device with the given identifier.
   You pass the device identifier as the `id` argument, the action name as the `actionName` argument, and the action-related arguments as the `args` argument.

   ```javascript
   entityManager.callAction(id, 'turnOn');
   ```

6. `callScene(id, args)`: This method is used to call a scene with the given identifier.
   You pass the scene identifier as the `id` argument and the scene-related arguments as the `args` argument.

   ```javascript
   entityManager.callScene(id, args);
   ```

The "Entity Manager" device manager simplifies the management of devices, settings, and state in the **Smart Home** project. It provides a convenient interface for retrieving information about devices, controlling them, and handling related events. This allows you to create more flexible and functional applications for your smart home.

----------------

### entityManager object

```coffeescript
# Available methods
entityManager
  .getEntity(id) -> Entity
  .setState(id, stateName)
  .setAttributes(id, attrs)
  .setMetric(id, name, value)
  .callAction(id, actionName, args)
  .callScene(id, args)
```

| Value      | Description                                              |
|------------|----------------------------------------------------------|
| id         | The unique identifier of the device.                      |
| getEntity  | Find a device, returns [Entity](#entity-object-methods).  |
| setState   | Set the current state, where **stateName** is the name of the future state. |
| setAttributes | Set the attribute of the current state, where **attrs** is a map of parameters {

key: val}. |
| setMetric | Update the device's metrics, where **name** is the name of the metric, and **value** is the metric value. |
| callAction | Execute a command on the device, where **actionName** is the name of the command, and **args** is a map of parameters {key: val}. |
| callScene  | Execute a scene with the specified **id**.                |


----------------

### Code Example

```coffeescript
# Entity Manager
# ##################################

# Find a device by its unique id, returns [Entity](#entity-object)
device = entityManager.getEntity('telegram.clavicus')

# Set the current state using the unique state name
entityManager.setState('telegram.clavicus', 'IN_WORK')

# Set attributes for the current state
entityManager.setAttributes('telegram.clavicus', {'temp': 23.3})

# Update device metrics
entityManager.setMetric('telegram.clavicus', 'cpu', {'all': 55})

# Execute a command on the device
entityManager.callAction('floor.light', 'ON', {'power': 30})

# Get settings
settings = device.getSettings()
headers = [{'apikey':settings['apikey']}]
url = 'https://webhook.site/2692a589-a5bc-4156-af7d-75875578798f'
res = http.headers(headers).get(url)
```

