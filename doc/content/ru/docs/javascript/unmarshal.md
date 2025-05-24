---
title: "unmarshal"
linkTitle: "unmarshal"
date: 2021-11-20
description: >

---

В JavaScript окружении доступна функция `unmarshal`, которая преобразует JSON-строку в объект с помощью
функции `JSON.parse`.

Пример реализации функции `unmarshal`:

```javascript
function unmarshal(j) {
  return JSON.parse(j);
}
```

Пример использования функции `unmarshal`:

```javascript
var jsonStr = '{"name":"John","age":30,"city":"New York"}';
var obj = unmarshal(jsonStr);
console.log(obj); // { name: "John", age: 30, city: "New York" }
```

В этом примере функция `unmarshal` принимает JSON-строку `jsonStr` и использует функцию `JSON.parse` для преобразования
строки в объект. Результатом выполнения функции будет объект, содержащий данные из JSON-строки.
