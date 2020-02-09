---
weight: 30
title: mqtt
groups:
    - javascript
---

## IC.Mqtt() {#ic_mqtt}

Объект mqtt клиента

```coffeescript
mqtt = IC.Mqtt()
return if !mqtt
```

**На выходе**

**Значение** | **Описание**
-------------|--------------
  `mqtt`     | type: Object, ссылка на экземпляр *Mqtt 
  
Доступные методы приведены далее:

### .publish(topic, payload, qos, retain) {#ic_mqtt_publish}

Отправить сообщение в канал

```coffeescript
mqtt = IC.Mqtt()
if mqtt
 mqtt.publish(topic, payload, qos, retain)
```

**На выходе**

**Значение** | **Описание**
-------------|--------------
  `topic`    | type: string
  `payload`  | type: []byte
  `qos`      | type: uint8
  `retain`   | type: bool


