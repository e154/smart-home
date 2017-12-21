---
weight: 3
title: flow
groups:
    - javascript
---

## IC.Flow() {#ic_flow}

Получить текущий процесс. Может вырнуть *null* если скрипт запускается 
вне контекста [Workflow](#ic_workflow)

```coffeescript
flow = IC.Flow()
return if !flow
```

**На выходе**

**Значение** | **Описание**
-------------|--------------
  `flow`     | type: Object, ссылка на экземпляр *Flow 
  
Доступные методы приведены далее:

### .getName() {#ic_flow_get_name}

Получить наименование текущего Flow процесса.

```coffeescript
flow = IC.Flow()
if flow
 name = flow.getName()
 print 'flow name', name
```

**На выходе**

**Значение** | **Описание**
-------------|--------------
  `name`     | type: string


### .setVar(key, value) {#ic_flow_set_var}

Запомнить переменнную в хранилище [Flow](#ic_flow). Хранилище позволяет 
сохранять состояния на время жизни [Flow](#ic_flow)

```coffeescript
flow = IC.Flow()
if flow
 flow.setVar("key", "value")
```

**На входе**

**Значение** | **Описание**
-------------|--------------
  `key`      | type: string
  `value`    | type: interface

### .getVar(key) {#ic_flow_get_var}

Получить ранее записанную переменную

```coffeescript
wf = IC.Flow()
if wf
  variable = wf.getVar(key)
```

**На входе**

**Значение** | **Описание**
-------------|--------------
  `key`      | type: string

**На выходе**

**Значение** | **Описание**
-------------|--------------
  `variable` | type: interface
  
### .node() {#ic_flow_node}

Получить ноду

```coffeescript
flow = IC.Flow()
if flow
 node = flow.node()
 print 'flow node', node
```

**На выходе**

**Значение** | **Описание**
-------------|--------------
  `node`     | type: Object, ссылка на экземпляр [Node](#ic_node) 
