---
weight: 20
title: workflow
groups:
    - javascript
---

## Workflow {#workflow}

Получить текущий Workflow. Выражение может вырнуть *null* если скрипт исполняется вне контекста Workflow 

**Синтаксис**

```coffeescript
workflow = Workflow
```

**На выходе**


**Значение** | **Описание**
-------------|--------------
  `workflow` | type: Object, ссылка на экземпляр *Workflow

Доступные методы приведены далее: 

### .GetName() {#workflow_get_name}

Получить наименование текущего Workflow

**Синтаксис**

```coffeescript
if Workflow
 name = Workflow.GetName()
 print 'wf name', name
```

### .GetDescription() {#workflow_get_description}

Получить описание текущего [Workflow](#workflow)

```coffeescript
if Workflow
 descr = Workflow.GetDescription()
 print 'wf descr', descr
```

### .SetVar(key, value) {#workflow_set_var}

Запомнить переменнную в хранилище [Workflow](#workflow). Хранилище позволяет 
сохрянять состояния на время жизни [Workflow](#workflow)  

```coffeescript
description = Workflow.SetVar(key, value)
```

**На входе**

**Значение** | **Описание**
-------------|--------------
  `key`      | type: string
  `value`    | type: interface

### .GetVar(key) {#workflow_get_var}

Получить ранее записанную переменную

```coffeescript
if Workflow
  variable = Workflow.GetVar(key)
```

**На входе**

**Значение** | **Описание**
-------------|--------------
  `key`      | type: string

**На выходе**

**Значение** | **Описание**
-------------|--------------
  `variable` | type: interface

### .GetScenario() {#workflow_get_scenario}

Получить активный сценарий для текущего [Workflow](#workflow)

```coffeescript
if Workflow
  scenario = Workflow.GetScenario()
```

**На выходе**

**Значение** | **Описание**
-------------|--------------
  `scenario` | type: string

### .SetScenario(name) {#workflow_set_scenario}

Переключить сценарий для текущего [Workflow](#workflow)

```coffeescript
if Workflow
  scenario = Workflow.SetScenario(name)
```

**На входе**

**Значение** | **Описание**
-------------|--------------
  `name`     | type: string
