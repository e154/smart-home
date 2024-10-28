---
title: "MQTT"
linkTitle: "mqtt"
date: 2021-10-20
description: >
  
---

The "mqtt" plugin is an extended version of the "sensor" plugin that provides MQTT protocol functionality. Here are some
settings available for the "mqtt" plugin:

1. `subscribe_topic`: The MQTT topic to subscribe to for receiving data.
2. `mqtt_login`: The login for authentication when connecting to the MQTT broker.
3. `mqtt_pass`: The password for authentication when connecting to the MQTT broker.

These settings allow you to specify the connection parameters for the MQTT broker and define the topic to subscribe to
for receiving data. Thus, the "mqtt" plugin enables integration with MQTT messages, allowing you to receive data from
devices using this protocol and use them in automation systems or other components.

The "sensor" plugin also implements a JavaScript handler called `entityAction`. This handler is designed to process
actions related to "entity" type devices based on the "sensor" plugin.

Here's an example implementation of the `entityAction` handler:

```javascript
entityAction = (entityId, actionName, args) => {
  // Action handling code
};
```

It also implements a JavaScript handler called `mqttEvent`. This handler is designed to process actions related to "
entity" type devices based on the "sensor" plugin.

Here's an example implementation of the `mqttEvent` handler:

```javascript
function mqttEvent(message) {
  // Action handling code
};
```

Example usage of the `mqttEvent` handler:

```coffeescript
mqttEvent = (message)->
#print '---mqtt new event from plug---'
  if !message || message.topic.includes('/set')
    return
  payload = unmarshal message.payload
  attrs =
    'consumption': payload.consumption
    'linkquality': payload.linkquality
    'power': payload.power
    'state': payload.state
    'temperature': payload.temperature
    'voltage': payload.voltage
  EntitySetState ENTITY_ID,
    'new_state': payload.state
    'attribute_values': attrs
```
