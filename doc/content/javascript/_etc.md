---
weight: 4
title: etc
groups:
    - javascript
---

## IC.CurrentNode() {#ic_current_node}

Получить текущую ноду. Выражение может вырнуть *null* если скрипт исполняется вне контекста [Workflow](#ic_workflow) 

```coffeescript
node = IC.CurrentNode()
print 'current node', node
```

**На выходе**

**Значение** | **Описание**
-------------|--------------
  `node`     | type: Object, ссылка на экземпляр [Node](#ic_node) 
  
## IC.CurrentDevice() {#ic_current_device}

Получить текущее утройство.

```coffeescript
device = IC.CurrentDevice()
print 'current device', device
```

**На выходе**

**Значение** | **Описание**
-------------|--------------
  `device`   | type: Object, ссылка на экземпляр [Device](#ic_device) 
  
## IC.Runmode {#ic_runmode}

Текущий режим, *prod* | *dev*

```coffeescript
runmode = IC.Runmode
print 'runmode', runmode
```

**На выходе**

**Значение** | **Описание**
-------------|--------------
 `runmode`   | type: string 