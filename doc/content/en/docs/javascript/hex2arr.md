---
title: "hex2arr"
linkTitle: "hex2arr"
date: 2021-11-20 
description: >

---

The `hex2arr(hexString)` function is available in a JavaScript environment, and it converts a hexadecimal string to a byte array (`[]byte`).

Here's an example implementation of the `hex2arr` function:

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
more examples:

```javascript
function hex2arr(hexString) {
    // Remove all spaces from the string
    hexString = hexString.replace(/\s/g, '');

    // Split the string into pairs of characters (each pair represents a byte)
    var hexPairs = hexString.match(/.{1,2}/g);

    // Convert each pair of characters to a numeric value
    var byteArr = hexPairs.map(function (hex) {
        return parseInt(hex, 16);
    });

    return byteArr;
}
```

You can use the `hex2arr` function as follows:

```javascript
var hexString = "48656C6C6F20576F726C64"; // Hexadecimal representation of "Hello World"
var byteArr = hex2arr(hexString);
console.log(byteArr); // [ 72, 101, 108, 108, 111, 32, 87, 111, 114, 108, 100 ]
```

In this example, the `hex2arr` function takes the hexadecimal string `"48656C6C6F20576F726C64"`, which represents "Hello World" in hexadecimal format. The result of the function is a byte array containing the corresponding ASCII values of the characters.
