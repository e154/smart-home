---
weight: 50
title: device
groups:
    - javascript
---

## Device{} {#device}

Получить текущее устройство. Может вырнуть *null* если скрипт запускается 
вне контекста [Workflow](#ic_workflow)
  
Доступные методы приведены далее:

### .getName() {#deice_get_name}

Получить наименование устройства.

```coffeescript
name = dev.getName()
```

**На выходе**

**Значение** | **Описание**
-------------|--------------
  `name`     | type: string


### .getDescription() {#deice_get_description}

Получить описание устройства.

```coffeescript
description = dev.getDescription()
```

**На выходе**

**Значение**    | **Описание**
----------------|--------------
  `description` | type: string

### .runCommand(name, args) {#deice_run_command}

Выполнить комманду на устройстве.

```coffeescript
description = dev.runCommand(name, ['foo','bar'])
```

**На выходе**

**Значение**    | **Описание**
----------------|--------------
  `name`        | type: string
  `args`        | type: array

### .modBus(func, address, count, command) {#deice_mod_bus}

Выполнить комманду на modbus устройстве.

```coffeescript
COMMAND = []
FUNC = 'ReadHoldingRegisters'
ADDRESS = 0
COUNT = 16
res = device.modBus(FUNC, ADDRESS, COUNT, COMMAND)
if res.error
  print 'error: ', res.error
```

**На входе**

**Значение**    | **Описание**
----------------|--------------
  `func`        | type: string, вызываемая функция
  `address`     | type: int, адрес устройства
  `count`       | type: int
  `command`     | type: array(int) 

варианты значения поля `func`

**1 битные функции**|
--------------------|
ReadCoils           |
ReadDiscreteInputs  |
WriteSingleCoil     |
WriteMultipleCoils  |

варианты значения поля `func`

**16 битные функции**       |
----------------------------|
ReadInputRegisters          |
ReadHoldingRegisters        |
ReadWriteMultipleRegisters  |
WriteSingleRegister         |
WriteMultipleRegisters      |

### .smartBus(command) {#deice_smart_bus}

Выполнить комманду на smartbus устройстве.

```coffeescript
COMMAND = []
res = device.smartBus(COMMAND)
if res.error
  print 'error: ', res.error
```

**На выходе**

**Значение**    | **Описание**
----------------|--------------
  `command`     | type: array(int) 
