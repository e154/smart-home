---
weight: 30
title: mqtt
groups:
    - javascript
---

## Mqtt {#mqtt}

Объект mqtt клиента

```coffeescript
return if !Mqtt
```

**На выходе**

**Значение** | **Описание**
-------------|--------------
  `mqtt`     | type: Object, ссылка на экземпляр *Mqtt 
  
Доступные методы приведены далее:

### .Publish(topic, payload, qos, retain) {#mqtt_publish}

Отправить сообщение в канал

```coffeescript
if Mqtt
 Mqtt.publish(topic, payload, qos, retain)
```

**На выходе**

**Значение** | **Описание**
-------------|--------------
  `topic`    | type: string
  `payload`  | type: []byte
  `qos`      | type: uint8
  `retain`   | type: bool


