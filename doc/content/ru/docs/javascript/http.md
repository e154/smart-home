
---
title: "http запросы"
linkTitle: "http"
date: 2021-10-20
description: >

---

В проекте **Smart Home** имеется возможность выполнения произвольных HTTP запросов синхронно к сторонним ресурсам. 

Объект `http` позволяет выполнять синхронные HTTP запросы к сторонним ресурсам, таким как API-сервисы,
и получать ответы. Вы можете использовать этот метод для интеграции с другими системами и получения или отправки данных
через HTTP протокол в вашем проекте **Smart Home**.

поддерживаемые методы:
* GET
* POST
* PUT
* DELETE

{{< alert color="success" >}}Функция доступна в любом скрипте системы.{{< /alert >}}

----------------

Для этого доступен соответствующий метод:

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
# auth
# ##################################

res = http.digestAuth('user','password').download(uri);

res = http.basicAuth('user','password').download(uri);

res = http.download(uri);


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
