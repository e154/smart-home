---
title: "unmarshal"
linkTitle: "unmarshal"
date: 2021-11-20
description: >

---

In the JavaScript environment, there is a `unmarshal` function available that converts a JSON string into an object using the `JSON.parse` function.

Example implementation of the `unmarshal` function:

```javascript
function unmarshal(j) {
  return JSON.parse(j);
}
```

Example of using the `unmarshal` function:

```javascript
var jsonStr = '{"name":"John","age":30,"city":"New York"}';
var obj = unmarshal(jsonStr);
console.log(obj); // { name: "John", age: 30, city: "New York" }
```

In this example, the `unmarshal` function takes the JSON string `jsonStr` and uses the `JSON.parse` function to convert the string into an object. The result of the function will be an object containing the data from the JSON string.
