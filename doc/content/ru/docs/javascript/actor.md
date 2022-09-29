---
title: "Actor"
linkTitle: "actor"
date: 2021-11-20
description: >

---

Объект устройства

{{< alert color="warning" >}}Объект доступен не во всех областях видимости.{{< /alert >}}

----------------

### объект Actor

```coffeescript
# доступные метода
Actor
  .setState(entityStateParams)
  .getSettings() -> map[string]any
```

### объект entityStateParams 

|  значение | type  | описание  |
|-----------|--|---------|
| new_state | string | наименование нового состояния (не обязательное поле) |
| attribute_values | map[string]any | описание атрибутов |
| settings_value | map[string]any | описание настоек |
| storage_save | bool | признак записи в бд |

----------------

### пример кода

```coffeescript
settings = Actor.getSettings()
headers = [{'apikey':settings['apikey']}]
url = 'https://webhook.site/2692a589-a5bc-4156-af7d-75875578798f'
res = http.headers(headers).get(url)
```
