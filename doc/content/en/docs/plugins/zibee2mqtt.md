
---
title: "Zigbee2mqtt"
linkTitle: "zigbee2mqtt"
date: 2021-10-20
description: >
  
---

The "Zigbee2mqtt" plugin is part of the system and provides integration between Zigbee devices and an MQTT broker. This plugin allows you to control and manage Zigbee devices through MQTT messages.

The "Zigbee2mqtt" plugin implements a JavaScript handler called `zigbee2mqttEvent`, which doesn't take any parameters. Inside the handler, there is an accessible `message` object that represents the information of the received MQTT message.

The properties of the `message` object include:

1. `payload`: The value of the message, represented as a dictionary (map) where the keys are strings and the values can be of any type.
2. `topic`: The MQTT message topic, indicating the source or destination of the message.
3. `qos`: The Quality of Service level of the MQTT message.
4. `duplicate`: A flag indicating whether the message is a duplicate.
5. `storage`: A `Storage` object that provides access to a data store for caching and retrieving arbitrary values.
6. `error`: A string containing information about an error if one occurred during message processing.
7. `success`: A boolean value indicating the successful completion of an operation or message processing.
8. `new_state`: An `Actor.StateParams` object representing the new state of an actor after an operation is performed.

Here's an example of using the `zigbee2mqttEvent` handler:

```javascript
function zigbee2mqttEvent() {
    console.log("Received MQTT message:");
    console.log("Payload:", message.payload);
    console.log("Topic:", message.topic);
    console.log("QoS:", message.qos);
    console.log("Duplicate:", message.duplicate);

    if (message.error) {
        console.error("Error:", message.error);
    } else if (message.success) {
        console.log("Operation successful!");
        console.log("New state:", message.new_state);
    }

    // Accessing the storage
    const value = message.storage.getByName("key");
    console.log("Value from storage:", value);
}
```

The `zigbee2mqttEvent` handler can be used to process incoming MQTT messages and perform additional operations based on the received data.

Here's an example in CoffeeScript:

```coffeescript
zigbee2mqttEvent = ->
  # Print '---mqtt new event from plug---'
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
  Actor.setState
    'new_state': payload.state
    'attribute_values': attrs
```

```coffeescript
zigbee2mqttEvent = ->
  # Print '---mqtt new event from button---'
  if !message
    return
  payload = unmarshal message.payload
  attrs =
    'battery': payload.battery
    'linkquality': payload.linkquality
    'voltage': payload.voltage
  state = ''
  if payload.action
    attrs.action = payload.action
    state = payload.action + "_action"
  if payload.click
    attrs.click = payload.click
    attrs.action = ""
    state = payload.click + "_click"
  Actor.setState
    'new_state': state.toUpperCase()
    'attribute_values': attrs
```
