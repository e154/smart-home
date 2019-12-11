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

### .setVar(key, value)

```coffeescript
message.setVar('foo', 'bar')
```

определение переменной в хранилище

**На входе**

**Значение** | **Описание**
-------------|--------------
  `key`      | type: string
  `value`    | type: interface

### .getVar(key)

получение переменной из хранилища

```coffeescript
value = message.getVar('foo')
```

**На входе**

**Значение** | **Описание**
-------------|--------------
  `key`      | type: string
  
**На выходе**

**Значение** | **Описание**
-------------|--------------
  `value`    | type: interface

### .setError(error)

сохранение состояния ошибки

```coffeescript
message.setError('some error')
```

**На входе**

**Значение** | **Описание**
-------------|--------------
  `error`    | type: string

### .error

переменная хранения состояния отрицательного результата

```coffeescript
err = message.error
```
**На выходе**

**Значение** | **Описание**
-------------|--------------
  `err`    | type: string

### .ok()

сохронение состояния положительного результата

```coffeescript
message.ok()
```

### .success

флаг положительного результата

```coffeescript
value = message.success
```

**На выходе**

**Значение** | **Описание**
-------------|--------------
  `value`    | type: bool

### .clear()

очистка состояний ошибки и хранилища переменных

```coffeescript
message.clear()
```

### .setdir(dir)

установка флага для указания желаемго направления движения процесса в **Flow engine**

```coffeescript
message.setdir(dir)
```

**На входе**

**Значение** | **Описание**
-------------|--------------
  `dir`    | type: bool
  
### .direction

получение состояние направления движения

```coffeescript
value = message.direction
```

**На выходе**

**Значение** | **Описание**
-------------|--------------
  `value`    | type: bool
  
### .mqtt

флаг вызывающей стороны текущего процесса **Flow**

```coffeescript
value = message.mqtt
topic = message.getVar('mqtt_topic')
payload = JSON.parse(message.getVar('mqtt_payload'))
```

**На выходе**

**Значение** | **Описание**
-------------|--------------
  `value`    | type: bool
  `topic`    | type: string
  `payload`    | type: Object

