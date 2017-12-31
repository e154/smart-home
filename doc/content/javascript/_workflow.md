---
weight: 20
title: workflow
groups:
    - javascript
---

## IC.Workflow() {#ic_workflow}

Получить текущий Workflow. Выражение может вырнуть *null* если скрипт исполняется вне контекста Workflow 

**Синтаксис**

```coffeescript
workflow = IC.Workflow()
```

**На выходе**


**Значение** | **Описание**
-------------|--------------
  `workflow` | type: Object, ссылка на экземпляр *Workflow

Доступные методы приведены далее: 

### .getName() {#ic_workflow_get_name}

Получить наименование текущего Workflow

**Синтаксис**

```coffeescript
wf = IC.Workflow()
if wf
 name = wf.getName()
 print 'wf name', name
```

### .getDescription() {#ic_workflow_get_description}

Получить описание текущего [Workflow](#ic_workflow)

```coffeescript
wf = IC.Workflow()
if wf
 descr = wf.getDescription()
 print 'wf descr', descr
```

### .setVar(key, value) {#ic_workflow_set_var}

Запомнить переменнную в хранилище [Workflow](#ic_workflow). Хранилище позволяет 
сохрянять состояния на время жизни [Workflow](#ic_workflow)  

```coffeescript
description = IC.Workflow().setVar(key, value)
```

**На входе**

**Значение** | **Описание**
-------------|--------------
  `key`      | type: string
  `value`    | type: interface

### .getVar(key) {#ic_workflow_get_var}

Получить ранее записанную переменную

```coffeescript
wf = IC.Workflow()
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

### .getScenario() {#ic_workflow_get_scenario}

Получить активный сценарий для текущего [Workflow](#ic_workflow)

```coffeescript
wf = IC.Workflow()
if wf
  scenario = wf.getScenario()
```

**На выходе**

**Значение** | **Описание**
-------------|--------------
  `scenario` | type: string

### .setScenario(name) {#ic_workflow_set_scenario}

Переключить сценарий для текущего [Workflow](#ic_workflow)

```coffeescript
wf = IC.Workflow()
if wf
  scenario = wf.setScenario(name)
```

**На входе**

**Значение** | **Описание**
-------------|--------------
  `name`     | type: string
