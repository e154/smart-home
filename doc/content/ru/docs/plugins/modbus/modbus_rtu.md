---
title: "Modbus rtu"
linkTitle: "modbus rtu"
date: 2021-10-20 
description: >

---

Modbus RTU — компактный двоичный вариант. Сообщения разделяются по паузе в линии. Сообщение должно начинаться и
заканчиваться интервалом тишины, длительностью не менее 3,5 символов при данной скорости передачи. Во время передачи
сообщения не должно быть пауз длительностью более 1,5 символов. Для скоростей более 19200 бод допускается использовать
интервалы 1,75 и 0,75 мс, соответственно. Проверка целостности осуществляется с помощью CRC.

{{< alert color="warning" >}}Для работы с modbus rtu устройством требуется настроенная **нода**{{< /alert >}}

### Настройка:

* slave_id `1-32`
* baud `9600, 19200`
* data_bits `5-9`
* timeout `milliseconds`
* stop_bits `1-2`
* sleep `milliseconds`
* parity `none, odd, even`

### Команды:

* произвольный набор

### Атрибуты

* произвольный набор

### Статус

* произвольный набор

----------------

### Доступные функции

**1 битные функции**

ReadCoils           
ReadDiscreteInputs  
WriteSingleCoil     
WriteMultipleCoils

**16 битные функции**

ReadInputRegisters          
ReadHoldingRegisters        
ReadWriteMultipleRegisters  
WriteSingleRegister         
WriteMultipleRegisters


----------------

### javascript свойства

* ENTITY_ID
* ModbusRtu
* entityAction
* Action

```coffeescript
# константа с уникальным id устройства
const ENTITY_ID
````

```coffeescript
# выполнение команды(функции) на устройсте:
result = ModbusRtu func, addr, count, command
```

|  значение  | описание  |
|-------------|---------|
| func | функция для выполнения на устройстве  |
| addr | Адрес первого регистра (40108-40001 = 107 =6B hex)  |
| count | Количество требуемых регистров (чтение 3-х регистров с 40108 по 40110) |
| command | Команда |


```coffeescript
# функция-обработчик события действий:
entityAction = (entityId, actionName)->
```

|  значение  | описание  |
|-------------|---------|
| entityId | уникальное id устройства  |
| actionName | системное наименование действия  |

{{< alert color="warning" >}}Объект **Action** доступен в скриптах действий и скриптах закрепленных за устройством.{{< /alert >}}
```coffeescript
state = {
  new_state: 'ENABLED',
  attribute_values: {
    heat: false
  },
  settings_value: {},
  storage_save: true
} 
# сохранение состояния 
Actor.setState(state)
```

|  значение  | описание  |
|-------------|---------|
| new_state | уникальное системное наименование состояния |
| attribute_values | значения атрибутов ранее определенных для устройства |
| settings_value | значения настроек ранее определенных для устройства |
| storage_save | признак записи среза состояния в БД |

----------------

### пример coffeescript кода

```coffeescript
# ModbusRtu
# ##################################
"use strict";

checkStatus = ->
  COMMAND = []
  FUNC = 'ReadHoldingRegisters'
  ADDRESS = 0
  COUNT = 16
  
  # выполнение команды(функции) на устройсте:
  res = ModbusRtu FUNC, ADDRESS, COUNT, COMMAND
  print res.error
  print res.result
  print res.time

# функция-обработчик события действий:
entityAction = (entityId, actionName)->
  switch actionName
    when 'ON' then doOnAction()
    when 'OFF' then doOffAction()
    when 'CHECK' then checkStatus()
    when 'ON_WITH_ERR' then doOnErrAction()

```


