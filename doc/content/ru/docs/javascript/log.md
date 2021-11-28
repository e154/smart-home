
---
title: "Логирование"
linkTitle: "log"
date: 2021-11-19
description: >

---

Предостовляет возможность логировать события. Позволяет упростить отладку, и последующий анализ поведения системы.

{{< alert color="success" >}}Функция доступна в любом скрипте системы.{{< /alert >}}

----------------

### Логирование
```coffeescript
Log
  .info(txt)
  .warn(txt)
  .error(txt)
  .debug(txt)
```

----------------

### пример кода

```coffeescript
# Log
# ##################################

Log.info 'some text'
Log.warn 'some text'
Log.error 'some text'
Log.debug 'some text'
```
