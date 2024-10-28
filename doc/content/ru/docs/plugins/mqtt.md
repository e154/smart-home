---
title: "MQTT"
linkTitle: "mqtt"
date: 2021-10-20
description: >
  
---

Плагин "mqtt" является расширенной версией плагина "sensor" и предоставляет возможность работы с протоколом MQTT.
Вот некоторые настройки, доступные для плагина "mqtt":

1. `subscribe_topic`: Топик MQTT, на который выполняется подписка для получения данных.
2. `mqtt_login`: Логин для аутентификации при подключении к MQTT-брокеру.
3. `mqtt_pass`: Пароль для аутентификации при подключении к MQTT-брокеру.

Эти настройки позволяют задать параметры подключения к MQTT-брокеру и определить топик, на который будет выполняться
подписка
для получения данных. Таким образом, плагин "mqtt" обеспечивает интеграцию с MQTT-сообщениями, что позволяет получать
данные
с устройств, использующих этот протокол, и использовать их в системе автоматизации или других компонентах.

Плагин "sensor" также реализует JavaScript-обработчик (handler) под названием `entityAction`. Этот обработчик
предназначен
для обработки действий, связанных с устройствами типа "entity" на основе плагина "sensor".

Пример реализации обработчика `entityAction`:

```javascript
entityAction = (entityId, actionName, args) => {
  // Код обработки действия
};
```

Также реализует JavaScript-обработчик (handler) под названием `mqttEvent`. Этот обработчик предназначен
для обработки действий, связанных с устройствами типа "entity" на основе плагина "sensor".

Пример реализации обработчика `mqttEvent`:

```javascript
function mqttEvent(message) {
  // Код обработки действия
};
```

Пример использования обработчика mqttEvent:

```coffeescript
mqttEvent = (message)->
#print '---mqtt new event from plug---'
  if !message || message.topic.includes('/set')
    return
  payload = unmarshal message.payload
  attrs =
    'consumption': payload.consumption
    'linkquality': payload.linkquality
    'power': payload.power
    'state': payload.state
    'temperature': payload.temperature
    'voltage': payload.voltage
  EntitySetState ENTITY_ID,
    'new_state': payload.state
    'attribute_values': attrs
```
