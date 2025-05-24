---
title: "Zigbee2mqtt"
linkTitle: "zigbee2mqtt"
date: 2021-10-20
description: >
  
---

Плагин "Zigbee2mqtt" является частью системы и предоставляет интеграцию между Zigbee-устройствами и MQTT-брокером. Этот
плагин позволяет контролировать и управлять Zigbee-устройствами через MQTT-сообщения.

В плагине "Zigbee2mqtt" реализован JavaScript-обработчик `zigbee2mqttEvent`, который не принимает параметры. Внутри
обработчика доступен объект `message`, представляющий информацию о полученном MQTT-сообщении.

Свойства объекта `message` включают:

1. `payload`: Значение сообщения, представленное в виде словаря (map), где ключи являются строками, а значения могут
   быть любого типа.
2. `topic`: Топик MQTT-сообщения, указывающий на источник или назначение сообщения.
3. `qos`: Уровень качества обслуживания (Quality of Service) MQTT-сообщения.
4. `duplicate`: Флаг, указывающий, является ли сообщение дубликатом.
5. `storage`: Объект `Storage`, предоставляющий доступ к хранилищу данных для кеширования и получения произвольных
   значений.
6. `error`: Строка, содержащая информацию об ошибке, если таковая возникла при обработке сообщения.
7. `success`: Булево значение, указывающее на успешное выполнение операции или обработку сообщения.
8. `new_state`: Объект `StateParams`, представляющий новое состояние актора после выполнения операции.

Пример использования обработчика `zigbee2mqttEvent`:

```javascript
function zigbee2mqttEvent(message) {
  console.log("Received MQTT message:");
  console.log("Payload:", message.payload);
  console.log("Topic:", message.topic);
  console.log("QoS:", message.qos);
  console.log("Duplicate:", message.duplicate);

  if (message.error) {
    console.error("Error:", message.error);
  } else if (message.success) {
    console.log("Operation successful!");
    console.log("New state:", message.new_state);
  }

  // Доступ к хранилищу
  const value = message.storage.getByName("key");
  console.log("Value from storage:", value);
}
```

Обработчик `zigbee2mqttEvent` может быть использован для обработки входящих MQTT-сообщений, а также для выполнения
дополнительных операций на основе полученных данных.

```coffeescript
zigbee2mqttEvent = (message)->
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

```coffeescript
zigbee2mqttEvent = (message)->
#print '---mqtt new event from button---'
  if !message
    return
  payload = unmarshal message.payload
  attrs =
    'battery': payload.battery
    'linkquality': payload.linkquality
    'voltage': payload.voltage
  state = ''
  if payload.action
    attrs.action = payload.action
    state = payload.action + "_action"
  if payload.click
    attrs.click = payload.click
    attrs.action = ""
    state = payload.click + "_click"
  EntitySetState ENTITY_ID,
    'new_state': state.toUpperCase()
    'attribute_values': attrs
```
