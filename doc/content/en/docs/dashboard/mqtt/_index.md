---
title: "MQTT"
linkTitle: "MQTT"
date: 2017-01-03
description: >
  
---

The **Smart Home** project system includes a built-in MQTT broker, which allows for message exchange and device control
in the smart home using the MQTT (Message Queuing Telemetry Transport) protocol.

MQTT is a lightweight messaging protocol designed for exchanging data between devices in networks with limited bandwidth
or unreliable connections. It is based on the publisher-subscriber model and provides reliable message delivery between
devices.

The built-in MQTT broker in the **Smart Home** project provides the server-side component for handling MQTT messages. It
performs the following functions:

1. Publisher and Subscriber Support: The MQTT broker accepts messages from publishers and forwards them to subscribers.
   This enables devices in the smart home to exchange data and control each other.

2. Topics and Filtering: The MQTT broker uses the concept of "topics" to classify messages. Subscribers can subscribe to
   specific topics to receive only the data they are interested in. The broker filters messages and forwards them only
   to subscribers subscribed to the corresponding topics.

3. Guaranteed Message Delivery: The MQTT broker ensures reliable message delivery, guaranteeing that messages will be
   delivered to subscribers even in case of temporary connection issues or device unavailability.

4. Access Control: The MQTT broker provides authentication and authorization capabilities, allowing control over access
   to different topics and restricting access to devices and data in the smart home.

The built-in MQTT broker in the **Smart Home** project simplifies the interaction and integration of devices in the
smart home. It provides a stable and efficient environment for data exchange and device control using the MQTT protocol.
This enables the creation of complex automation scenarios, monitoring device states, and exchanging data between
different devices.
