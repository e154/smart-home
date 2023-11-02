
---
title: "Bitmine"
linkTitle: "bitmine"
date: 2021-10-20
description: >

---

Working with Bitmine Antminer ASIC devices.

Supported devices:
* S9
* S7
* L3
* L3
* D3
* T9

### JavaScript Properties

----------------

### Miner Base Entity

```javascript
// Constant with the unique device ID
const ENTITY_ID;
```

```javascript
const result = Miner
    .stats()
    .devs()
    .summary()
    .pools()
    .addPool(url)
    .version()
    .enable(poolId)
    .disable(poolId)
    .delete(poolId)
    .switchPool(poolId)
    .restart();
```

|  Property  | Description  |
|-------------|---------|
| result |    Type: Object (Result)   |

|  Property  | Description  |
|-------------|---------|
| stats |   Type: Method, statistics    |
| devs |   Type: Method, statistics    |
| summary |   Type: Method, statistics    |
| pools |   Type: Method, server list   |
| addPool |   Type: Method, add server   |
| version |   Type: Method, firmware version  |
| enable |   Type: Method, enable server  |
| disable |   Type: Method, disable server  |
| delete |   Type: Method, delete server  |
| switchPool |   Type: Method, switch server  |
| restart |   Type: Method, soft device restart  |

----------------

### Result Object

```javascript
const result = {
  error: false,
  errMessage: "",
  result: ""
};
``` 
|  Property  | Description  |
|-------------|---------|
| error |    Type: boolean, error indicator   |
| errMessage |   Type: string, human-readable error message  |
| result | Type: string, JSON object in string format, if the request is completed without errors |

----------------

### entityAction Function

```javascript
entityAction = (entityId, actionName, args) => {
}
```
| Property  | Description  |
|-------------|---------|
| entityId |    Type: string, ID of the entity sending the message   |
| actionName |   Type: string, name of the action, in uppercase without the '/' symbol  |
| args | Type: map[string]any |

----------------

### EntityStateParams Object

```javascript
const entityStateParams = {
  "new_state": "",
  "attribute_values": {},
  "settings_value": {},
  "storage_save": true
};
``` 
|  Property  | Description  |
|-------------|---------|
| new_state |    Type: *string, name of the new state   |
| attribute_values |   Type: Object, new state of the entity's attributes |
| settings_value | Type: Object, new state of the entity's settings |
| storage_save | Type: boolean, indicator to save the new state to the database |

----------------

### Code Example

```coffeescript
# cgminer
# ##################################

ifError =(res)->
  return !res || res.error || res.Error

checkStatus =(args)->
  stats = Miner.stats()
  if ifError(stats)
    EntitySetState ENTITY_ID,
      'new_state': 'ERROR'
    return
  p = JSON.parse(stats.result)
  print p

checkSum =(args)->
  summary = Miner.summary()
  if ifError(summary)
    EntitySetState ENTITY_ID,
      'new_state': 'ERROR'
    return
  p = JSON.parse(summary.result)
  print p

entityAction = (entityId, actionName, args)->
  switch actionName
    when 'CHECK' then checkStatus(args)
    when 'SUM' then checkSum(args)
```
