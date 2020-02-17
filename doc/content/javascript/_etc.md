---
weight: 40
title: etc
groups:
    - javascript
---

## CurrentNode() {#current_node}

Получить текущую ноду. Выражение может вырнуть *null* если скрипт исполняется вне контекста [Workflow](#workflow) 

```coffeescript
node = CurrentNode()
print 'current node', node
```

**На выходе**

**Значение** | **Описание**
-------------|--------------
  `node`     | type: Object, ссылка на экземпляр [Node](#node) 
  
## CurrentDevice() {#current_device}

Получить текущее утройство. Выражение может вырнуть *null* если скрипт исполняется вне контекста [Workflow](#workflow)

```coffeescript
device = CurrentDevice()
print 'current device', device
```

**На выходе**

**Значение** | **Описание**
-------------|--------------
  `device`   | type: Object, ссылка на экземпляр [Device](#device) 
  
## Runmode {#runmode}

Текущий режим, *prod* | *dev*

```coffeescript
runmode = Runmode
print 'runmode', runmode
```

**На выходе**

**Значение** | **Описание**
-------------|--------------
 `runmode`   | type: string 
