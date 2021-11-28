
---
title: "Временное хранилище"
linkTitle: "storage"
date: 2021-11-19
description: >

---

Предоставляет `in memory` хранилище. Для хранения/кеширования произвольного значения в памяти. Раз в минуту значение 
архивируется на диск. Поиск значения происходит сначала в пямяти, затем в БД. 

{{< alert color="success" >}}Функция доступна в любом скрипте системы.{{< /alert >}}

----------------

### Работа с хранилищем
```coffeescript
Storage
  .push(key, value)
  .getByName(key)
  .search(key)
  .pop(key)
```

|  значение  | описание  |
|-------------|---------|
| push |    поместить значение в хранилище по ключу `key`  |
| getByName | получить значение из хранилища по ключу `key` |
| search | поиск значения по ключу `key` |
| pop | удаление записи |

----------------

### пример кода

```coffeescript
# Storage
# ##################################

foo =
  'bar': 'bar'

value = JSON.stringify foo

# save var
Storage.push 'foo', value

# get exist var
value = Storage.getByName 'foo'

# search
list = Storage.search 'bar'

Storage.pop 'foo'
```
