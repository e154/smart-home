---
title: "Storage"
linkTitle: "storage"
date: 2021-11-19
description: >

---

In the **Smart Home** project, there is an `in-memory` storage provided, which allows you to store and cache arbitrary
values in memory. The values are also periodically archived to disk every minute.

{{< alert color="success" >}}This function is available in any system script.{{< /alert >}}

The following methods are available through the "Storage" object for working with the storage:

1. `push(key, value)`: This method is used to add a value to the storage. You pass the `key` and `value` as arguments,
   which will be stored in memory. Example usage:

```javascript
Storage.push('temperature', 25.5);
```

2. `getByName(key)`: This method allows you to retrieve a value from the storage based on the specified `key`. If the
   value is found in memory, it will be returned. Example usage:

```javascript
const temperature = Storage.getByName('temperature');
console.log(temperature);
```

3. `search(key)`: This method performs a search for a value based on the `key`. It first searches for the value in
   memory, and if not found, it searches in the database. The result of the search is returned. Example usage:

```javascript
const result = Storage.search('temperature');
console.log(result);
```

4. `pop(key)`: This method removes a value from the storage based on the specified `key`. If the value is found and
   removed, the method returns the removed value. Example usage:

```javascript
const removedValue = Storage.pop('temperature');
console.log(removedValue);
```

The methods of the "Storage" object provide the ability to add, retrieve, search, and remove values from the `in-memory`
storage. You can use this storage for temporary data storage in memory, with the option to save to disk and perform
efficient value search.

----------------

### Working with the storage

```coffeescript
Storage
  .push(key, value)
  .getByName(key)
  .search(key)
  .pop(key)
```

| Value     | Description                                                    |
|-----------|----------------------------------------------------------------|
| push      | Store a value in the storage with the specified `key`          |
| getByName | Retrieve a value from the storage based on the specified `key` |
| search    | Search for a value based on the specified `key`                |
| pop       | Remove a record                                                |

----------------

### Code example

```coffeescript
# Storage
# ##################################

foo =
  'bar': 'bar'

value = JSON.stringify foo

# save var
Storage.push 'foo', value

# get exist var
value = Storage.getByName 'foo'

# search
list = Storage.search 'bar'

Storage.pop 'foo'
```
