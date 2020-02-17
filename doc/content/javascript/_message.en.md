---
weight: 57
title: message
groups:
    - javascript
---

## message {#message}

**message** - транспортный объект для передачи состояния(переменные и объекты вызова) между сущностями в
конструкторе процессов **Flow editor**

Доступные свойства и методы:

### .SetVar(key, value)

```coffeescript
message.SetVar('foo', 'bar')
```

определение переменной в хранилище

**На входе**

**Значение** | **Описание**
-------------|--------------
  `key`      | type: string
  `value`    | type: interface

### .GetVar(key)

получение переменной из хранилища

```coffeescript
value = message.GetVar('foo')
```

**На входе**

**Значение** | **Описание**
-------------|--------------
  `key`      | type: string
  
**На выходе**

**Значение** | **Описание**
-------------|--------------
  `value`    | type: interface

### .SetError(error)

сохранение состояния ошибки

```coffeescript
message.SetError('some error')
```

**На входе**

**Значение** | **Описание**
-------------|--------------
  `error`    | type: string

### .Error

переменная хранения состояния отрицательного результата

```coffeescript
err = message.Error
```
**На выходе**

**Значение** | **Описание**
-------------|--------------
  `err`    | type: string

### .Ok()

сохронение состояния положительного результата

```coffeescript
message.Ok()
```

### .Success

флаг положительного результата

```coffeescript
value = message.Success
```

**На выходе**

**Значение** | **Описание**
-------------|--------------
  `value`    | type: bool

### .Clear()

очистка состояний ошибки и хранилища переменных

```coffeescript
message.Clear()
```

### .Setdir(dir)

установка флага для указания желаемго направления движения процесса в **Flow engine**

```coffeescript
message.Setdir(dir)
```

**На входе**

**Значение** | **Описание**
-------------|--------------
  `dir`    | type: bool
  
### .Direction

получение состояние направления движения

```coffeescript
value = message.Direction
```

**На выходе**

**Значение** | **Описание**
-------------|--------------
  `value`    | type: bool
  
### .Mqtt

флаг вызывающей стороны текущего процесса **Flow**

```coffeescript
value = message.Mqtt
topic = message.GetVar('mqtt_topic')
payload = JSON.parse(message.GetVar('mqtt_payload'))
qos = message.GetVar('mqtt_qos')
duplicate = message.GetVar('mqtt_duplicate')
```

**На выходе**

**Значение** | **Описание**
-------------|--------------
  `value`    | type: bool
  `topic`    | type: string
  `payload`  | type: Object
  `qos`      | type: byte
  `duplicate`| type: bool

