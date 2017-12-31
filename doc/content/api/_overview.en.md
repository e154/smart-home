---
weight: 1
title: overview
groups:
    - api
---

## Обзор {#overview}

**Smart home** **API** позволяет управлять системой **Smart home**, производить настройку конфигурации. Исполнять комманды на устройствах подключенных к системе.

**Smart home** **API** организован на основе **REST** запросов по верх HTTP протокола. 
Такая связка была выбрана исключительно из простоты поддержания, и разработки нового функционала. 
Все запросы должны быть в формате JSON и кодировке UTF-8 *(Content-Type: application/json; charset=utf-8)*. 
Заголовок Accept будут игнорироваться для всех запросов.

```bash
$ curl -i http://localhost:3000/api/v1/signin
```

## Безопасность {#security}

Сервер **Smart home** имеет систему контроля доступа основанную на ролях. Любое действие или команда по **API** будут проходить проверку на право доступа
Идентификация пользователя и роли происходит по [JWT](https://jwt.io) токену указанному в **access_token** заголовке запроса.
Все запросы кроме, запросов связанных с авторизацией, должны содержать в теле заглоловок **access_token**:

> Пример токена авторизации

```bash
access_token: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

# Ошибки {#errors}

При возникновении ошибки Вы получите её код в теле ответа. 

Поле | Тип | Значение
-----|-----|---------
`message` | строка | Текст ошибки
`status` | строка | Системное название

                                                                                                                                                                                                           
```bash
HTTP/1.1 404 Not Found
Date: Mon, 21 Apr 2014 13:26:48 GMT
Content-Type: application/json; charset=utf-8:
```

```json
{
    "message": "Описание ошибки",
    "status": "error"
}
```

Код ошибки | Описание                                                                                                                                                                                      
---------- | -------                                                                                                                                                                                       
`400` | Bad Request -- Плохой запрос                                                                                                                                                                       
`401` | Unauthorized -- **API** токен не верный                                                                                                                                                            
`403` | Forbidden -- Доступ запрещен, возможно не хватает прав                                                                                                                                             
`404` | Not Found -- Контент не найден                                                                                                                                                                     
`500` | Internal Server Error -- Ошибка на сервере                                                                                                                                                         
                                                                                                                                                                                                           
## HTTP методы {#http-methods}

Сервер **Smart home** **REST API** поддерживает данный список HTTP методов для каждого действия:

Метод | Действие
------|---------
`GET`| Получить ресурс
`POST`| Создать ресурс
`PUT`| Обновить ресурс
`DELETE`| Удалить ресурс

## Cross Origin Resource Sharing (CORS) {#cors}

Сервер **Smart home** API version 1.0 поддерживает Cross Origin Resource Sharing (CORS) запросы для AJAX.

```bash
$ curl -i http://localhost:3000/api/v1/signin
HTTP/1.1 200 OK
Access-Control-Allow-Credentials: true
Access-Control-Allow-Headers: Content-Type,access_token
Access-Control-Allow-Methods: GET,POST,PUT,DELETE,OPTIONS,UPDATE
Access-Control-Allow-Origin: *
```
