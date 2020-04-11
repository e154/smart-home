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

## DoAction(action_id) {#do_action}

Фукция, выполнить действие

**На входе**

**Значение**        | **Описание**
--------------------|--------------
  `action_id`       | type: number, id действия (action)

Пример:

```coffeescript
DoAction(1)
```
  
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
