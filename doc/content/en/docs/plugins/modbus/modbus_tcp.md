
---
title: "Modbus tcp"
linkTitle: "modbus tcp"
date: 2021-10-20
description: >

---

In the **Smart Home** system, there is a Modbus TCP plugin implemented that enables interaction with devices using the Modbus TCP protocol. The plugin provides the `ModbusTcp(func, addr, count, command)` method, which allows sending commands and receiving data using the specified protocol.

Arguments of the `ModbusTcp` method:

1. `func`: Modbus function code indicating the type of operation to perform with the device.
2. `addr`: Modbus device address to which the command is sent.
3. `count`: Number of registers or bits to read or write.
4. `command`: Command to be executed, indicating additional parameters and settings.

Usage example:

```javascript
const COMMAND = []
const FUNC = 'ReadHoldingRegisters'
const ADDRESS = 0
const COUNT = 16

const response = ModbusTcp(FUNC, ADDRESS, COUNT, COMMAND);
console.log(response);
```

The `ModbusTcp` method allows sending Modbus TCP commands, performing register reads or writes, and retrieving data from devices. You can use this method in your **Smart Home** project to interact with devices that support the Modbus TCP protocol.

{{< alert color="warning" >}}To work with a Modbus RTU device, a configured **node** is required.{{< /alert >}}

### Configuration:

* slave_id `1-32`
* address_port `localhost:502`

### Commands:

* Custom set

### Attributes:

* Custom set

### Status:

* Custom set

----------------

### Available Functions

**1-bit functions**

ReadCoils
ReadDiscreteInputs
WriteSingleCoil
WriteMultipleCoils

**16-bit functions**

ReadInputRegisters
ReadHoldingRegisters
ReadWriteMultipleRegisters
WriteSingleRegister
WriteMultipleRegisters

----------------

### JavaScript Properties

* ENTITY_ID
* ModbusRtu
* entityAction
* Action

```coffeescript
# Constant with a unique device ID
const ENTITY_ID
````

```coffeescript
# Execute a command (function) on the device:
result = ModbusTcp(func, addr, count, command)
```

|  Value   | Description |
|-------------|---------|
| func | Function to be executed on the device |
| addr | Address of the first register (40108-40001 = 107 = 6B hex) |
| count | Number of registers to read (reading 3 registers from 40108 to 40110) |
| command | Command |


```coffeescript
# Event handler function for actions:
entityAction = (entityId, actionName, args)->
```

|  Value   | Description |
|-------------|---------|
| entityId | Unique ID of the device |
| actionName | System name of the action |
| args | Type: map[string]any |

{{< alert color="warning" >}}The **Action** object is available in action scripts and scripts attached to the device.{{< /alert >}}
```coffeescript
state = {
  new_state: 'ENABLED',
  attribute_values: {
    heat: false
  },
  settings_value: {},
  storage_save: true
} 
# Save the state 
Actor.setState(state)
```

|  Value   | Description |
|-------------|---------|
| new_state | Unique system name of the state |
| attribute_values | Values of attributes previously defined for the device |
| settings_value | Values of settings previously defined for the device |
| storage_save | Flag indicating whether to save the state

----------------

### Example CoffeeScript code

```coffeescript
# ModbusTcp
# ##################################
"use strict";

checkStatus = ->
  COMMAND = []
  FUNC = 'ReadHoldingRegisters'
  ADDRESS = 0
  COUNT = 16

  res = ModbusTcp FUNC, ADDRESS, COUNT, COMMAND
  print res.error
  print res.result
  print res.time

entityAction = (entityId, actionName, args)->
  switch actionName
    when 'ON' then doOnAction()
    when 'OFF' then doOffAction()
    when 'CHECK' then checkStatus()
    when 'ON_WITH_ERR' then doOnErrAction()

```


