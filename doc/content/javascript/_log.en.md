---
weight: 3
title: log
groups:
    - javascript
---

## IC.Log {#ic_log}

Логирование состояния.

```coffeescript
log = IC.Log
```

**На выходе**

**Значение** | **Описание**
-------------|--------------
  `log`     | type: Object, ссылка на экземпляр *Log
  
Доступные методы приведены далее:

### .info() {#ic_log_info}

```coffeescript
IC.Log.info('Info message')
```

### .warn() {#ic_log_warn}

```coffeescript
IC.Log.warn('Warning message')
```

### .error() {#ic_log_error}

```coffeescript
IC.Log.error('Error message')
```

### .debug() {#ic_log_debug}

```coffeescript
IC.Log.debug('Debug message')
```