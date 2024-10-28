---
title: "MQTT bridge"
linkTitle: "mqtt_bridge"
date: 2024-01-04
description: >
  
---

The MQTT Bridge plugin provides functionality for integrating devices through the MQTT protocol. This plugin allows
efficient data exchange between devices and the Smart Home system using an MQTT broker.

#### Device Settings

Each device created using the MQTT Bridge plugin has the following settings:

- **`keepAlive` (type: Int)**: The time in seconds after which the device sends a ping to maintain an active connection
  with the broker.

- **`pingTimeout` (type: Int)**: The time in seconds expected to receive a response to the ping from the broker.

- **`broker` (type: String)**: The address of the MQTT broker to which the device will connect.

- **`clientID` (type: String)**: The client identifier used when connecting to the broker.

- **`connectTimeout` (type: Int)**: The time in seconds allocated for establishing a connection with the broker.

- **`cleanSession` (type: Bool)**: A flag indicating whether to use a "clean" session when connecting.

- **`username` (type: String)**: The username for authentication when connecting to the broker.

- **`password` (type: Encrypted)**: The encrypted password for authentication when connecting to the broker.

- **`qos` (type: Int)**: The Quality of Service level for interacting with the broker.

- **`direction` (type: String)**: The direction of interaction (e.g., "inbound" or "outbound").

- **`topics` (type: String)**: A list of topics with which the device will interact.

#### Device Statuses

Each device created using the MQTT Bridge plugin can have the following statuses:

- **`connected`**: The device is successfully connected to the MQTT broker and ready for data exchange.

- **`offline`**: The device is not connected to the broker or has lost the connection.

These settings and statuses provide flexibility in integrating devices through the MQTT protocol, allowing easy
configuration and monitoring of their state within the Smart Home system.
