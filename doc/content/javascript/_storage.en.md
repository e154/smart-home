---
weight: 30
title: storage
groups:
    - javascript
---

## Storage {#storage}

Быстрая база данных ключ-значение. Ключ и значение хранит в строковом виде. Для быстрого доступа данные кешируются в памяти.
Измененные значения сохраняются в БД, с определенным интервалом. 

### .Push(key, value) {#storage_push}

Сохранение объекта в БД

**На входе**

**Значение** | **Описание**
-------------|--------------
  `key`      | type: string
  `value`    | type: string

Пример: 

```coffeescript
foo = 
	'bar': 'bar'

value = JSON.stringify foo
Storage.Push 'foo', value
```

### .Pop(key) {#storage_pop}

Получить копию объекта и удалить из БД

**На входе**

**Значение** | **Описание**
-------------|--------------
  `key`      | type: string

**На выходе**

**Значение** | **Описание**
-------------|--------------
  `value`    | type: string

Пример:

```coffeescript
value = Storage.Pop 'foo'
foo = JSON.parse value
```

### .Search(key) {#storage_search}

Поиск в БД по ключу

```coffeescript
value = Storage.Search 'foo'
```

**На входе**

**Значение** | **Описание**
-------------|--------------
  `key`      | type: string

**На выходе**

**Значение** | **Описание**
-------------|--------------
  `value`    | type: map[string]string


### .GetByName(key) {#storage_get_by_name}

Получить объект по ключу

```coffeescript
value = Storage.GetByName 'foo'
foo = JSON.parse value
print foo
```

**На входе**

**Значение** | **Описание**
-------------|--------------
  `key`      | type: string

**На выходе**

**Значение** | **Описание**
-------------|--------------
  `value`    | type: string


