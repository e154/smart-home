
---
title: "Выполнение системных команд"
linkTitle: "execute"
date: 2021-11-19
description: >

---

Запуск произвольных файлов/скриптов. Запрос происходит в синхронном режиме. 

{{< alert color="success" >}}Функция доступна в любом скрипте системы.{{< /alert >}}

----------------

### Вызов команды
```coffeescript
ExecuteSync(file, args)

ExecuteAsync(file, args)
```

|  значение  | описание  |
|-------------|---------|
| file |    путь до исполняемого файла  |
| args | не обязательные аргументы передаваемые файлу  |


----------------

### пример кода

```coffeescript
# ExecuteSync
# ##################################
"use strict";

r = ExecuteSync "data/scripts/ping.sh", "google.com"
if r.out == 'ok'
    print "site is available ^^"
```

ping.sh
```bash
#!/usr/bin/env bash

ping -c1 $1 > /dev/null 2> /dev/null; [[ $? -eq 0 ]] && echo ok || echo "err"
```
