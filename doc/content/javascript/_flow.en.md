---
weight: 30
title: flow
groups:
    - javascript
---

## Flow {#flow}

Получить текущий процесс. Может вырнуть *null* если скрипт запускается 
вне контекста [Workflow](#workflow)

```coffeescript
flow = Flow
return if !flow
```

**На выходе**

**Значение** | **Описание**
-------------|--------------
  `flow`     | type: Object, ссылка на экземпляр *Flow 
  
Доступные методы приведены далее:

### .GetName() {#flow_get_name}

Получить наименование текущего Flow процесса.

```coffeescript
flow = Flow
if flow
 name = flow.GetName()
 print 'flow name', name
```

**На выходе**

**Значение** | **Описание**
-------------|--------------
  `name`     | type: string


### .SetVar(key, value) {#flow_set_var}

Запомнить переменнную в хранилище [Flow](#flow). Хранилище позволяет 
сохранять состояния на время жизни [Flow](#flow)

```coffeescript
flow = Flow
if flow
 flow.SetVar("key", "value")
```

**На входе**

**Значение** | **Описание**
-------------|--------------
  `key`      | type: string
  `value`    | type: interface

### .GetVar(key) {#flow_get_var}

Получить ранее записанную переменную

```coffeescript
wf = Flow
if wf
  variable = wf.GetVar(key)
```

**На входе**

**Значение** | **Описание**
-------------|--------------
  `key`      | type: string

**На выходе**

**Значение** | **Описание**
-------------|--------------
  `variable` | type: interface
  
### .node() {#flow_node}

Получить ноду

```coffeescript
flow = Flow
if flow
 node = flow.node()
 print 'flow node', node
```

**На выходе**

**Значение** | **Описание**
-------------|--------------
  `node`     | type: Object, ссылка на экземпляр [Node](#node) 
