# Авторизация

Авторизация - получения JWT токена, происходит `POST` запросом по паре email/пароль.

> пример запроса на авторизацию

```json
{
    "email": "user@server.name",
    "password": "******"
}
```

* Токен регистрируется на 24 часа, по истечении которых становится не действительным, затем снова требуется авторизация.
* Пароль должен состоять не менне чем из 6 символов.
* Первоночальная регистрация возможна только, по одной из двух учетных записей

> Учетные записи существующие после установки сервера  

```shell
email: demo@e154.ru
password: demodemo

email: admin@e154.ru
password: adminadmin
```

## Вход

> `POST` http://localhost/api/v1/signin

```shell
#!/usr/bin/env bash

EMAIL="admin@e154.ru"
PASSWORD="adminadmin"
URL=http://localhost:3000/api/v1/signin

curl -i -H "Content-Type: application/json" -X POST -d "{\"email\":\"${EMAIL}\",\"password\":\"${PASSWORD}\"}" ${URL}
```

Значение полей в результате положительного ответа 

Поле | Тип | Значение
-----|-----|---------
`token` | строка | JWT токен
`current_user` | объект | Параметры текущеко пользователя

> `200` положительный ответ

```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "current_user": {
    "avatar": null,
    "email": "admin@e154.ru",
    "first_name": "",
    "history": "[{\"ip\":\"127.0.0.1\",\"time\":\"2017-01-27T11:40:59.018428538+07:00\"},
    {\"ip\":\"127.0.0.1\",\"time\":\"2017-01-29T16:50:25.846523732+07:00\"}    
    "id": 1,
    "last_name": "",
    "meta": [
      {
        "id": 49,
        "key": "phone1",
        "value": ""
      },
      {
        "id": 50,
        "key": "phone2",
        "value": ""
      },
      {
        "id": 51,
        "key": "phone3",
        "value": ""
      },
      {
        "id": 52,
        "key": "autograph",
        "value": ""
      }
    ],
    "nickname": "admin",
    "role": {
      "name": "admin",
      "description": "",
      "parent": {
        "name": "user",
        "description": "",
        "parent": null,
        "children": null,
        "access_list": null,
        "created_at": "0001-01-01T00:00:00Z",
        "update_at": "0001-01-01T00:00:00Z"
      },
      "children": null,
      "access_list": {
        "dashboard": [
          "create",
          "update",
          "delete",
          "read"
        ],
        "device": [
          "update",
          "create",
          "delete",
          "read"
        ],
        "device_action": [
          "create",
          "read",
          "delete"
        ],
        "device_state": [
          "create",
          "delete",
          "update",
          "read"
        ],
        "flow": [
          "create",
          "delete",
          "update",
          "read"
        ],
        "image": [
          "read",
          "upload",
          "create",
          "delete",
          "update"
        ],
        "log": [
          "read"
        ],
        "map": [
          "read_map_layer",
          "read_map",
          "read_map_element"
        ],
        "node": [
          "read"
        ],
        "notifr": [
          "preview_notifr",
          "read_notifr_item",
          "read_notifr_template"
        ],
        "script": [
          "update",
          "create",
          "delete",
          "exec_script",
          "read"
        ],
        "user": [
          "read_role",
          "read_role_access_list",
          "read"
        ],
        "worker": [
          "delete",
          "create",
          "update",
          "read"
        ],
        "workflow": [
          "create",
          "delete",
          "read",
          "update"
        ],
        "ws": [
          "read"
        ]
      },
      "created_at": "2017-01-15T12:20:58+07:00",
      "update_at": "2017-01-15T12:21:29+07:00"
    },
    "sign_in_count": 121
  }
}
```

> `401` ошибка

```json
{
  "message": "Пользователь не найден",
  "status": "error"
}
```

## Выход

> `POST` http://localhost/api/v1/signout

