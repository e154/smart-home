---
weight: 35
title: log
groups:
    - javascript
---

## Log {#log}

Логирование состояния.

```coffeescript
log = Log
```

**На выходе**

**Значение** | **Описание**
-------------|--------------
  `log`     | type: Object, ссылка на экземпляр *Log
  
Доступные методы приведены далее:

### .info() {#log_info}

```coffeescript
Log.Info('Info message')
```

### .warn() {#log_warn}

```coffeescript
Log.Warn('Warning message')
```

### .error() {#log_error}

```coffeescript
Log.Error('Error message')
```

### .debug() {#log_debug}

```coffeescript
Log.Debug('Debug message')
```
