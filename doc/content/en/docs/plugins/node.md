---
title: "Node"
linkTitle: "node"
date: 2021-10-20
description: >
  
---

The Node plugin is designed for the integration of external Node agents, providing the ability to work with Modbus
devices on remote servers in the network or from other subnets. This plugin simplifies the connection and interaction
with remote Modbus agents, expanding the functionality of the Smart Home system.

#### Device Properties

- **`thread` (type: Int)**: Identifier of the thread in which the Node agent is operating.

- **`rps` (type: Int)**: Number of requests per second (Requests Per Second) processed by the agent.

- **`min` (type: Int)**: Minimum response time from the agent in milliseconds.

- **`max` (type: Int)**: Maximum response time from the agent in milliseconds.

- **`latency` (type: Int)**: Delay between the request and response from the agent in milliseconds.

- **`startedAt` (type: Time)**: Time when the Node agent was started.

#### Device Settings

- **`node_login` (type: String)**: Login for authentication on the Node agent.

- **`node_pass` (type: String)**: Password for authentication on the Node agent.

#### Device Statuses

- **`wait` (type: String)**: The agent is waiting for connections and requests from the Smart Home system.

- **`connected` (type: String)**: Successful connection established with the Node agent.

- **`error` (type: String)**: An error occurred during the interaction with the Node agent.
