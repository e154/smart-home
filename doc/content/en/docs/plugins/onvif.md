---
title: "Onvif"
linkTitle: "onvif"
date: 2024-01-04
description: >
  
---

The Smart Home project provides the ONVIF plugin, enabling interaction with surveillance cameras using the ONVIF
protocol. This plugin implements several JavaScript methods for camera control and snapshot retrieval.

#### JavaScript Methods

1. **`Camera.continuousMove(X, Y)`**: This method allows for smooth camera movement to the specified coordinates X and
   Y.

2. **`Camera.stopContinuousMove()`**: This method stops continuous camera movement.

3. **`OnvifGetSnapshotUri(entityId)`**: Method for obtaining the snapshot URI for the specified device identifier.

4. **`DownloadSnapshot(entityId)`**: Method for downloading a snapshot from the device based on its identifier.

#### Device Status

- **`motion` (type: Boolean)**: Status indicating the presence of motion detected by the surveillance camera.

#### Device Settings

- **`address` (type: String)**: IP address of the surveillance camera.

- **`onvifPort` (type: Int)**: Port for connecting to the camera via the ONVIF protocol.

- **`rtspPort` (type: Int)**: Port for video transmission via the RTSP protocol.

- **`userName` (type: String)**: User name for authentication when accessing the camera.

- **`password` (type: Encrypted)**: Encrypted password for authentication.

- **`requireAuthorization` (type: Bool)**: Flag indicating whether authorization is required for interacting with the
  camera.

#### Control Commands

- **`continuousMove`**: Command to initiate continuous camera movement.

- **`stopContinuousMove`**: Command to stop continuous camera movement.

#### Device Statuses

- **`connected`**: The device is successfully connected and ready to operate.

- **`offline`**: The device is unavailable or disconnected from the system.

These functions allow the integration of surveillance cameras into the Smart Home system and efficient management
through the ONVIF plugin.

Example of using the "onvif" plugin to implement camera control:

```javascript
continuousMove = function (args) {
  var X, Y;
  X = args['X'] || 0;
  Y = args['Y'] || 0;
  if (Math.abs(X) > Math.abs(Y)) {
    Y = 0;
  } else {
    X = 0;
  }
  return Camera.continuousMove(X, Y);
};

stopStop = function (args) {
  return Camera.stopContinuousMove();
};

```
