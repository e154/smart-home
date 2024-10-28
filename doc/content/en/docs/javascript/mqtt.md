---
title: "Mqtt"
linkTitle: "mqtt"
date: 2021-11-20 
description: >

---

The "Mqtt" object in the **Smart Home** project provides the ability to publish messages to an MQTT broker.

#### Quality of Service in MQTT (QoS)

QoS 0: At most once. At this level, the publisher sends the message to the broker once and does not wait for
confirmation. It means "fire and forget."

QoS 1: At least once. This level guarantees that the message will be delivered to the broker, but there is a possibility
of message duplication from the publisher. When the broker receives a duplicate message, it resends it to the
subscribers and sends an acknowledgment to the publisher. If the publisher doesn't receive a PUBACK message from the
broker, it resends the packet with the DUP flag set to "1."

QoS 2: Exactly once. This level guarantees the delivery of messages to the subscriber and eliminates the possibility of
duplicating sent messages.

{{< alert color="success" >}}This function is available in any system script.{{< /alert >}}

The "Mqtt" object in the **Smart Home** project provides the following method:

1. `publish(topic, payload, qos, retain)`: This method is used to publish messages to an MQTT broker. You pass the
   arguments `topic` (the topic), `payload` (the message content), `qos` (Quality of Service level), and `retain` (the "
   retain" flag).

- `topic`: A string representing the topic to which the message will be published.
- `payload`: The content of the message, which can be a string, an object, or binary data.
- `qos`: The Quality of Service level determines the message delivery guarantees. Possible values are 0 (at most once),
  1 (at least once with acknowledgment), or 2 (exactly once with retransmission if needed).
- `retain`: A flag indicating whether the message should be retained on the broker and sent to new subscribers. Values
  can be true (retain) or false (do not retain).

Usage example:

```javascript
const topic = 'smart-home/living-room/light';
const payload = 'on';
const qos = 1;
const retain = true;

Mqtt.publish(topic, payload, qos, retain);
```

The `publish` method allows you to send messages to an MQTT broker for transmitting information and controlling devices
or the system as a whole in your **Smart Home** project.

----------------

### Mqtt object

```coffeescript
# Publish a message
Mqtt
  .publish(topic, payload, qos, retain)
```

| Value   | Description                                                                                                                                                                                                |
|---------|------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| topic   | Channel for the message                                                                                                                                                                                    |
| payload | JSON object intended for transmission                                                                                                                                                                      |
| qos     | Quality of Service (1, 2, 3)                                                                                                                                                                               |
| retain  | When publishing data with the retain flag set, the broker will store it. Upon the next subscription to this topic, the broker immediately sends the message with this flag. Only used in PUBLISH messages. |

----------------

### Code Example

```coffeescript
# mqtt
# ##################################

payload = JSON.stringify({"state": actionName})
Mqtt.publish"zigbee2mqtt/device/set", payload, 0, false)

```
