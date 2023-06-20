---
title: "Entity"
linkTitle: "Entity"
date: 2021-11-20
description: >

---


```bash
+---------------------+
|       Entity        |
+---------------------+
|      Actions        |
|      State          |
|    Attributes       |
|     Settings        |
|      Metrics        |
|      Storage        |
+---------------------+
```

This is a simple scheme that represents an "Entity" object. Each component has its role and functionality in managing and monitoring the "Entity" object in a smart home.

The "Entity" object is a central element of the Smart Home system and combines various aspects of the object, its states, attributes, actions, settings, and metrics. This allows the system to efficiently control and monitor objects in a smart home.

Here is a detailed description of the components of the "Entity" object:

Actions: The "Entity" object can receive and process various actions or commands. Actions are operations that can be performed on the object, such as turning on, turning off, changing parameters, etc.

State: The "Entity" object has a list of states it can transition into. At any given moment, the object can be in only one state. Examples of states may include "On," "Off," "Standby," "Playing," etc.

Attributes: Attributes are a storage of the object's state. It is a predefined list of fields and properties that contain information about the current state of the object. Attributes can be represented as a map[string]any, where the key is the attribute name, and the value is the corresponding attribute value.

Settings: The settings of the "Entity" object are a pre-defined list of immutable fields and properties. They define the configuration parameters of the object that can be set during its setup. Settings can also be represented as a map[string]any.

Metrics: The metrics of the "Entity" object provide information about its attributes or state, which is used for monitoring and measuring the performance or behavior of the object. Metrics can include data such as the average value of an attribute, the number of state changes, etc.

Storage: The storage of the "Entity" object provides a history of changes to its state or attributes. It records and stores previous values, allowing for tracking and analyzing the object's change history. The storage can be used for displaying graphs, analytics, or performing other operations with historical data of the object.

The "Entity" object brings all these components together, providing a unified and flexible approach to managing and monitoring various devices and systems in a smart home.

In the Smart Home system, each "Entity" object is implemented based on a specific plugin. Plugins provide different functionalities and capabilities for the "Entity" objects in the system. Some popular plugins that can be used to create "Entity" objects include sensor, mqtt, weather, automation, and others.

When a new "Entity" object is created, it is associated with a specific plugin that determines its functionality and capabilities. For example, if the "Entity" object represents a sensor, the sensor plugin can be used for its implementation. If the "Entity" object is intended to interact with an MQTT broker, the mqtt plugin can be used for it.

Each plugin provides its own handlers and methods that allow "Entity" objects to perform specific actions, retrieve data, send messages, etc.

```bash
              +-----------------------+
              |                       |
              |     Smart Home        |
              |                       |
              +-----------------------+
                         |
                         |
                         v
              +-----------------------+
              |                       |
              |       Entities        |
              |                       |
              +-----------------------+
                         |
                         |
                         v
   +----------------------------------------+
   |                                        |
   |              Plugins                   |
   |                                        |
   +----------------------------------------+
   |             |           |              |
   |             |           |              |
   v             v           v              v
+----------+ +----------+ +----------+ +----------+
|          | |          | |          | |          |
|  Plugin  | |  Plugin  | |  Plugin  | |  Plugin  |
|          | |          | |          | |          |
+----------+ +----------+ +----------+ +----------+
   |             |           |              |
   |             |           |              |
   v             v           v              v
+----------+ +----------+ +----------+ +----------+
|          | |          | |          | |          |
| Entity 1 | | Entity 2 | | Entity 3 | | Entity 4 |
|          | |          | |          | |          |
+----------+ +----------+ +----------+ +----------+

                     ^
                     |
                     |
              +------------------+
              |                  |
              |  Automation      |
              |                  |
              +------------------+
```

On the diagram, the general structure of the relationship between "Entity" objects, plugins, and the automation component in the Smart Home system is presented.

Smart Home is the central part of the system and coordinates the interaction between all the components.
"Entity" objects represent specific devices, sensors, or other components in the system. Each "Entity" object can use a specific plugin for its implementation.
Plugins provide functionality and capabilities for the "Entity" objects. They contain logic and methods that allow "Entity" objects to interact with external devices, gather data, send commands, etc.
The automation component is responsible for creating scenarios and triggering specific actions based on conditions and triggers. It can use "Entity" objects and their plugins to define conditions, triggers, and actions in the scenarios.
This structure allows the Smart Home system to be flexible and extensible, as new plugins and "Entity" objects can be added, and the automation component can use them to create various scenarios and automate actions.
