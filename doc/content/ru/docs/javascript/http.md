
---
title: "http запросы"
linkTitle: "http"
date: 2021-10-20
description: >

---

Выполнение произвольного http запроса. Позволяет выполнять синхронные запросы к сторонним ресурсам.

поддерживаемые методы:
* GET
* POST
* PUT
* DELETE

{{< alert color="success" >}}Функция доступна в любом скрипте системы.{{< /alert >}}

----------------

### GET запрос 
```coffeescript
response = http.get(url)
```

|  значение  | описание  |
|-------------|---------|
| url |    адрес запроса   |
| response | рекзультат запроса |


----------------

### POST запрос 
```coffeescript
response = http.post(url, body)
```

|  значение  | описание  |
|-------------|---------|
| url |    адрес запроса   |
| body |    тело запроса   |
| response | рекзультат запроса |

----------------

### Headers запрос 
```coffeescript
response = http.headers(headers).post(url, body)
```

|  значение  | описание  |
|-------------|---------|
| headers |    адрес запроса   |
| url |    адрес запроса   |
| body |    тело запроса   |
| response | рекзультат запроса |

----------------

### пример кода

```coffeescript
# GET http
# ##################################

res = http.get("%s")
if res.error
  return
p = JSON.parse(res.body)


# POST http
# ##################################

res = http.post("%s", {'foo':'bar'})
if res.error
  return
p = JSON.parse(res.body)


# PUT http
# ##################################

res = http.put("%s", {'foo':'bar'})
if res.error
  return
p = JSON.parse(res.body)


# GET http + custom headers
# ##################################

res = http.headers([{'apikey':'some text'}]).get("%s")
if res.error
return
p = JSON.parse(res.body)

# DELETE http
# ##################################

res = http.delete("%s")
if res.error
  return
p = JSON.parse(res.body)

```
