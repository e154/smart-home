---
title: "Speedtest"
linkTitle: "speedtest"
date: 2024-01-04
description: >
  
---

The Speedtest plugin provides the capability to measure internet connection speed directly from the Smart Home system.
This plugin allows users to check the quality of their internet connection and use the obtained data to make decisions
regarding network optimization.

#### Device Properties

Each device created using the Speedtest plugin has the following properties:

- **`Point`**: Geographical coordinates of the Speedtest node.

- **`Name` (type: String)**: Name of the conducted Speedtest.

- **`Country` (type: String)**: Country where the test is conducted.

- **`Sponsor` (type: String)**: Name of the Speedtest sponsor.

- **`Distance` (type: Float)**: Distance to the Speedtest server in kilometers.

- **`Latency` (type: String)**: Latency (ping) to the Speedtest server.

- **`Jitter` (type: String)**: Indicator of internet connection stability.

- **`DLSpeed` (type: Float)**: Download Speed in megabits per second.

- **`ULSpeed` (type: Float)**: Upload Speed in megabits per second.

- **`location` (type: Point)**: Geographical coordinates of the device location.

#### Device Settings

Each device created using the Speedtest plugin contains the following settings:

- **`location` (type: Point)**: Geographical coordinates indicating the device location for conducting the Speedtest.

- **`city` (type: String)**: Name of the city where the testing will take place.

#### Device Statuses

Each device created using the Speedtest plugin can have the following statuses:

- **`completed`**: Testing has been successfully completed, and the data is available for analysis.

- **`in process`**: Testing is currently in progress.
