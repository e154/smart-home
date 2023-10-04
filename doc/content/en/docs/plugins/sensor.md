
---
title: "Sensor"
linkTitle: "sensor"
date: 2021-10-20
description: >
  
---

This "sensor" plugin allows you to create and manage various types of sensors that can measure and track various environmental parameters such as temperature, humidity, light intensity, and other physical quantities.

Creating an abstract device based on the "sensor" plugin simplifies the process of integrating different sensors into a smart home system. It provides a unified interface for working with different types of sensors, making it easy to retrieve data from them, set their state, and configure their operation parameters.

Thanks to the flexibility and versatility of the "sensor" plugin, it is suitable for most use cases where working with sensors and obtaining data about the environment's state is required.

The "sensor" plugin also implements a JavaScript handler called `entityAction`. This handler is designed to handle actions related to "entity" type devices based on the "sensor" plugin.

Example implementation of the `entityAction` handler:

```javascript
entityAction = (entityId, actionName, args) => {
  // Action handling code
};
```

In this example, the `entityAction` handler takes two parameters: `entityId`, representing the identifier of the "entity" type device, and `actionName`, representing the name of the action to be performed.

The `entityAction` handler allows you to perform the necessary logic in response to actions related to "entity" type devices. Inside the handler, you can access functions and methods of the "sensor" plugin to perform specific operations, update the devices' state, or interact with other components of the smart home system.

Example usage of the `entityAction` handler:

```javascript
entityManager.callAction('sensor.sensor123', 'turnOn');
```

In this example, the `entityAction` handler is called with the device identifier "sensor123" and the action name "turnOn". Inside the handler, you can define the logic associated with executing the required action for the specified device.
