---
title: "marshal"
linkTitle: "marshal"
date: 2021-11-20
description: >

---

В JavaScript окружении доступна функция `marshal`, которая преобразует объект в его JSON-представление с помощью
функции `JSON.stringify`.

Пример реализации функции `marshal`:

```javascript
function marshal(obj) {
  return JSON.stringify(obj);
}
```

Пример использования функции `marshal`:

```javascript
var obj = { name: "John", age: 30, city: "New York" };
var jsonStr = marshal(obj);
console.log(jsonStr); // {"name":"John","age":30,"city":"New York"}
```

В этом примере функция `marshal` принимает объект `obj` и использует функцию `JSON.stringify` для преобразования объекта
в его JSON-представление. Результатом выполнения функции будет строка JSON, содержащая сериализованные данные объекта.
