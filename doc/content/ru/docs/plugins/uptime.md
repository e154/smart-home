
---
title: "Uptime"
linkTitle: "uptime"
date: 2021-10-20
description: >
  
---

{{< figure src="/smart-home/uptime1.png" >}}

&nbsp;

&nbsp;

В системе **Smart Home** присутствует плагин "uptime", который позволяет отображать параметры времени работы системы и приложения. Этот плагин предоставляет доступ к следующим параметрам:

1. `total`: Параметр `total` отображает общее время работы системы или приложения. Он указывает на количество времени, прошедшее с момента запуска.

2. `app_started`: Параметр `app_started` указывает на время запуска приложения. Он показывает дату и время, когда приложение было запущено или перезапущено.

Пример использования плагина "uptime" для получения параметров времени работы:

```javascript
const entity = entityManager.getEntity('uptime.uptime')
const uptimeParams = entity.getAttributes()
console.log(uptimeParams.total);
console.log(uptimeParams.app_started);
```

Вы можете использовать эти параметры для отображения информации о времени работы системы или приложения в вашем проекте **Smart Home**.