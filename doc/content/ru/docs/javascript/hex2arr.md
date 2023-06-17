---
title: "hex2arr"
linkTitle: "hex2arr"
date: 2021-11-20 
description: >

---

В JavaScript окружении доступна функция `hex2arr(hexString)`, которая преобразует строку шестнадцатеричных значений в
массив байт (`[]byte`).

```javascript
hex2arr = function (hexString) {
    var result = [];
    while (hexString.length >= 2) {
        result.push(parseInt(hexString.substring(0, 2), 16));
        hexString = hexString.substring(2, hexString.length);
    }
    return result;
};
```

Пример реализации функции `hex2arr`:

```javascript
function hex2arr(hexString) {
    // Удаляем все пробелы из строки
    hexString = hexString.replace(/\s/g, '');

    // Разбиваем строку на пары символов (каждая пара представляет байт)
    var hexPairs = hexString.match(/.{1,2}/g);

    // Преобразуем каждую пару символов в числовое значение
    var byteArr = hexPairs.map(function (hex) {
        return parseInt(hex, 16);
    });

    return byteArr;
}
```

Пример использования функции `hex2arr`:

```javascript
var hexString = "48656C6C6F20576F726C64"; // "Hello World" в шестнадцатеричном формате
var byteArr = hex2arr(hexString);
console.log(byteArr); // [ 72, 101, 108, 108, 111, 32, 87, 111, 114, 108, 100 ]
```

В этом примере функция `hex2arr` принимает строку `"48656C6C6F20576F726C64"`, которая представляет "Hello World" в
шестнадцатеричном формате. Результатом выполнения функции будет массив байт, содержащий соответствующие значения ASCII
символов.
